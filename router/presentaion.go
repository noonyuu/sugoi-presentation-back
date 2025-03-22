package router

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noonyuu/websocket-server/db"
)

type PresentationHandler struct {
	DB *db.Database
}

func NewPresentationHandler(database *db.Database) *mux.Router {
	handler := &PresentationHandler{DB: database}
	router := mux.NewRouter()
	router.HandleFunc("/presentation/create", handler.createPresentation).Methods("POST")
	router.HandleFunc("/presentation/get", handler.getPresentationUrls).Methods("GET")
	router.HandleFunc("/presentation/get/selected", handler.getSelectedPresentationUrl).Methods("GET")

	return router
}

func (handler *PresentationHandler) createPresentation(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserId string `json:"user_id"`
		Url    string `json:"url"`
	}
	// リクエストボディのデコード
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "リクエスト解析エラー", http.StatusBadRequest)
		return
	}

	// 存在確認
	if _, err := handler.DB.GetUser(data.UserId); err != nil {
		http.Error(w, "ユーザーが存在しません", http.StatusBadRequest)
		return
	}

	if _, err := handler.DB.GetSelectedPresentationUrl(data.UserId, data.Url); err == nil {
		http.Error(w, "すでに登録されています", http.StatusBadRequest)
		return
	}

	presentationUrl := &db.PresentationUrl{
		UserId: data.UserId,
		Url:    data.Url,
	}

	// 保存
	if err := handler.DB.InsertPresentationUrl(presentationUrl); err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}

	// 成功
	json.NewEncoder(w).Encode(presentationUrl)
}

func (handler *PresentationHandler) getPresentationUrls(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserId string `json:"user_id"`
	}
	// リクエストボディのデコード
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "リクエスト解析エラー", http.StatusBadRequest)
		return
	}

	presentationUrls, err := handler.DB.GetAllPresentationUrls(data.UserId)
	if err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(presentationUrls)
}

func (handler *PresentationHandler) getSelectedPresentationUrl(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserId string `json:"user_id"`
		Url    string `json:"url"`
	}
	// リクエストボディのデコード
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "リクエスト解析エラー", http.StatusBadRequest)
		return
	}

	presentationUrl, err := handler.DB.GetSelectedPresentationUrl(data.UserId, data.Url)
	if err != nil {
		http.Error(w, "データベースエラー", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(presentationUrl)
}
