package entity

import (
	"github.com/joeljunstrom/go-luhn"
)

type OrderNumber string

func (on OrderNumber) Valid() bool {
	return luhn.Valid(string(on))
}
