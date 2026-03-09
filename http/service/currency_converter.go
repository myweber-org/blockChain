
package main

import (
	"fmt"
	"os"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

type CurrencyConverter struct {
	rates []ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: []ExchangeRate{
			{"USD", "EUR", 0.92},
			{"EUR", "USD", 1.09},
			{"USD", "GBP", 0.79},
			{"GBP", "USD", 1.27},
			{"USD", "JPY", 148.50},
			{"JPY", "USD", 0.0067},
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	for _, rate := range c.rates {
		if rate.FromCurrency == fromCurrency && rate.ToCurrency == toCurrency {
			return amount * rate.Rate, nil
		}
	}

	return 0, fmt.Errorf("no exchange rate found for %s to %s", fromCurrency, toCurrency)
}

func (c *CurrencyConverter) AddRate(fromCurrency, toCurrency string, rate float64) {
	c.rates = append(c.rates, ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
	})
}

func main() {
	converter := NewCurrencyConverter()

	converter.AddRate("USD", "CAD", 1.35)
	converter.AddRate("CAD", "USD", 0.74)

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	conversions := []struct {
		amount float64
		from   string
		to     string
	}{
		{50, "EUR", "USD"},
		{200, "GBP", "JPY"},
		{1000, "USD", "CAD"},
		{5000, "JPY", "GBP"},
	}

	for _, conv := range conversions {
		result, err := converter.Convert(conv.amount, conv.from, conv.to)
		if err != nil {
			fmt.Printf("Failed to convert %.2f %s to %s: %v\n", conv.amount, conv.from, conv.to, err)
			continue
		}
		fmt.Printf("%.2f %s = %.2f %s\n", conv.amount, conv.from, result, conv.to)
	}
}