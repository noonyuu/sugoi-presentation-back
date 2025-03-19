package main

import (
	"fmt"

	"github.com/noonyuu/websocket-server/router"
)

func main() {
	fmt.Println("Hello, World!")

	router.Server()

	// サーバーの終了を待機する
	select {}
}
