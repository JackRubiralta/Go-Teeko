package main

import (
    "fmt"
)
func printProgress(current, total int, changes uint) {
	const PBWIDTH = 50
    // Fraction done
    fraction := float64(current) / float64(total)
    // Number of '=' to display
    filled := int(fraction * PBWIDTH)
    // Print the bar
    fmt.Printf("\r[")
    for i := 0; i < filled; i++ {
        fmt.Print("=")
    }
    for i := filled; i < PBWIDTH; i++ {
        fmt.Print(" ")
    }
    fmt.Printf("] %.2f%% (%d/%d) Changes: %d", fraction*100.0, current, total, changes)
}

func comb(n, k int) int {
	if k < 0 || k > n || n > 25 {
		return 0
	}
	if k == 0 || k == n {
		return 1
	}
	// For efficiency, do k = min(k, n-k)
	if k > n - k {
		k = n - k
	}
	result := 1
	for i := 0; i < k; i++ {
		result = result * (n - i) / (i + 1)
	}
	return result
}

func arrayToBitboard(pos []int) bitboard {
	var bb bitboard
	for _, p := range pos {
		bb |= (1 << p)
	}
	return bb
}

func bitboardToArray(bb bitboard) []int {
	var positions []int
	for i := 0; i < 25; i++ {
		if (bb & (1 << i)) != 0 {
			positions = append(positions, i)
		}
	}
	return positions
}

func popCount(bb bitboard) int {
	count := 0
	for bb != 0 {
		bb &= (bb - 1)
		count++
	}
	return count
}
