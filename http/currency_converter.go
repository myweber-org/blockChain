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

func (er *ExchangeRates) AddRate(currency Currency, rate float64) {
	er.rates[currency] = rate
}

func main() {
	rates := NewExchangeRates()
	
	converted, err := rates.Convert(100.0, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("100 USD = %.2f EUR\n", converted)
	
	rates.AddRate("CAD", 1.25)
	cadAmount, _ := rates.Convert(50.0, USD, "CAD")
	fmt.Printf("50 USD = %.2f CAD\n", cadAmount)
}
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
	rates []ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: []ExchangeRate{
			{USD, EUR, 0.92},
			{USD, GBP, 0.79},
			{USD, JPY, 149.50},
			{EUR, USD, 1.09},
			{EUR, GBP, 0.86},
			{EUR, JPY, 162.50},
			{GBP, USD, 1.27},
			{GBP, EUR, 1.16},
			{GBP, JPY, 189.24},
			{JPY, USD, 0.0067},
			{JPY, EUR, 0.0062},
			{JPY, GBP, 0.0053},
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from Currency, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	for _, rate := range c.rates {
		if rate.From == from && rate.To == to {
			result := amount * rate.Rate
			return math.Round(result*100) / 100, nil
		}
	}

	return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
}

func (c *CurrencyConverter) AddRate(from Currency, to Currency, rate float64) {
	c.rates = append(c.rates, ExchangeRate{from, to, rate})
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	converted, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, EUR)

	converted, err = converter.Convert(amount, EUR, GBP)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, EUR, converted, GBP)

	converter.AddRate(USD, CAD, 1.38)
	converted, err = converter.Convert(amount, USD, CAD)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, CAD)
}package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	LastUpdated  time.Time
}

type CurrencyConverter struct {
	rates map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	key := from + "_" + to
	c.rates[key] = ExchangeRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		LastUpdated:  time.Now(),
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	currencies := make(map[string]bool)
	for _, rate := range c.rates {
		currencies[rate.FromCurrency] = true
		currencies[rate.ToCurrency] = true
	}

	result := make([]string, 0, len(currencies))
	for currency := range currencies {
		result = append(result, currency)
	}
	return result
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	
	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	fmt.Printf("Supported currencies: %v\n", converter.GetSupportedCurrencies())
}package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float64
	LastUpdated  time.Time
}

type CurrencyConverter struct {
	rates map[string]map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]float64),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
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
	if from == to {
		return amount, nil
	}
	
	if rate, exists := c.rates[from][to]; exists {
		return amount * rate, nil
	}
	
	return 0, fmt.Errorf("no conversion rate available from %s to %s", from, to)
}

func (c *CurrencyConverter) GetSupportedCurrencies() []string {
	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("USD", "JPY", 110.0)
	converter.AddRate("EUR", "GBP", 0.86)
	
	amount := 100.0
	from := "USD"
	to := "EUR"
	
	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
	
	fmt.Println("Supported currencies:", converter.GetSupportedCurrencies())
}