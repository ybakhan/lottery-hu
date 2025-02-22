package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_NUMBER_OF_PICKS          = 5
	DEFAULT_MIN_MATCHES              = 2
	DEFAULT_MIN_LOTTERY_PICK         = 1
	DEFAULT_MAX_LOTTERY_PICK         = 90
	DEFAULT_PLAYER_NUMBERS_FILE_PATH = "10m-v2.txt"
)

type Lottery interface {
	ProcessPlayerNumbers(file *os.File) ([][]int, error)
	ParseLotteryPicks(lotteryPicksEntry string, allPlayersNumbers [][]int) error
}

func initializeLottery() Lottery {
	numberOfPicks, err := strconv.Atoi(os.Getenv("NUMBER_OF_PICKS"))
	if err != nil {
		numberOfPicks = DEFAULT_NUMBER_OF_PICKS
	}

	minMatches, err := strconv.Atoi(os.Getenv("MIN_MATCHES"))
	if err != nil {
		minMatches = DEFAULT_MIN_MATCHES
	}

	minLotteryPick, err := strconv.Atoi(os.Getenv("MIN_LOTTERY_PICK"))
	if err != nil {
		minLotteryPick = DEFAULT_MIN_LOTTERY_PICK
	}

	maxLotteryPick, err := strconv.Atoi(os.Getenv("MAX_LOTTERY_PICK"))
	if err != nil {
		maxLotteryPick = DEFAULT_MAX_LOTTERY_PICK
	}

	return lottery{
		numberOfPicks,
		minMatches,
		minLotteryPick,
		maxLotteryPick,
	}
}

func main() {
	fmt.Println("Hungarian Lottery system. Press CTRL+C to exit")

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	playerNumbersFile, err := initializePlayerNumbersFile()

	if err != nil {
		fmt.Println("Error reading player numbers file:", err)
		os.Exit(1)
	}

	defer playerNumbersFile.Close()

	l := initializeLottery()

	allPlayersNumbers, err := l.ProcessPlayerNumbers(playerNumbersFile)

	if err != nil {
		fmt.Println("Error processing player numbers:", err)
		os.Exit(1)
	}

	fmt.Println("\nEnter lottery picks")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		if err := l.ParseLotteryPicks(scanner.Text(), allPlayersNumbers); err != nil {
			continue
		}

		fmt.Println("\nEnter lottery picks")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading lottery picks:", err)
		os.Exit(1)
	}
}

func initializePlayerNumbersFile() (*os.File, error) {

	playerNumbersFilePath := os.Getenv("PLAYER_NUMBERS_FILE_PATH")

	if playerNumbersFilePath == "" {
		playerNumbersFilePath = DEFAULT_PLAYER_NUMBERS_FILE_PATH
	}

	return os.Open(playerNumbersFilePath)
}
