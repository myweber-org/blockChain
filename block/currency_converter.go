package main

import (
	"fmt"
	"math"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
)

type ExchangeRate struct {
	From Currency
	To   Currency
	Rate float64
}

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	rates := map[string]float64{
		"USD_EUR": 0.92,
		"USD_GBP": 0.79,
		"USD_JPY": 148.50,
		"EUR_USD": 1.09,
		"EUR_GBP": 0.86,
		"EUR_JPY": 161.20,
		"GBP_USD": 1.27,
		"GBP_EUR": 1.16,
		"GBP_JPY": 187.50,
		"JPY_USD": 0.0067,
		"JPY_EUR": 0.0062,
		"JPY_GBP": 0.0053,
	}
	return &CurrencyConverter{rates: rates}
}

func (c *CurrencyConverter) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	key := fmt.Sprintf("%s_%s", from, to)
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not available for %s to %s", from, to)
	}

	converted := amount * rate
	return math.Round(converted*100) / 100, nil
}

func (c *CurrencyConverter) AddRate(from, to Currency, rate float64) {
	key := fmt.Sprintf("%s_%s", from, to)
	c.rates[key] = rate
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	from := USD
	to := EUR

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	converter.AddRate(USD, GBP, 0.78)
	testResult, _ := converter.Convert(50.0, USD, GBP)
	fmt.Printf("50.00 USD = %.2f GBP (with custom rate)\n", testResult)
}