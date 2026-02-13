package main

import (
	"fmt"
)

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			"USD_EUR": 0.85,
			"EUR_USD": 1.18,
			"USD_GBP": 0.73,
			"GBP_USD": 1.37,
		},
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	key := from + "_" + to
	rate, exists := c.rates[key]
	if !exists {
		return 0, fmt.Errorf("conversion rate not available for %s to %s", from, to)
	}
	return amount * rate, nil
}

func (c *CurrencyConverter) AddRate(from, to string, rate float64) {
	c.rates[from+"_"+to] = rate
}

func main() {
	converter := NewCurrencyConverter()
	
	// Convert 100 USD to EUR
	result, err := converter.Convert(100, "USD", "EUR")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("100 USD = %.2f EUR\n", result)
	
	// Add new rate and convert
	converter.AddRate("EUR_JPY", 130.5)
	jpyResult, _ := converter.Convert(50, "EUR", "JPY")
	fmt.Printf("50 EUR = %.2f JPY\n", jpyResult)
}package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	apiEndpoint string
	client      *http.Client
	cache       map[string]float64
	lastUpdated time.Time
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		apiEndpoint: "https://api.exchangerate-api.com/v4/latest/USD",
		client:      &http.Client{Timeout: 10 * time.Second},
		cache:       make(map[string]float64),
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastUpdated) < 30*time.Minute && len(c.cache) > 0 {
		return nil
	}

	resp, err := c.client.Get(c.apiEndpoint)
	if err != nil {
		return fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	var rateResponse ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&rateResponse); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	c.cache = rateResponse.Rates
	c.lastUpdated = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, fromCurrency, toCurrency string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, fromExists := c.cache[fromCurrency]
	toRate, toExists := c.cache[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency: %s or %s", fromCurrency, toCurrency)
	}

	usdAmount := amount / fromRate
	convertedAmount := usdAmount * toRate
	return convertedAmount, nil
}

func main() {
	converter := NewCurrencyConverter()
	
	amount := 100.0
	result, err := converter.Convert(amount, "EUR", "JPY")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f EUR = %.2f JPY\n", amount, result)
}