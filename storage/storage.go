package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const NUM_OF_QUESTIONS = 3283

// Question represents a question with all necessary details.
type Question struct {
	Number     int
	Type       string
	Difficulty string
	Prompt     string
	Answer     string
}

// QuestionManager manages available and completed questions using a SQLite database.
type QuestionManager struct {
	db *sql.DB
}

// NewQuestionManager initializes a new QuestionManager and sets up the database.
func NewQuestionManager(dbFile string) (*QuestionManager, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	// Create tables if they don't exist
	createTablesSQL := `
    CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY,
    number INTEGER UNIQUE,
    type TEXT,
    difficulty TEXT,
    prompt TEXT
	);

	CREATE TABLE IF NOT EXISTS completed (
		id INTEGER PRIMARY KEY,
		number INTEGER UNIQUE,
		type TEXT,
		difficulty TEXT,
		prompt TEXT,
		answer TEXT
	);
	`
	_, err = db.Exec(createTablesSQL)
	if err != nil {
		return nil, err
	}

	return &QuestionManager{db: db}, nil
}

// Close closes the database connection.
func (qm *QuestionManager) Close() error {
	return qm.db.Close()
}

// AddQuestion adds a new question or updates an existing one in the questions table.
func (qm *QuestionManager) AddQuestion(q Question) error {
	// Check if the question exists
	var exists bool
	err := qm.db.QueryRow("SELECT EXISTS(SELECT 1 FROM questions WHERE number = ?)", q.Number).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking existence of question: %v", err)
	}

	if exists {
		// Update existing question
		_, err := qm.db.Exec(
			"UPDATE questions SET type = ?, difficulty = ?, prompt = ? WHERE number = ?",
			q.Type, q.Difficulty, q.Prompt, q.Number,
		)
		if err != nil {
			return fmt.Errorf("error updating question: %v", err)
		}
	} else {
		// Insert new question
		_, err := qm.db.Exec(
			"INSERT INTO questions (number, type, difficulty, prompt) VALUES (?, ?, ?, ?)",
			q.Number, q.Type, q.Difficulty, q.Prompt,
		)
		if err != nil {
			return fmt.Errorf("error adding new question: %v", err)
		}
	}

	return nil
}

// DeleteQuestion deletes a question from the questions table by its number.
func (qm *QuestionManager) DeleteQuestion(number int) error {
	_, err := qm.db.Exec("DELETE FROM questions WHERE number = ?", number)
	return err
}

// LoadQuestion retrieves a question by its number from the questions table.
func (qm *QuestionManager) LoadQuestion(number int) (Question, error) {
	var q Question
	err := qm.db.QueryRow("SELECT number, type, difficulty, prompt FROM questions WHERE number = ?", number).Scan(&q.Number, &q.Type, &q.Difficulty, &q.Prompt)
	if err != nil {
		return q, err
	}
	return q, nil
}

// LoadAllQuestions loads all available questions from the database.
func (qm *QuestionManager) LoadAllQuestions() ([]Question, error) {
	rows, err := qm.db.Query("SELECT number, type, difficulty, prompt FROM questions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.Number, &q.Type, &q.Difficulty, &q.Prompt); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// LoadCompletedQuestions loads completed questions from the database.
func (qm *QuestionManager) LoadCompletedQuestions() ([]Question, error) {
	rows, err := qm.db.Query("SELECT number, type, difficulty, prompt, answer FROM completed")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.Number, &q.Type, &q.Difficulty, &q.Prompt, &q.Answer); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, nil
}

// AddCompletedQuestion adds a question to the completed questions table.
func (qm *QuestionManager) AddCompletedQuestion(q Question) error {
	_, err := qm.db.Exec("INSERT INTO completed (number, type, difficulty, prompt, answer) VALUES (?, ?, ?, ?, ?)", q.Number, q.Type, q.Difficulty, q.Prompt, q.Answer)
	return err
}

// ShowCompletedQuestionDetails shows details of a completed question by its number.
func (qm *QuestionManager) ShowCompletedQuestionDetails(number int) (Question, error) {
	var q Question
	err := qm.db.QueryRow("SELECT number, type, difficulty, prompt, answer FROM completed WHERE number = ?", number).Scan(&q.Number, &q.Type, &q.Difficulty, &q.Prompt, &q.Answer)
	if err != nil {
		return q, fmt.Errorf("error loading completed question details: %v", err)
	}
	return q, nil
}
