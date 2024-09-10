package business

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/khoaphungnguyen/go-leetcode/storage"
)

// BusinessManager handles business logic related to managing questions.
type BusinessManager struct {
	qm *storage.QuestionManager
}

// NewBusinessManager creates a new instance of BusinessManager.
func NewBusinessManager(qm *storage.QuestionManager) *BusinessManager {
	return &BusinessManager{qm: qm}
}

// GenerateQuestions generates random questions and adds them to the available list.
func (bm *BusinessManager) GenerateQuestions(count int) ([]storage.Question, error) {
	existingQuestions, err := bm.qm.LoadAllQuestions()
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	var newQuestions []storage.Question

	for len(newQuestions) < count {
		qNumber := rand.Intn(storage.NUM_OF_QUESTIONS) + 1
		unique := true
		for _, q := range existingQuestions {
			if q.Number == qNumber {
				unique = false
				break
			}
		}
		if unique {
			q := storage.Question{
				Number:     qNumber,
				Type:       "array", // Example: Set type, difficulty, prompt dynamically if available
				Difficulty: "medium",
				Prompt:     fmt.Sprintf("This is the prompt for question %d", qNumber),
			}
			newQuestions = append(newQuestions, q)
			if err := bm.qm.AddQuestion(q); err != nil {
				return nil, err
			}
		}
	}

	return newQuestions, nil
}

// CompleteQuestion marks a question as completed by moving it to the completed table with an answer.
func (bm *BusinessManager) CompleteQuestion(number int, answer string) error {
	q, err := bm.qm.LoadQuestion(number)
	if err != nil {
		return fmt.Errorf("question not found: %v", err)
	}

	q.Answer = answer

	// Delete from available and add to completed
	if err := bm.qm.DeleteQuestion(number); err != nil {
		return fmt.Errorf("failed to delete question: %v", err)
	}
	if err := bm.qm.AddCompletedQuestion(q); err != nil {
		return fmt.Errorf("failed to add completed question: %v", err)
	}

	return nil
}

// DeleteRandomQuestions deletes a specified number of random questions from the available list.
func (bm *BusinessManager) DeleteRandomQuestions(count int) error {
	availableQuestions, err := bm.qm.LoadAllQuestions()
	if err != nil {
		return err
	}

	if len(availableQuestions) < count {
		count = len(availableQuestions)
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		idx := rand.Intn(len(availableQuestions))
		q := availableQuestions[idx]
		if err := bm.qm.DeleteQuestion(q.Number); err != nil {
			return fmt.Errorf("failed to delete question: %v", err)
		}
		availableQuestions = append(availableQuestions[:idx], availableQuestions[idx+1:]...)
	}

	return nil
}

// ShowCompletedQuestionDetails shows detailed information for a completed question.
func (bm *BusinessManager) ShowCompletedQuestionDetails(number int) (storage.Question, error) {
	q, err := bm.qm.ShowCompletedQuestionDetails(number)
	if err != nil {
		return q, fmt.Errorf("failed to load completed question details: %v", err)
	}
	return q, nil
}
