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

func (c *CurrencyConverter) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	for _, rate := range c.rates {
		if rate.From == from && rate.To == to {
			return math.Round((amount*rate.Rate)*100) / 100, nil
		}
	}

	return 0, fmt.Errorf("no exchange rate found for %s to %s", from, to)
}

func (c *CurrencyConverter) AddRate(from, to Currency, rate float64) {
	c.rates = append(c.rates, ExchangeRate{from, to, rate})
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	result, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	result, err = converter.Convert(amount, EUR, GBP)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, EUR, result, GBP)

	converter.AddRate(USD, CAD, 1.36)
	result, err = converter.Convert(amount, USD, CAD)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, CAD)
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

	return 0, fmt.Errorf("conversion rate not available from %s to %s", from, to)
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
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	
	supported := converter.GetSupportedCurrencies()
	fmt.Println("Supported currencies:", supported)
}