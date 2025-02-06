package repository

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type QuotesRepo struct {
	FilePath string
}

func NewQuotesRepo(filePath string) *QuotesRepo {
	return &QuotesRepo{
		FilePath: filePath,
	}
}

// Для больших файлов хранение, конечно же, лучше будет вынести в какую-то БД.
// Текущий алгоритм позволит более оптимально использовать память даже для больших
// файлов.
func (q *QuotesRepo) GetOneQuote() (string, error) {
	file, err := os.Open(q.FilePath)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rand.Seed(time.Now().UnixNano())

	var quote string
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		if rand.Intn(lineCount) == 0 {
			quote = scanner.Text()
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	if lineCount == 0 {
		return "", fmt.Errorf("file is empty")
	}

	return quote, nil
}
