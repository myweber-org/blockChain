
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
	rates map[string]map[string]ExchangeRate
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: make(map[string]map[string]ExchangeRate),
	}
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	if c.rates[from] == nil {
		c.rates[from] = make(map[string]ExchangeRate)
	}
	c.rates[from][to] = ExchangeRate{
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

	if rates, ok := c.rates[from]; ok {
		if rate, ok := rates[to]; ok {
			return amount * rate.Rate, nil
		}
	}

	return 0, fmt.Errorf("conversion rate not available from %s to %s", from, to)
}

func (c *CurrencyConverter) GetRate(from, to string) (ExchangeRate, bool) {
	if rates, ok := c.rates[from]; ok {
		if rate, ok := rates[to]; ok {
			return rate, true
		}
	}
	return ExchangeRate{}, false
}

func main() {
	converter := NewCurrencyConverter()
	
	converter.AddRate("USD", "EUR", 0.92)
	converter.AddRate("EUR", "USD", 1.09)
	converter.AddRate("USD", "JPY", 147.50)
	
	amount := 100.0
	result, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, result)
	
	if rate, ok := converter.GetRate("USD", "JPY"); ok {
		fmt.Printf("Current USD to JPY rate: %.2f (updated: %v)\n", 
			rate.Rate, rate.LastUpdated.Format("2006-01-02 15:04:05"))
	}
}