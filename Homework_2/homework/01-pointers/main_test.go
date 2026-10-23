package main

import "testing"

func TestIncrement(t *testing.T) {
	for _, initial := range []int{10, 0, -2} {
		number := initial
		Increment(&number)
		if number != initial+1 {
			t.Errorf("Increment(%d): число стало %d, ожидается %d", initial, number, initial+1)
		}
		Increment(&number)
		if number != initial+2 {
			t.Errorf("после двух вызовов для %d: получено %d, ожидается %d", initial, number, initial+2)
		}
	}
}

func TestIncrementNil(t *testing.T) {
	Increment(nil)
}
