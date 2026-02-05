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
	c.rates[to][from] = 1 / rate
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	if from == to {
		return amount, nil
	}
	
	if rate, exists := c.rates[from][to]; exists {
		return amount * rate, nil
	}
	
	return 0, fmt.Errorf("no conversion rate found from %s to %s", from, to)
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
	if rates, found := c.cache[baseCurrency]; found {
		if c.isCacheValid(rates.Date) {
			return &rates, nil
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

func (c *CurrencyConverter) isCacheValid(dateStr string) bool {
	cachedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return false
	}
	return time.Since(cachedDate) < c.cacheTTL
}

func (c *CurrencyConverter) Convert(amount float64, from, to string) (float64, error) {
	rates, err := c.GetRates(from)
	if err != nil {
		return 0, err
	}

	rate, exists := rates.Rates[to]
	if !exists {
		return 0, fmt.Errorf("currency %s not found in exchange rates", to)
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

	rates, err := converter.GetRates("JPY")
	if err != nil {
		fmt.Printf("Failed to get JPY rates: %v\n", err)
		return
	}

	fmt.Printf("JPY exchange rates as of %s:\n", rates.Date)
	for currency, rate := range rates.Rates {
		if currency == "USD" || currency == "EUR" || currency == "GBP" {
			fmt.Printf("  1 JPY = %.6f %s\n", rate, currency)
		}
	}
}
package main

import (
	"fmt"
)

const usdToEurRate = 0.85

func ConvertUSDToEUR(amount float64) float64 {
	return amount * usdToEurRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
}
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type ExchangeRates struct {
	Rates map[string]float64 `json:"rates"`
	Base  string             `json:"base"`
	Date  string             `json:"date"`
}

func fetchExchangeRates(apiKey string) (*ExchangeRates, error) {
	url := fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/USD")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rates ExchangeRates
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return nil, err
	}

	return &rates, nil
}

func convertCurrency(amount float64, fromCurrency, toCurrency string, rates *ExchangeRates) (float64, error) {
	if fromCurrency == toCurrency {
		return amount, nil
	}

	fromRate, fromExists := rates.Rates[fromCurrency]
	toRate, toExists := rates.Rates[toCurrency]

	if !fromExists || !toExists {
		return 0, fmt.Errorf("unsupported currency")
	}

	amountInUSD := amount / fromRate
	return amountInUSD * toRate, nil
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
		fmt.Println("Example: currency_converter 100 USD EUR")
		os.Exit(1)
	}

	amount, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Printf("Invalid amount: %v\n", err)
		os.Exit(1)
	}

	fromCurrency := os.Args[2]
	toCurrency := os.Args[3]

	rates, err := fetchExchangeRates("")
	if err != nil {
		fmt.Printf("Failed to fetch exchange rates: %v\n", err)
		os.Exit(1)
	}

	convertedAmount, err := convertCurrency(amount, fromCurrency, toCurrency, rates)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s (as of %s)\n", amount, fromCurrency, convertedAmount, toCurrency, rates.Date)
}