
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
    Base  string             `json:"base"`
    Rates map[string]float64 `json:"rates"`
}

func fetchExchangeRates(apiKey string) (*ExchangeRates, error) {
    url := fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", apiKey)
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

func convertCurrency(amount float64, from, to string, rates *ExchangeRates) (float64, error) {
    if from == rates.Base {
        rate, exists := rates.Rates[to]
        if !exists {
            return 0, fmt.Errorf("currency %s not found", to)
        }
        return amount * rate, nil
    }

    fromRate, exists := rates.Rates[from]
    if !exists {
        return 0, fmt.Errorf("currency %s not found", from)
    }

    toRate, exists := rates.Rates[to]
    if !exists {
        return 0, fmt.Errorf("currency %s not found", to)
    }

    return amount * (toRate / fromRate), nil
}

func main() {
    apiKey := os.Getenv("OPENEXCHANGE_API_KEY")
    if apiKey == "" {
        fmt.Println("Please set OPENEXCHANGE_API_KEY environment variable")
        os.Exit(1)
    }

    rates, err := fetchExchangeRates(apiKey)
    if err != nil {
        fmt.Printf("Error fetching exchange rates: %v\n", err)
        os.Exit(1)
    }

    if len(os.Args) != 4 {
        fmt.Println("Usage: currency_converter <amount> <from_currency> <to_currency>")
        fmt.Println("Example: currency_converter 100 USD EUR")
        os.Exit(1)
    }

    amount, err := strconv.ParseFloat(os.Args[1], 64)
    if err != nil {
        fmt.Printf("Invalid amount: %v\n", err)
        os.Exit(1)
    }

    from := os.Args[2]
    to := os.Args[3]

    converted, err := convertCurrency(amount, from, to, rates)
    if err != nil {
        fmt.Printf("Conversion error: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("%.2f %s = %.2f %s\n", amount, from, converted, to)
}