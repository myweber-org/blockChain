
package main

import (
	"fmt"
)

func ConvertUSDToEUR(amount float64) float64 {
	const exchangeRate = 0.85
	return amount * exchangeRate
}

func main() {
	usdAmount := 100.0
	eurAmount := ConvertUSDToEUR(usdAmount)
	fmt.Printf("%.2f USD = %.2f EUR\n", usdAmount, eurAmount)
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
	if len(os.Args) < 4 {
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

	result, err := convertCurrency(amount, fromCurrency, toCurrency, rates)
	if err != nil {
		fmt.Printf("Conversion error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%.2f %s = %.2f %s\n", amount, fromCurrency, result, toCurrency)
}