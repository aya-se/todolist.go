package service

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	database "todolist.go/db"
	. "todolist.go/db"
	"golang.org/x/crypto/bcrypt"
	"time"
	"regexp"
	"strings"
	"log"
)

var LoginInfo = User{}

// Home renders index.html
func Home(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	t := time.Now().Local().Format("2006-01-02T15:04")
	var categories []database.Category
	user_id := LoginInfo.UserID
	err = db.Select(&categories, "SELECT * FROM categories WHERE user_id=?", user_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME", "Now": t, "Categories": categories, "User": LoginInfo.UserName})
}

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// タスクの検索条件指定に対応
	var searchForm SearchForm
	ctx.Bind(&searchForm)
	status := searchForm.Status
	priority := searchForm.Priority
	category_id := searchForm.CategoryID
	order := searchForm.Order
	// 空の場合(指定なし)
	if order == "" {
		order = "deadline"
	}
	// 指定の部分文字列
	substring := searchForm.Substring
	query_status := ""
	query_priority := ""
	query_category_id := ""
	// 完了状態の絞り込みクエリを追加
	if status == "incomplete" {
		query_status = " AND is_done=0"
	} else if status == "completed" {
		query_status = " AND is_done=1"
	}
	// 優先度の絞り込みクエリを追加
	if priority == "high" {
		query_priority = " AND priority<=1"
	}
	// カテゴリIDの絞り込みクエリの追加
	if category_id != 0 {
		query_category_id = " AND T.category_id=" + strconv.FormatUint(category_id, 10)
	}
	// Get tasks in DB
	var tasks []database.Task
	query := "SELECT id, T.user_id AS 'user_id', title, detail, priority, T.category_id AS 'category_id', category_name, T.created_at AS 'created_at', deadline, is_done FROM tasks AS T LEFT JOIN categories AS C ON T.category_id=C.category_id WHERE T.user_id="+ strconv.FormatUint(LoginInfo.UserID, 10) + " AND title LIKE '%"+ substring + "%'" + query_status + query_priority + query_category_id + " ORDER BY " + order
	log.Println(query);
	err = db.Select(&tasks, query)
	var categories []database.Category
	user_id := LoginInfo.UserID
	err = db.Select(&categories, "SELECT * FROM categories WHERE user_id=?", user_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks, "User": LoginInfo.UserName, "Status": status, "Priority": priority, "Substring": substring, "Order": order, "CategoryID": category_id, "Categories": categories})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Render task
	ctx.String(http.StatusOK, task.Title)
}

//　新しいタスクの追加
func InsertTask(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// 新しいタスクの追加
	var data TaskForm
	ctx.Bind(&data)
	user_id := LoginInfo.UserID
	title := data.Title
	detail := data.Detail
	priority := data.Priority
	category_id := data.CategoryID
	deadline := data.Deadline
	_, err = db.Query("INSERT INTO tasks (user_id, title, detail, priority, category_id, deadline) VALUES (?, ?, ?, ?, ?, ?)", user_id, title, detail, priority, category_id, deadline)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

//　指定タスクの編集
func UpdateTask(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// 指定タスクの編集
	var data TaskForm
	ctx.Bind(&data)
	title := data.Title
	detail := data.Detail
	priority := data.Priority
	category_id := data.CategoryID
	deadline := data.Deadline
	_, err = db.Query("UPDATE tasks SET title=? , detail=?, priority=?, category_id=?, deadline=? WHERE id=?", title, detail, priority, category_id, deadline, id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

// 指定タスクの削除
func DeleteTask(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Get a task with given ID
	_, err = db.Query("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

// 指定タスクの完了・再開
func CompleteTask(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	status, err := strconv.Atoi(ctx.Param("status"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Get a task with given ID
	_, err = db.Query("UPDATE tasks SET is_done=? WHERE id=?", status, id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

// 指定タスクの編集ページ
func EditTask(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Get a task with given ID
	var task Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	
	var categories []database.Category
	user_id := LoginInfo.UserID
	err = db.Select(&categories, "SELECT * FROM categories WHERE user_id=?", user_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Render task
	deadline := task.Deadline.Format("2006-01-02T15:04")
	ctx.HTML(http.StatusOK, "edit_task.html", gin.H{"ID": task.ID, "Title": task.Title, "Detail": task.Detail, "Priority": task.Priority, "CategoryID": task.CategoryID, "Deadline": deadline, "Categories": categories, "User": LoginInfo.UserName })
}

// ユーザー登録ページ
func Signup(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", gin.H{"Title": "SignUp", "User": LoginInfo.UserName})
}

// ユーザーログインページ
func Signin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signin.html", gin.H{"Title": "SignIn", "User": LoginInfo.UserName})
}

// ユーザー編集ページ
func EditUser(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	ctx.HTML(http.StatusOK, "edit_user.html", gin.H{"Title": "EditUser", "User": LoginInfo.UserName, "UserName": LoginInfo.UserName})
}

// ユーザー登録
func InsertUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// 新しいユーザーの追加
	var data UserForm
	ctx.Bind(&data)
	user_name := data.UserName
	password := data.Password
	confirm := data.Confirm

	// パスワードの一致確認
	if password != confirm {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "2つのフォームに入力されたパスワードが異なっています！再度パスワードを登録してください。"})
		return
	}
	// バリデーション
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if len(user_name)<4 || len(user_name)>16 || !re.MatchString(user_name) {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "ユーザー名には4文字以上16文字以下の英数字とアンダーバー(_)を入力してください。"})
		return
	}
	if len(password)<4 || len(password)>16 || !re.MatchString(password) {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "パスワードには4文字以上16文字以下の英数字とアンダーバー(_)を入力してください。"})
		return
	}
	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),12)
	_, err = db.Query("INSERT INTO users (user_name, password) VALUES (?, ?)", user_name, hashedPassword)
	if err != nil {
		if strings.Count(err.Error(),"Duplicate")>0 {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "このユーザー名は既に登録されています！異なるユーザー名を指定してください。"})
			return
		} else {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
			return
		}
	}

	// リダイレクト
	ctx.Redirect(303, "/signout-user")
}

