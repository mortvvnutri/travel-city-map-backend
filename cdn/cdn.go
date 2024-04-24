package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

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
	"BECAUSE I'M A POTATO",
	"Aperture Science computer-aided enrichment center",
	"Fifteen Hundred Megawatt Aperture Science Heavy Duty Super-Colliding Super Button",
}

var http_port string

func fileserver(fs http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "inline")
		w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Key, API-usr, Token, ref-key, lu-key, X-Identity")
		fs.ServeHTTP(w, r)
	}
}

// Blackhole
func denyIncoming(w http.ResponseWriter, r *http.Request) {
	rd, e := rand.Int(rand.Reader, big.NewInt(int64(len(pb))))
	if e != nil {
		rd = big.NewInt(int64(0))
	}
	w.Header().Set("X-Powered-By", pb[rd.Int64()])
	w.Header().Set("content-type", "text/plain")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Access-Key, API-usr, Token, ref-key, lu-key")
	w.WriteHeader(403)
	fmt.Fprintf(w, "403: Access denied")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(false)
	fs := http.StripPrefix("/userfiles/", http.FileServer(http.Dir("./userfiles/")))
	router.PathPrefix(`/userfiles/{user:.+}/{folder:.+}/{file:.+\.\w+}`).Handler(fileserver(fs))

	router.NotFoundHandler = router.NewRoute().HandlerFunc(denyIncoming).GetHandler()

	if http_port == "443" {
		log.Fatal(http.ListenAndServeTLS(":"+http_port, "cert.crt", "priv.key", router))
	} else {
		log.Fatal(http.ListenAndServe(":"+http_port, router))
	}
}

func main() {
	// Keep the app running
	log.Println("Starting up TCM CDN...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http_port = os.Getenv("PORT")
	if http_port == "" {
		log.Fatal("Env var PORT is not defined")
	}

	go func() {
		handleRequests()
	}()

	<-sc
	log.Println("Shutting down TCM CDN")
}
