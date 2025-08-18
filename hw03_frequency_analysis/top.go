package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const count = 10

type WordList struct {
	Word     string
	Quantity int
}

func Top10(text string) []string {
	words := strings.Fields(text)

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
