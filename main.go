package main

import (
	"fmt"
	"net/http"

	"goslots/game"
)

var gameState *game.GameState

func main() {
	gameState = game.NewGameState()
	gameState.InitSymbolArray()

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/spin", handleSpin)
	http.HandleFunc("/reset", handleReset)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
