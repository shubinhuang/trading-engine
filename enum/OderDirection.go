package enum

type OrderDirection string

const (
	BUY  OrderDirection = "buy"
	SELL OrderDirection = "sell"
)

func (o OrderDirection) String() string {
	switch o {
	case BUY:
		return "buy"
	case SELL:
		return "sell"
	default:
		return "unknown"
	}
}
