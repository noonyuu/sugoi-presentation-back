package router

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// WebSocketの接続を管理するマップとミューテックス
	activeConnections = make(map[*websocket.Conn]bool)
	connMutex         = sync.Mutex{}
	upgrader          = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 同一オリジンのチェックを無効化（必要に応じて強化）
		},
	}
)

// WebSocket接続を登録
func addConnection(conn *websocket.Conn) {
	connMutex.Lock()
	defer connMutex.Unlock()
	activeConnections[conn] = true
}

// WebSocket接続を解除
func removeConnection(conn *websocket.Conn) {
	connMutex.Lock()
	defer connMutex.Unlock()
	delete(activeConnections, conn)
}

// WebSocketのハンドラー関数
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to set WebSocket upgrade:", err)
		return
	}
	defer conn.Close()

	addConnection(conn)
	defer removeConnection(conn)

	// 初回のメッセージを送信
	err = conn.WriteMessage(websocket.TextMessage, []byte("Server: Hello, Client!"))
	if err != nil {
		log.Println("Failed to send initial message:", err)
		return
	}

	// クライアントからのメッセージを待機
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error receiving message:", err)
			break
		}
		fmt.Printf("Received message from client: %s\n", msg)
	}
}

// リアクション更新をクライアントに通知する関数
func NotifyReactionUpdate(wordID string, newCount int) {
	connMutex.Lock()
	defer connMutex.Unlock()

	message := fmt.Sprintf(`{"wordID": "%s", "count": %d}`, wordID, newCount)

	for conn := range activeConnections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error sending message to client: %v\n", err)
			conn.Close()
			delete(activeConnections, conn) // エラーが発生した接続を削除
		}
	}
}

// 投稿次にwsで通知する
func NotifyPost(word string) {
	fmt.Println("NotifyPost")
	connMutex.Lock()
	defer connMutex.Unlock()

	message := fmt.Sprint(word)

	for conn := range activeConnections {
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error sending message to client: %v\n", err)
			conn.Close()
			delete(activeConnections, conn) // エラーが発生した接続を削除
		}
	}
}
