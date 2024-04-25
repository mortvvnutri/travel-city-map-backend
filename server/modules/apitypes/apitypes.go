package apitypes

import (
	"time"

	"github.com/lib/pq"
)

type ErrorStruct struct {
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type API_obj struct {
	Status *string `json:"status,omitempty"`
	Token  *string `json:"token,omitempty"`

	Error *ErrorStruct `json:"error,omitempty"`

	Initiator *User_Obj   `json:"user_initiator,omitempty"`
	User      *User_Obj   `json:"user,omitempty"`
	Users     *[]User_Obj `json:"users,omitempty"`

	Route  *Route_Obj   `json:"route,omitempty"`
	Routes *[]Route_Obj `json:"routes,omitempty"`

	RHJ  *Routes_History_Join   `json:"rhj,omitempty"`
	RHJs *[]Routes_History_Join `json:"rhjs,omitempty"`

	Feedback  *Feedback_Obj   `json:"feedback,omitempty"`
	Feedbacks *[]Feedback_Obj `json:"feedbacks,omitempty"`

	Place  *Place_Obj   `json:"place,omitempty"`
	Places *[]Place_Obj `json:"places,omitempty"`

	Category   *Category_Obj   `json:"category,omitempty"`
	Categories *[]Category_Obj `json:"categories,omitempty"`

	PosReq  *PosReq_Obj   `json:"pos_req,omitempty"`
	PosReqs *[]PosReq_Obj `json:"pos_reqs,omitempty"`

	Weather *OWM_Weather `json:"weather,omitempty"`

	File *File_Obj `json:"file,omitempty"`
}

type File_Obj struct {
	FullUrl *string `json:"full_url,omitempty"`
	Href    *string `json:"href,omitempty"`
	Name    *string `json:"name,omitempty"`
}

type User_Obj struct {
	Token *string `json:"token,omitempty"`

	Id             *int           `json:"id,omitempty"`
	DisplayName    *string        `json:"display_name,omitempty"`
	Email          *string        `json:"email,omitempty"`
	Pwd            *string        `json:"pwd,omitempty"`
	OldPwd         *string        `json:"oldpwd,omitempty"`
	Pic            *string        `json:"pic,omitempty"`
	YandexId       *string        `json:"yandex_id,omitempty"`
	PreferredCats  *pq.Int32Array `json:"preferred_cats,omitempty"`
	DefCustomPlace *int           `json:"def_custom_place,omitempty"`
	FirebaseToken  *string        `json:"firebase_token,omitempty"`
	Meta           *string        `json:"meta,omitempty"`
	CreatedAt      *time.Time     `json:"created_at,omitempty"`
	UpdatedAt      *time.Time     `json:"updated_at,omitempty"`
}

type Route_Obj struct {
	Id             *int           `json:"id,omitempty"`
	UserId         *int           `json:"user_id,omitempty"`
	Places         *pq.Int32Array `json:"places,omitempty"`
	Categories     *pq.Int32Array `json:"categories,omitempty"`
	TimesCompleted *int           `json:"times_completed,omitempty"`
	TotalDistance  *string        `json:"total_distance,omitempty"`
	StartP         *string        `json:"start_p,omitempty"`
	EndP           *string        `json:"end_p,omitempty"`
	DisplayName    *string        `json:"display_name,omitempty"`
	RouteData      *string        `json:"route_data,omitempty"`
	TimeTook       *string        `json:"time_took,omitempty"`
	Meta           *string        `json:"meta,omitempty"`
	CreatedAt      *time.Time     `json:"created_at,omitempty"`
	UpdatedAt      *time.Time     `json:"updated_at,omitempty"`
}

type Routes_History_Join struct {
	Id        *int       `json:"id,omitempty"`
	UserId    *int       `json:"user_id,omitempty"`
	RouteId   *int       `json:"route_id,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Feedback_Obj struct {
	UserId    *int       `json:"user_id,omitempty"`
	PlaceId   *int       `json:"place_id,omitempty"`
	Comment   *string    `json:"comment,omitempty"`
	Rating    *float64   `json:"rating,omitempty"`
	Meta      *string    `json:"meta,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type Place_Obj struct {
	Id          *int           `json:"id,omitempty"`
	Name        *string        `json:"name,omitempty"`
	Description *string        `json:"description,omitempty"`
	Lat         *float64       `json:"lat,omitempty"`
	Long        *float64       `json:"long,omitempty"`
	CategoryId  *int           `json:"category_id,omitempty"`
	POptions    *pq.Int32Array `json:"p_options,omitempty"`
	Meta        *string        `json:"meta,omitempty"`
	CreatedAt   *time.Time     `json:"created_at,omitempty"`
	UpdatedAt   *time.Time     `json:"updated_at,omitempty"`
}

type Category_Obj struct {
	Id        *int       `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	ParentId  *int       `json:"parent_id,omitempty"`
	Meta      *string    `json:"meta,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type CustomPlace_Obj struct {
	Id        *int       `json:"id,omitempty"`
	UserId    *int       `json:"user_id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Lat       *float64   `json:"lat,omitempty"`
	Long      *float64   `json:"long,omitempty"`
	Meta      *string    `json:"meta,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type PosReq_Obj struct {
	MyLat  *float64       `json:"my_lat,omitempty"`
	MyLong *float64       `json:"my_long,omitempty"`
	Radius *float64       `json:"radius,omitempty"`
	Cats   *pq.Int32Array `json:"cats,omitempty"`
}

// OpenWeatherMap

type OWM_CFG struct {
	ApiKey *string `json:"api_key,omitempty"`
}

type OWM_Weather struct {
	Coord      *OWM_W_Coord         `json:"coord,omitempty"`
	Weather    *[]OWM_W_Weather     `json:"weather,omitempty"`
	Base       *string              `json:"base,omitempty"`
	Main       *OWM_W_Main          `json:"main,omitempty"`
	Visibility *float32             `json:"visibility,omitempty"`
	Wind       *OWM_W_Wind          `json:"wind,omitempty"`
	Clouds     *OWM_W_Clouds        `json:"clouds,omitempty"`
	Rain       *OWM_W_Precipitation `json:"rain,omitempty"`
	Snow       *OWM_W_Precipitation `json:"snow,omitempty"`
	Dt         *int64               `json:"dt,omitempty"`
	Sys        *OWM_W_Sys           `json:"sys,omitempty"`
	Timezone   *int                 `json:"timezone,omitempty"`
	Id         *int64               `json:"id,omitempty"`
	Name       *string              `json:"name,omitempty"`
	Cod        *int64               `json:"cod,omitempty"`
}

type OWM_W_Sys struct {
	Type    *int64  `json:"type,omitempty"`
	Id      *int64  `json:"id,omitempty"`
	Country *string `json:"country,omitempty"`
	Sunrise *int64  `json:"sunrise,omitempty"`
	Sunset  *int64  `json:"sunset,omitempty"`
}

type OWM_W_Clouds struct {
	All *int `json:"all,omitempty"`
}

type OWM_W_Precipitation struct {
	H1  *float32 `json:"1h,omitempty"`
	H2  *float32 `json:"2h,omitempty"`
	H3  *float32 `json:"3h,omitempty"`
	H4  *float32 `json:"4h,omitempty"`
	H5  *float32 `json:"5h,omitempty"`
	H6  *float32 `json:"6h,omitempty"`
	H7  *float32 `json:"7h,omitempty"`
	H8  *float32 `json:"8h,omitempty"`
	H9  *float32 `json:"9h,omitempty"`
	H10 *float32 `json:"10h,omitempty"`
	H11 *float32 `json:"11h,omitempty"`
	H12 *float32 `json:"12h,omitempty"`
	H13 *float32 `json:"13h,omitempty"`
	H14 *float32 `json:"14h,omitempty"`
	H15 *float32 `json:"15h,omitempty"`
	H16 *float32 `json:"16h,omitempty"`
	H17 *float32 `json:"17h,omitempty"`
	H18 *float32 `json:"18h,omitempty"`
	H19 *float32 `json:"19h,omitempty"`
	H20 *float32 `json:"20h,omitempty"`
	H21 *float32 `json:"21h,omitempty"`
	H22 *float32 `json:"22h,omitempty"`
	H23 *float32 `json:"23h,omitempty"`
	H24 *float32 `json:"24h,omitempty"`
}

type OWM_W_Wind struct {
	Speed *float32 `json:"speed,omitempty"`
	Deg   *float32 `json:"deg,omitempty"`
	Gust  *float32 `json:"gust,omitempty"`
}

type OWM_W_Main struct {
	Temp        *float32 `json:"temp,omitempty"`
	FeelsLike   *float32 `json:"feels_like,omitempty"`
	Pressure    *float32 `json:"pressure,omitempty"`
	Humidity    *float32 `json:"humidity,omitempty"`
	TempMin     *float32 `json:"temp_min,omitempty"`
	TempMax     *float32 `json:"temp_max,omitempty"`
	SeaLevel    *float32 `json:"sea_level,omitempty"`
	GroundLevel *float32 `json:"grnd_level,omitempty"`
}

type OWM_W_Weather struct {
	Id          *int64  `json:"id,omitempty"`
	Main        *string `json:"main,omitempty"`
	Description *string `json:"description,omitempty"`
	Icon        *string `json:"icon,omitempty"`
}

type OWM_W_Coord struct {
	Lat *float64 `json:"lat,omitempty"`
	Lon *float64 `json:"lon,omitempty"`
}
