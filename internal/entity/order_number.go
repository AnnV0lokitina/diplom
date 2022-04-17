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
	if n >= 0 && n <= 9 {
		return int(n), nil
	}
	return 0, errors.New("illegal order number")
}

func checkLuhn(cardNumber OrderNumber) (bool, error) {
	sum := 0
	length := len(cardNumber)
	for i := 1; i < length; i++ {
		number, err := byteToDigit(cardNumber[i])
		if err != nil {
			return false, err
		}
		if i%2 == 0 {
			number *= 2
			if number > 9 {
				number -= 9
			}
		}
		sum += number
		if sum >= 10 {
			sum -= 10
		}
	}
	return sum == 0, nil
}
