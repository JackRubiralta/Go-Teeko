package main

// bitboard for storing piece positions
type bitboard uint32
type GameMode int

const (
	Regular GameMode = iota
	Advanced
)

// Board is 5x5 => 25 bits
const TOTAL_MARKER int = 8
const BOARD_LENGTH bitboard = 5
const BOARD_SIZE int = int(BOARD_LENGTH * BOARD_LENGTH)
const BOARD_MASK bitboard = 0b1111111111111111111111111 // make this be set

const GAME_MODE = Advanced

// PHASE enum
type Phase int

const (
	DropPhase Phase = iota
	MovePhase
)

// Player enum
type Player int

const (
	BlackToMove Player = iota
	RedToMove
)

// GameMode enum

// Teeko struct
// changes from player_positions to player_positions and occupied_positions
type Teeko struct {
	player_positions   bitboard // squares for "current player"
	occupied_positions bitboard // squares occupied by both sides
	current_player     Player
}

// Constructor
func makeTeeko() Teeko {
	var game Teeko
	game.player_positions = 0
	game.occupied_positions = 0
	game.current_player = BlackToMove
	return game
}

// Figure out the current phase (drop or move)
func (game *Teeko) phase() Phase {
	var n bitboard
	n = game.player_positions
	// Count how many bits are set; if >= 4 => move phase
	// This trick unsets 3 bits and checks if there's something left
	n = n & (n - 1)
	n = n & (n - 1)
	n = n & (n - 1)

	if n != 0 {
		return MovePhase
	} else {
		return DropPhase
	}
}

// Drop a marker onto the board
func (game *Teeko) dropMarker(drop bitboard) {
	// Toggle current_player
	if game.current_player == BlackToMove {
		game.current_player = RedToMove
	} else {
		game.current_player = BlackToMove
	}

	// Swap ownership (original design)
	game.player_positions ^= game.occupied_positions

	// Add the new drop bit
	game.occupied_positions |= drop
}

// Move a marker on the board using a single parameter with two bits set
func (game *Teeko) moveMarker(move bitboard) {
	// Toggle current_player
	if game.current_player == BlackToMove {
		game.current_player = RedToMove
	} else {
		game.current_player = BlackToMove
	}

	// Swap ownership (original design)
	game.player_positions ^= game.occupied_positions

	// XOR the old and new bits in or out of the occupied_positions
	game.occupied_positions ^= move
}

// Check if the opponent has a winning shape
func (game *Teeko) isWin() bool {
	// Opponent is everything in occupied_positions except for current player's bits
	var opponent_positions bitboard = game.player_positions ^ game.occupied_positions

	// horizontal 4 in a row
	if ((opponent_positions & (opponent_positions >> 1) & (opponent_positions >> 2) & (opponent_positions >> 3)) &
		0b0001100011000110001100011) != 0 {
		return true
	}

	// vertical 4 in a row
	if ((opponent_positions &
		(opponent_positions >> (BOARD_LENGTH * 1)) &
		(opponent_positions >> (BOARD_LENGTH * 2)) &
		(opponent_positions >> (BOARD_LENGTH * 3))) &
		0b0000000000000001111111111) != 0 {
		return true
	}

	// diagonal (one direction) 4 in a row
	if ((opponent_positions &
		(opponent_positions >> ((BOARD_LENGTH + 1) * 1)) &
		(opponent_positions >> ((BOARD_LENGTH + 1) * 2)) &
		(opponent_positions >> ((BOARD_LENGTH + 1) * 3))) &
		0b0001100011000110001100011) != 0 {
		return true
	}

	// diagonal (other direction) 4 in a row
	if ((opponent_positions &
		(opponent_positions >> ((BOARD_LENGTH - 1) * 1)) &
		(opponent_positions >> ((BOARD_LENGTH - 1) * 2)) &
		(opponent_positions >> ((BOARD_LENGTH - 1) * 3))) &
		0b1100011000110001100011000) != 0 {
		return true
	}

	// 2x2 square
	if ((opponent_positions &
		(opponent_positions >> 1) &
		(opponent_positions >> BOARD_LENGTH) &
		(opponent_positions >> (BOARD_LENGTH + 1))) &
		0b011110111101111011110111101111) != 0 {
		return true
	}

	// -------------------------------------------------------------
	// ADVANCED checks for 3x3, 4x4, 5x5 squares
	// (Bitmasks might need adjusting to match your board's bit layout)
	// -------------------------------------------------------------
	if GAME_MODE == Advanced {
		// We'll define a helper variable 'm' for these checks
		var m bitboard

		// 3x3 square
		m = opponent_positions & (opponent_positions >> 2) & (opponent_positions >> 10) & (opponent_positions >> 12)
		if (m & 0b001110011100111001110011100111) != 0 {
			return true
		}

		// 4x4 square
		m = opponent_positions & (opponent_positions >> 3) & (opponent_positions >> 15) & (opponent_positions >> 18)
		if (m & 0b000110001100011000110001100011) != 0 {
			return true
		}

		// 5x5 square (the entire board)
		// Checking corners = 17825809 is from your snippet
		// You might want to do a more precise bitmask check
		if (opponent_positions & 17825809) == 17825809 {
			return true
		}
	}

	return false
}

