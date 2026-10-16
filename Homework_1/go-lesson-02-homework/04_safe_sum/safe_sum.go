package safe_sum

import "math"

// SafeSum последовательно суммирует транзакции, пропуская те, добавление
// которых привело бы к переполнению int64. skipped — индексы пропущенных
// элементов исходного среза, в порядке появления.
func SafeSum(transactions []int64) (total int64, skipped []int) {

	if len(transactions) == 0 {
		return 0, nil
	}

	for i, n := range transactions {
		if (n >= 0 && total <= 0) || (n <= 0 && total >= 0) {
			total += n
		} else if (n >= 0 && total >= 0 && math.MaxInt64-total >= n) || (n <= 0 && total <= 0 && math.MinInt64-total <= n) {
			total += n
		} else {
			skipped = append(skipped, i)
		}
	}

	return total, skipped

}
