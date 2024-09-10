package main

import (
	"fmt"

	"github.com/khoaphungnguyen/go-leetcode/business"
	"github.com/khoaphungnguyen/go-leetcode/storage"
)

const (
	DB_FILE         = "questions.db"
	QUESTIONS_COUNT = 5
)

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
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Generate 5 random questions")
		fmt.Println("2. Show available questions")
		fmt.Println("3. Delete 5 random questions")
		fmt.Println("4. Complete a question")
		fmt.Println("5. Show completed questions")
		fmt.Println("6. Show details of a completed question")
		fmt.Println("7. Exit")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			newQuestions, err := bm.GenerateQuestions(QUESTIONS_COUNT)
			if err != nil {
				fmt.Println("Error generating questions:", err)
			} else {
				fmt.Println("Generated questions:")
				for _, q := range newQuestions {
					fmt.Printf("Number: %d, Type: %s, Difficulty: %s, Prompt: %s\n", q.Number, q.Type, q.Difficulty, q.Prompt)
				}
			}

		case 2:
			availableQuestions, err := qm.LoadAllQuestions()
			if err != nil {
				fmt.Println("Error loading available questions:", err)
			} else {
				fmt.Println("Available questions:")
				for _, q := range availableQuestions {
					fmt.Printf("Number: %d, Type: %s, Difficulty: %s, Prompt: %s\n", q.Number, q.Type, q.Difficulty, q.Prompt)
				}
			}

		case 3:
			if err := bm.DeleteRandomQuestions(QUESTIONS_COUNT); err != nil {
				fmt.Println("Error deleting random questions:", err)
			}

		case 4:
			fmt.Println("Enter the question number to complete:")
			var number int
			fmt.Scanln(&number)
			fmt.Println("Enter the answer to the question:")
			var answer string
			fmt.Scanln(&answer)

			if err := bm.CompleteQuestion(number, answer); err != nil {
				fmt.Println("Error completing question:", err)
			}

		case 5:
			completedQuestions, err := qm.LoadCompletedQuestions()
			if err != nil {
				fmt.Println("Error loading completed questions:", err)
			} else {
				fmt.Println("Completed questions:")
				for _, q := range completedQuestions {
					fmt.Printf("Number: %d, Type: %s, Difficulty: %s, Prompt: %s, Answer: %s\n", q.Number, q.Type, q.Difficulty, q.Prompt, q.Answer)
				}
			}

		case 6:
			fmt.Println("Enter the question number to show details:")
			var number int
			fmt.Scanln(&number)

			q, err := bm.ShowCompletedQuestionDetails(number)
			if err != nil {
				fmt.Println("Error showing completed question details:", err)
			} else {
				fmt.Printf("Number: %d, Type: %s, Difficulty: %s, Prompt: %s, Answer: %s\n", q.Number, q.Type, q.Difficulty, q.Prompt, q.Answer)
			}

		case 7:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice, please select again.")
		}
	}
}
