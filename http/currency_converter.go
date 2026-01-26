package main

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
	key := from + "->" + to
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

	key := from + "->" + to
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
	converter.AddRate("GBP", "USD", 1.38)

	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
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
	converter.AddRate("JPY", "USD", 0.00905)

	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported currencies: %v\n", converter.GetSupportedCurrencies())
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

type ExchangeRates map[Currency]float64

type CurrencyConverter struct {
	rates ExchangeRates
}

func NewCurrencyConverter(rates ExchangeRates) *CurrencyConverter {
	return &CurrencyConverter{rates: rates}
}

func (c *CurrencyConverter) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}

	fromRate, ok1 := c.rates[from]
	toRate, ok2 := c.rates[to]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("unsupported currency")
	}

	usdAmount := amount / fromRate
	converted := usdAmount * toRate

	return math.Round(converted*100) / 100, nil
}

func (c *CurrencyConverter) UpdateRates(newRates ExchangeRates) {
	for currency, rate := range newRates {
		c.rates[currency] = rate
	}
}

func main() {
	rates := ExchangeRates{
		USD: 1.0,
		EUR: 0.85,
		GBP: 0.73,
		JPY: 110.5,
	}

	converter := NewCurrencyConverter(rates)

	amount := 100.0
	result, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, result, EUR)

	newRates := ExchangeRates{EUR: 0.88, GBP: 0.75}
	converter.UpdateRates(newRates)

	result, _ = converter.Convert(amount, USD, EUR)
	fmt.Printf("After update: %.2f %s = %.2f %s\n", amount, USD, result, EUR)
}package main

import (
	"fmt"
	"time"
)

type ExchangeRate struct {
	BaseCurrency    string
	TargetCurrency  string
	Rate            float64
	LastUpdated     time.Time
}

type CurrencyConverter struct {
	rates map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(base, target string, rate float64) {
	key := base + "_" + target
	c.rates[key] = ExchangeRate{
		BaseCurrency:   base,
		TargetCurrency: target,
		Rate:           rate,
		LastUpdated:    time.Now(),
	}
}

func (c *CurrencyConverter) Convert(amount float64, base, target string) (float64, error) {
	if base == target {
		return amount, nil
	}

	key := base + "_" + target
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("exchange rate not found for %s to %s", base, target)
	}

	return amount * rate.Rate, nil
}

func (c *CurrencyConverter) GetSupportedPairs() []string {
	pairs := make([]string, 0, len(c.rates))
	for key := range c.rates {
		pairs = append(pairs, key)
	}
	return pairs
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.85)
	converter.AddRate("EUR", "USD", 1.18)
	converter.AddRate("USD", "JPY", 110.5)
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported pairs: %v\n", converter.GetSupportedPairs())
}package main

import (
	"fmt"
	"os"
	"strconv"
)

const usdToEurRate = 0.92

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run currency_converter.go <amount_in_usd>")
		os.Exit(1)
	}

	usdAmount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	if usdAmount < 0 {
		fmt.Println("Amount cannot be negative")
		os.Exit(1)
	}

	eurAmount := usdAmount * usdToEurRate
	fmt.Printf("%.2f USD = %.2f EUR (Rate: %.2f)\n", usdAmount, eurAmount, usdToEurRate)
}