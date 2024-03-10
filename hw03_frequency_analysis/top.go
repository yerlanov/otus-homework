package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	words := strings.Fields(text)

	wordCounts := make(map[string]int)

	for _, word := range words {
		wordCounts[word]++
	}

	// Создаем структуру для сортировки
	type wordCount struct {
		Word  string
		Count int
	}

	wordCountsSlice := make([]wordCount, 0, len(wordCounts))

	for word, count := range wordCounts {
		wordCountsSlice = append(wordCountsSlice, wordCount{word, count})
	}

	// Сортируем слова
	sort.Slice(wordCountsSlice, func(i, j int) bool {
		if wordCountsSlice[i].Count == wordCountsSlice[j].Count {
			return wordCountsSlice[i].Word < wordCountsSlice[j].Word // Лексикографическая сортировка
		}
		return wordCountsSlice[i].Count > wordCountsSlice[j].Count // Сортировка по убыванию частоты
	})

	// Выбираем топ-10 слов
	topWords := make([]string, 0, 10)
	for i := 0; i < 10 && i < len(wordCountsSlice); i++ {
		topWords = append(topWords, wordCountsSlice[i].Word)
	}

	return topWords
}
