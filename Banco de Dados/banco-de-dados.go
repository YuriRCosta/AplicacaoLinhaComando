package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	stringConnect := "root:admin@/first?charset=utf8&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", stringConnect)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conectado com sucesso!")

	linhas, err := db.Query("select * from usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer linhas.Close()

	for linhas.Next() {
		var (
			id    int
			nome  string
			email string
		)
		if err := linhas.Scan(&id, &nome, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Println(id, nome, email)
	}
}
