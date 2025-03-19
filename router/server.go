package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noonyuu/websocket-server/db"
)

func Server(database *db.Database) {
	mainRouter := mux.NewRouter()

	// 疎通確認用
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received a request")
		w.Write([]byte("Hello, World!"))
	})

	mainRouter.HandleFunc("/ws", handleWebSocket)
	mainRouter.HandleFunc("/comment-form/{sessionId}", commentFormHandler).Methods("GET")
	mainRouter.HandleFunc("/comment", commentHandler).Methods("POST")

	go func() {
		log.Println("Server is running on port 8080")
		if err := http.ListenAndServe(":8080", mainRouter); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
}
