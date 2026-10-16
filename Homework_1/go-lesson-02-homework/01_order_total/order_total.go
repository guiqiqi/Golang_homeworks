package order_total

// CalculateOrder возвращает стоимость заказа до скидки, размер скидки и итог
// (всё в копейках). Скидка задаётся целым числом процентов от 0 до 100.
// При некорректном входе функция возвращает три нуля.
// Расчёт скидки выполняется целочисленно с округлением вниз до копейки.
// Можно считать, что суммы и произведения не переполняют int64.
func CalculateOrder(prices []int64, quantities []int, discount int) (subtotal, discountAmount, total int64) {
	if (len(prices) != len(quantities)) ||
		(len(prices) == len(quantities) && len(prices) == 0) ||
		(discount > 100 || discount < 0) {
		return 0, 0, 0
	}

	subtotal = 0

	for i := 0; i < len(prices); i++ {
		price := prices[i]
		quantity := quantities[i]

		if price < 0 || quantity < 0 {
			return 0, 0, 0
		}

		subtotal += price * int64(quantity)
	}

	discountAmount = (subtotal * int64(discount)) / 100

	return subtotal, discountAmount, subtotal - discountAmount
}
