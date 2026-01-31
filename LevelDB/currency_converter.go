
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
	rates := make(map[Currency]map[Currency]float64)
	
	rates[USD] = map[Currency]float64{
		EUR: 0.92,
		GBP: 0.79,
		JPY: 148.50,
	}
	
	rates[EUR] = map[Currency]float64{
		USD: 1.09,
		GBP: 0.86,
		JPY: 161.20,
	}
	
	rates[GBP] = map[Currency]float64{
		USD: 1.27,
		EUR: 1.16,
		JPY: 187.50,
	}
	
	rates[JPY] = map[Currency]float64{
		USD: 0.0067,
		EUR: 0.0062,
		GBP: 0.0053,
	}
	
	return &ExchangeRates{rates: rates}
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}
	
	rate, exists := er.rates[from][to]
	if !exists {
		return 0, fmt.Errorf("exchange rate not available from %s to %s", from, to)
	}
	
	converted := amount * rate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) AddRate(from, to Currency, rate float64) {
	if er.rates[from] == nil {
		er.rates[from] = make(map[Currency]float64)
	}
	er.rates[from][to] = rate
	
	reciprocalRate := 1.0 / rate
	if er.rates[to] == nil {
		er.rates[to] = make(map[Currency]float64)
	}
	er.rates[to][from] = reciprocalRate
}

func main() {
	converter := NewExchangeRates()
	
	amount := 100.0
	
	result, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)
	
	result, err = converter.Convert(amount, GBP, JPY)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, GBP, result, JPY)
	
	converter.AddRate(USD, Currency("CAD"), 1.35)
	result, err = converter.Convert(amount, USD, Currency("CAD"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f CAD\n", amount, USD, result)
}
package main

import (
	"fmt"
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
			{"USD", "EUR", 0.85},
			{"EUR", "USD", 1.18},
			{"USD", "GBP", 0.73},
			{"GBP", "USD", 1.37},
			{"USD", "JPY", 110.25},
			{"JPY", "USD", 0.0091},
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

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)

	converter.AddRate("EUR", "CAD", 1.47)
	converted, _ := converter.Convert(50.0, "EUR", "CAD")
	fmt.Printf("50.00 EUR = %.2f CAD\n", converted)
}