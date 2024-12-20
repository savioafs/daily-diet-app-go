package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func LoadConfigs() (*sql.DB, *jwtauth.JWTAuth, int, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpireIn := os.Getenv("JWT_EXPIRESIN")

	expireIn, err := strconv.Atoi(jwtExpireIn)
	if err != nil {
		log.Fatalf("Erro ao converter JWT_EXPIRESIN para int: %v", err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	configToken := jwtauth.New("HS256", []byte(jwtSecret), nil)

	return db, configToken, expireIn, nil
}
