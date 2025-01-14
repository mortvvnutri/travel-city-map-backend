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
	KM           = 0.012685281829263324
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

func (db *DBwrap) ChangeName(initiator *apitypes.User_Obj, to_profile *apitypes.User_Obj) (*apitypes.User_Obj, error) {
	// quick validation
	if to_profile == nil || to_profile.DisplayName == nil || initiator == nil || initiator.Id == nil {
		return nil, errors.New("missing required parametes")
	}
	if len(*to_profile.DisplayName) > 64 || len(*to_profile.DisplayName) < 2 {
		return nil, errors.New("invalid name length")
	}
	uo := &apitypes.User_Obj{}

	err := db.db.QueryRow(`UPDATE users SET display_name=$1 WHERE id=$2 
	RETURNING id, email, pic, preferred_cats, def_custom_place, display_name, meta, created_at, updated_at `,
		to_profile.DisplayName, initiator.Id).Scan(
		&uo.Id,
		&uo.Email,
		&uo.Pic,
		&uo.PreferredCats,
		&uo.DefCustomPlace,
		&uo.DisplayName,
		&uo.Meta,
		&uo.CreatedAt,
		&uo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return uo, nil
}

func (db *DBwrap) ChangePic(initiator *apitypes.User_Obj, to_profile *apitypes.User_Obj) (*apitypes.User_Obj, error) {
	// quick validation
	if to_profile == nil || to_profile.Pic == nil || initiator == nil || initiator.Id == nil {
		return nil, errors.New("missing required parametes")
	}
	if len(*to_profile.Pic) > 128 {
		return nil, errors.New("invalid url length")
	}
	uo := &apitypes.User_Obj{}

	err := db.db.QueryRow(`UPDATE users SET pic=$1 WHERE id=$2 
	RETURNING id, email, pic, preferred_cats, def_custom_place, display_name, meta, created_at, updated_at `,
		to_profile.Pic, initiator.Id).Scan(
		&uo.Id,
		&uo.Email,
		&uo.Pic,
		&uo.PreferredCats,
		&uo.DefCustomPlace,
		&uo.DisplayName,
		&uo.Meta,
		&uo.CreatedAt,
		&uo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return uo, nil
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

/*
db.CatList() - Returns a list of categories that exist in the database
*/
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

// Builds a route from the selected categories and a current user position
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

func (db *DBwrap) PlacesNearby(pos_req *apitypes.PosReq_Obj) (*[]apitypes.Place_Obj, error) {
	// public, noauth, nopage

	// Basic validation
	if pos_req == nil || pos_req.MyLat == nil || pos_req.MyLong == nil {
		return nil, errors.New("required request parameters are not present")
	}

	ret := []apitypes.Place_Obj{}

	rows, err := db.db.Query(`WITH pre AS (
		SELECT
		id, name, description,
		lat, long, p_options,
		category_id, created_at, updated_at,
		meta, rating,
		distance(lat, long, $1::double precision, $2::double precision) dist
		FROM places
		ORDER BY dist ASC
		LIMIT 100
	)
	SELECT * from pre
	WHERE dist<=$3::double precision`, pos_req.MyLat, pos_req.MyLong, KM)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var val apitypes.Place_Obj
		rows.Scan(&val.Id, &val.Name, &val.Description,
			&val.Lat, &val.Long, &val.POptions,
			&val.CategoryId, &val.CreatedAt, &val.UpdatedAt,
			&val.Meta, &val.Rating,
			&val.Distance,
		)
		ret = append(ret, val)
	}

	return &ret, nil
}

func (db *DBwrap) SaveRoute(initiator *apitypes.User_Obj, route *apitypes.Route_Obj) (*apitypes.Route_Obj, error) {
	if initiator == nil || initiator.Id == nil || route == nil || route.DisplayName == nil {
		return nil, errors.New("missing required parameters")
	}
	ret := &apitypes.Route_Obj{}

	err := db.db.QueryRow(`INSERT INTO
		routes(user_id, places, categories,
			 total_distance, start_p, end_p,
			 display_name, route_data, time_took,
			 meta)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING
		id, user_id, places,
		categories, times_completed, total_distance,
		start_p, end_p, display_name,
		route_data, time_took, meta,
		created_at, updated_at`,
		initiator.Id, route.Places, route.Categories,
		route.TotalDistance, route.StartP, route.EndP,
		route.DisplayName, route.RouteData, route.TimeTook,
		route.Meta,
	).Scan(
		&ret.Id, &ret.UserId, &ret.Places,
		&ret.Categories, &ret.TimesCompleted, &ret.TotalDistance,
		&ret.StartP, &ret.EndP, &ret.DisplayName,
		&ret.RouteData, &ret.TimeTook, &ret.Meta,
		&ret.CreatedAt, &ret.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *DBwrap) ListRoutes(initiator *apitypes.User_Obj) (*[]apitypes.Route_Obj, error) {
	if initiator == nil || initiator.Id == nil {
		return nil, errors.New("missing required parameters")
	}
	ret := &[]apitypes.Route_Obj{}

	rows, err := db.db.Query(`SELECT
	id, user_id, places, categories,
	times_completed, total_distance, start_p, end_p,
	display_name, route_data, time_took, meta,
	created_at, updated_at
	FROM routes
	WHERE user_id=$1
	ORDER BY updated_at DESC`, initiator.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var val apitypes.Route_Obj
		rows.Scan(&val.Id, &val.UserId, &val.Places, &val.Categories,
			&val.TimesCompleted, &val.TotalDistance, &val.StartP, &val.EndP,
			&val.DisplayName, &val.RouteData, &val.TimeTook, &val.Meta,
			&val.CreatedAt, &val.UpdatedAt)
		*ret = append(*ret, val)
	}

	return ret, nil
}

func (db *DBwrap) CompleteRoute(initiator *apitypes.User_Obj, route *apitypes.Route_Obj) (*apitypes.Route_Obj, error) {
	if initiator == nil || initiator.Id == nil || route == nil || route.Id == nil {
		return nil, errors.New("missing required parameters")
	}
	ret := &apitypes.Route_Obj{}

	err := db.db.QueryRow(`UPDATE routes 
		SET times_completed = times_completed+1
		WHERE user_id=$1 AND id=$2
		RETURNING
		id, user_id, places,
		categories, times_completed, total_distance,
		start_p, end_p, display_name,
		route_data, time_took, meta,
		created_at, updated_at`,
		initiator.Id, route.Id,
	).Scan(
		&ret.Id, &ret.UserId, &ret.Places,
		&ret.Categories, &ret.TimesCompleted, &ret.TotalDistance,
		&ret.StartP, &ret.EndP, &ret.DisplayName,
		&ret.RouteData, &ret.TimeTook, &ret.Meta,
		&ret.CreatedAt, &ret.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (db *DBwrap) CreateCustomPlace(initiator *apitypes.User_Obj, place *apitypes.CustomPlace_Obj) (*apitypes.CustomPlace_Obj, error) {
	if initiator == nil || initiator.Id == nil {
		return nil, errors.New("authorization is required")
	}

	if place == nil || place.Name == nil || place.Lat == nil || place.Lat == nil {
		return nil, errors.New("missing required parameters")
	}

	ret := apitypes.CustomPlace_Obj{}
	err := db.db.QueryRow(`INSERT INTO
		custom_places(user_id, name, lat,
			 long, meta)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
		id, user_id, name,
		lat, long, meta,
		created_at, updated_at;`,
		initiator.Id, place.Name, place.Lat,
		place.Long, place.Meta,
	).Scan(
		&ret.Id, &ret.UserId, &ret.Name,
		&ret.Lat, &ret.Long, &ret.Meta,
		&ret.CreatedAt, &ret.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (db *DBwrap) ListMyCustomPlaces(initiator *apitypes.User_Obj) (*[]apitypes.CustomPlace_Obj, error) {
	if initiator == nil || initiator.Id == nil {
		return nil, errors.New("missing required parameters")
	}
	ret := &[]apitypes.CustomPlace_Obj{}

	rows, err := db.db.Query(`SELECT
	id, user_id, name,
	lat, long, meta,
	created_at, updated_at
	FROM custom_places
	WHERE user_id=$1
	ORDER BY updated_at DESC`, initiator.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var val apitypes.CustomPlace_Obj
		rows.Scan(&val.Id, &val.UserId, &val.Name,
			&val.Lat, &val.Long, &val.Meta,
			&val.CreatedAt, &val.UpdatedAt)
		*ret = append(*ret, val)
	}

	return ret, nil
}

func (db *DBwrap) DeleteCustomPlace(initiator *apitypes.User_Obj, place *apitypes.CustomPlace_Obj) error {
	if initiator == nil || initiator.Id == nil {
		return errors.New("authorization is required")
	}

	if place == nil || place.Id == nil {
		return errors.New("missing required parameters")
	}
	_, err := db.db.Exec(`DELETE FROM
		custom_places
		WHERE
		user_id=$1 AND id=$2
		`,
		initiator.Id, place.Id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *DBwrap) SetDefaultPlace(initiator *apitypes.User_Obj, place *apitypes.CustomPlace_Obj) (*apitypes.User_Obj, error) {
	if initiator == nil || initiator.Id == nil {
		return nil, errors.New("authorization is required")
	}

	if place == nil || place.Id == nil {
		return nil, errors.New("missing required parameters")
	}

	// Check if the place actually belongs to out user and yeet them otherwise
	tgt_id := -1
	err := db.db.QueryRow(`SELECT user_id FROM custom_places WHERE id=$1`, place.Id).Scan(&tgt_id)
	if err != nil {
		return nil, err
	}

	if tgt_id != *initiator.Id {
		return nil, errors.New("you do not have required permissions to perform that action")
	}

	uo := &apitypes.User_Obj{}

	err = db.db.QueryRow(`UPDATE users SET def_custom_place=$1 WHERE id=$2 
	RETURNING id, email, pic, preferred_cats, def_custom_place, display_name, meta, created_at, updated_at `,
		place.Id, initiator.Id).Scan(
		&uo.Id,
		&uo.Email,
		&uo.Pic,
		&uo.PreferredCats,
		&uo.DefCustomPlace,
		&uo.DisplayName,
		&uo.Meta,
		&uo.CreatedAt,
		&uo.UpdatedAt,
	)
	return uo, err
}

func (db *DBwrap) MyReviews(initiator *apitypes.User_Obj) (*[]apitypes.Review_Obj, error) {
	if initiator == nil || initiator.Id == nil {
		return nil, errors.New("authorization is required")
	}

	ret := &[]apitypes.Review_Obj{}

	rows, err := db.db.Query(`SELECT 
	user_id, place_id,
	reviews.rating, comment,
	reviews.created_at, reviews.updated_at,
	places.id, places.name,
	places.description, places.lat,
	places.long, places.rating
	FROM reviews
	LEFT JOIN places on places.id=reviews.place_id
	WHERE reviews.user_id=$1
	ORDER BY reviews.updated_at DESC`, initiator.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var val apitypes.Review_Obj
		val.Place = &apitypes.Place_Obj{}
		rows.Scan(&val.UserId, &val.PlaceId,
			&val.Rating, &val.Comment,
			&val.CreatedAt, &val.UpdatedAt,
			&val.Place.Id, &val.Place.Name,
			&val.Place.Description, &val.Place.Lat,
			&val.Place.Long, &val.Place.Rating,
		)
		*ret = append(*ret, val)
	}

	return ret, nil
}

func (db *DBwrap) CreateReview(initiator *apitypes.User_Obj, review *apitypes.Review_Obj, pos_req *apitypes.PosReq_Obj) (*apitypes.Review_Obj, error) {
	if initiator == nil || initiator.Id == nil {
		return nil, errors.New("authorization is required")
	}

	if review == nil || review.PlaceId == nil || review.Rating == nil || pos_req == nil || pos_req.MyLat == nil || pos_req.MyLong == nil {
		return nil, errors.New("missing required parameters")
	}

	if *review.Rating < 0 || *review.Rating > 5 {
		return nil, errors.New("review score is invalid")
	}

	// check distance
	dist := 100000.0
	err := db.db.QueryRow(`WITH pl AS (
		SELECT lat, long FROM places pl WHERE id=$1
	)
	SELECT distance($2::real, $3::real, lat, long) FROM pl`, review.PlaceId, pos_req.MyLat, pos_req.MyLong).Scan(&dist)
	if err != nil {
		return nil, err
	}

	if dist > KM {
		return nil, errors.New("you are too far away from the place to leave a review")
	}

	// Check if the review already exists
	chk_rate := -1
	err = db.db.QueryRow(`SELECT COUNT(rating) FROM reviews WHERE user_id=$1 AND place_id=$2`, initiator.Id, review.PlaceId).Scan(&chk_rate)
	if err != nil {
		return nil, err
	}

	ret := &apitypes.Review_Obj{}

	if chk_rate >= 1 {
		// update
		err = db.db.QueryRow(`UPDATE reviews SET rating=$1, comment=$2 WHERE user_id=$3 AND place_id=$4 
		RETURNING rating, comment, created_at, updated_at, user_id, place_id `,
			review.Rating, review.Comment, initiator.Id, review.PlaceId).Scan(
			&ret.Rating,
			&ret.Comment,
			&ret.CreatedAt,
			&ret.UpdatedAt,
			&ret.UserId,
			&ret.PlaceId,
		)
		if err != nil {
			return nil, err
		}
	} else {
		// insert
		err = db.db.QueryRow(`INSERT INTO reviews (rating, comment, user_id, place_id) VALUES($1, $2, $3, $4)
		RETURNING rating, comment, created_at, updated_at, user_id, place_id `,
			review.Rating, review.Comment, initiator.Id, review.PlaceId).Scan(
			&ret.Rating,
			&ret.Comment,
			&ret.CreatedAt,
			&ret.UpdatedAt,
			&ret.UserId,
			&ret.PlaceId,
		)
		if err != nil {
			return nil, err
		}
	}

	return ret, err
}

func (db *DBwrap) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}
