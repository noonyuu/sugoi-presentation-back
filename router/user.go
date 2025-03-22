package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noonyuu/websocket-server/db"
	"github.com/noonyuu/websocket-server/lib"
)

type UserHandler struct {
	DB *db.Database
}

func NewUserHandler(database *db.Database) *mux.Router {
	handler := &UserHandler{DB: database}
	router := mux.NewRouter()
	router.HandleFunc("/user/create", handler.createUser).Methods("POST")
	router.HandleFunc("/user/get", handler.getUserInfo).Methods("POST")

	return router
}

func (h *UserHandler) getUserInfo(w http.ResponseWriter, r *http.Request) {
	var userInfo struct {
		UserId   string `json:"user_id"`
		Password string `json:"password"`
	}
	// リクエストボディーのデコード
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	user, err := h.DB.GetUser(userInfo.UserId)
	if err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "ユーザーが存在しません", http.StatusUnauthorized)
		return
	}

	if err := lib.PasswordCompare(user.Password, userInfo.Password); err != nil {
		http.Error(w, "パスワードが違います", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {

	var userInfo struct {
		UserId   string `json:"user_id"`
		Name     string `json:"name"`
		Password string `json:"password"`
		// Email    string `json:"email"`
	}

	// リクエストボディーのデコード
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		http.Error(w, `{"error": "リクエストの解析に失敗しました"}`, http.StatusBadRequest)
		return
	}

	// 必須項目のバリデーション
	if userInfo.UserId == "" || userInfo.Name == "" || userInfo.Password == "" {
		http.Error(w, `{"error": "全ての項目を入力してください"}`, http.StatusBadRequest)
		return
	}

	// 既存ユーザー確認
	userPtr, err := h.DB.GetUser(userInfo.UserId)
	if err != nil {
		http.Error(w, `{"error": "データベースエラー"}`, http.StatusInternalServerError)
		return
	}

	if userPtr != nil {
		http.Error(w, `{"error": "既に存在するユーザーです"}`, http.StatusConflict)
		return
	}

	// パスワードの暗号化
	encryptPassword, err := lib.PasswordEncrypt(userInfo.Password)
	if err != nil {
		http.Error(w, `{"error": "パスワードの暗号化に失敗しました"}`, http.StatusInternalServerError)
		return
	}

	// ユーザー登録
	user := db.User{
		UserId:   userInfo.UserId,
		Name:     userInfo.Name,
		Password: encryptPassword,
		// Email:    userInfo.Email,
	}

	err = h.DB.CreateUser(&user)
	if err != nil {
		http.Error(w, `{"error": "ユーザー登録に失敗しました"}`, http.StatusInternalServerError)
		return
	}

	// 成功レスポンス
	json.NewEncoder(w).Encode(user)
}
