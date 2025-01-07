package main

import (
	"log"
)



func rankCombination(subset []int, n int) int {
	rank := 0
	k := len(subset)
	if k == 0 {
		return 0
	}
	previous := -1
	for i := 0; i < k; i++ {
		x := subset[i]
		for v := previous + 1; v < x; v++ {
			rank += comb(n-1-v, k-1-i)
		}
		previous = x
	}
	return rank
}

func unrankCombination(rank, k, n int) []int {
	subset := make([]int, 0, k)
	current := 0
	for i := 0; i < k; i++ {
		for v := current; v < n; v++ {
			c := comb(n-1-v, k-1-i)
			if rank < c {
				subset = append(subset, v)
				current = v + 1
				break
			} else {
				rank -= c
			}
		}
	}
	return subset
}

// ------------------------------------------------------------------- //
// Offsets for valid (b, r) pairs: b=#player, r=#opponent, with b=r or b=r+1, b+r <=8
var (
    offset_po = func() [5][5]int {
        var po [5][5]int
        accum := 0
        for total := 0; total <= 8; total++ {
            for o := 0; o <= 4; o++ {
                p := total - o
                if p < 0 || p > 4 {
                    continue
                }
                if !(o == p || o == p+1) {
                    continue
                }
                po[o][p] = accum

                accum += comb(25, o) * comb(25-o, p)
            }
        }
        MAX_KEY = accum // Note: MAX_KEY must be a global variable for this to work
        return po
    }()
    MAX_KEY int
)





// ------------------------------------------------------------------- //
// encodeTeeko: Input is a Teeko struct. We figure out the "opponent" bits
// as  player_positions ^ occupied_positions. Then we apply the standard combination logic
// to produce a unique integer in [0..MAX_KEY).



func encodeTeeko(game Teeko) int {
    // The "player" bitboard:
    player_positions := game.player_positions
    // The "opponent" bitboard:
    opponent_mask := game.player_positions ^ game.occupied_positions

    // We keep these exactly the same:
    // we can get ri
    player_count := popCount(player_positions)
    opponent_count := popCount(opponent_mask)

    // Changed from offset_po[player_count][opponent_count] to:
    base := offset_po[opponent_count][player_count]

    //
    // Now we swap the names in the arrays, because previously
    // 'player_positions' was the second line; with roles reversed,
    // we rename local arrays so the code “makes sense.”
    //
    // Old code had:
    //    player_pos := bitboardToArray(player_positions)
    //    opponent_pos := bitboardToArray(opponent_mask)
    // Here, we reverse who gets called “player_pos” vs “opponent_pos”
    // to stay consistent with the swapped bitboards above.
    //

    // "opponent_pos" now comes from 'opponent_mask'
    opponent_pos := bitboardToArray(opponent_mask)
    // The “opponent” rank (formerly player_rank in the old code)
    opponent_rank := rankCombination(opponent_pos, 25)

    // "player_pos" now comes from 'player_positions'
    player_pos := bitboardToArray(player_positions)

    // Build the 'used' array based on opponent_pos
    used := make([]bool, 25)
    for _, idx := range opponent_pos {
        used[idx] = true
    }

    var leftover []int
    for i := 0; i < 25; i++ {
        if !used[i] {
            leftover = append(leftover, i)
        }
    }

    // Convert the 'player_pos' squares into relative indices w.r.t. leftover
    rel_player := make([]int, player_count)
    for i := 0; i < player_count; i++ {
        val := player_pos[i]
        j := 0
        for j < len(leftover) && leftover[j] != val {
            j++
        }
        rel_player[i] = j
    }

    // The “player” rank (formerly opponent_rank in old code)
    player_rank := rankCombination(rel_player, 25 - opponent_count)

    // Number of ways to place 'player_count' among leftover squares
    ways_for_player := comb(25 - opponent_count, player_count)

    // Final local rank
    local_rank := opponent_rank*ways_for_player + player_rank

    return base + local_rank
}

