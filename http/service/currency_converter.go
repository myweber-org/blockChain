package main

import (
	"fmt"
	"sync"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
}

type CurrencyConverter struct {
	rates map[string]map[string]float64
	mu    sync.RWMutex
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.rates[from] == nil {
		c.rates[from] = make(map[string]float64)
	}
	c.rates[from][to] = rate

	if c.rates[to] == nil {
		c.rates[to] = make(map[string]float64)
	}
	c.rates[to][from] = 1.0 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if from == to {
		return amount, nil
	}

	if _, exists := c.rates[from]; !exists {
		return 0, fmt.Errorf("no exchange rate found for currency: %s", from)
	}

	rate, exists := c.rates[from][to]
	if !exists {
		return 0, fmt.Errorf("no exchange rate found from %s to %s", from, to)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()

	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "GBP", 0.73)
	converter.AddRate("EUR", "JPY", 130.0)

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	currencies := converter.GetSupportedCurrencies()
	fmt.Println("Supported currencies:", currencies)
}package main

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
	rates map[Currency]float64
}

func NewExchangeRates() *ExchangeRates {
	return &ExchangeRates{
		rates: map[Currency]float64{
			USD: 1.0,
			EUR: 0.85,
			GBP: 0.73,
			JPY: 110.0,
		},
	}
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	fromRate, fromOk := er.rates[from]
	toRate, toOk := er.rates[to]

	if !fromOk || !toOk {
		return 0, fmt.Errorf("unsupported currency")
	}

	baseAmount := amount / fromRate
	converted := baseAmount * toRate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) UpdateRate(currency Currency, rate float64) {
	er.rates[currency] = rate
}

func main() {
	rates := NewExchangeRates()

	converted, err := rates.Convert(100.0, USD, EUR)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	fmt.Printf("100 USD = %.2f EUR\n", converted)

	rates.UpdateRate(JPY, 115.5)
	converted, _ = rates.Convert(5000.0, JPY, GBP)
	fmt.Printf("5000 JPY = %.2f GBP\n", converted)
}package main

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
			{"USD", "JPY", 110.5},
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

	return 0, fmt.Errorf("conversion rate not found for %s to %s", fromCurrency, toCurrency)
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
	cadResult, _ := converter.Convert(50.0, "EUR", "CAD")
	fmt.Printf("50.00 EUR = %.2f CAD\n", cadResult)
}