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
	rates map[Currency]map[Currency]float64
}

func NewExchangeRates() *ExchangeRates {
	rates := make(map[Currency]map[Currency]float64)
	
	rates[USD] = map[Currency]float64{
		EUR: 0.85,
		GBP: 0.73,
		JPY: 110.15,
	}
	
	rates[EUR] = map[Currency]float64{
		USD: 1.18,
		GBP: 0.86,
		JPY: 129.47,
	}
	
	rates[GBP] = map[Currency]float64{
		USD: 1.37,
		EUR: 1.16,
		JPY: 150.89,
	}
	
	rates[JPY] = map[Currency]float64{
		USD: 0.0091,
		EUR: 0.0077,
		GBP: 0.0066,
	}
	
	return &ExchangeRates{rates: rates}
}

func (er *ExchangeRates) Convert(amount float64, from, to Currency) (float64, error) {
	if from == to {
		return amount, nil
	}
	
	rate, exists := er.rates[from][to]
	if !exists {
		return 0, fmt.Errorf("conversion rate from %s to %s not available", from, to)
	}
	
	converted := amount * rate
	return math.Round(converted*100) / 100, nil
}

func (er *ExchangeRates) AddRate(from, to Currency, rate float64) {
	if er.rates[from] == nil {
		er.rates[from] = make(map[Currency]float64)
	}
	er.rates[from][to] = rate
	
	if er.rates[to] == nil {
		er.rates[to] = make(map[Currency]float64)
	}
	er.rates[to][from] = 1 / rate
}

func main() {
	converter := NewExchangeRates()
	
	amount := 100.0
	converted, err := converter.Convert(amount, USD, EUR)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f %s = %.2f %s\n", amount, USD, converted, EUR)
	
	converter.AddRate(USD, CAD, 1.25)
	cadAmount, _ := converter.Convert(50.0, USD, CAD)
	fmt.Printf("50.00 %s = %.2f CAD\n", USD, cadAmount)
}
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
			{"USD", "JPY", 149.50},
			{"JPY", "USD", 0.0067},
			{"GBP", "USD", 1.27},
			{"USD", "GBP", 0.79},
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

	converter.AddRate("EUR", "JPY", 162.50)
	result2, err := converter.Convert(50.0, "EUR", "JPY")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f EUR = %.2f JPY\n", 50.0, result2)
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
	converted, err := converter.Convert(amount, "USD", "EUR")
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		return
	}
	
	fmt.Printf("%.2f USD = %.2f EUR\n", amount, converted)
	fmt.Printf("Supported currencies: %v\n", converter.GetSupportedCurrencies())
}