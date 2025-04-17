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

	var v ulid.ULID
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
		err = v.UnmarshalBinary(users[0].ID)
		if err != nil {
			slog.Error("failed to marshal ID", "error", err)
			return
		}
		for _, user := range users {
			fmt.Println(v.String())
			fmt.Println(user.Name)
			fmt.Println(user.Email)
			fmt.Println(user.Password)
			fmt.Println(user.CreatedAt)
			fmt.Println(user.UpdatedAt)
		}
	} else {
		fmt.Println("No users found")
		CreateUser(ctx, &db)
	}
}

func CreateUser(ctx context.Context, db *bob.DB) {
	johnID := ulid.Make()

	slog.Info("johnID", slog.String("johnID before binary", johnID.String()))
	// トランザクションを開始
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("failed to begin transaction", "error", err)
		return
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			slog.Error("failed to rollback transaction", "error", err)
		}
	}()

	// Insertでuserを追加
	setter := &models.UserSetter{
		ID:        omit.From(johnID.Bytes()),
		Name:      omit.From("John Doe"),
		Email:     omit.From("john.doe@example.com"),
		Password:  omit.From("password"),
		CreatedAt: omit.From(time.Now()),
	}

	m, err := models.Users.Insert(setter).Exec(ctx, db)

	if err != nil {
		slog.Error("failed to insert user",
			"error", err,
			"id", johnID.String(),
			"name", "John Doe",
			"email", "john.doe@example.com")
		return
	}

	fmt.Println(m)

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		slog.Error("failed to commit transaction", "error", err)
		return
	}

	john, err := models.FindUser(ctx, db, johnID.Bytes())
	if err != nil {
		slog.Error("failed to find user", "error", err)
		return
	}
	var id ulid.ULID

	err = id.UnmarshalBinary(john.ID)
	if err != nil {
		slog.Error("failed to unmarshal ID", "error", err)
		return
	}
	slog.Info("IDs", slog.String("ID", id.String()), slog.String("Name", john.Name), slog.String("Email", john.Email), slog.String("Password", john.Password), slog.String("CreatedAt", john.CreatedAt.Format(time.RFC3339)), slog.String("UpdatedAt", john.UpdatedAt.Format(time.RFC3339)))
}
