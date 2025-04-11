package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aarondl/opt/omit"
	_ "github.com/go-sql-driver/mysql"
	"github.com/junsazanami430u/test-bob/pkg/gen/models"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/stephenafamo/bob"
)

func main() {
	ctx := context.Background()

	// テーブルモデルを取得
	userTable := models.Users

	// if err := godotenv.Load(); err != nil {
	// 	slog.Error("failed to load environment variables", "error", err)
	// }
	dsn := os.Getenv("MYSQL_DSN")

	slog.Info(dsn)
	db, err := bob.Open("mysql", dsn)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}()

	// user一覧を取得(select * from users)
	users, err := userTable.View.Query().All(ctx, db)
	if err != nil {
		slog.Error("failed to get users", "error", err)
		return
	}

	if len(users) > 0 {
		for _, user := range users {
			fmt.Println(ulid.Parse(string(user.ID)))
			fmt.Println(user.Name)
			fmt.Println(user.Email)
			fmt.Println(user.Password)
			fmt.Println(user.CreatedAt)
			fmt.Println(user.UpdatedAt)
		}
	}

	johnID := ulid.Make()
	//Insertでuserを追加
	_, err = userTable.Insert(&models.UserSetter{
		ID:        omit.From(johnID.Bytes()),
		Name:      omit.From("John Doe"),
		Email:     omit.From("john.doe@example.com"),
		Password:  omit.From("password"),
		CreatedAt: omit.From(time.Now()),
	}).One(ctx, db)

	if err != nil {
		slog.Error("failed to insert user", "error", err)
		return
	}

	john, err := models.FindUser(ctx, db, johnID.Bytes())
	if err != nil {
		slog.Error("failed to find user", "error", err)
		return
	}

	fmt.Printf("ID: %s\n, Name: %s\n, Email: %s\n, Password: %s\n, CreatedAt: %s\n, UpdatedAt: %s\n", string(john.ID), john.Name, john.Email, john.Password, john.CreatedAt, john.UpdatedAt)

	r := echo.New()
	r.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello World!")
	})
	r.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "pong",
		})
	})
	if err := r.Start(":8080"); err != nil {
		slog.Error("Error starting server", "error", err)
	}
}
