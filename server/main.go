package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tcm/apitypes"
	"tcm/dblink"
	"tcm/weblink"

	"github.com/joho/godotenv"
)

type env_cfg struct {
	db  *dblink.DBconfig
	owm *apitypes.OWM_CFG
}

/*
Loads .env file and compiles variables into struct
*/
func loadEnv() (*env_cfg, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	ret := env_cfg{}

	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_pwd := os.Getenv("DB_PWD")
	db_name := os.Getenv("DB_DBNAME")
	owm_api_key := os.Getenv("OWM_API_KEY")
	ret.db = &dblink.DBconfig{
		Host:   &db_host,
		Port:   &db_port,
		User:   &db_user,
		Pwd:    &db_pwd,
		Dbname: &db_name,
	}
	ret.owm = &apitypes.OWM_CFG{
		ApiKey: &owm_api_key,
	}

	return &ret, nil
}

func main() {
	fmt.Println("Please wait while the server is starting...")

	// Keep the app rinning, listen for interrupts
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	// Load ENV file into env_cfg struct
	env_vars, err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	// web and db
	err = weblink.Init(env_vars.db, env_vars.owm)
	if err != nil {
		log.Fatal(err)
	}

	// Block execution until interrupt is received
	<-sc

	// Start graceful shutdown
	fmt.Println("Stopping server and saving data...")
	weblink.Close()
}
