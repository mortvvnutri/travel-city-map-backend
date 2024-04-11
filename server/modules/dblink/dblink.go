package dblink

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tcm/apitypes"
	"tcm/utils"

	"github.com/lib/pq"
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
	rows, err := db.db.Query(`SELECT
	id, name, parent_id, meta, created_at, updated_at 
	FROM categories ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := &[]apitypes.Category_Obj{}
	for rows.Next() {
		var val apitypes.Category_Obj
		rows.Scan(&val.Id, &val.Name, &val.ParentId, &val.Meta, &val.CreatedAt, &val.UpdatedAt)
		*ret = append(*ret, val)
	}
	return ret, nil
}

func (db *DBwrap) BuildRoute(pos_req *apitypes.PosReq_Obj) (*[]apitypes.Place_Obj, error) {
	// public, noauth, nopage

	// Basic validation
	if pos_req == nil || pos_req.MyLat == nil || pos_req.MyLong == nil {
		return nil, errors.New("required request parameters are not present")
	}

	ret := []apitypes.Place_Obj{}
	if pos_req.Cats == nil || len([]int32(*pos_req.Cats)) == 0 {
		// no route to build, return empty array
		return &ret, nil
	}

	local_cats := []int32(*pos_req.Cats)

	if len(local_cats) > 15 {
		// expensive operation, deny
		return nil, errors.New("maximum allowed points: 15")
	}

	var catset []int32
	seen_places := make(map[int32]bool)
	rolling_lat := *pos_req.MyLat
	rolling_long := *pos_req.MyLong

	slen := len(local_cats)
	for i := 0; i < slen; i++ {
		catset = utils.ToSetInt32(local_cats)
		rows, err := db.db.Query(`SELECT 
		id, name, description,
		lat, long, p_options,
		category_id, created_at, updated_at,
		meta 
		FROM places
		WHERE
		category_id = ANY($1::int[])
		ORDER BY distance(lat, long, $2::double precision, $3::double precision) ASC
		LIMIT 30`, pq.Array(catset), rolling_lat, rolling_long)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var val apitypes.Place_Obj
			rows.Scan(&val.Id, &val.Name, &val.Description,
				&val.Lat, &val.Long, &val.POptions,
				&val.CategoryId, &val.CreatedAt, &val.UpdatedAt,
				&val.Meta)
			if seen_places[int32(*val.Id)] {
				// the place already exists in the route
				continue
			} else {
				// the closest place is not yet in our route
				// add it to our route and mark as seen
				ret = append(ret, val)
				seen_places[int32(*val.Id)] = true
				break
			}
		}

		// edge case where no places could be found
		// on first iteration. Thus the route couldn't
		// be built
		if len(ret) == 0 {
			break
		}

		// Get the last added category
		last_found_cat := int32(*ret[len(ret)-1].CategoryId)
		// Remove the category ONCE from local copy
		local_cats = utils.RemoveSingleInt32(local_cats, last_found_cat)
		// redefine our position to the place's coordinates
		rolling_lat = *ret[len(ret)-1].Lat
		rolling_long = *ret[len(ret)-1].Long
	}

	return &ret, nil
}

func (db *DBwrap) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
