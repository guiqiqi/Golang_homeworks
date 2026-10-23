package main

import (
	"slices"
	"testing"
)

func TestSum(t *testing.T) {
	tests := []struct {
		name    string
		numbers []int
		want    int
	}{
		{name: "nil", numbers: nil, want: 0},
		{name: "пустой срез", numbers: []int{}, want: 0},
		{name: "один элемент", numbers: []int{7}, want: 7},
		{name: "положительные", numbers: []int{2, 4, 6, 8, 10}, want: 30},
		{name: "отрицательные", numbers: []int{-5, -3, -1}, want: -9},
		{name: "смешанные", numbers: []int{-3, 10, -2}, want: 5},
		{name: "нулевая сумма", numbers: []int{-2, -1, 0, 1, 2}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := slices.Clone(tt.numbers)
			if got := Sum(tt.numbers); got != tt.want {
				t.Errorf("Sum(%v) = %d, ожидается %d", before, got, tt.want)
			}
			if !slices.Equal(tt.numbers, before) {
				t.Errorf("Sum изменила срез: было %v, стало %v", before, tt.numbers)
			}
		})
	}
}

func TestSumAfterAppendAndSharedChange(t *testing.T) {
	numbers := []int{2, 4, 6, 8, 10}
	numbers = append(numbers, 12)
	if got := Sum(numbers); got != 42 {
		t.Errorf("сумма после append = %d, ожидается 42", got)
	}
	shared := numbers
	shared[0] = 100
	if got := Sum(numbers); got != 140 {
		t.Errorf("сумма после изменения через shared = %d, ожидается 140", got)
	}
}
