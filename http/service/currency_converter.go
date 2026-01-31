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
}