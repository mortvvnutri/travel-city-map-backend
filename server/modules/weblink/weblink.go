package weblink

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"tcm/apitypes"

	"tcm/dblink"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/gorilla/mux"
)

type CustomClaims struct {
	UserInfo apitypes.User_Obj `json:"user_info"`
	Login    string            `json:"login"`
	Note     string            `json:"note"`
	jwt.StandardClaims
}

var pb []string = []string{
	"The machine with a base-plate of prefabulated aluminite, surmounted by a malleable logarithmic casing in such a way that the two main spurving bearings were in a direct line with the pentametric fan",
	"IKEA battery supplies",
	"Probably not you...",
	"php 4.0.1",
	"The smallest brainfuck interpreter written using Piet",
	"8192 monkeys with typewriters",
	"16 dumplings and one chicken nuggie",
	"Imaginary cosmic duck",
	"13 space chickens",
	" // TODO: Fill this field in",
	"Marshmallow on a stick",
	"Two sticks and a duct tape",
	"Multipolygonal eternal PNGs",
	"40 potato batteries. Embarrassing. Barely science, really.",
	"Aperture Science computer-aided enrichment center",
	"A cluster*** of protogens",
	"Fifteen Hundred Megawatt Aperture Science Heavy Duty Super-Colliding Super Button",
}

var JWT_PRIV_KEY *[]byte = &[]byte{}
var dbl = dblink.DBwrap{}

const (
	WEB_CRYPTO_KEY_PATH = "./key.mkey"
)

func checkTokenAndGetInfo(initiator *apitypes.User_Obj) (*apitypes.User_Obj, error) {
	if initiator == nil || initiator.Token == nil {
		return nil, errors.New("initiator token is not present")
	}
	return verifyJWT(initiator.Token)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	// _, err := checkTokenAndGetInfo(r)
	// if err != nil {
	// 	http.Redirect(w, r, "/web/login", http.StatusSeeOther)
	// 	return
	// }
	// http.Redirect(w, r, "/web/dashboard", http.StatusSeeOther)
	// ThrowApiErr()
	return
}

func inarr(s1 string, arr []string) bool {
	for _, v := range arr {
		if v == s1 {
			return true
		}
	}
	return false
}

// Blackhole
func denyIncoming(w http.ResponseWriter, r *http.Request) {
	rd, e := rand.Int(rand.Reader, big.NewInt(int64(len(pb))))
	if e != nil {
		rd = big.NewInt(int64(0))
	}
	w.Header().Add("X-Powered-By", pb[rd.Int64()])
	w.Header().Add("content-type", "text/plain")
	// w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Key, API-usr, Token, ref-key, lu-key")
	w.WriteHeader(403)
	fmt.Fprintf(w, "403: Access denied")
}

/*
Answers preflight OPTIONS requests
*/
func preflight(w http.ResponseWriter, r *http.Request) {
	rd, e := rand.Int(rand.Reader, big.NewInt(int64(len(pb))))
	if e != nil {
		rd = big.NewInt(int64(0))
	}
	w.Header().Add("X-Powered-By", pb[rd.Int64()])
	w.Header().Add("content-type", "text/plain")
	w.Header().Add("access-control-allow-origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Key, API-usr, Token, ref-key, lu-key")
	w.WriteHeader(204)
	fmt.Fprintf(w, "204: Access denied, but with love to the poor browser that for some reason wanted to access this page.")
}