// ------------------------------------------------------------------- //
// decodeTeeko: from an integer key -> a new Teeko struct.
// We'll interpret the "player" bits vs "opponent" bits, then
// set game.current_player = BlackToMove if total # markers is even, else RedToMove.

func decodeTeeko(key int) Teeko {
    // After flipping roles, we’ll find opponent_count first, then player_count.
    var opponent_count, player_count, base int
    found := false

outer:
    // We keep the same total loop, but now treat `o` as the opponent_count 
    // and `p` as the player_count.
    for total := 0; total <= 8; total++ {
        for o := 0; o <= 4; o++ {
            p := total - o
            if p < 0 || p > 4 {
                continue
            }
            // Flip the logic that used to check (p == o || p == o+1):
            // now it’s (o == p || o == p+1) because we swapped roles.
            if !(o == p || o == p+1) {
                continue
            }
            // Also flip offset_po[p][o] → offset_po[o][p]
            off := offset_po[o][p]
            // Likewise flip comb(25, p)*comb(25-p, o) → comb(25, o)*comb(25-o, p)
            ways := comb(25, o) * comb(25-o, p)
            if key >= off && key < off+ways {
                opponent_count = o
                player_count   = p

                base  = off
                found = true
                break outer
            }
        }
    }

    if !found {
        log.Fatalf("decodeTeeko: key=%d out of range (MAX_KEY=%d)", key, MAX_KEY)
    }

    local_rank := key - base

    // Previously was ways_for_opp = comb(25 - player_count, opponent_count).
    // Now we flip to ways_for_player = comb(25 - opponent_count, player_count).
    ways_for_player := comb(25 - opponent_count, player_count)

    // Flip player_rank ↔ opponent_rank:
    //   old: player_rank = local_rank / ways_for_opp
    //        opponent_rank = local_rank % ways_for_opp
    //   new: opponent_rank = local_rank / ways_for_player
    //        player_rank   = local_rank % ways_for_player
    opponent_rank := local_rank / ways_for_player
    player_rank   := local_rank % ways_for_player

    //
    // --- Reconstruct Opponent Squares (formerly "player" squares) ---
    //
    // Old code:
    //   player_pos = unrankCombination(player_rank, player_count, 25)
    // We flip that to:
    opponent_pos   := unrankCombination(opponent_rank, opponent_count, 25)
    opponent_mask  := arrayToBitboard(opponent_pos)

    // leftover squares used by opponent
    used := make([]bool, 25)
    for _, idx := range opponent_pos {
        used[idx] = true
    }
    var leftover []int
    for i := 0; i < 25; i++ {
        if !used[i] {
            leftover = append(leftover, i)
        }
    }

    //
    // --- Reconstruct Player Squares (formerly "opponent" squares) ---
    //
    // Old code:
    //   rel_opp = unrankCombination(opponent_rank, opponent_count, 25 - player_count)
    //   for i in opponent_count ...
    // Flip that to:
    rel_player := unrankCombination(player_rank, player_count, 25 - opponent_count)
    player_pos := make([]int, player_count)
    for i := 0; i < player_count; i++ {
        player_pos[i] = leftover[rel_player[i]]
    }
    player_positions := arrayToBitboard(player_pos)

    //
    // --- Figure out current_player as before, but with swapped references ---
    //
    // The old code used if (player_count + opponent_count < 8) { ... } else { ... }
    // That logic remains the same; we’re just calling them “opponent_count + player_count.”
    var current_player Player
    if opponent_count+player_count < 8 {
        // drop phase
        if opponent_count == player_count {
            current_player = BlackToMove
        } else {
            current_player = RedToMove
        }
    } else {
        // move phase
        current_player = BlackToMove
    }

    //
    // --- Final assembly of the Teeko struct ---
    //
    var game Teeko
    occupied_positions := opponent_mask | player_positions

    // Now we do the final swap so that:
    //    game.player_positions = player_positions
    // instead of opponent_mask.
    game.player_positions   = player_positions
    game.occupied_positions = occupied_positions
    game.current_player     = current_player

    return game
}

