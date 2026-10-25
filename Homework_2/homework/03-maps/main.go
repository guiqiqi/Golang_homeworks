package main

import "fmt"

func StudentScores() map[string]int {
	scores := make(map[string]int)
	scores["Анна"] = 3
	scores["Борис"] = 0
	scores["Вера"] = 5
	scores["Анна"] += 2
	delete(scores, "Вера")

	return scores
}

// 1. В Go поиск по карте возвращает null в случае неудачи, что равно 0 для целого числа. Поэтому необходима отдельная переменная проверки, чтобы убедиться, что ключ не найден.
// 2. `var` только объявляет переменную, но не выделяет базовое пространство памяти для карты (оно равно nil). Поэтому запись в карту приведет к сбою. Для инициализации всегда требуется `make`.

func main() {
	scores := StudentScores()
	fmt.Println(scores["Анна"])
	for _, name := range []string{"Борис", "Глеб"} {
		score, ok := scores[name]
		fmt.Println(score, ok)
	}
	fmt.Println(len(scores))
}
