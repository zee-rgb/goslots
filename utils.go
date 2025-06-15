package main

import (
	"fmt"
	"math/rand"
)

func GetName() string {
	name := ""

	fmt.Println("Welcome to the Casino...")
	fmt.Printf("Enter your name: ")
	_, err := fmt.Scanln(&name)
	if err != nil {
		return ""
	}

	fmt.Printf("Welcome %s, let's play!\n", name)
	return name
}

func GetBet(balance uint) uint {
	var bet uint

	for {
		fmt.Printf("Enter your bet, or 0 to quit (balance: $%d): ", balance)
		fmt.Scanln(&bet)

		if bet > balance {
			fmt.Println("Bet can not exceed balance")
		} else {
			break
		}
	}
	return bet
}

func GetRandNum(min int, max int) int {
	randNum := rand.Intn(max-min+1) + min
	return randNum
}

func GetSpin(reel []string, rows int, cols int) [][]string {
	resultSpin := [][]string{}

	for i := 0; i < rows; i++ {
		resultSpin = append(resultSpin, []string{})
	}

	for col := 0; col < cols; col++ {
		selected := map[int]bool{}
		for row := 0; row < rows; row++ {
			for {
				index := GetRandNum(0, len(reel)-1)
				_, exists := selected[index]
				if !exists {
					selected[index] = true
					resultSpin[row] = append(resultSpin[row], reel[index])
					break
				}
			}
		}

	}
	return resultSpin
}