/*
A global entry to the API endpoints.

Routes requests to their corresponding functions and actions
*/
func apiGlobalRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		preflight(w, r)
		return
	}

	vars := mux.Vars(r)
	group, ok := vars["group"]
	if !ok {
		fmt.Println("group is missing in parameters")
		denyIncoming(w, r)
		return
	}
	endpoint, ok := vars["endpoint"]
	if !ok {
		fmt.Println("endpoint is missing in parameters")
		denyIncoming(w, r)
		return
	}
	operation, ok := vars["operation"]
	if !ok {
		fmt.Println("operation is missing in parameters")
		denyIncoming(w, r)
		return
	}
	fmt.Println(endpoint + ":" + operation)

	w.Header().Set("access-control-allow-origin", "*")

	if r.Method != "POST" {
		denyIncoming(w, r)
		return
	}
	w.Header().Add("content-type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		apiRespond(w, &apitypes.API_obj{Error: &apitypes.ErrorStruct{Error: err.Error(), Message: "Could not read request body", Code: 500}})
		return
	}
	var apireq apitypes.API_obj
	err = json.Unmarshal(body, &apireq)
	if err != nil {
		// invalid params
		apiRespond(w, &apitypes.API_obj{Error: &apitypes.ErrorStruct{Error: err.Error(), Message: "Could not parse JSON from the request body", Code: 500}})
		return
	}

	switch group {
	case "public":
		apiPublic(w, r, endpoint, operation, &apireq)
	default:
		fmt.Println("invalid API group")
		ThrowApiErr(w, "Invalid API group", nil, 400)
	}

	return
}

/*
A request router for public API

Routes via endpoint parameter and passes through the operation to the
routed function
*/
func apiPublic(w http.ResponseWriter, r *http.Request, endpoint string, operation string, apireq *apitypes.API_obj) {
	switch endpoint {
	case "user":
		apiPubUser(w, r, operation, apireq)
	case "category":
		apiPubCats(w, r, operation, apireq)
	case "route":
		apiPubRoute(w, r, operation, apireq)
	default:
		fmt.Println("invalid API endpoint")
		ThrowApiErr(w, "Invalid API endpoint", nil, 400)
	}
}

/*
Defines a router for public user actions.

Acts according to the operation provided
*/
func apiPubUser(w http.ResponseWriter, r *http.Request, operation string, apireq *apitypes.API_obj) {
	var uo *apitypes.User_Obj
	var err error
	if operation != "login" && operation != "register" {
		uo, err = checkTokenAndGetInfo(apireq.Initiator)
		if err != nil {
			ThrowApiErr(w, "token is invalid", nil, 403)
			return
		}
	}
	switch operation {
	case "register":
		if apireq == nil || apireq.Initiator == nil || apireq.Initiator.Email == nil || apireq.Initiator.Pwd == nil || apireq.Initiator.DisplayName == nil {
			ThrowApiErr(w, "initiator email, password and display name must be present", nil, 401)
			return
		}

		if len(*apireq.Initiator.Email) < 5 {
			ThrowApiErr(w, "Email is invalid", nil, 400)
			return
		}
		if len(*apireq.Initiator.Pwd) < 8 {
			ThrowApiErr(w, "Password is too short", nil, 400)
			return
		}
		if len(*apireq.Initiator.DisplayName) < 2 {
			ThrowApiErr(w, "Name is too short", nil, 400)
			return
		}

		newuser, err := dbl.RegisterUser(apireq.Initiator)
		if err != nil {
			ThrowApiErr(w, "Email already exists", err, 403)
			return
		}

		dt := time.Hour * 24 * 3

		tk, err := genJWTtokenByUserObj(newuser, dt)
		if err != nil {
			ThrowApiErr(w, "Failed to generate token", err, 500)
			return
		}
		resp := apitypes.API_obj{Token: &tk, User: newuser}
		// w.Header().Set("Set-Cookie", fmt.Sprintf("mc_token=%s; SameSite=Strict; Path=/; HttpOnly; Max-Age=%d;", tk, int(dt.Seconds())))
		apiRespond(w, &resp)
	case "login":
		if apireq == nil || apireq.Initiator == nil || apireq.Initiator.Email == nil || apireq.Initiator.Pwd == nil {
			ThrowApiErr(w, "Not enough data to auth", nil, 401)
			return
		}

		usr, err := dbl.CheckAuth(*apireq.Initiator.Email, *apireq.Initiator.Pwd)
		if err != nil {
			ThrowApiErr(w, "Failed to login", err, 403)
			return
		}

		dt := time.Hour * 24 * 3

		tk, err := genJWTtoken(*apireq.Initiator.Email, dt)
		if err != nil {
			ThrowApiErr(w, "Failed to generate token", err, 500)
			return
		}
		resp := apitypes.API_obj{Token: &tk, User: usr}
		// w.Header().Set("Set-Cookie", fmt.Sprintf("mc_token=%s; SameSite=Strict; Path=/; HttpOnly; Max-Age=%d;", tk, int(dt.Seconds())))
		apiRespond(w, &resp)
	case "logout":
		// w.Header().Set("Set-Cookie", "mc_token=none; SameSite=Strict; Path=/; HttpOnly; Max-Age=-1;")
		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
	case "me":
		if uo.Email == nil {
			ThrowApiErr(w, "Username was not parsed", errors.New("uo.email is nil"), 500)
			return
		}
		uo, err = dbl.GetUserInfoByEmail(*uo.Email, true)
		if err != nil {
			ThrowApiErr(w, "Failed to get userinfo", err, 500)
			return
		}
		ap := apitypes.API_obj{User: uo}
		apiRespond(w, &ap)
	default:
		fmt.Println("invalid API operation")
		ThrowApiErr(w, "Invalid API operation", nil, 400)
	}
}

