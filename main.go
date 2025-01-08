package main

import (
    "fmt"
    "os"

    "github.com/eiannone/keyboard"
)

// Basic arrow key constants:
const (
    KeyArrowUp = iota
    KeyArrowDown
    KeyArrowLeft
    KeyArrowRight
    KeyEnter
    KeyOther
)

// readKey uses eiannone/keyboard to return a simpler integer representing arrow keys or Enter.
func readKey() int {
    _, key, err := keyboard.GetKey()
    if err != nil {
        return KeyOther
    }
    switch key {
    case keyboard.KeyArrowUp:
        return KeyArrowUp
    case keyboard.KeyArrowDown:
        return KeyArrowDown
    case keyboard.KeyArrowLeft:
        return KeyArrowLeft
    case keyboard.KeyArrowRight:
        return KeyArrowRight
    case keyboard.KeyEnter:
        return KeyEnter
    }
    return KeyOther
}

// navigateBoard modifies x,y based on arrow keys. If user presses ENTER, we return true.
func navigateBoard(x, y *int) bool {
    switch readKey() {
    case KeyArrowUp:
        if *y < 4 {
            *y++
            // Move cursor up visually (2 lines).
            fmt.Print("\x1b[A\x1b[A")
        }
    case KeyArrowDown:
        if *y > 0 {
            *y--
            fmt.Print("\x1b[B\x1b[B")
        }
    case KeyArrowLeft:
        if *x > 0 {
            *x--
            // Move cursor ~6 columns left
            fmt.Print("\x1b[D\x1b[D\x1b[D\x1b[D\x1b[D\x1b[D")
        }
    case KeyArrowRight:
        if *x < 4 {
            *x++
            // Move cursor ~6 columns right
            fmt.Print("\x1b[C\x1b[C\x1b[C\x1b[C\x1b[C\x1b[C")
        }
    case KeyEnter:
        return true
    }
    return false
}

// -------------------------------------------------------------------
// Minimal helper to move the cursor from "center" (2,2) to another (x,y).
// Because after printing, your code repositions the cursor near (2,2),
// we just do little arrow steps to get from (2,2) to (x,y).
// -------------------------------------------------------------------
func moveCursorFromCenterTo(x, y int) {
    // Starting at center => (2,2)
    dx := x - 2
    dy := y - 2

    // If dy > 0 => we need to move up
    //   because in your code, y=0 is bottom, y=4 is top
    //   so "increasing y" means going up on the board
    if dy > 0 {
        for i := 0; i < dy; i++ {
            fmt.Print("\x1b[A\x1b[A") // same logic as navigateBoard
        }
    } else {
        // negative => we move down
        for i := 0; i < -dy; i++ {
            fmt.Print("\x1b[B\x1b[B")
        }
    }

    // If dx > 0 => move left or right?
    // Actually, x=0 is left, x=4 is right, so
    //   if dx>0 => we move right
    //   if dx<0 => we move left
    if dx > 0 {
        for i := 0; i < dx; i++ {
            fmt.Print("\x1b[C\x1b[C\x1b[C\x1b[C\x1b[C\x1b[C")
        }
    } else {
        for i := 0; i < -dx; i++ {
            fmt.Print("\x1b[D\x1b[D\x1b[D\x1b[D\x1b[D\x1b[D")
        }
    }
}

