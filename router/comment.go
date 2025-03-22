package router

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/noonyuu/websocket-server/db"
)

type PageData struct {
	SessionID string
}

type CommentHandler struct {
	DB *db.Database
}

func NewCommentHandler(database *db.Database) *mux.Router {
	handler := &CommentHandler{DB: database}
	router := mux.NewRouter()
	router.HandleFunc("/comment", handler.commentHandler).Methods("POST")
	router.HandleFunc("/comment/get/{sessionId}", handler.commentListHandler).Methods("GET")
	router.HandleFunc("/comment/form/{sessionId}", commentFormHandler).Methods("GET")

	return router
}

func (h *CommentHandler) commentHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Name      string `json:"name"`
		Comment   string `json:"comment"`
		SessionId string `json:"sessionId"`
	}
	// リクエストボディのデコード
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "リクエストの解析エラー", http.StatusBadRequest)
		return
	}
	// ニックネームのバリデーション
	if len(data.Name) > 10 {
		http.Error(w, "ニックネームは10字以内で入力してください", http.StatusBadRequest)
		return
	} else if data.Name == "" {
		data.Name = "匿名"
	}
	// コメントのバリデーション
	if data.Comment == "" {
		http.Error(w, "コメントを入力してください", http.StatusBadRequest)
		return
	}
	if utf8.RuneCountInString(data.Comment) > 20 { // ✅ 文字数をカウント
		http.Error(w, "コメントは20字以内で入力してください", http.StatusBadRequest)
		return
	}

	comment := db.Comment{
		Name:      data.Name,
		Comment:   data.Comment,
		SessionId: data.SessionId,
	}

	// データベースへコメントを保存
	if err := h.DB.InsertComment(&comment); err != nil {
		http.Error(w, "コメントの保存に失敗しました", http.StatusInternalServerError)
		return
	}

	// ログ出力（デバッグ用）
	fmt.Printf("セッションID: %s, コメント: %s\n", comment.SessionId, comment.Comment)

	// 成功レスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	NotifyPost(comment.Comment)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "コメントを受け付けました",
	})
}

