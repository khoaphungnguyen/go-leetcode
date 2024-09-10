package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/rand"
)

const NUM_OF_QUESTIONS = 3283

// GenerateRandQuestions generates 5 random questions from the pool of LeetCode questions
func GenerateRandQuestions() []int {
	// Read existing questions from the file, if any
	results, err := ReadFromFile("random_questions.txt") // returns map[int]struct{}
	if err != nil {
		fmt.Println("No existing questions found, starting fresh.")
		results = map[int]struct{}{}
	}

	rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator with the current time

	// Convert map to a slice for easy handling
	resultsSlice := convertMapToSlice(results)

	// Target number of questions
	targetCount := len(resultsSlice) + 5

	// Add new unique questions until we reach the target count
	for len(resultsSlice) < targetCount {
		question := rand.Intn(NUM_OF_QUESTIONS) + 1
		if _, exists := results[question]; !exists {
			results[question] = struct{}{}
			resultsSlice = append(resultsSlice, question)
		}
	}

	return resultsSlice
}

// DeleteRandQuestions deletes 5 random questions from the map of questions
func DeleteRandQuestions(questions map[int]struct{}) map[int]struct{} {
	if len(questions) == 0 {
		fmt.Println("The list is empty, nothing to delete.")
		return questions
	}

	rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator with the current time

	// Number of questions to delete is the minimum of 5 or the length of the questions list
	numToDelete := 5
	if len(questions) < 5 {
		numToDelete = len(questions)
	}

	// Convert keys to a slice to select random items
	keys := make([]int, 0, len(questions))
	for k := range questions {
		keys = append(keys, k)
	}

	// Randomly delete questions
	for i := 0; i < numToDelete; i++ {
		idx := rand.Intn(len(keys))
		delete(questions, keys[idx])
		keys = append(keys[:idx], keys[idx+1:]...) // Remove the deleted element from the keys slice
	}

	return questions
}

// CompleteQuestion removes a specific question by its value from the map
func CompleteQuestion(number int, questions map[int]struct{}) map[int]struct{} {
	if len(questions) == 0 {
		fmt.Println("The list is empty, nothing to delete.")
		return questions
	}

	// Check if the question exists
	if _, exists := questions[number]; !exists {
		fmt.Printf("The question %d does not exist in the list.\n", number)
		return questions
	}

	// Remove the question
	delete(questions, number)
	return questions
}

// SaveToFile saves the generated questions to a text file
func SaveToFile(questions map[int]struct{}, filename string) error {
	file, err := os.Create(filename) // Create a file
	if err != nil {
		return err
	}
	defer file.Close()

	for question := range questions {
		_, err := file.WriteString(fmt.Sprintf("%d\n", question)) // Write each question on a new line
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadFromFile reads the questions from the text file and returns a map of questions
func ReadFromFile(filename string) (map[int]struct{}, error) {
	file, err := os.Open(filename) // Open the file
	if err != nil {
		return nil, err
	}
	defer file.Close()

	questions := make(map[int]struct{})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var question int
		_, err := fmt.Sscanf(scanner.Text(), "%d", &question)
		if err != nil {
			return nil, err
		}
		questions[question] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

func main() {
	for {
		fmt.Println("Choose an option:")
		fmt.Println("1. Generate 5 random questions")
		fmt.Println("2. Read saved questions from file")
		fmt.Println("3. Delete random 5 questions from the list")
		fmt.Println("4. Complete a question from the list")
		fmt.Println("5. Exit")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			// Generate random questions and convert the result to a map
			questions := GenerateRandQuestions()
			questionsMap := convertSliceToMap(questions)
			fmt.Println("Randomly selected questions:", questions)

			// Save the result to a file
			err := SaveToFile(questionsMap, "random_questions.txt")
			if err != nil {
				fmt.Println("Error saving to file:", err)
			} else {
				fmt.Println("Questions saved to random_questions.txt")
			}

		case 2:
			// Read the questions from the file
			questions, err := ReadFromFile("random_questions.txt")
			if err != nil {
				fmt.Println("Error reading from file:", err)
			} else {
				fmt.Println("Questions read from file:", convertMapToSlice(questions))
			}

		case 3:
			// Read the questions from the file
			questions, err := ReadFromFile("random_questions.txt")
			if err != nil {
				fmt.Println("Error reading from file:", err)
				continue
			}

			// Delete 5 random questions
			questions = DeleteRandQuestions(questions)
			fmt.Println("Updated questions after deletion:", convertMapToSlice(questions))

			// Save the updated list back to the file
			err = SaveToFile(questions, "random_questions.txt")
			if err != nil {
				fmt.Println("Error saving to file:", err)
			} else {
				fmt.Println("Updated questions saved to random_questions.txt")
			}

		case 4:
			// Read the questions from the file
			questions, err := ReadFromFile("random_questions.txt")
			if err != nil {
				fmt.Println("Error reading from file:", err)
				continue
			}

			// Prompt for the question number to complete
			fmt.Println("Enter the question number to complete:")
			var number int
			fmt.Scanln(&number)

			// Complete (remove) the question
			questions = CompleteQuestion(number, questions)
			fmt.Println("Updated questions after completing:", convertMapToSlice(questions))

			// Save the updated list back to the file
			err = SaveToFile(questions, "random_questions.txt")
			if err != nil {
				fmt.Println("Error saving to file:", err)
			} else {
				fmt.Println("Updated questions saved to random_questions.txt")
			}

		case 5:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice, please select again.")
		}
	}
}

// Helper function to convert a slice to a map
func convertSliceToMap(slice []int) map[int]struct{} {
	m := make(map[int]struct{}, len(slice))
	for _, v := range slice {
		m[v] = struct{}{}
	}
	return m
}

// Helper function to convert a map to a slice
func convertMapToSlice(m map[int]struct{}) []int {
	slice := make([]int, 0, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	return slice
}