func apiPubWeather(w http.ResponseWriter, r *http.Request, operation string, apireq *apitypes.API_obj) {
	switch operation {
	case "now":
		// weather now at the specified location or Moscow Center

	default:
		fmt.Println("invalid API operation")
		ThrowApiErr(w, "Invalid API operation", nil, 400)
	}
}

func apiPubCats(w http.ResponseWriter, r *http.Request, operation string, apireq *apitypes.API_obj) {
	switch operation {
	case "list":
		cats, err := dbl.CatList()
		if err != nil || cats == nil {
			ThrowApiErr(w, "Failed to fetch categories", err, 500)
			return
		}

		apiRespond(w, &apitypes.API_obj{Categories: cats})
	default:
		ThrowApiErr(w, "Invalid API operation", nil, 400)
	}
}

func apiPubRoute(w http.ResponseWriter, r *http.Request, operation string, apireq *apitypes.API_obj) {
	switch operation {
	case "build":
		// Basic input validation
		if apireq == nil ||
			apireq.PosReq == nil ||
			apireq.PosReq.MyLat == nil ||
			apireq.PosReq.MyLong == nil ||
			apireq.PosReq.Cats == nil ||
			len(*apireq.PosReq.Cats) == 0 {
			ThrowApiErr(w, "my_lat, my_long and cats must be present", nil, 400)
			return
		}

		places, err := dbl.BuildRoute(apireq.PosReq)
		if err != nil {
			ThrowApiErr(w, "Unable to build a route", err, 500)
			return
		}
		apiRespond(w, &apitypes.API_obj{Places: places})
	default:
		ThrowApiErr(w, "Invalid API operation", nil, 400)
	}
}

func apiRespond(w http.ResponseWriter, apiobj *apitypes.API_obj) {
	j, err := json.Marshal(apiobj)
	if err != nil {
		ThrowApiErr(w, "Failed to respond with data", err, 500)
		return
	}
	fmt.Fprint(w, string(j))
}

func fileserver(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/web/dashboard") {
			c, err := r.Cookie("mc_token")
			if err != nil {
				http.Redirect(w, r, "/web/login", http.StatusSeeOther)
				return
			}
			arr := strings.Split(c.Value, "mc_token=")
			if len(arr) == 0 {
				http.Redirect(w, r, "/web/login", http.StatusSeeOther)
				return
			}
			_, err = verifyJWT(&arr[0])
			if err != nil {
				http.Redirect(w, r, "/web/login", http.StatusSeeOther)
				return
			}
		}
		w.Header().Set("Content-Disposition", "inline")
		w.Header().Add("access-control-allow-origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Key, API-usr, Token, ref-key, lu-key")
		fs.ServeHTTP(w, r)
	}
}

