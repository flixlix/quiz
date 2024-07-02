package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFileName := flag.String("csv", "./problems.csv", "A csv file in the format of 'question, answer'")
	timeLimit := flag.Int("limit", 30, "The maximum amount of time a user is allowed to complete the quiz (in seconds).")
	flag.Parse()

	file, err := os.Open(*csvFileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open %s", *csvFileName))
		os.Exit(1)
	}

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	defer file.Close()

	if err != nil {
		exit("Failed to parse the provided csv file.")
	}

	problems := parseLines(records)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correctAnswers := 0

problemLoop:

	for index, problem := range problems {
		fmt.Printf("Problem #%d: %s = ", index+1, problem.q)

		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\nYour time has run out!")
			break problemLoop
		case answer := <-answerCh:
			if answer != problem.a {
				continue
			}
			correctAnswers++
		}
	}

	fmt.Printf("You scored %d out of %d answers correctly.\n", correctAnswers, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
