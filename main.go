package main

import (
	"fmt"	
)

func printBoardWithInfo(game Teeko) {
	game.printBoard()
	position_score := evaluate(game)
	fmt.Printf("\nCurrent Board Evaluation: %d\n", position_score)

	if position_score == WIN {
		fmt.Println("Victory is imminent for the current player!")
	} else if position_score == LOSE {
		fmt.Println("The current player is in a losing position!")
	} else {
		fmt.Printf("It will take %d moves for the current player to reach the best outcome.\n", abs(position_score))
	}
}

func abs(value int8) int8 {
	if value < 0 {
		return -value
	}
	return value
}

func main() {
	// Load the precomputed book
	loadTable("book.txt")
	// Initialize the game
	var game Teeko
	game = makeTeeko()

	for game.isWin() == false {
		// Clear the screen and display the board with info
		fmt.Print("\033[H\033[2J")
		printBoardWithInfo(game)

		var player_text string
		if game.current_player == BlackToMove {
			player_text = "\u001b[30;1mBlack\u001b[0m"
		} else {
			player_text = "\u001b[31;1mRed\u001b[0m"
		}

		if game.phase() == DropPhase {
			fmt.Printf("%s, enter drop (e.g., 12 for X=1, Y=2): ", player_text)
			var input string
			fmt.Scanln(&input)

			// Convert input chars '1'..'5' to 0-based indexes
			position_x := uint32(input[0] - '0' - 1)
			position_y := uint32(input[1] - '0' - 1)

			var drop bitboard = 1 << (position_x*uint32(BOARD_LENGTH) + position_y)
			game.dropMarker(drop)
		} else {
			fmt.Printf("%s, enter move (e.g., 0102 for marker (0,1) -> (0,2)): ", player_text)
			var input string
			fmt.Scanln(&input)

			marker_y := uint32(input[0] - '0' - 1)
			marker_x := uint32(input[1] - '0' - 1)
			destination_y := uint32(input[2] - '0' - 1)
			destination_x := uint32(input[3] - '0' - 1)

			var old_pos bitboard = 1 << (marker_y*uint32(BOARD_LENGTH) + marker_x)
			var new_pos bitboard = 1 << (destination_y*uint32(BOARD_LENGTH) + destination_x)

			game.moveMarker(old_pos | new_pos)
		}
		fmt.Println("Great move! Let's see what happens next...")
	}

	// Final board and winner announcement
	fmt.Print("\033[H\033[2J")
	printBoardWithInfo(game)
	fmt.Printf("Game Over! Winner is: %s\n", func() string {
		if game.current_player == BlackToMove {
			return "\u001b[31;1mRed\u001b[0m"
		}
		return "\u001b[30;1mBlack\u001b[0m"
	}())
}
