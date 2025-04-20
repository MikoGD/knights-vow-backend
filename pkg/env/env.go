package env

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// Loads env from the path provided
//
// Function will exit with code 1 for any failures
func LoadEnv(envPath string) error {
	file, err := os.Open(envPath)
	if err != nil {
		log.Fatalf("Unable to read %s\n", envPath)
	}

	scanner := bufio.NewScanner(file)
	scanner.Text()

	for scanner.Scan() {
		line := scanner.Text()

		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			log.Fatalf("\"%s\" is missing a key or value\n", line)
		}

		key := strings.Trim(tokens[0], " ")
		value := strings.Trim(tokens[1], " ")

		if err := os.Setenv(key, value); err != nil {
			log.Fatalf("Error setting env variable\nKey = %s\nValue = %s", key, value)
		}
	}

	return nil
}
