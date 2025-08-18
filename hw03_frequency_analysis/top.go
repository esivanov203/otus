package hw03frequencyanalysis

import (
	"errors"
	"log"
	"sort"
	"strings"
)

const count = 10

var ErrEmptyString = errors.New("empty string")

var ErrNoWords = errors.New("no words in empty string")

type WordList struct {
	Word     string
	Quantity int
}

func split(text string) ([]string, error) {
	if strings.TrimSpace(text) == "" {
		// Эта операция по сути не нужна, применил ее,
		// чтобы показать обработку разных ошибок (для учебных целей)
		return nil, ErrEmptyString
	}
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil, ErrNoWords
	}

	return words, nil
}

func Top10(text string) []string {
	words, err := split(text)
	if errors.Is(err, ErrNoWords) || errors.Is(err, ErrEmptyString) {
		return []string{}
	} else if err != nil {
		log.Fatalf("Unexpected fail to split words: %v", err)
	}

	uniqueWords := map[string]int{}
	for _, word := range words {
		uniqueWords[word]++
	}

	sortedWords := make([]WordList, len(uniqueWords))
	i := 0
	for k, v := range uniqueWords {
		sortedWords[i] = WordList{k, v}
		i++
	}
	sort.Slice(sortedWords, func(i, j int) bool {
		if sortedWords[i].Quantity == sortedWords[j].Quantity {
			return sortedWords[i].Word < sortedWords[j].Word
		}
		return sortedWords[i].Quantity > sortedWords[j].Quantity
	})

	result := make([]string, 0, count)
	for i, w := range sortedWords {
		if i >= count {
			break
		}
		result = append(result, w.Word)
		// Возможен вывод с количеством: fmt.Sprintf("%s (%d)", w.Word, w.Quantity), но ТЗ требует только слова
	}

	return result
}
