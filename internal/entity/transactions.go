package entity

import (
	"math"
	"time"
)

const defaultPrecision = 100

type PointValue int

type OperationType int

const (
	OperationAdd OperationType = iota
	OperationSub
)

type Withdrawal struct {
	OrderNumber string
	Sum         PointValue
	ProcessedAt time.Time
}

type UserBalance struct {
	Current   PointValue
	Withdrawn PointValue
}

func NewPointValue(value float64) PointValue {
	sign := math.Round(value * defaultPrecision)
	return PointValue(sign)
}

func (pv PointValue) ToFloat() float64 {
	if pv == 0 {
		return 0
	}
	r := float64(pv) / defaultPrecision
	return r
}

func (ot OperationType) String() string {
	switch ot {
	case OperationAdd:
		return "ADD"
	case OperationSub:
		return "SUB"
	default:
		return "UNKNOWN"
	}
}
