package enum

type OrderStatus string

const (
	FULLY_FILLED   OrderStatus = "fully_filled"
	PARTIAL_FILLED OrderStatus = "partial_filled"
	PENDING        OrderStatus = "pending"
)

func (o OrderStatus) String() string {
	switch o {
	case FULLY_FILLED:
		return "fully_filled"
	case PARTIAL_FILLED:
		return "partial_filled"
	case PENDING:
		return "pending"
	default:
		return "unknown"
	}
}