// ユーザーログイン
func SigninUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	var data UserForm
	ctx.Bind(&data)
	user_name := data.UserName
	password := data.Password

	var user User
	err = db.Get(&user, "SELECT * FROM users WHERE user_name=?", user_name)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "ユーザー名もしくはパスワードが間違っています。"})
		return
	}
	// 退会確認
	if user.IsDeleted {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "既に退会したユーザーです。別のアカウントでログインしてください。"})
		return
	}
	// パスワード確認
	hash := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "ユーザー名もしくはパスワードが間違っています。"})
		return
	}
	// リダイレクト
	LoginInfo = user
	ctx.Redirect(303, "/list")
}

// ユーザー更新
func UpdateUser(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// ユーザー情報更新
	var data UserForm
	ctx.Bind(&data)
	user_id := LoginInfo.UserID
	user_name := data.UserName
	password := data.Password
	confirm := data.Confirm

	// パスワードの一致確認
	if password != confirm {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "2つのフォームに入力されたパスワードが異なっています！再度パスワードを登録してください。"})
		return
	}
	// バリデーション
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if len(user_name)<4 || len(user_name)>16 || !re.MatchString(user_name) {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "ユーザー名には4文字以上16文字以下の英数字とアンダーバー(_)を入力してください。"})
		return
	}
	if len(password)<4 || len(password)>16 || !re.MatchString(password) {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "パスワードには4文字以上16文字以下の英数字とアンダーバー(_)を入力してください。"})
		return
	}
	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),12)
	_, err = db.Query("UPDATE users SET user_name=?, password=? WHERE user_id=?", user_name, hashedPassword, user_id)
	if err != nil {
		if strings.Count(err.Error(),"Duplicate")>0 {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "このユーザー名は既に登録されています！異なるユーザー名を指定してください。"})
			return
		} else {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
			return
		}
	}
	// リダイレクト
	LoginInfo.UserName = user_name
	ctx.Redirect(303, "/list")
}

// ユーザーログアウト
func SignoutUser(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// リダイレクト
	var emptyUser User
	LoginInfo = emptyUser
	ctx.Redirect(303, "/signin")
}

// ユーザー退会
func DeleteUser(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// ユーザーの退会処理
	user_id := LoginInfo.UserID
	_, err = db.Query("UPDATE users SET is_deleted=1 WHERE user_id=?", user_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// リダイレクト
	var emptyUser User
	LoginInfo = emptyUser
	ctx.Redirect(303, "/signin")
}

// カテゴリ管理ページ
func EditCategories(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	
	var categories []database.Category
	user_id := LoginInfo.UserID
	err = db.Select(&categories, "SELECT * FROM categories WHERE user_id=?", user_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "edit_categories.html", gin.H{"Title": "Categories", "Categories": categories, "User": LoginInfo.UserName})
}

//　新しいカテゴリの追加
func InsertCategory(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// 新しいカテゴリの追加
	var data CategoryForm
	ctx.Bind(&data)
	user_id := LoginInfo.UserID
	category_name := data.CategoryName
	
	// バリデーション
	if len(category_name)<1 {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "カテゴリ名は1文字以上にしてください。"})
		return
	}
	_, err = db.Query("INSERT INTO categories (user_id, category_name) VALUES (?, ?)", user_id, category_name)
	if err != nil {
		if strings.Count(err.Error(),"Duplicate")>0 {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "このカテゴリ名は既に登録されています！異なるカテゴリ名を指定してください。"})
			return
		} else {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
			return
		}
	}

	// リダイレクト
	ctx.Redirect(303, "/edit-categories")
}

//　指定カテゴリの編集
func UpdateCategory(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	category_id, err := strconv.Atoi(ctx.Param("category_id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// 指定カテゴリの編集
	var data CategoryForm
	ctx.Bind(&data)
	category_name := data.CategoryName

	// バリデーション
	if len(category_name)<1 {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "カテゴリ名は1文字以上にしてください。"})
		return
	}
	_, err = db.Query("UPDATE categories SET category_name=? WHERE category_id=?", category_name, category_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/edit-categories")
}

// 指定カテゴリの削除
func DeleteCategory(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	user_id := LoginInfo.UserID
	category_id, err := strconv.Atoi(ctx.Param("category_id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	_, err = db.Query("DELETE FROM categories WHERE category_id=?", category_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// 削除するカテゴリに登録されているタスクをカテゴリ未登録に変更
	_, err = db.Query("UPDATE tasks SET category_id=1 WHERE user_id=? AND category_id=?", user_id, category_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/edit-categories")
}