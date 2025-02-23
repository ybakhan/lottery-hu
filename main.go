package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	DEFAULT_NUMBER_OF_PICKS        = 5            // Number of picks per player
	DEFAULT_MIN_MATCHES            = 2            // Minimum matches required to win
	DEFAULT_MIN_PICK               = 1            // Minimum valid lottery number
	DEFAULT_MAX_PICK               = 90           // Maximum valid lottery number
	DEFAULT_PLAYER_PICKS_FILE_PATH = "10m-v2.txt" // Path to player picks file
)

// Lottery defines the interface for lottery functionality.
type Lottery interface {
	ProcessPlayerPicks(file *os.File) ([]LotteryPick, error)
	MatchPicks(winningEntry string, playerPicks []LotteryPick) error
}

// main is the entry point for the Hungarian Lottery system.
// Loads environment variables, processes player picks from a file, and matches lottery picks from stdin.
func main() {
	fmt.Println("Hungarian Lottery system. Press CTRL+C to exit")

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	file, err := initializePicksFile()
	if err != nil {
		fmt.Println("Error reading player picks file:", err)
		os.Exit(1)
	}
	defer file.Close()

	l := initializeLottery()

	// Read player picks from file
	playerPicks, err := l.ProcessPlayerPicks(file)
	if err != nil {
		fmt.Println("Error processing player picks:", err)
		os.Exit(1)
	}

	fmt.Println("\nEnter winning lottery pick")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		// Process each lottery pick entry
		_ = l.MatchPicks(scanner.Text(), playerPicks)
		fmt.Println("\nEnter winning lottery pick")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading lottery picks:", err)
		os.Exit(1)
	}
}

// initializeLottery constructs a Lottery instance with configuration from environment variables
// or uses default settings if environment variables unset or invalid.
func initializeLottery() Lottery {
	numberOfPicks, err := strconv.Atoi(os.Getenv("NUMBER_OF_PICKS"))
	if err != nil {
		numberOfPicks = DEFAULT_NUMBER_OF_PICKS
	}

	minMatches, err := strconv.Atoi(os.Getenv("MIN_MATCHES"))
	if err != nil || minMatches <= 0 {
		minMatches = DEFAULT_MIN_MATCHES
	}

	minPick, err := strconv.Atoi(os.Getenv("MIN_PICK"))
	if err != nil || minPick <= 0 {
		minPick = DEFAULT_MIN_PICK
	}

	maxPick, err := strconv.Atoi(os.Getenv("MAX_PICK"))
	if err != nil || maxPick > 128 {
		maxPick = DEFAULT_MAX_PICK
	}

	return lottery{
		numberOfPicks,
		minMatches,
		minPick,
		maxPick,
	}
}

// initializePicksFile opens the file containing player lottery picks.
// Returns the opened file or an error if the file cannot be opened.
func initializePicksFile() (*os.File, error) {
	picksFilePath := os.Getenv("PLAYER_PICKS_FILE_PATH")
	if picksFilePath == "" {
		picksFilePath = DEFAULT_PLAYER_PICKS_FILE_PATH
	}
	return os.Open(picksFilePath)
}
