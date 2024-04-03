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
	Id          *int       `json:"id,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Description *string    `json:"description,omitempty"`
	Lat         *float64   `json:"lat,omitempty"`
	Long        *float64   `json:"long,omitempty"`
	CategoryId  *int       `json:"category_id,omitempty"`
	Meta        *string    `json:"meta,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
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