// possibleMoves returns all legal "move" bit positions for each marker of current_player
func (game *Teeko) possibleMoves() []bitboard {
	
	var player_positions bitboard = game.player_positions

	var unoccupied_positions bitboard = game.occupied_positions ^ BOARD_MASK

	var possible_moves []bitboard

	for player_positions != 0 {
		// isolate the least significant set bit
		var current_marker bitboard
		current_marker = player_positions ^ (player_positions & (player_positions - 1))
		player_positions = player_positions ^ current_marker

		// We’ll try each possible shift from your code:
		var move bitboard

		// left
		move = ((current_marker << 1) & 0b1111011110111101111011110) & unoccupied_positions
		if move != 0 {
			possible_moves = append(possible_moves, move|current_marker)
		}

		// up-left
		move = ((current_marker << 6) & 0b1111011110111101111011110) & unoccupied_positions
		if move != 0 {
			possible_moves = append(possible_moves, move|current_marker)
		}

		// down-right
		move = ((current_marker >> 4) & 0b1111011110111101111011110) & unoccupied_positions
		if move != 0 {
			possible_moves = append(possible_moves, move|current_marker)
		}

		// right
		move = ((current_marker >> 1) & 0b0111101111011110111101111) & unoccupied_positions
		if move != 0 {
			possible_moves = append(possible_moves, move|current_marker)
		}

		// up-right
		move = ((current_marker >> 6) & 0b0111101111011110111101111) & unoccupied_positions
		if move != 0 {
			possible_moves = append(possible_moves, move|current_marker)
		}

		// down-left
		move = ((current_marker << 4) & 0b0111101111011110111101111) & unoccupied_positions
		if move != 0 {
			possible_moves = append(possible_moves, move|current_marker)
		}

		// south
		move = (current_marker >> 5) & unoccupied_positions
		if move != 0 {
			possible_moves = append(possible_moves, move|current_marker)
		}

		// north
		move = (current_marker << 5) & unoccupied_positions
		if move != 0 {
			possible_moves = append(possible_moves, move|current_marker)
		}
	}
	return possible_moves
}

// possibleDrops returns all empty squares for dropping a new piece
func (game *Teeko) possibleDrops() []bitboard {

	var empty_positions bitboard
	empty_positions = game.occupied_positions ^ BOARD_MASK

	var possible_drops []bitboard

	for empty_positions != 0 {
		var current_position bitboard
		current_position = empty_positions ^ (empty_positions & (empty_positions - 1))
		empty_positions = empty_positions ^ current_position
		possible_drops = append(possible_drops, current_position)
	}
	return possible_drops
}

