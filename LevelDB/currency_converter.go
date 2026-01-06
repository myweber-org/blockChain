
package main

import (
	"fmt"
	"os"
	"strconv"
)

type ExchangeRate struct {
	Currency string
	Rate     float64
}

var rates = []ExchangeRate{
	{"USD", 1.0},
	{"EUR", 0.85},
	{"GBP", 0.73},
	{"JPY", 110.0},
	{"CAD", 1.25},
}

func convertCurrency(amount float64, fromCurrency, toCurrency string) (float64, error) {
	var fromRate, toRate float64
	fromFound, toFound := false, false

	for _, rate := range rates {
		if rate.Currency == fromCurrency {
			fromRate = rate.Rate
			fromFound = true
		}
		if rate.Currency == toCurrency {
			toRate = rate.Rate
			toFound = true
		}
	}

	if !fromFound {
		return 0, fmt.Errorf("unsupported source currency: %s", fromCurrency)
	}
	if !toFound {
		return 0, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	convertedAmount := (amount / fromRate) * toRate
	return convertedAmount, nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
		fmt.Println("Example: currency_converter 100 USD EUR")
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	result, err := convertCurrency(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}
package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}