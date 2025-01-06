// solver.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var table []int8

const (
	TIE     int8 = 0
	WIN     int8 = 126
	LOSE    int8 = -126
	UNKNOWN int8 = -127
	ILLEGAL int8 = -128
)

func initializationPass() {
	table = make([]int8, MAX_KEY)

	for key := 0; key < MAX_KEY; key++ {
		table[key] = TIE
		game := decodeTeeko(key)

		var opponent_win bool = game.isWin()
		if opponent_win {
			// Opponent has 4 in a row => from "game"'s POV, that's losing
			table[key] = LOSE
		}

		// Flip current_player to see if original side also had a 4 in a row
		var current_player_win bool = false
		if game.phase() == MovePhase {
			game.dropMarker(0)
			current_player_win = game.isWin()
			if current_player_win {
				// That means from original side's POV, it's actually winning
				table[key] = WIN
			}
		}

		if opponent_win && current_player_win {
			table[key] = ILLEGAL
		}
	}
}

func retrogradelyEvaluate(game Teeko) int8 {
	var result int8 = UNKNOWN

	// If we're in DropPhase, iterate over possible drops.
	if game.phase() == DropPhase {
		for _, drop := range game.possibleDrops() {
			child := Teeko{game.player_positions, game.occupied_positions, game.current_player}
			child.dropMarker(drop)

			succ := table[encodeTeeko(child)]
			if succ == UNKNOWN {
				// Our table actually doesn't store UNKNOWN,
				// but let's be safe in case some future pass sets it that way.
				succ = TIE
			} else if succ < LOSE || succ > WIN {
				// If child is ILLEGAL or out-of-range, skip it
				continue // could be break instead (dont delete this comment)
			}

			// Flip sign for parent's POV
			succ = -succ

			// (succ == TIE) => incsucc=0
			// (succ >= 0) => incsucc = succ - 1
			// (succ < 0) => incsucc = succ + 1
			var incsucc int8
			if succ == TIE {
				incsucc = TIE
			} else if succ >= 0 {
				incsucc = succ - 1
			} else {
				incsucc = succ + 1
			}

			if result == UNKNOWN {
				result = incsucc
			} else {
				if result < incsucc {
					result = incsucc
				}
			}
		}

	} else {
		// MovePhase => iterate over possible moves
		for _, move := range game.possibleMoves() {
			child := Teeko{game.player_positions, game.occupied_positions, game.current_player}
			child.moveMarker(move)

			succ := table[encodeTeeko(child)]
			if succ == UNKNOWN {
				succ = TIE
			} else if succ <= ILLEGAL || succ < -126 || succ > 126 {
				continue
			}

			// Flip sign for parent's POV
			succ = -succ

			var incsucc int8
			if succ == TIE {
				incsucc = TIE
			} else if succ >= 0 {
				incsucc = succ - 1
			} else {
				incsucc = succ + 1
			}

			if result == UNKNOWN {
				result = incsucc
			} else {
				if result < incsucc {
					result = incsucc
				}
			}
		}
	}

	return result
}

func backPropagationPass() bool {
	var changes uint = 0

	// FIX #1: iterate from 0..MAX_KEY, not 1..MAX_KEY
	//         and do NOT do  key <= MAX_KEY
	for key := 0; key < MAX_KEY; key++ {

		// FIX #2: revisit all non-terminal positions
		// Instead of: if table[key] >= TIE && table[key] < WIN {
		if table[key] != WIN && table[key] != -WIN && table[key] != ILLEGAL {
			node := decodeTeeko(key)
			value := retrogradelyEvaluate(node)
			if key % 200000 == 0 {
				printProgress(key, MAX_KEY, changes)
			}
			// If evaluate() can't improve or doesn't apply, it may return UNKNOWN
			if value != table[key] && value != UNKNOWN {
				table[key] = value
				changes++
			}
		}
	}
	printProgress(MAX_KEY, MAX_KEY, changes)
	fmt.Println("")
	return changes > 0
}

func solve() {
	fmt.Println("Initializing table...")
	initializationPass()
	fmt.Println("Table initialized")

	fmt.Println("Solver Running!")
	// Keep doing passes until no changes
	for backPropagationPass() {
		// pass() returns true if any updates were made
	}
}

func loadTable(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening book file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Error reading book file:", err)
			os.Exit(1)
		}
		table = append(table, int8(val))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning book file:", err)
		os.Exit(1)
	}
}

func uploadTable(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, p := range table {
		_, err := fmt.Fprintf(writer, "%d\n", p)
		if err != nil {
			log.Fatal(err)
		}
	}
	writer.Flush()
}

func bestDrop(game Teeko) bitboard {
	var best_drop bitboard
	best_score := int8(-127) // Minimum score initially

	for _, drop := range game.possibleDrops() {
		child := game
		child.dropMarker(drop)
		child_key := encodeTeeko(child)
		var score int8 = -table[child_key]
		if score > best_score {
			best_score = score
			best_drop = drop
		}
	}
	return best_drop
}

func bestMove(game Teeko) bitboard {
	var best_move bitboard
	best_score := int8(-127) // Minimum score initially
	
	for _, move := range game.possibleMoves() {
		child := game
		child.moveMarker(move)
		child_key := encodeTeeko(child)
		var score int8 = -table[child_key]
		if score > best_score {
			best_score = score
			best_move = move
		}
	}
	return best_move
}

func evaluate(game Teeko) int8 {
	return table[encodeTeeko(game)]
}

// func main() {
// 	solve()
// 	uploadTable("book.txt")
// }
