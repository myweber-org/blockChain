package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	apiKey     string
	baseURL    string
	lastUpdate time.Time
	cache      map[string]float64
	cacheTTL   time.Duration
}

func NewCurrencyConverter(apiKey string) *CurrencyConverter {
	return &CurrencyConverter{
		apiKey:   apiKey,
		baseURL:  "https://api.exchangerate-api.com/v4/latest/",
		cache:    make(map[string]float64),
		cacheTTL: 30 * time.Minute,
	}
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if time.Since(c.lastUpdate) > c.cacheTTL {
		if err := c.updateRates(from); err != nil {
			return 0, err
		}
	}

	rate, exists := c.cache[to]
	if !exists {
		return 0, fmt.Errorf("exchange rate for %s not available", to)
	}

	return amount * rate, nil
}

func (c *CurrencyConverter) updateRates(baseCurrency string) error {
	url := fmt.Sprintf("%s%s", c.baseURL, baseCurrency)
	
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rateResponse ExchangeRateResponse
	if err := json.Unmarshal(body, &rateResponse); err != nil {
		return err
	}

	c.cache = rateResponse.Rates
	c.lastUpdate = time.Now()
	
	return nil
}

func main() {
	converter := NewCurrencyConverter("your-api-key-here")
	
	amount := 100.0
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
}