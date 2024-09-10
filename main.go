package main

import (
	"bufio"
	"fmt"

	"os"
	"strings"

	"github.com/khoaphungnguyen/go-leetcode/business"
	"github.com/khoaphungnguyen/go-leetcode/storage"
)

const (
	DB_FILE         = "questions.db"
	QUESTIONS_COUNT = 5
)

func printHeader(title string) {
	fmt.Println("\n==============================")
	fmt.Printf("  %s\n", title)
	fmt.Println("==============================")
}

func printDivider() {
	fmt.Println("---------------------------------")
}

func printMenu() {
	printHeader("LeetCode Question Manager")
	fmt.Println("1. Generate 5 random questions")
	fmt.Println("2. Show available questions")
	fmt.Println("3. Delete 5 random questions")
	fmt.Println("4. Complete a question")
	fmt.Println("5. Show completed questions")
	fmt.Println("6. Show details of a completed question")
	fmt.Println("7. Add details to a question")
	fmt.Println("8. Exit")
	printDivider()
	fmt.Print("Enter your choice: ")
}

func printQuestion(q storage.Question) {
	fmt.Printf("Number: %d\n", q.Number)
	fmt.Printf("Type: %s\n", q.Type)
	fmt.Printf("Difficulty: %s\n", q.Difficulty)
	fmt.Printf("Prompt: %s\n", q.Prompt)
	if q.Answer != "" {
		fmt.Println("Answer:")
		printDivider()
		// Split the answer into lines for better readability
		lines := strings.Split(q.Answer, "\n")
		for _, line := range lines {
			// Print each line with indentation for better visibility
			fmt.Printf("    %s\n", line)
		}
		printDivider()
	}
}

func addQuestionDetails(qm *storage.QuestionManager, number int) error {
	// Check if the question exists in the available questions list
	availableQuestions, err := qm.LoadAllQuestions()
	if err != nil {
		return fmt.Errorf("error loading available questions: %v", err)
	}

	// Check if the question is available
	var found bool
	for _, q := range availableQuestions {
		if q.Number == number {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("question number %d is not available to add details", number)
	}

	// Proceed to add details to the question
	reader := bufio.NewReader(os.Stdin)

	// Prompt for question type
	fmt.Print("Enter the type of the question (e.g., array, graph): ")
	qType, _ := reader.ReadString('\n')
	qType = strings.TrimSpace(qType)

	// Prompt for difficulty level
	fmt.Print("Enter the difficulty level (easy, medium, hard): ")
	difficulty, _ := reader.ReadString('\n')
	difficulty = strings.TrimSpace(difficulty)

	// Prompt for the question prompt
	fmt.Println("Enter the prompt for the question (end with a blank line):")
	var promptLines []string
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" { // End input on a blank line
			break
		}
		promptLines = append(promptLines, line)
	}
	prompt := strings.Join(promptLines, " ")

	// Prompt for the answer (multiline input)
	fmt.Println("Enter the answer to the question (type 'END' on a new line to finish):")
	var answerLines []string
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "END" {
			break
		}
		answerLines = append(answerLines, line)
	}
	answer := strings.Join(answerLines, "\n")

	// Add the details to the question
	question := storage.Question{
		Number:     number,
		Type:       qType,
		Difficulty: difficulty,
		Prompt:     prompt,
		Answer:     answer,
	}

	// Update the question in the database
	if err := qm.AddQuestion(question); err != nil {
		return fmt.Errorf("error adding question details: %v", err)
	}

	fmt.Println("Details added successfully.")
	return nil
}

func main() {
	// Initialize QuestionManager
	qm, err := storage.NewQuestionManager(DB_FILE)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer qm.Close()

	// Initialize BusinessManager
	bm := business.NewBusinessManager(qm)

	for {
		printMenu()

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			printHeader("Generate Random Questions")
			newQuestions, err := bm.GenerateQuestions(QUESTIONS_COUNT)
			if err != nil {
				fmt.Println("Error generating questions:", err)
			} else {
				fmt.Println("Generated Questions:")
				printDivider()
				for _, q := range newQuestions {
					printQuestion(q)
				}
			}

		case 2:
			printHeader("Available Questions")
			availableQuestions, err := qm.LoadAllQuestions()
			if err != nil {
				fmt.Println("Error loading available questions:", err)
			} else if len(availableQuestions) == 0 {
				fmt.Println("No available questions found.")
			} else {
				for _, q := range availableQuestions {
					printQuestion(q)
				}
			}

		case 3:
			printHeader("Delete Random Questions")
			if err := bm.DeleteRandomQuestions(QUESTIONS_COUNT); err != nil {
				fmt.Println("Error deleting random questions:", err)
			} else {
				fmt.Println("Successfully deleted 5 random questions.")
			}

		case 4:
			printHeader("Complete a Question")
			fmt.Print("Enter the question number to complete: ")
			var number int
			fmt.Scanln(&number)
			fmt.Println("Enter the answer to the question (type 'END' on a new line to finish):")
			reader := bufio.NewReader(os.Stdin)
			var answerLines []string
			for {
				line, _ := reader.ReadString('\n')
				line = strings.TrimSpace(line)
				if line == "END" {
					break
				}
				answerLines = append(answerLines, line)
			}
			answer := strings.Join(answerLines, "\n")

			if err := bm.CompleteQuestion(number, answer); err != nil {
				fmt.Println("Error completing question:", err)
			} else {
				fmt.Println("Question completed successfully.")
			}

		case 5:
			printHeader("Completed Questions")
			completedQuestions, err := qm.LoadCompletedQuestions()
			if err != nil {
				fmt.Println("Error loading completed questions:", err)
			} else if len(completedQuestions) == 0 {
				fmt.Println("No completed questions found.")
			} else {
				for _, q := range completedQuestions {
					printQuestion(q)
				}
			}

		case 6:
			printHeader("Show Details of a Completed Question")
			fmt.Print("Enter the question number to show details: ")
			var number int
			fmt.Scanln(&number)

			q, err := bm.ShowCompletedQuestionDetails(number)
			if err != nil {
				fmt.Println("Error showing completed question details:", err)
			} else {
				printQuestion(q)
			}

		case 7:
			printHeader("Add Details to a Question")
			fmt.Print("Enter the question number to add details: ")
			var number int
			fmt.Scanln(&number)

			if err := addQuestionDetails(qm, number); err != nil {
				fmt.Println(err)
			}

		case 8:
			printHeader("Exiting")
			fmt.Println("Thank you for using LeetCode Question Manager!")
			return

		default:
			fmt.Println("Invalid choice, please select again.")
		}
	}
}
