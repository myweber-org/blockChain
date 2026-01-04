
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
	{"EUR", 0.92},
	{"GBP", 0.79},
	{"JPY", 149.5},
	{"CAD", 1.36},
}

func convertCurrency(amount float64, fromCurrency, toCurrency string) (float64, error) {
	var fromRate, toRate float64
	foundFrom, foundTo := false, false

	for _, rate := range rates {
		if rate.Currency == fromCurrency {
			fromRate = rate.Rate
			foundFrom = true
		}
		if rate.Currency == toCurrency {
			toRate = rate.Rate
			foundTo = true
		}
	}

	if !foundFrom {
		return 0, fmt.Errorf("unsupported source currency: %s", fromCurrency)
	}
	if !foundTo {
		return 0, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	return amount * (toRate / fromRate), nil
}

func listSupportedCurrencies() {
	fmt.Println("Supported currencies:")
	for _, rate := range rates {
		fmt.Printf("  %s (rate: %.4f)\n", rate.Currency, rate.Rate)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
		fmt.Println("Example: currency_converter 100 USD EUR")
		listSupportedCurrencies()
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
}package main

import (
	"fmt"
)

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			"USD_EUR": 0.85,
			"EUR_USD": 1.18,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("conversion rate not available for %s to %s", from, to)
	}
	return amount * rate, nil
}

func main() {
	converter := NewCurrencyConverter()
	
	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
}