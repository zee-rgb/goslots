package main

import (
	"fmt"
)

func generateSymbolArray(symbols map[string]uint) []string {
	symbolArray := []string{}
	for symbol, count := range symbols {
		for i := uint(0); i < count; i++ {
			symbolArray = append(symbolArray, symbol)
		}
	}
	return symbolArray
}

func main() {

	symbols := map[string]uint{
		"A": 4,
		"B": 7,
		"C": 12,
		"D": 20,
	}
	multipliers := map[string]uint{
		"A": 20,
		"B": 10,
		"C": 5,
		"D": 2,
	}

	symbolArr := generateSymbolArray(symbols)

	balance := uint(200)
	GetName()

	for balance > 0 {
		bet := GetBet(balance)
		if bet == 0 {
			break
		}

		balance -= bet
		spin := GetSpin(symbolArr, 3, 3)
		PrintSpin(spin)
		winningLines := CheckWin(spin, multipliers)
		fmt.Println(winningLines)
		for i, multi := range winningLines {
			win := multi * bet
			balance += win
			if multi > 0 {
				fmt.Printf("You won $%d (%dx) on line %d\n", win, multi, i+1)
			}
		}

	}

	fmt.Printf("You have $%d left to play with\n", balance)
}