// -------------------------------------------------------------------
// printTeeko prints the board, optionally highlighting a specific bit (for a selected marker).
// We use "\u001b[46;1m" (cyan background) for the highlighted marker (Black or Red).
// -------------------------------------------------------------------
func printTeeko(game Teeko, highlight bitboard) {
    const vertical_separator = "\u001b[36m|"
    const horizontal_separator = "\u001b[36m-------------------------------------"

    var black bitboard
    var red bitboard

    if game.current_player == BlackToMove {
        black = game.player_positions
        red = game.player_positions ^ game.occupied_positions
    } else {
        black = game.player_positions ^ game.occupied_positions
        red = game.player_positions
    }

    var row bitboard = BOARD_LENGTH
    for row != 0 {
        fmt.Print(horizontal_separator, "\n")
        fmt.Print(vertical_separator, "  ", row, "  ")
        var column bitboard = 1
        for column <= BOARD_LENGTH {
            fmt.Print(vertical_separator)

            index := (column - 1)*BOARD_LENGTH + (row - 1)
            mask := bitboard(1) << index
            is_highlight := (highlight & mask) != 0

            switch {
            case (black & mask) != 0:
                if is_highlight {
                    // Selected black marker => cyan background
                    fmt.Print(" \u001b[46;1m \u001b[46;1mB \u001b[0m ")
                } else {
                    fmt.Print("  \u001b[30;1mB\u001b[0m  ")
                }
            case (red & mask) != 0:
                if is_highlight {
                    // Selected red marker => cyan background
					
					fmt.Print(" \u001b[46;1m \u001b[31;1mR \u001b[0m ")
				} else {
                    fmt.Print("  \u001b[31;1mR\u001b[0m  ")
                }
            default:
                if is_highlight {
                    // For empty highlighted square
                    fmt.Print("  \u001b[43;1m-\u001b[0m  ")
                } else {
                    fmt.Print("  -  ")
                }
            }
            column++
        }
        fmt.Print(vertical_separator, "\n")
        row--
    }
    fmt.Print(horizontal_separator, "\n")

    // Print column headers
    fmt.Print(vertical_separator, "  0  ")
    for colIdx := bitboard(1); colIdx <= BOARD_LENGTH; colIdx++ {
        fmt.Print(vertical_separator)
        fmt.Print("  ", colIdx, "  ")
    }
    fmt.Print(vertical_separator, "\n")
    fmt.Print(horizontal_separator, "\n")
    fmt.Print("\u001b[0m")
}

// -------------------------------------------------------------------
// printBoardWithInfo
// 1) Clears screen
// 2) Prints board (no highlight).
// 3) Prints evaluation and instructions
// 4) Moves cursor up & right near the center (like your original).
// -------------------------------------------------------------------
func printBoardWithInfo(game Teeko, marker bitboard) {
    fmt.Print("\033[H\033[2J") // Clear screen

    // Print board with no highlight
    printTeeko(game, marker)
    // Print evaluation
    // Phase-based instructions
    var player_text string
    if game.current_player == BlackToMove {
        player_text = "\u001b[30;1mBlack\u001b[0m"
    } else {
        player_text = "\u001b[31;1mRed\u001b[0m"
    }
	score := evaluate(game)
	switch {
		case score > 0:
			// 1..125 => how many more moves until 126
			fmt.Printf("Current player can force a win in %d moves.\n", WIN - score)
		case score < 0:
			// -125..-1 => how many more moves until opponent hits 126
			fmt.Printf("Opponent can force a win in %d moves.\n", WIN + score)
		default:
			// score == 0 => no forced result either way
			fmt.Println("No forced win.")
	}
	

    if game.phase() == DropPhase {
        fmt.Printf("%s, use arrow-keys to pick a drop; ENTER to confirm.\n", player_text)
    } else {
        fmt.Printf("%s, arrow-keys to pick marker & destination; ENTER to confirm.\n", player_text)
    }

    // Move cursor up ~10 lines
    fmt.Print("\x1b[A\x1b[A\x1b[A\x1b[A\x1b[A\x1b[A\x1b[A\x1b[A\x1b[A\x1b[A")
    // Then move right ~20 columns
    for i := 0; i < 21; i++ {
        fmt.Print("\x1b[C")
    }
}

func computerMove(game *Teeko) {
	if game.phase() == DropPhase {
		var drop bitboard = bestDrop(*game)
		game.dropMarker(drop) 
	} else {
		var move bitboard = bestMove(*game)
		game.moveMarker(move) 
	}
}

