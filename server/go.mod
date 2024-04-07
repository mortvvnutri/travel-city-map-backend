module tcm

go 1.22.4

require tcm/weblink v0.0.0

replace tcm/weblink => ./modules/weblink

require tcm/dblink v0.0.0

require (
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	golang.org/x/crypto v0.24.0 // indirect
)

replace tcm/dblink => ./modules/dblink

require (
	github.com/joho/godotenv v1.5.1
	tcm/apitypes v0.0.0
)

replace tcm/apitypes => ./modules/apitypes


require tcm/utils v0.0.0
replace tcm/utils => ./modules/utils