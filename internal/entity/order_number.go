package entity

import (
	"errors"
)

const code0 = 48

type OrderNumber string

func (on OrderNumber) Valid() bool {
	valid, err := checkLuhn(on)
	if err != nil {
		return false
	}
	return valid
}

func byteToDigit(b byte) (int, error) {
	n := b - code0
	if n <= 9 {
		return int(n), nil
	}
	return 0, errors.New("illegal order number")
}

func checkLuhn(cardNumber OrderNumber) (bool, error) {
	sum := 0
	length := len(cardNumber)
	for i := 0; i < length; i++ {
		number, err := byteToDigit(cardNumber[i])
		if err != nil {
			return false, err
		}
		var val int
		if i%2 == 0 {
			val = number
		} else {
			val = number * 2
			if val > 9 {
				val -= 9
			}
		}
		sum += val
	}
	return sum%10 == 0, nil
}
