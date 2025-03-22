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
	router.HandleFunc("/user/get", handler.getUserInfo).Methods("GET")

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
		json.NewEncoder(w).Encode(err)
		return
	}

	if user == nil {
		json.NewEncoder(w).Encode("ユーザーが存在しません")
		return
	}

	err = lib.PasswordCompare(user.Password, userInfo.Password)
	if err != nil {
		json.NewEncoder(w).Encode("パスワードが違います")
		return
	}



	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var userInfo struct {
		UserId   string `json:"user_id"`
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	// リクエストボディーのデコード
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	var user db.User

	userPtr, err := h.DB.GetUser(userInfo.UserId)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	if userPtr != nil {
		user = *userPtr
	}

	if user.UserId != "" {
		json.NewEncoder(w).Encode("既に存在するユーザーです")
		return
	}

	encryptPassword, err := lib.PasswordEncrypt(userInfo.Password)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	// 登録
	user = db.User{
		UserId:   userInfo.UserId,
		Name:     userInfo.Name,
		Password: encryptPassword,
		Email:    userInfo.Email,
	}

	err = h.DB.CreateUser(&user)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	json.NewEncoder(w).Encode("ユーザー登録完了")
}
