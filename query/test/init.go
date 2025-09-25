package test

import (
	"database/sql"
	"log"
	"time"

	"github.com/Cooooing/cutil/common/logger"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("connect to database success")
}

type User struct {
	Id        *int       `json:"id"`
	Name      *string    `json:"name"`
	Age       *int       `json:"age"`
	Email     *string    `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
}

type PostTitle struct {
	UserId    *int    `json:"id"`
	UserName  *string `json:"name"`
	PostTitle *string `json:"title"`
}
