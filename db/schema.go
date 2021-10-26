package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
type Task struct {
	ID        uint64    `db:"id"`
	UserID		uint64		`db:"user_id"`
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