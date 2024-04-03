package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tcm/dblink"
	"tcm/weblink"

	"github.com/joho/godotenv"
)

type env_cfg struct {
	db *dblink.DBconfig
}

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
	ret.db = &dblink.DBconfig{
		Host:   &db_host,
		Port:   &db_port,
		User:   &db_user,
		Pwd:    &db_pwd,
		Dbname: &db_name,
	}

	return &ret, nil
}

func main() {
	fmt.Println("Please wait while the server is starting...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	env_vars, err := loadEnv()
	if err != nil {
		log.Fatal(err)
	}

	// web and db
	err = weblink.Init(env_vars.db)
	if err != nil {
		log.Fatal(err)
	}

	<-sc

	fmt.Println("Stopping server and saving data...")
	weblink.Close()
	// player.SaveRec()
}
