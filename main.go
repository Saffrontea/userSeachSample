package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Priority int    `db:"priority"`
}

type UserData []User

func handler(w http.ResponseWriter, r *http.Request) {
	var userData UserData
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("B_USER"),
		os.Getenv("B_PASS"),
		os.Getenv("B_HOST"),
		os.Getenv("B_PORT"),
		os.Getenv("B_DBNAME"),
	)
	db, err := sqlx.Open("mysql", dataSource)
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
	db.SetConnMaxLifetime(time.Minute * 3)
	// 接続数はとりあえず10に設定
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err != nil {
		log.Fatal(err)
	}
	t := strings.Replace(
		r.URL.Path,
		"/",
		"",
		1,
	)
	err = db.Select(
		&userData,
		"select * from test where name like ? order by priority limit 10",
		fmt.Sprintf("%s%%", t),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(userData)
	fmt.Fprint(w, userData)
	fmt.Fprintln(w)
}

func main() {
	http.HandleFunc("/", handler) // ハンドラを登録してウェブページを表示させる
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