func playerMove(game *Teeko) {
	printBoardWithInfo(*game, 0)

	if game.phase() == DropPhase {
		// ---------------- DROP PHASE ----------------
		cursorX, cursorY := 2, 2
		for {
			pressedEnter := navigateBoard(&cursorX, &cursorY)
			if pressedEnter {
				dropMask := bitboard(1) << (uint32(cursorX)*uint32(BOARD_LENGTH) + uint32(cursorY))
				var allDrops bitboard
				for _, d := range game.possibleDrops() {
					allDrops |= d
				}
				if (allDrops & dropMask) != 0 {
					game.dropMarker(dropMask)
					break
				}
			}
		}

	} else {
		// --------------- MOVE PHASE -----------------
		// 1) Select marker
		cursorX, cursorY := 2, 2
		markerX, markerY := -1, -1

		for {
			pressedEnter := navigateBoard(&cursorX, &cursorY)
			if pressedEnter {
				mask := bitboard(1) << (uint32(cursorX)*uint32(BOARD_LENGTH) + uint32(cursorY))
				cpPositions := game.player_positions
				if (cpPositions & mask) != 0 {
					// Valid marker => store coords
					markerX, markerY = cursorX, cursorY

					// === NEW PART: Re-print board (highlight the marker),
					// then move cursor back to that same marker. ===
					// 1) Re-print:
					fmt.Print("\033[H\033[2J")
					highlightMask := bitboard(1) << (uint32(markerX)*uint32(BOARD_LENGTH) + uint32(markerY))
					printBoardWithInfo(*game, highlightMask)

					// Print a quick line about next step

					// 2) The default printing logic tries to place the cursor near center again.
					//    So move from (2,2) => (markerX,markerY).
					//    We just do that directly:
					moveCursorFromCenterTo(markerX, markerY)

					break
				}
			}
		}

		// 2) Select destination from the same screen
		cursorX, cursorY = markerX, markerY

		for {
			pressedEnter := navigateBoard(&cursorX, &cursorY)
			if pressedEnter {
				oldMask := bitboard(1) << (uint32(markerX)*uint32(BOARD_LENGTH) + uint32(markerY))
				newMask := bitboard(1) << (uint32(cursorX)*uint32(BOARD_LENGTH) + uint32(cursorY))
				moveMask := oldMask | newMask

				possibleMoves := game.possibleMoves()
				valid := false
				for _, mv := range possibleMoves {
					if mv == moveMask {
						valid = true
						break
					}
				}
				if valid {
					game.moveMarker(moveMask)
					break
				}
				// else do nothing
			}
		}
	}
}

// -------------------------------------------------------------------
// main
// -------------------------------------------------------------------
func main() {
	loadTable("book.txt")

    // 1) Clear screen
    fmt.Print("\033[H\033[2J\u001b[0m")

    // 2) Initialize keyboard for arrow-key usage
    if err := keyboard.Open(); err != nil {
        fmt.Println("Cannot open keyboard:", err)
        os.Exit(1)
    }
    defer keyboard.Close()

    // 3) Prepare two menu items for arrow-key selection
    menuItems := []string{"Player vs Player", "Player vs AI"}
    selectedIndex := 0

    // A small function to print the menu with highlight
    printMenu := func() {
        fmt.Print("\033[H\033[2J\u001b[0m") // Clear screen each time we reprint
        fmt.Println("Select Game Mode (use arrow keys, ENTER to confirm):\n")
        for i, item := range menuItems {
            if i == selectedIndex {
                // Print highlighted item => e.g. "[ Player vs Player ]"
                fmt.Printf("\u001b[46;1m%s\u001b[0m\n", item)
            } else {
                fmt.Println(item)
            }
        }
        fmt.Println()
    }

    // Print the menu initially
    printMenu()

MenuLoop:
    for {
        _, key, err := keyboard.GetKey()
        if err != nil {
            fmt.Println("Error reading key:", err)
            break
        }
        switch key {
        case keyboard.KeyArrowUp:
            if selectedIndex > 0 {
                selectedIndex--
            }
        case keyboard.KeyArrowDown:
            if selectedIndex < len(menuItems)-1 {
                selectedIndex++
            }
        case keyboard.KeyEnter:
            // Confirm selection
            break MenuLoop
        }
        // Reprint menu after changing selection
        printMenu()
    }

    // Now we have either index=0 => "Player vs Player" or index=1 => "Player vs AI"
    mode := selectedIndex // 0 or 1

    // 4) Load the table for your Teeko game

    // 5) Initialize the game
    game := makeTeeko()

    // 6) Game loop
    for !game.isWin() {
        if mode == 0 {
            // 0 => Player vs Player
            playerMove(&game) // both sides always use playerMove
        } else {
            // 1 => Player vs AI
            if game.current_player == BlackToMove {
                playerMove(&game)   // Black is human
            } else {
                computerMove(&game) // Red is AI
            }
        }
    }

    // 7) Game finished => final output
    fmt.Print("\033[H\033[2J\u001b[0m")
    printTeeko(game, 0)
    fmt.Printf("Game Over! Winner is: %s\n", func() string {
        if game.current_player == BlackToMove {
            return "\u001b[31;1mRed\u001b[0m"
        }
        return "\u001b[30;1mBlack\u001b[0m"
    }())

    // Final reset
    fmt.Print("\u001b[0m")
}

