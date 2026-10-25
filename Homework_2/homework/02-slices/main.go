package main

import "fmt"

func Sum(numbers []int) int {
	var sum = 0
	if len(numbers) != 0 {
		for i := range numbers {
			sum += numbers[i]
		}
	}
	return sum
}

func main() {
	// 1. В Go при передаче среза (slice) фактически передается структура, содержащая указатель (адрес), длину и емкость.
	numbers := []int{2, 4, 6, 8, 10}
	fmt.Println(Sum(numbers))                // 30
	fmt.Println(Sum([]int{}))                // 0
	fmt.Println(Sum(nil))                    // 0
	fmt.Println(Sum([]int{-2, -1, 0, 1, 2})) // 0

	// 3. Операция `append` изменяет структуру среза (в частности, его длину) и может привести к расширению базового массива, если текущей емкости недостаточно; поэтому результирующий срез необходимо сохранить (присвоить переменной).
	numbers = append(numbers, 12)
	fmt.Println(len(numbers), Sum(numbers))

	// 2. Поскольку `shared` и `numbers` используют одну и ту же область памяти, любые изменения, внесенные в один из них, будут видны в другом.
	shared := numbers
	shared[0] = 100
	fmt.Println(numbers[0], shared[0]) // 100 100
}
