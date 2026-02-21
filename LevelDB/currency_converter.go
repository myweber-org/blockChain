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
			{"USD", "EUR", 0.92},
			{"EUR", "USD", 1.09},
			{"USD", "GBP", 0.79},
			{"GBP", "USD", 1.27},
			{"USD", "JPY", 150.25},
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
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	converter.AddRate("EUR", "CAD", 1.46)
	result2, err := converter.Convert(50.0, "EUR", "CAD")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f EUR = %.2f CAD\n", 50.0, result2)
}package main

import (
	"fmt"
)

func main() {
	const usdToEurRate = 0.85
	var usdAmount float64

	fmt.Print("Enter amount in USD: ")
	fmt.Scan(&usdAmount)

	eurAmount := usdAmount * usdToEurRate
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
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
	if _, exists := c.rates[from]; !exists {
		c.rates[from] = make(map[string]float64)
	}
	c.rates[from][to] = rate
	
	if _, exists := c.rates[to]; !exists {
		c.rates[to] = make(map[string]float64)
	}
	c.rates[to][from] = 1.0 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}

	if rates, exists := c.rates[from]; exists {
		if rate, ok := rates[to]; ok {
			return amount * rate, nil
		}
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
	converter.AddRate("USD", "GBP", 0.73)
	converter.AddRate("EUR", "JPY", 130.0)
	
	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	
	currencies := converter.GetSupportedCurrencies()
	fmt.Println("Supported currencies:", currencies)
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
			{"USD", "JPY", 110.0},
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
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)

	converter.AddRate("EUR", "GBP", 0.86)
	result, err = converter.Convert(50.0, "EUR", "GBP")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f EUR = %.2f GBP\n", 50.0, result)
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
	converter.AddRate("USD", "JPY", 110.0)
	converter.AddRate("JPY", "USD", 0.0091)

	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported currencies: %v\n", converter.GetSupportedCurrencies())
}