package main

import "fmt"

func Increment(value *int) {
	if value == nil {
		return
	}
	(*value)++
}

func main() {

	// 1. `&` and `*` are the address-of and dereference operators, respectively—equivalent to the same operations in the C language.
	number := 10
	Increment(&number)
	fmt.Println(number)

	// 2. The reason for passing a pointer is that Go uses "pass-by-value" for function arguments; therefore, to modify an external variable, you must pass its pointer.
	Increment(&number)
	fmt.Println(number)

	Increment(nil)
}
