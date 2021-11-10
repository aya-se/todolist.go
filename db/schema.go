package db

// schema.go provides data models in DB
import (
	"time"
)

// タスク
type Task struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	Title     string    `db:"title"`
	Detail    string    `db:"detail"`
	Priority  int       `db:"priority"`
	CategoryID   uint64 `db:"category_id"`
	CategoryName string `db:"category_name"`
	CreatedAt time.Time `db:"created_at"`
	Deadline	time.Time `db:"deadline"`
	IsDone    bool      `db:"is_done"`
}

// タスク編集フォーム
type TaskForm struct {
	Title     string    `form:"title"`
	Detail    string    `form:"detail"`
	Priority  int       `form:"priority"`
	CategoryID   uint64 `form:"category_id"`
	Deadline	string    `form:"deadline"`
}

// ユーザー
type User struct {
	UserID    uint64    `db:"user_id"`
	UserName  string    `db:"user_name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	IsDeleted	bool      `db:"is_deleted"`
}

// ユーザー編集フォーム
type UserForm struct {
	UserName  string    `form:"user_name"`
	Password  string    `form:"password"`
	Confirm   string    `form:"confirm"`
}

// タスク検索フォーム
type SearchForm struct {
	Substring string    `form:"substring"`
	Status    string    `form:"status"`
	Priority  string    `form:"priority"`
	CategoryID   uint64 `form:"category_id"`
	Order     string    `form:"order"`
}

// カテゴリ
type Category struct {
	CategoryID   uint64 `db:"category_id"`
	UserID       uint64 `db:"user_id"`
	CategoryName string `db:"category_name"`
	CreatedAt time.Time `db:"created_at"`
}

// カテゴリ編集フォーム
type CategoryForm struct {
	CategoryName string `form:"category_name"`
}