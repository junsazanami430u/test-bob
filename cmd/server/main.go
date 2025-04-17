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

type User struct {
	ID       ulid.ULID
	Name     string
	Email    string
	Password string
}

func main() {

	// テーブルモデルを取得
	userTable := models.Users
	dsn := os.Getenv("MYSQL_DSN")
	ctx := context.Background()

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
		// UnmarshalBinaryでIDを取得
		var v ulid.ULID
		err = v.UnmarshalBinary(users[0].ID)
		if err != nil {
			slog.Error("failed to marshal ID", "error", err)
			return
		}
		for _, user := range users {
			slog.Info(fmt.Sprintf(`
				ID: %s
				Name: %s
				Email: %s
				Password: %s
				CreatedAt: %s
				UpdatedAt: %s
			`,
				v.String(),
				user.Name,
				user.Email,
				user.Password,
				user.CreatedAt.Format(time.RFC3339),
				user.UpdatedAt.Format(time.RFC3339)))
		}
	} else {
		fmt.Println("No users found")
		CreateUser(ctx, &User{
			ID:       ulid.Make(),
			Name:     "John Doe",
			Email:    "john.doe@example.com",
			Password: "password",
		}, &db)
	}
}

func CreateUser(ctx context.Context, user *User, db *bob.DB) {
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
		Name:      omit.From(user.Name),
		Email:     omit.From(user.Email),
		Password:  omit.From(user.Password),
		CreatedAt: omit.From(time.Now()),
	}

	m, err := models.Users.Insert(setter).Exec(ctx, db)

	if err != nil {
		slog.Error(fmt.Sprintf(`failed to insert user: %s\n
			ID: %s
			Name: %s
			Email: %s
			Password: %s
			CreatedAt: %s
			UpdatedAt: %s
		`, err,
			johnID.String(),
			user.Name,
			user.Email,
			user.Password,
			time.Now().Format(time.RFC3339),
			time.Now().Format(time.RFC3339)))
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
	// UnmarshalBinaryでIDを取得
	var id ulid.ULID
	err = id.UnmarshalBinary(john.ID)
	if err != nil {
		slog.Error("failed to unmarshal ID", "error", err)
		return
	}

	slog.Info(fmt.Sprintf(`
		ID: %s
		Name: %s
		Email: %s
		Password: %s
		CreatedAt: %s
		UpdatedAt: %s
	`, id.String(),
		john.Name,
		john.Email,
		john.Password,
		john.CreatedAt.Format(time.RFC3339),
		john.UpdatedAt.Format(time.RFC3339)))
	return
}
