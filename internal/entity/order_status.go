package entity

type OrderStatus int

const (
	New OrderStatus = iota
	Processing
	Invalid
	Processed
)

func (os OrderStatus) String() string {
	return [...]string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}[os]
}