func ThrowApiErr(w http.ResponseWriter, custom string, err error, code int) {
	if custom != "" && err == nil {
		apiRespond(w, &apitypes.API_obj{Error: &apitypes.ErrorStruct{Message: custom, Error: custom, Code: code}})
		return
	}

	if err != nil && custom != "" {
		apiRespond(w, &apitypes.API_obj{Error: &apitypes.ErrorStruct{Message: custom, Error: err.Error(), Code: 403}})
		return
	}

	custom = "You cannot perform this action. Access Denied."
	apiRespond(w, &apitypes.API_obj{Error: &apitypes.ErrorStruct{Message: custom, Error: custom, Code: 403}})
	return
}

func setupRoutes() {
	router := mux.NewRouter().StrictSlash(false)

	// router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { wsHandler(w, r) })
	router.HandleFunc("/api/{group}/{endpoint}/{operation}", apiGlobalRouter)

	fs := http.StripPrefix("", http.FileServer(http.Dir("./web/dist/")))
	router.PathPrefix("").Handler(fileserver(fs))

	router.NotFoundHandler = router.NewRoute().HandlerFunc(homePage).GetHandler()
	log.Fatal(http.ListenAndServe(":5012", router))
}

func genJWTkey() error {
	k := make([]byte, 128)
	_, err := rand.Read(k)
	if err != nil {
		fmt.Println(err)
		return errors.New("failed to generate JWT private key")
	}
	JWT_PRIV_KEY = &k
	return nil
}

func genJWTtokenByUserObj(for_user *apitypes.User_Obj, dt time.Duration) (string, error) {
	e := ""
	if for_user.Email != nil {
		e = *for_user.Email
	}
	cl := CustomClaims{
		*for_user,
		e,
		string([]byte{0x75, 0x5f, 0x72, 0x5f, 0x63, 0x75, 0x74, 0x65, 0x5f, 0x3c, 0x33}),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(dt).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "tcm_weblink_module",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)

	var token_s string
	token_s, err := token.SignedString(*JWT_PRIV_KEY)
	if err != nil {
		fmt.Println(err)
		return "", errors.New("failed to sign JWT for web")
	}

	return token_s, nil
}

func genJWTtoken(for_user string, dt time.Duration) (string, error) {
	// get data
	usr, err := dbl.GetUserInfoByEmail(for_user, true)
	if err != nil {
		return "", err
	}

	return genJWTtokenByUserObj(usr, dt)
}

func verifyJWT(token_in *string) (*apitypes.User_Obj, error) {
	if token_in == nil {
		return nil, errors.New("invalid token")
	}

	uo := apitypes.User_Obj{}
	cl := CustomClaims{}
	_, err := jwt.ParseWithClaims(*token_in, &cl, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		if cl.Note == "" || cl.Issuer != "tcm_weblink_module" || cl.ExpiresAt <= time.Now().Unix() {
			return nil, errors.New("invalid token")
		}
		uo = cl.UserInfo

		return *JWT_PRIV_KEY, nil
	})
	if err != nil {
		return nil, err
	}

	return &uo, nil
}

func getBinFile(f string) ([]byte, error) {
	fl, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer fl.Close()

	st, e := fl.Stat()
	if e != nil {
		return nil, e
	}

	var s int64 = st.Size()
	bs := make([]byte, s)

	b := bufio.NewReader(fl)
	_, err = b.Read(bs)

	return bs, err
}

func Init(db_cfg *dblink.DBconfig) error {
	fmt.Println("Starting TCM API server...")
	// init db
	dbl = dblink.DBwrap{}
	err := dbl.Init(db_cfg)
	if err != nil {
		return err
	}

	// do we need to regen keys?
	if _, err := os.Stat(WEB_CRYPTO_KEY_PATH); err == nil {
		// load key from file
		var err2 error
		*JWT_PRIV_KEY, err2 = getBinFile(WEB_CRYPTO_KEY_PATH)
		if err2 != nil {
			return err2
		}

	} else if errors.Is(err, os.ErrNotExist) {
		// regen keys
		err = genJWTkey()
		if err != nil {
			fmt.Println(err)
			return err
		}
	} else {
		// Uhhh... Panic I guess
		panic("Key file access failed. Failed to check for file existence. OS error.")
	}

	go func() {
		setupRoutes()
	}()
	fmt.Println("Web ready")
	return nil
}

func Close() error {
	return dbl.Close()
}
