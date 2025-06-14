package main

import "fmt"

func getName() string {
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

func getBet(balance uint) uint {
	var bet uint

	for true {
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

func main() {
	balance := uint(200)
	getName()

	for balance > 0 {
		bet := getBet(balance)
		if bet == 0 {
			break
		}
		balance -= bet
	}
	fmt.Printf("You have $%d left to play with\n", balance)

}
