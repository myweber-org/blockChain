
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
			{"USD", "EUR", 0.85},
			{"EUR", "USD", 1.18},
			{"USD", "GBP", 0.73},
			{"GBP", "USD", 1.37},
			{"EUR", "GBP", 0.86},
			{"GBP", "EUR", 1.16},
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

	converter.AddRate("JPY", "USD", 0.0091)
	jpyResult, err := converter.Convert(5000, "JPY", "USD")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("%.2f JPY = %.2f USD\n", 5000.0, jpyResult)
}package main

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
	if fromCurrency == rates.Base {
		rate, exists := rates.Rates[toCurrency]
		if !exists {
			return 0, fmt.Errorf("currency %s not found", toCurrency)
		}
		return amount * rate, nil
	}

	if toCurrency == rates.Base {
		rate, exists := rates.Rates[fromCurrency]
		if !exists {
			return 0, fmt.Errorf("currency %s not found", fromCurrency)
		}
		return amount / rate, nil
	}

	fromRate, exists := rates.Rates[fromCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not found", fromCurrency)
	}

	toRate, exists := rates.Rates[toCurrency]
	if !exists {
		return 0, fmt.Errorf("currency %s not found", toCurrency)
	}

	return amount * (toRate / fromRate), nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run currency_converter.go <amount> <from_currency> <to_currency>")
		fmt.Println("Example: go run currency_converter.go 100 USD EUR")
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