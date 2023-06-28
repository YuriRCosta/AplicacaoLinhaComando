package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Driver de conexão com o MySQL
)

func Conectar() (*sql.DB, error) {
	stringConnect := "root:admin@/first?charset=utf8&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", stringConnect)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
