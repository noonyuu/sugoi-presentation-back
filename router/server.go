package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/noonyuu/websocket-server/db"
)

func Server(database *db.Database) {
	mainRouter := mux.NewRouter()

	commentRouter := NewCommentHandler(database)
	userRouter := NewUserHandler(database)
	presentationRouter := NewPresentationHandler(database)

	// 疎通確認用
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received a request")
		w.Write([]byte("Hello, World!"))
	})

	mainRouter.HandleFunc("/ws", handleWebSocket)

	mainRouter.PathPrefix("/comment").Handler(commentRouter)
	mainRouter.PathPrefix("/user").Handler(userRouter)
	mainRouter.PathPrefix("/presentation").Handler(presentationRouter)

	go func() {
		log.Println("Server is running on port 8080")
		if err := http.ListenAndServe(":8080", mainRouter); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()
}
