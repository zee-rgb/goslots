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
		w.Write([]byte(`<div class="bg-red-500 text-white p-2 rounded text-center">Please enter a valid bet amount</div>`))
		return
	}

	// Check if bet is more than balance
	if uint(bet) > gameState.Balance {
		w.Header().Set(contentTypeHeader, contentTypeHTML)
		w.Write([]byte(fmt.Sprintf(`
            <div class="bg-red-500 text-white p-2 rounded text-center">
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

	// Prepare slot grid update
	var slotHTML strings.Builder
	for _, row := range spin {
		for _, symbol := range row {
			slotHTML.WriteString(fmt.Sprintf(`
                <div class="text-2xl font-bold p-4 slot-cell rounded-lg text-center">%s</div>`, symbol))
		}
	}

	// Prepare response
	var resultHTML strings.Builder

	// Win/loss message
	if totalWin > 0 {
		multiplier := totalWin / uint(bet)
		resultHTML.WriteString(fmt.Sprintf(`
            <div class="bg-green-500 text-white p-2 rounded text-center">
                Bet: $%d | Win: $%d | Multiplier: %dx
            </div>`, bet, totalWin, multiplier))
	} else {
		resultHTML.WriteString(fmt.Sprintf(`
            <div class="bg-red-500 text-white p-2 rounded text-center">
                Bet: $%d | No win
            </div>`, bet))
	}

	// HTMX out-of-band swaps
	resultHTML.WriteString(fmt.Sprintf(`
        <div id="slot-grid" hx-swap-oob="true" class="grid grid-cols-3 gap-4 mb-6">
            %s
        </div>`, slotHTML.String()))

	resultHTML.WriteString(fmt.Sprintf(`
        <div id="balance" hx-swap-oob="true" class="text-2xl font-bold">
            $%d
        </div>`, gameState.Balance))

	// Game over check
	if gameState.Balance == 0 {
		resultHTML.WriteString(`
            <div class="mt-4">
                <div class="bg-red-500 text-white p-2 rounded text-center mb-2">Game Over!</div>
                <button hx-post="/reset" 
                        hx-target="body"
                        class="w-full bg-green-600 text-white font-bold py-2 rounded">
                    Reset Game ($200)
                </button>
            </div>`)
	}

	w.Header().Set(contentTypeHeader, contentTypeHTML)
	w.Write([]byte(resultHTML.String()))
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameState.Balance = 200
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
