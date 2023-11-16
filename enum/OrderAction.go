package enum

type OrderAction string

const (
	CREATE_ORDER OrderAction = "create"
	CANCEL_ORDER OrderAction = "cancel"
)

func (o OrderAction) String() string {
	switch o {
	case CREATE_ORDER:
		return "create"
	case CANCEL_ORDER:
		return "cancel"
	default:
		return "unknown"
	}
}
