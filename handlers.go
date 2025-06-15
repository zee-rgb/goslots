package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeHTML   = "text/html"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Create initial slots display
	initialSlots := make([]string, 9)
	for i := range initialSlots {
		initialSlots[i] = "?"
	}

	data := map[string]interface{}{
		"Balance": gameState.Balance,
		"Slots":   initialSlots,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}

func handleSpin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get and validate bet
	betStr := r.FormValue("bet")
	bet, err := strconv.ParseUint(betStr, 10, 32)
	if err != nil || bet == 0 {
		w.Header().Set(contentTypeHeader, contentTypeHTML)
		w.Write([]byte(`
            <div class="text-red-500 bg-[#1e3a8a] p-2 rounded text-center mb-4">
                Please enter a valid bet amount
            </div>`))
		return
	}

	// Check if bet is more than balance
	if uint(bet) > gameState.Balance {
		w.Header().Set(contentTypeHeader, contentTypeHTML)
		w.Write([]byte(fmt.Sprintf(`
            <div class="text-red-500 bg-[#1e3a8a] p-2 rounded text-center mb-4">
                Cannot bet more than current balance ($%d)
            </div>`, gameState.Balance)))
		return
	}

	// Process spin
	gameState.Balance -= uint(bet)
	spin := GetSpin(gameState.SymbolArr, 3, 3)
	winningLines := CheckWin(spin, gameState.Multipliers)

	// Calculate winnings
	totalWin := uint(0)
	for _, multiplier := range winningLines {
		totalWin += multiplier * uint(bet)
	}
	gameState.Balance += totalWin

	// Build slot grid HTML
	var slotsHtml strings.Builder
	for _, row := range spin {
		for _, symbol := range row {
			slotsHtml.WriteString(fmt.Sprintf(`
                <div class="aspect-square bg-[#1e3a8a] rounded flex items-center justify-center text-xl font-bold">
                    %s
                </div>`, symbol))
		}
	}

	// Prepare response
	var resultHtml strings.Builder

	// Win/loss message
	if totalWin > 0 {
		multiplier := totalWin / uint(bet)
		resultHtml.WriteString(fmt.Sprintf(`
            <div class="bg-green-500 text-white p-2 rounded text-center mb-4">
                Bet: $%d | Win: $%d | Multiplier: %dx
            </div>`, bet, totalWin, multiplier))
	} else {
		resultHtml.WriteString(fmt.Sprintf(`
            <div class="bg-red-500 text-white p-2 rounded text-center mb-4">
                Bet: $%d | No win
            </div>`, bet))
	}

	// Out-of-band swaps for updates
	resultHtml.WriteString(fmt.Sprintf(`
        <div id="slot-grid" hx-swap-oob="true" class="grid grid-cols-3 gap-2 mb-6">
            %s
        </div>`, slotsHtml.String()))

	resultHtml.WriteString(fmt.Sprintf(`
        <div id="balance-display" hx-swap-oob="true" class="text-2xl font-bold">
            $%d
        </div>`, gameState.Balance))

	// Game over handling
	if gameState.Balance == 0 {
		resultHtml.WriteString(`
            <div class="text-red-500 bg-[#1e3a8a] p-2 rounded text-center mb-4">
                Game Over! Balance: $0
            </div>
            <button hx-post="/reset" 
                    hx-target="#game-container"
                    class="w-full bg-green-600 text-white font-bold py-2 rounded mt-2">
                Reset Game ($200)
            </button>`)
	}

	w.Header().Set(contentTypeHeader, contentTypeHTML)
	w.Write([]byte(resultHtml.String()))
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameState.Balance = 200

	// Render initial game state
	initialSlots := make([]string, 9)
	for i := range initialSlots {
		initialSlots[i] = "?"
	}

	data := map[string]interface{}{
		"Balance": gameState.Balance,
		"Slots":   initialSlots,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, data)
}
