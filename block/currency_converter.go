
package main

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
	var pairs []string
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
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ExchangeRates struct {
	Base  string             `json:"base"`
	Rates map[string]float64 `json:"rates"`
	Date  string             `json:"date"`
}

type CurrencyConverter struct {
	rates     map[string]float64
	lastFetch time.Time
	apiKey    string
}

func NewCurrencyConverter(apiKey string) *CurrencyConverter {
	return &CurrencyConverter{
		rates:  make(map[string]float64),
		apiKey: apiKey,
	}
}

func (c *CurrencyConverter) fetchRates() error {
	if time.Since(c.lastFetch) < 30*time.Minute && len(c.rates) > 0 {
		return nil
	}

	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD")
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rates ExchangeRates
	if err := json.Unmarshal(body, &rates); err != nil {
		return err
	}

	c.rates = rates.Rates
	c.lastFetch = time.Now()
	return nil
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if err := c.fetchRates(); err != nil {
		return 0, err
	}

	fromRate, fromExists := c.rates[from]
	toRate, toExists := c.rates[to]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency")
	}

	if from == "USD" {
		return amount * toRate, nil
	}

	usdAmount := amount / fromRate
	return usdAmount * toRate, nil
}

func (c *CurrencyConverter) SupportedCurrencies() []string {
	if err := c.fetchRates(); err != nil {
		return []string{}
	}

	currencies := make([]string, 0, len(c.rates))
	for currency := range c.rates {
		currencies = append(currencies, currency)
	}
	return currencies
}

func main() {
	converter := NewCurrencyConverter("")

	amount := 100.0
	from := "EUR"
	to := "JPY"

	result, err := converter.Convert(amount, from, to)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, from, result, to)
	fmt.Printf("Supported currencies: %v\n", converter.SupportedCurrencies())
}
package main

import (
	"fmt"
	"os"
	"strconv"
)

type ExchangeRate struct {
	Currency string
	Rate     float64
}

var rates = []ExchangeRate{
	{"USD", 1.0},
	{"EUR", 0.85},
	{"GBP", 0.73},
	{"JPY", 110.0},
	{"CAD", 1.25},
}

func convertCurrency(amount float64, fromCurrency, toCurrency string) (float64, error) {
	var fromRate, toRate float64
	foundFrom, foundTo := false, false

	for _, rate := range rates {
		if rate.Currency == fromCurrency {
			fromRate = rate.Rate
			foundFrom = true
		}
		if rate.Currency == toCurrency {
			toRate = rate.Rate
			foundTo = true
		}
	}

	if !foundFrom {
		return 0, fmt.Errorf("unsupported source currency: %s", fromCurrency)
	}
	if !foundTo {
		return 0, fmt.Errorf("unsupported target currency: %s", toCurrency)
	}

	return (amount / fromRate) * toRate, nil
}

func listSupportedCurrencies() {
	fmt.Println("Supported currencies:")
	for _, rate := range rates {
		fmt.Printf("  %s (rate: %.4f)\n", rate.Currency, rate.Rate)
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
		fmt.Println("Example: currency_converter 100 USD EUR")
		listSupportedCurrencies()
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	result, err := convertCurrency(amount, fromCurrency, toCurrency)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}