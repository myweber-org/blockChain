package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ExchangeRates struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	client     *http.Client
	apiBaseURL string
	cache      map[string]ExchangeRates
	cacheTTL   time.Duration
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		client:     &http.Client{Timeout: 10 * time.Second},
		apiBaseURL: "https://api.exchangerate.host/latest",
		cache:      make(map[string]ExchangeRates),
		cacheTTL:   30 * time.Minute,
	}
}

func (c *CurrencyConverter) GetRates(baseCurrency string) (*ExchangeRates, error) {
	if cached, found := c.cache[baseCurrency]; found {
		if c.isCacheValid(cached.Date) {
			return &cached, nil
		}
	}

	url := fmt.Sprintf("%s?base=%s", c.apiBaseURL, baseCurrency)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	var rates ExchangeRates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.cache[baseCurrency] = rates
	return &rates, nil
}

func (c *CurrencyConverter) isCacheValid(cacheDate string) bool {
	cachedTime, err := time.Parse("2006-01-02", cacheDate)
	if err != nil {
		return false
	}
	return time.Since(cachedTime) < c.cacheTTL
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	rates, err := c.GetRates(fromCurrency)
	if err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not supported", toCurrency)
	}

	return amount * rate, nil
}

func main() {
	converter := NewCurrencyConverter()

	amount := 100.0
	from := "USD"
	to := "EUR"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
}