package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type lottery struct {
	numberOfPicks  int
	minMatches     int
	minLotteryPick int
	maxLotteryPick int
}

func (l lottery) ProcessPlayerNumbers(file *os.File) ([][]int, error) {
	playerIndex := 1
	var allPlayersNumbers [][]int
	scanner := bufio.NewScanner(file)

	// Read the player numbers line by line and process each player's numbers
	for scanner.Scan() {

		line := scanner.Text()
		numbersStr := strings.Fields(line)

		if len(numbersStr) != l.numberOfPicks {

			fmt.Printf("Player %d ignored expected %d numbers got %d\n", playerIndex, l.numberOfPicks, len(numbersStr))
			continue
		}

		numbers := make([]int, l.numberOfPicks)
		validLine := true

		// Convert strings to integers
		for i, numberStr := range numbersStr {

			num, err := strconv.Atoi(numberStr)

			if err != nil {

				fmt.Printf("Player %d ignored error converting %s to integer\n", playerIndex, numberStr)

				validLine = false
				break
			}

			numbers[i] = num
		}

		if validLine {

			sort.Ints(numbers[:])

			allPlayersNumbers = append(allPlayersNumbers, numbers[:])
		}

		//fmt.Printf("Player scores %d sorted: %v\n", playerIndex, numbers)
		playerIndex++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return allPlayersNumbers, nil
}

func (l lottery) parsePlayerNumbers(playerNumbersEntry string) ([]int, error) {

	lotteryPicksStr := strings.Fields(playerNumbersEntry)

	if len(lotteryPicksStr) != l.numberOfPicks {

		fmt.Printf("Enter %d lottery picks\n", l.numberOfPicks)

		return nil, fmt.Errorf("invalid number of lottery picks")
	}

	lotteryPicks := make([]int, l.numberOfPicks)

	for i, lotteryPickStr := range lotteryPicksStr {

		lotteryPick, err := strconv.Atoi(lotteryPickStr)

		if err != nil {

			fmt.Printf("Error parsing lottery pick %s : %v", lotteryPickStr, err)

			return nil, err
		}

		if lotteryPick < l.minLotteryPick || lotteryPick > l.maxLotteryPick {

			fmt.Printf("Lottery pick entry out of range %d\n", lotteryPick)

			return nil, fmt.Errorf("lottery pick entry out of range %d", lotteryPick)
		}

		lotteryPicks[i] = lotteryPick
	}

	fmt.Printf("lottery pick: %v \n", lotteryPicks)

	sort.Ints(lotteryPicks)

	return lotteryPicks, nil
}

func (l lottery) ParseLotteryPicks(lotteryPicksEntry string, allPlayersNumbers [][]int) error {
	lotteryPicks, err := l.parsePlayerNumbers(lotteryPicksEntry)

	if err != nil {
		return err
	}

	start := time.Now()

	winners := l.matchLotteryPicks(lotteryPicks, allPlayersNumbers)

	elapsed := time.Since(start)

	fmt.Printf("matchLotteryPicks took %v\n", elapsed)

	fmt.Printf("%-16s\t%s\n", "Numbers matching", "Winners")

	for i := l.numberOfPicks; i >= l.minMatches; i-- {

		fmt.Printf("%-16d\t%d\n", i, winners[i])

	}

	return nil
}

func (l lottery) matchLotteryPicks(lotteryPicks []int, allPlayersNumbers [][]int) map[int]int {
	winners := make(map[int]int)

	if len(lotteryPicks) != l.numberOfPicks {
		return winners
	}

	for _, playerNumbers := range allPlayersNumbers {

		if len(playerNumbers) != l.numberOfPicks {
			continue
		}

		matches := l.countMatches(lotteryPicks, playerNumbers)

		if matches >= l.minMatches {
			winners[matches]++
		}
	}

	return winners
}

func (l lottery) countMatches(lotteryPicks, playerNumbers []int) int {
	matches := 0

	for i, j := 0, 0; i < l.numberOfPicks && j < l.numberOfPicks; {

		if playerNumbers[i] == lotteryPicks[j] {

			matches++
			i++
			j++

		} else if playerNumbers[i] < lotteryPicks[j] {

			i++

		} else {

			j++
		}
	}

	return matches
}
