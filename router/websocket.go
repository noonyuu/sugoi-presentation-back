package router

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// WebSocketの接続を管理するマップとミューテックス
	activeSessions = make(map[string]map[*websocket.Conn]bool)
	connMutex      = sync.Mutex{}
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 同一オリジンのチェックを無効化（必要に応じて強化）
		},
	}
)

// WebSocket接続を登録
func addConnection(sessionId string, conn *websocket.Conn) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if _, exists := activeSessions[sessionId]; !exists {
		activeSessions[sessionId] = make(map[*websocket.Conn]bool)
	}
	activeSessions[sessionId][conn] = true
}

// WebSocket接続を解除
func removeConnection(sessionId string, conn *websocket.Conn) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if connections, exists := activeSessions[sessionId]; exists {
		delete(connections, conn)
		if len(connections) == 0 {
			delete(activeSessions, sessionId)
		}
	}
}

// WebSocketのハンドラー関数
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	sessionId := r.URL.Query().Get("sessionId") // クエリパラメータから取得
	if sessionId == "" {
		http.Error(w, "Missing sessionId", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade:", err)
		return
	}
	defer conn.Close()

	addConnection(sessionId, conn)
	defer removeConnection(sessionId, conn)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error receiving message:", err)
			break
		}
		fmt.Printf("Received message from client [%s]: %s\n", sessionId, msg)
	}
}

// リアクション更新をクライアントに通知する関数
func NotifySession(sessionId string, message string) {
	connMutex.Lock()
	defer connMutex.Unlock()

	if connections, exists := activeSessions[sessionId]; exists {
		for conn := range connections {
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("Error sending message to client: %v\n", err)
				conn.Close()
				delete(connections, conn)
			}
		}
	}
}

// コメント投稿時に通知
func NotifyPost(sessionId string, word string) {
	fmt.Println("NotifyPost to session:", sessionId)

	message := fmt.Sprintf(`{"word": "%s", "timestamp": "%s"}`, word, time.Now().Format(time.RFC3339))
	NotifySession(sessionId, message)
}