func (h *CommentHandler) commentListHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["sessionId"]
	comments, err := h.DB.GetComments(sessionID)
	if err != nil {
		http.Error(w, "コメントの取得に失敗しました", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

func commentFormHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := mux.Vars(r)["sessionId"]
	fmt.Println("セッションID:", sessionID)
	data := PageData{SessionID: sessionID}
	tmpl, err := template.New("chat").Parse(chatPageHTML)
	if err != nil {
		http.Error(w, "テンプレートの解析エラー", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

var chatPageHTML = `
    <!DOCTYPE html>
    <html lang="ja">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>コメント入力</title>
				<style>
					:root {
						--primary-color: #4361ee;
						--text-color: #333333;
						--border-color: #e0e0e0;
						--success-color: #38b000;
						--error-color: #e63946;
						--bg-color: #f9f9f9;
					}
					
					* {
						margin: 0;
						padding: 0;
						box-sizing: border-box;
						font-family: 'Helvetica Neue', Arial, sans-serif;
					}
					
					body {
						background-color: var(--bg-color);
						display: flex;
						justify-content: center;
						align-items: center;
						min-height: 100vh;
						padding: 20px;
					}
					
					.comment-container {
						width: 100%;
						max-width: 500px;
						background-color: white;
						border-radius: 12px;
						box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
						padding: 24px;
					}
					
					.comment-title {
						font-size: 18px;
						font-weight: 600;
						color: var(--text-color);
						margin-bottom: 16px;
						text-align: center;
					}
					
					.connection-status {
						text-align: center;
						font-size: 14px;
						padding: 6px;
						margin-bottom: 16px;
						border-radius: 4px;
						background-color: #f8f9fa;
					}
					
					.status-connected {
						color: var(--success-color);
					}
					
					.status-disconnected {
						color: var(--error-color);
					}
					
					.comment-form {
						display: flex;
						flex-direction: column;
						gap: 16px;
					}

					.name-input {
						width: 100%;
						padding: 16px;
						border: 1px solid var(--border-color);
						border-radius: 8px;
						font-size: 16px;
						transition: border-color 0.3s, box-shadow 0.3s;
					}

					.name-input:focus {
						outline: none;
						border-color: var(--primary-color);
						box-shadow: 0 0 0 2px rgba(67, 97, 238, 0.2);
					}
					
					.comment-input {
						width: 100%;
						min-height: 120px;
						padding: 16px;
						border: 1px solid var(--border-color);
						border-radius: 8px;
						font-size: 16px;
						resize: none;
						transition: border-color 0.3s, box-shadow 0.3s;
					}
					
					.comment-input:focus {
						outline: none;
						border-color: var(--primary-color);
						box-shadow: 0 0 0 2px rgba(67, 97, 238, 0.2);
					}
					
					.submit-button {
						background-color: var(--primary-color);
						color: white;
						border: none;
						border-radius: 8px;
						padding: 14px;
						font-size: 16px;
						font-weight: 500;
						cursor: pointer;
						transition: background-color 0.2s, transform 0.1s;
						display: block;
						width: 100%;
					}
					
					.submit-button:hover {
						background-color: #3651d1;
					}
					
					.submit-button:active {
						transform: translateY(1px);
					}
					
					.status-message {
						height: 20px;
						font-size: 14px;
						text-align: center;
						margin-top: 12px;
						transition: opacity 0.3s;
					}
					
					.status-success {
						color: var(--success-color);
					}
					
					.status-error {
						color: var(--error-color);
					}

					.char-counter {
						font-size: 14px;
						text-align: right;
						color: var(--text-color);
					}

					.char-counter.error {
						color: var(--error-color);
					}
				</style>
			</head>
			<body>
				<div class="comment-container">
					<div class="comment-title">コメントを入力してください</div>
					<div class="connection-status" id="connectionStatus">接続中...</div>
					<form class="comment-form" id="commentForm">
						<input class="name-input" placeholder="ニックネーム (任意)" id="nameInput" />
						<div class="char-counter" id="nameCharCounter">0 / 10</div>
						<textarea class="comment-input" placeholder="コメントを入力..." id="commentInput"></textarea>
						<div class="char-counter" id="commentCharCounter">0 / 20</div>
						<button type="submit" class="submit-button" id="submitButton">送信</button>
					</form>
					<div class="status-message" id="statusMessage"></div>
				</div>

				<script>
					// WebSocket接続
					let socket;
					const statusMessage = document.getElementById("statusMessage");
					const connectionStatus = document.getElementById("connectionStatus");
					const commentForm = document.getElementById("commentForm");
					const nameInput = document.getElementById("nameInput");
					const commentInput = document.getElementById("commentInput");
					const submitButton = document.getElementById("submitButton");
					const sessionId = "{{.SessionID}}";

					// カウンター
					const nameCharCounter = document.getElementById("nameCharCounter");
					const commentCharCounter = document.getElementById("commentCharCounter");

					nameInput.addEventListener("input", () => {
						const length = nameInput.value.length;
						nameCharCounter.textContent = ` + "`${length} / 20`" + `;
						if (length > 10) {
							nameCharCounter.classList.add("error");
						} else {
							nameCharCounter.classList.remove("error");
						}
					});

					commentInput.addEventListener("input", () => {
						const length = commentInput.value.length;
						commentCharCounter.textContent = ` + "`${length} / 20`" + `;
						if (length > 20) {
							commentCharCounter.classList.add("error");
						} else {
							commentCharCounter.classList.remove("error");
						}
					});

					// WebSocketの接続関数
					function connectWebSocket() {
						// 現在のURLからWebSocketのURLを生成（httpをwsに変換）
						const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
						const wsUrl = ` + "`${protocol}//${window.location.host}/app/ws`" + `;
						
						socket = new WebSocket(wsUrl);
						
						socket.onopen = function() {
							connectionStatus.textContent = '接続しました';
							connectionStatus.classList.add('status-connected');
							connectionStatus.classList.remove('status-disconnected');
							submitButton.disabled = false;
						};
						
						socket.onclose = function() {
							connectionStatus.textContent = '接続が切断されました。再接続中...';
							connectionStatus.classList.add('status-disconnected');
							connectionStatus.classList.remove('status-connected');
							submitButton.disabled = true;
							
							// 再接続を試みる
							setTimeout(connectWebSocket, 3000);
						};
						
						socket.onerror = function() {
							connectionStatus.textContent = '接続エラーが発生しました';
							connectionStatus.classList.add('status-disconnected');
							submitButton.disabled = true;
						};
						
						socket.onmessage = function(event) {
							// メッセージ受信時の処理（必要に応じて）
							console.log('受信したメッセージ:', event.data);
						};
					}
					
					// WebSocketに接続
					connectWebSocket();
					console.log('コメント送信処理前');
					// コメント送信処理
					commentForm.addEventListener('submit', async function(e) {
						console.log('コメント送信処理');
						e.preventDefault();
						const name = nameInput.value.trim();
        		const comment = commentInput.value.trim();  

						if (!comment) {
							showStatusMessage('コメントを入力してください', 'error');
							return;
						}

						if (name.length > 10) {
							statusMessage.textContent = "ニックネームは10字以内で入力してください";
							statusMessage.className = "status-message status-visible status-error";
							return;
						}

						switch (comment.length) {
							case 0:
								statusMessage.textContent = "コメントを入力してください";
								statusMessage.className = "status-message status-visible status-error";
								return;
							case 1:
								statusMessage.textContent = "コメントが短すぎます";
								statusMessage.className = "status-message status-visible status-error";
								return;
							case comment.length > 20:
								statusMessage.textContent = "コメントは20字以内で入力してください";
								statusMessage.className = "status-message status-visible status-error";
								return;
						}
						
						try {
							const response = await fetch('/app/comment', {
								method: 'POST',
								headers: {
									'Content-Type': 'application/json',
								},
								body: JSON.stringify({ name, comment, sessionId : ` + "`${sessionId}`" + ` }),
							});
							
							const data = await response.json();
							
							if (data.success) {
								showStatusMessage('コメントを送信しました', 'success');
								commentInput.value = '';
							} else {
								showStatusMessage(data.message || 'エラーが発生しました', 'error');
							}
						} catch (error) {
							console.error('エラー:', error);
							showStatusMessage('通信エラーが発生しました', 'error');
						}
					});
					
					// ステータスメッセージの表示処理
					function showStatusMessage(message, type) {
						statusMessage.textContent = message;
						statusMessage.className = 'status-message';
						statusMessage.classList.add(` + "`status-${type}`" + `);
						
						// 3秒後にメッセージをクリア
						setTimeout(() => {
							statusMessage.textContent = '';
						}, 3000);
					}
				</script>
			</body>
    </html>
  `
