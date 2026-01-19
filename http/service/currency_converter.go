
package main

import (
	"fmt"
)

func convertUSDToEUR(amount float64) float64 {
	const exchangeRate = 0.85
	return amount * exchangeRate
}

func main() {
	usdAmount := 100.0
	eurAmount := convertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}