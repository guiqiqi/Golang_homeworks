package order_status

// CountStatuses возвращает количество new, paid, shipped, cancelled,
// unknown и active (new + paid + shipped) в указанном порядке.
// Любой статус, не равный одному из четырёх известных, считается unknown.
func CountStatuses(statuses []string) (newCount, paidCount, shippedCount, cancelledCount, unknownCount, activeCount int) {

	for _, status := range statuses {
		switch status {
		case "new":
			newCount += 1
			activeCount += 1
		case "paid":
			paidCount += 1
			activeCount += 1
		case "shipped":
			shippedCount += 1
			activeCount += 1
		case "cancelled":
			cancelledCount += 1
		default:
			unknownCount += 1
		}
	}

	return newCount, paidCount, shippedCount, cancelledCount, unknownCount, activeCount
}
