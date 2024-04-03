package dblink

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tcm/apitypes"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DBwrap_int interface {
	Init()
	Close()

	// API
	// PwdReset()
}

const (
	DB_FILE_PATH = "./storage.db"
)

type DBwrap struct {
	DBwrap_int
	db *sql.DB
}

type DBconfig struct {
	Host   *string
	Port   *string
	User   *string
	Pwd    *string
	Dbname *string
}

func (db *DBwrap) Init(cfg *DBconfig) error {

	if cfg == nil || cfg.Host == nil || cfg.Port == nil || cfg.User == nil || cfg.Pwd == nil || cfg.Dbname == nil {
		return errors.New("database configuration: missing parameters")
	}

	var err error
	db.db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", *cfg.Host, *cfg.Port, *cfg.User, *cfg.Pwd, *cfg.Dbname))

	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	err = db.db.Ping()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (db *DBwrap) GetUserInfoByEmail(email string, redact_pwd bool) (*apitypes.User_Obj, error) {
	uo := apitypes.User_Obj{}
	if email == "" {
		return &uo, errors.New("no username supplied")
	}

	err := db.db.QueryRow(`SELECT
	id, email, pic, yandex_id, pwd,
	preferred_cats, def_custom_place, display_name, meta,
	created_at, updated_at
	FROM users WHERE email=$1 LIMIT 1`, email).Scan(
		&uo.Id,
		&uo.Email,
		&uo.Pic,
		&uo.YandexId,
		&uo.Pwd,
		&uo.PreferredCats,
		&uo.DefCustomPlace,
		&uo.DisplayName,
		&uo.Meta,
		&uo.CreatedAt,
		&uo.UpdatedAt,
	)
	if err != nil {
		return &uo, errors.New("user does not exist")
	}

	if redact_pwd {
		uo.Pwd = nil
		uo.OldPwd = nil
	}

	return &uo, nil
}

func (db *DBwrap) CheckAuth(email string, pwd string) (*apitypes.User_Obj, error) {
	if email == "" || pwd == "" {
		return nil, errors.New("email or password is empty")
	}

	usr, err := db.GetUserInfoByEmail(email, false)

	if err != nil {
		return nil, errors.New("login credentials are incorrect")
	}

	if usr.Pwd == nil {
		return nil, errors.New("failed to check the password")
	}

	if CheckPasswordHash(pwd, *usr.Pwd) {
		// redact password manually before answering
		usr.Pwd = nil
		usr.OldPwd = nil
		return usr, nil
	}
	return nil, errors.New("login credentials are incorrect")
}

func (db *DBwrap) RegisterUser(initiator *apitypes.User_Obj) (*apitypes.User_Obj, error) {
	if initiator == nil || initiator.Email == nil || initiator.Pwd == nil || initiator.DisplayName == nil {
		return nil, errors.New("initiator email, password or display name is empty")
	}

	// hash password before continuing
	hashed, err := HashPassword(*initiator.Pwd)
	if err != nil {
		return nil, err
	}
	usr := &apitypes.User_Obj{}

	err = db.db.QueryRow(`INSERT INTO 
				users(email, pwd, preferred_cats, display_name, meta) VALUES ($1,$2,$3,$4,$5)
				RETURNING id, email, pic,
				preferred_cats, display_name, meta, created_at,
				updated_at`,
		initiator.Email, hashed, initiator.PreferredCats, initiator.DisplayName, initiator.Meta).Scan(
		&usr.Id, &usr.Email, &usr.Pic, &usr.PreferredCats, &usr.DisplayName, &usr.Meta, &usr.CreatedAt, &usr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (db *DBwrap) CatList() (*[]apitypes.Category_Obj, error) {
	// public, noauth, nopage
	rows, err := db.db.Query(`SELECT id, name, parent_id, meta, created_at, updated_at FROM categories`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := &[]apitypes.Category_Obj{}
	for rows.Next() {
		var val apitypes.Category_Obj
		rows.Scan(&val.Id, &val.Name)
		*ret = append(*ret, val)
	}
	return ret, nil
}

func (db *DBwrap) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
