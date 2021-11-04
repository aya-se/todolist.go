package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
type Task struct {
	ID        uint64    `db:"id"`
	UserID		string		`db:"user_id"`
	Title     string    `db:"title"`
	Detail    string    `db:"detail"`
	Priority  string    `db:"priority"`
	CreatedAt time.Time `db:"created_at"`
	Deadline	time.Time `db:"deadline"`
	IsDone    bool      `db:"is_done"`
}

// タスク編集フォーム
type TaskForm struct {
	Title     string    `form:"title"`
	Detail    string    `form:"detail"`
	Priority  string    `form:"priority"`
	Deadline	string    `form:"deadline"`
}

// ユーザー
type User struct {
	UserID    string    `db:"user_id"`
	UserName  string    `db:"user_name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	IsDeleted	bool      `db:"is_deleted"`
}

// ユーザー編集フォーム
type UserForm struct {
	UserID    string    `form:"user_id"`
	UserName  string    `form:"user_name"`
	Password  string    `form:"password"`
}