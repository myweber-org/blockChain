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

type ExchangeRates struct {
	rates map[Currency]map[Currency]float64
}

func NewExchangeRates() *ExchangeRates {
	er := &ExchangeRates{
		rates: make(map[Currency]map[Currency]float64),
	}

	rates := map[Currency]float64{
		USD: 1.0,
		EUR: 0.85,
		GBP: 0.73,
		JPY: 110.0,
	}

	for from, rateFrom := range rates {
		er.rates[from] = make(map[Currency]float64)
		for to, rateTo := range rates {
			er.rates[from][to] = rateTo / rateFrom
		}
	}

	return er
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	if _, exists := er.rates[from]; !exists {
		return 0, fmt.Errorf("unsupported source currency: %s", from)
	}
	if _, exists := er.rates[from][to]; !exists {
		return 0, fmt.Errorf("unsupported target currency: %s", to)
	}

	converted := amount * er.rates[from][to]
	return math.Round(converted*100) / 100, nil
}

func main() {
	converter := NewExchangeRates()

	amount := 100.0
	result, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	result, err = converter.Convert(amount, EUR, JPY)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, EUR, result, JPY)
}