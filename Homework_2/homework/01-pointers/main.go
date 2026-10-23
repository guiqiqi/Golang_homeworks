package main

import "fmt"

func Increment(value *int) {
	if value == nil {
		return
	}
	(*value)++
}

func main() {
    number := 10
    Increment(&number)
    fmt.Println(number)

    Increment(&number)
    fmt.Println(number)

    Increment(nil)
}
