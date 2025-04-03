package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/stephenafamo/bob"
)

func NewDB(ctx context.Context) (*bob.DB, error) {
	dbname := os.Getenv("MYSQL_DATABASE")
	password := os.Getenv("MYSQL_PASSWORD")
	user := os.Getenv("MYSQL_USER")
	db, err := bob.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", user, password, dbname))
	if err != nil {
		return nil, err
	}

	return &db, nil
}
