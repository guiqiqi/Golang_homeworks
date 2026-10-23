package main

import "fmt"

type Status string

func (s Status) String() string {
	m := map[Status]string{
		"new":      "Новый",
		"paid":     "Оплачен",
		"canceled": "Отменен",
	}

	if v, ok := m[s]; ok {
		return v
	}

	return "Неизвестный статус"
}

type Item struct {
	Name  string
	Price int64
}

type Order struct {
	ID     int
	Items  []Item
	Status Status
}

func (o *Order) AddItem(item Item) {
	o.Items = append(o.Items, item)
}

func (o Order) Total() int64 {
	var total int64
	for _, item := range o.Items {
		total += item.Price
	}
	return total
}

func NewOrder(id int) Order {
	return Order{
		ID:     id,
		Status: "new",
	}
}

func main() {
	o := NewOrder(1)
	o.AddItem(Item{Name: "Кофе", Price: 199_99})
	fmt.Println(o, o.Total())

	// s := make([]int, 0, 3)
	// s = append(s, []int{1, 2, 3}...)
	// a := s[0:3]
	// a = append(a, 4)
	// s[1] = 100
	// fmt.Println(a, len(a), cap(a))
}
