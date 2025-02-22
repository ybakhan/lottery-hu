package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// LotteryPick is a 128-bit mask representing numbers of a lottery pick.
// It contains two uint64 numbers. Each bit of both numbers represents a lottery number.
// Bits of the first number represent lottery numbers 1-64
// Bits of the second number represent lottery numbers 65-128
type LotteryPick [2]uint64

// lottery represents the configuration of a Hungarian lottery
type lottery struct {
	numberOfPicks int // number of picks per player (e.g. 5)
	minMatches    int // minimum matches required to win (e.g. 2) cannot be negative or zero
	minPick       int // minimum valid pick number (e.g. 1)  cannot be negative or zero
	maxPick       int // maximum valid pick number (e.g. 90) cannot be more than 128
}

// ProcessPlayerPicks reads player lottery picks from a file and converts them to bit masks.
// Each line in the file should contain space separated integers matching lottery numberOfPicks.
// Returns a slice of LotteryPicks bit masks or an error if file reading fails.
func (l lottery) ProcessPlayerPicks(file *os.File) ([]LotteryPick, error) {
	playerIndex := 1
	var picks []LotteryPick
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		pickEntry := scanner.Text()
		pickStr := strings.Fields(pickEntry)

		if len(pickStr) != l.numberOfPicks {
			// Skip pick entry with incorrect number of picks
			fmt.Printf("Player %d ignored expected %d numbers got %d\n", playerIndex, l.numberOfPicks, len(pickStr))
			continue
		}

		var pick LotteryPick
		var valid = true

		for _, numStr := range pickStr {

			num, err := strconv.Atoi(numStr)
			if err != nil {
				valid = false
				fmt.Printf("Player %d ignored error converting %s to integer\n", playerIndex, numStr)
				break
			}

			if num < l.minPick || num > l.maxPick {
				valid = false
				fmt.Printf("Player %d ignored number %d out of lottery range\n", playerIndex, num)
				break
			}

			if num <= 64 {
				pick[0] |= 1 << (num - 1) // Bits 0-63 for picks 1-64
			} else {
				pick[1] |= 1 << (num - 65) // Bits 0-127 for picks 65-128
			}
		}

		if valid {
			picks = append(picks, pick)
		}
		playerIndex++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return picks, nil
}

// MatchPicks processes a winning lottery pick entry and matches it against all player picks.
// It computes the number of winners for each match count ranging from minMatches to numberOfPicks
// and prints the results with execution time.
// Returns an error if lottery pick entry is invalid.
func (l lottery) MatchPicks(winningEntry string, playerPicks []LotteryPick) error {
	winningPick, err := l.parseWinningEntry(winningEntry)
	if err != nil {
		return err
	}

	start := time.Now()
	winners := l.matchPicks(winningPick, playerPicks)

	elapsed := time.Since(start)
	fmt.Printf("\nMatching lottery picks took %dms\n", int64(float64(elapsed)/float64(time.Millisecond)+0.5))

	fmt.Printf("\n%-16s\t%s\n", "Numbers matching", "Winners")
	for i := l.numberOfPicks; i >= l.minMatches; i-- {
		fmt.Printf("%-16d\t%d\n", i, winners[i])
	}

	return nil
}

// parseWinningEntry converts a space separated string of a winning lottery pick into a LotteryPick bit mask.
// The input must contain exactly numberOfPicks integers within range [minPick, maxPick].
// Returns LotteryPick bit mask or an error if the input is invalid.
func (l lottery) parseWinningEntry(winningEntry string) (LotteryPick, error) {

	winningPick := strings.Fields(winningEntry)
	if len(winningPick) != l.numberOfPicks {
		fmt.Printf("Enter %d lottery picks\n", l.numberOfPicks)
		return LotteryPick{}, fmt.Errorf("invalid number of lottery picks")
	}

	var pick LotteryPick
	for _, numStr := range winningPick {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			fmt.Printf("Lottery pick entry not a number %s\n", numStr)
			return LotteryPick{}, err
		}

		if num < l.minPick || num > l.maxPick {
			fmt.Printf("Lottery pick entry out of range %d\n", num)
			return LotteryPick{}, fmt.Errorf("lottery pick entry out of range %d", num)
		}

		if num <= 64 {
			pick[0] |= 1 << (num - 1) // Bits 0-63 for picks 1-64
		} else {
			pick[1] |= 1 << (num - 65) // Bits 0-127 for picks 65-128
		}
	}
	return pick, nil
}

// matchPicks computes number of matches between a winning lottery pick and picks of all players.
// Uses go routines, bitwise AND, and popcount to efficiently count matches.
// Returns a map where keys are match counts and values are the number of players with that count.
// Match counts less than minMatches are not included in result
func (l lottery) matchPicks(winningPick LotteryPick, playerPicks []LotteryPick) map[int]int {
	// Use available CPU cores for parallelization
	numRoutines := runtime.NumCPU()
	if numRoutines > runtime.GOMAXPROCS(0) {
		numRoutines = runtime.GOMAXPROCS(0) // Cap at available OS threads
	}

	if numRoutines < 1 {
		numRoutines = 1
	}

	chunkSize := (len(playerPicks) + numRoutines - 1) / numRoutines
	resultChan := make(chan map[int]int, numRoutines)

	// Split playerPicks into chunks and process in parallel
	for i := 0; i < numRoutines; i++ {

		start := i * chunkSize

		end := start + chunkSize
		if end > len(playerPicks) {
			end = len(playerPicks) // Last chunk takes remainder
		}

		if start >= end {
			break // No more work to do
		}

		go func(picks []LotteryPick) {
			winners := make(map[int]int)

			// Process this chunk of player picks
			for _, pick := range picks {
				matches := bits.OnesCount64(pick[0]&winningPick[0]) + bits.OnesCount64(pick[1]&winningPick[1])
				if matches >= l.minMatches {
					winners[matches]++
				}
			}

			resultChan <- winners

		}(playerPicks[start:end])
	}

	// Aggregate results from all goroutines
	winners := make(map[int]int, l.numberOfPicks-l.minMatches+1)

	for i := 0; i < numRoutines && i*chunkSize < len(playerPicks); i++ {
		winnersChunk := <-resultChan
		for matches, count := range winnersChunk {
			winners[matches] += count
		}
	}
	return winners
}
