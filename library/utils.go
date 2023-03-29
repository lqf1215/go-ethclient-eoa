package library

import (
	"math"
	"strconv"
)

func EtherConvertAmount(amountVal string) float64 {
	pow := math.Pow(10, 18)
	amount, _ := strconv.ParseFloat(amountVal, 64)
	caleAmount := amount / pow
	return caleAmount
}
