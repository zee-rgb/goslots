package main

import "fmt"

func PrintSpin(spin [][]string) {
	for _, row := range spin {
		for j, symbol := range row {
			fmt.Printf("%s ", symbol)
			if j < len(row)-1 {
				fmt.Print(" | ")
			}
		}
		fmt.Println()
	}
}

func CheckWin(spin [][]string, multiplers map[string]uint) []uint {
	lines := []uint{}

	for _, row := range spin {
		win := true
		checkSymbol := row[0]
		for _, symbol := range row[1:] {
			if checkSymbol != symbol {
				win = false
				break
			}
		}

		if win {
			lines = append(lines, multiplers[checkSymbol])

		} else {
			lines = append(lines, 0)
		}

	}
	return lines

}
