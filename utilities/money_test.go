package utilities

import (
	"errors"
	"testing"
)

// TestCurrencyUnit tests the CurrencyUnit type
// which is a string type representing currency codes
func TestCurrencyUnit(t *testing.T) {
	tests := []struct {
		name     string
		currency CurrencyUnit
		expected string
	}{
		{
			name:     "USD currency",
			currency: USD,
			expected: "USD",
		},
		{
			name:     "VND currency",
			currency: VND,
			expected: "VND",
		},
		{
			name:     "custom currency",
			currency: CurrencyUnit("EUR"),
			expected: "EUR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := string(tt.currency)

			// Assert
			if result != tt.expected {
				t.Errorf("CurrencyUnit string = %q; expected %q", result, tt.expected)
			}
		})
	}
}

// TestDefaultLoadExchangeRate tests the default exchange rate map
// which contains predefined USD and VND exchange rates
func TestDefaultLoadExchangeRate(t *testing.T) {
	tests := []struct {
		name             string
		currency         CurrencyUnit
		expectedRates    int
		checkFirstRate   bool
		expectedCurrency CurrencyUnit
		expectedRate     float64
	}{
		{
			name:             "USD has 2 exchange rates",
			currency:         USD,
			expectedRates:    2,
			checkFirstRate:   true,
			expectedCurrency: USD,
			expectedRate:     1,
		},
		{
			name:             "VND has 2 exchange rates",
			currency:         VND,
			expectedRates:    2,
			checkFirstRate:   true,
			expectedCurrency: USD,
			expectedRate:     0.00004,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			rates, exists := DefaultLoadExchangeRate[tt.currency]

			// Assert
			if !exists {
				t.Fatalf("DefaultLoadExchangeRate[%s] does not exist", tt.currency)
			}
			if len(rates) != tt.expectedRates {
				t.Errorf("DefaultLoadExchangeRate[%s] has %d rates; expected %d",
					tt.currency, len(rates), tt.expectedRates)
			}
			if tt.checkFirstRate && len(rates) > 0 {
				if rates[0].Currency != tt.expectedCurrency {
					t.Errorf("First rate currency = %s; expected %s", rates[0].Currency, tt.expectedCurrency)
				}
				if rates[0].Rate != tt.expectedRate {
					t.Errorf("First rate = %f; expected %f", rates[0].Rate, tt.expectedRate)
				}
			}
		})
	}
}

// mockExchangeRateLoader is a mock implementation of ILoadExchangeRate for testing
type mockExchangeRateLoader struct {
	rates map[CurrencyUnit][]ExchangeRate
	err   error
}

func (m mockExchangeRateLoader) LoadExchangeRate() (map[CurrencyUnit][]ExchangeRate, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.rates, nil
}

// TestDefaultExchangeRate_LoadExchangeRate tests the defaultExchangeRate implementation
// which loads the default exchange rates
func TestDefaultExchangeRate_LoadExchangeRate(t *testing.T) {
	tests := []struct {
		name          string
		expectError   bool
		expectedCount int
	}{
		{
			name:          "loads default rates successfully",
			expectError:   false,
			expectedCount: 2, // USD and VND
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			loader := defaultExchangeRate{}

			// Act
			rates, err := loader.LoadExchangeRate()

			// Assert
			if (err != nil) != tt.expectError {
				t.Errorf("LoadExchangeRate() error = %v; expected error = %v", err, tt.expectError)
			}
			if !tt.expectError {
				if len(rates) != tt.expectedCount {
					t.Errorf("LoadExchangeRate() returned %d currencies; expected %d", len(rates), tt.expectedCount)
				}
			}
		})
	}
}

// TestNewMoneyTransform tests the NewMoneyTransform constructor
// which creates a new MoneyTransform instance with exchange rates loaded
func TestNewMoneyTransform(t *testing.T) {
	tests := []struct {
		name       string
		loader     ILoadExchangeRate
		expectNil  bool
		shouldLoad bool
	}{
		{
			name:       "create with nil loader uses default",
			loader:     nil,
			expectNil:  false,
			shouldLoad: true,
		},
		{
			name: "create with custom loader",
			loader: mockExchangeRateLoader{
				rates: map[CurrencyUnit][]ExchangeRate{
					USD: {{Currency: VND, Rate: 25000}},
				},
			},
			expectNil:  false,
			shouldLoad: true,
		},
		{
			name: "create with loader that returns error",
			loader: mockExchangeRateLoader{
				err: errors.New("load error"),
			},
			expectNil:  false,
			shouldLoad: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			mt := NewMoneyTransform(tt.loader)

			// Assert
			if (mt == nil) != tt.expectNil {
				t.Errorf("NewMoneyTransform() nil = %v; expected nil = %v", mt == nil, tt.expectNil)
			}
			if !tt.expectNil {
				if mt.loadExchangeRate == nil {
					t.Error("NewMoneyTransform().loadExchangeRate is nil")
				}
			}
		})
	}
}

// TestMoneyTransform_LoadExchange tests the LoadExchange method
// which populates the exchange rate map from the loader
func TestMoneyTransform_LoadExchange(t *testing.T) {
	tests := []struct {
		name          string
		loader        ILoadExchangeRate
		checkKey      string
		shouldExist   bool
		expectedValue float64
	}{
		{
			name: "loads rates successfully",
			loader: mockExchangeRateLoader{
				rates: map[CurrencyUnit][]ExchangeRate{
					USD: {
						{Currency: VND, Rate: 25000},
						{Currency: USD, Rate: 1},
					},
				},
			},
			checkKey:      "USD_VND",
			shouldExist:   true,
			expectedValue: 25000,
		},
		{
			name: "loads multiple currencies",
			loader: mockExchangeRateLoader{
				rates: map[CurrencyUnit][]ExchangeRate{
					VND: {
						{Currency: USD, Rate: 0.00004},
					},
				},
			},
			checkKey:      "VND_USD",
			shouldExist:   true,
			expectedValue: 0.00004,
		},
		{
			name: "handles error gracefully",
			loader: mockExchangeRateLoader{
				err: errors.New("load error"),
			},
			checkKey:    "USD_VND",
			shouldExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mt := &MoneyTransform{loadExchangeRate: tt.loader}

			// Act
			mt.LoadExchange()

			// Assert
			value, exists := mt.Load(tt.checkKey)
			if exists != tt.shouldExist {
				t.Errorf("After LoadExchange(), key %q exists = %v; expected %v",
					tt.checkKey, exists, tt.shouldExist)
			}
			if tt.shouldExist && exists {
				rate, ok := value.(float64)
				if !ok {
					t.Errorf("Value for key %q is not float64", tt.checkKey)
				} else if rate != tt.expectedValue {
					t.Errorf("Rate for key %q = %f; expected %f", tt.checkKey, rate, tt.expectedValue)
				}
			}
		})
	}
}

// TestMoneyTransform_ExchangeRate tests the ExchangeRate method
// which converts currency values using exchange rates
func TestMoneyTransform_ExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		setupLoader    ILoadExchangeRate
		fromCurrency   CurrencyUnit
		toCurrency     CurrencyUnit
		value          float64
		expectedResult float64
		expectedOk     bool
	}{
		{
			name:           "convert USD to VND",
			setupLoader:    nil, // Use default
			fromCurrency:   USD,
			toCurrency:     VND,
			value:          100,
			expectedResult: 2500000,
			expectedOk:     true,
		},
		{
			name:           "convert VND to USD",
			setupLoader:    nil,
			fromCurrency:   VND,
			toCurrency:     USD,
			value:          25000,
			expectedResult: 1.0,
			expectedOk:     true,
		},
		{
			name:           "convert USD to USD (same currency)",
			setupLoader:    nil,
			fromCurrency:   USD,
			toCurrency:     USD,
			value:          100,
			expectedResult: 100,
			expectedOk:     true,
		},
		{
			name:           "convert with non-existent rate",
			setupLoader:    nil,
			fromCurrency:   USD,
			toCurrency:     CurrencyUnit("EUR"),
			value:          100,
			expectedResult: 100,
			expectedOk:     false,
		},
		{
			name: "convert with custom rate",
			setupLoader: mockExchangeRateLoader{
				rates: map[CurrencyUnit][]ExchangeRate{
					CurrencyUnit("EUR"): {{Currency: CurrencyUnit("JPY"), Rate: 130}},
				},
			},
			fromCurrency:   CurrencyUnit("EUR"),
			toCurrency:     CurrencyUnit("JPY"),
			value:          10,
			expectedResult: 1300,
			expectedOk:     true,
		},
		{
			name:           "convert zero value",
			setupLoader:    nil,
			fromCurrency:   USD,
			toCurrency:     VND,
			value:          0,
			expectedResult: 0,
			expectedOk:     true,
		},
		{
			name:           "convert negative value",
			setupLoader:    nil,
			fromCurrency:   USD,
			toCurrency:     VND,
			value:          -50,
			expectedResult: -1250000,
			expectedOk:     true,
		},
		{
			name:           "convert decimal value",
			setupLoader:    nil,
			fromCurrency:   USD,
			toCurrency:     VND,
			value:          0.5,
			expectedResult: 12500,
			expectedOk:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mt := NewMoneyTransform(tt.setupLoader)

			// Act
			result, ok := mt.ExchangeRate(tt.fromCurrency, tt.toCurrency, tt.value)

			// Assert
			if ok != tt.expectedOk {
				t.Errorf("ExchangeRate() ok = %v; expected %v", ok, tt.expectedOk)
			}
			if result != tt.expectedResult {
				t.Errorf("ExchangeRate() result = %f; expected %f", result, tt.expectedResult)
			}
		})
	}
}

// TestMoneyTransform_ExchangeRate_ZeroRate tests handling of zero exchange rates
// which should return the original value and false
func TestMoneyTransform_ExchangeRate_ZeroRate(t *testing.T) {
	tests := []struct {
		name           string
		value          float64
		expectedResult float64
		expectedOk     bool
	}{
		{
			name:           "zero rate returns original value",
			value:          100,
			expectedResult: 100,
			expectedOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mt := &MoneyTransform{
				loadExchangeRate: mockExchangeRateLoader{
					rates: map[CurrencyUnit][]ExchangeRate{
						USD: {{Currency: VND, Rate: 0}},
					},
				},
			}
			mt.LoadExchange()

			// Act
			result, ok := mt.ExchangeRate(USD, VND, tt.value)

			// Assert
			if ok != tt.expectedOk {
				t.Errorf("ExchangeRate() with zero rate ok = %v; expected %v", ok, tt.expectedOk)
			}
			if result != tt.expectedResult {
				t.Errorf("ExchangeRate() with zero rate result = %f; expected %f", result, tt.expectedResult)
			}
		})
	}
}

// TestMoneyTransform_ConcurrentAccess tests thread-safety of MoneyTransform
// using sync.Map for concurrent read/write operations
func TestMoneyTransform_ConcurrentAccess(t *testing.T) {
	tests := []struct {
		name       string
		goroutines int
		operations int
	}{
		{
			name:       "concurrent exchange rate conversions",
			goroutines: 10,
			operations: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mt := NewMoneyTransform(nil)
			done := make(chan bool, tt.goroutines)

			// Act: Multiple goroutines performing exchanges concurrently
			for i := 0; i < tt.goroutines; i++ {
				go func() {
					for j := 0; j < tt.operations; j++ {
						// Perform various operations
						mt.ExchangeRate(USD, VND, 100)
						mt.ExchangeRate(VND, USD, 25000)
					}
					done <- true
				}()
			}

			// Wait for all goroutines to complete
			for i := 0; i < tt.goroutines; i++ {
				<-done
			}

			// Assert: Verify we can still perform operations after concurrent access
			result, ok := mt.ExchangeRate(USD, VND, 100)
			if !ok {
				t.Error("ExchangeRate() failed after concurrent operations")
			}
			if result <= 0 {
				t.Errorf("ExchangeRate() after concurrent access = %f; expected positive value", result)
			}
		})
	}
}

// TestExchangeRate_Struct tests the ExchangeRate struct
// which holds currency and rate information
func TestExchangeRate_Struct(t *testing.T) {
	tests := []struct {
		name             string
		currency         CurrencyUnit
		rate             float64
		expectedCurrency CurrencyUnit
		expectedRate     float64
	}{
		{
			name:             "create USD to VND rate",
			currency:         VND,
			rate:             25000,
			expectedCurrency: VND,
			expectedRate:     25000,
		},
		{
			name:             "create with decimal rate",
			currency:         USD,
			rate:             0.00004,
			expectedCurrency: USD,
			expectedRate:     0.00004,
		},
		{
			name:             "create with rate of 1",
			currency:         USD,
			rate:             1,
			expectedCurrency: USD,
			expectedRate:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			er := ExchangeRate{
				Currency: tt.currency,
				Rate:     tt.rate,
			}

			// Assert
			if er.Currency != tt.expectedCurrency {
				t.Errorf("ExchangeRate.Currency = %s; expected %s", er.Currency, tt.expectedCurrency)
			}
			if er.Rate != tt.expectedRate {
				t.Errorf("ExchangeRate.Rate = %f; expected %f", er.Rate, tt.expectedRate)
			}
		})
	}
}

// TestMoneyTransform_ChainConversions tests multiple currency conversions
// to verify consistency and correctness of exchange calculations
func TestMoneyTransform_ChainConversions(t *testing.T) {
	tests := []struct {
		name           string
		initialValue   float64
		fromCurrency   CurrencyUnit
		toCurrency     CurrencyUnit
		backToCurrency CurrencyUnit
		expectedFinal  float64
		toleranceDelta float64
	}{
		{
			name:           "USD to VND and back to USD",
			initialValue:   100,
			fromCurrency:   USD,
			toCurrency:     VND,
			backToCurrency: USD,
			expectedFinal:  100,
			toleranceDelta: 0.01, // Allow small floating point error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mt := NewMoneyTransform(nil)

			// Act: Convert forward
			intermediate, ok1 := mt.ExchangeRate(tt.fromCurrency, tt.toCurrency, tt.initialValue)
			if !ok1 {
				t.Fatalf("First conversion failed")
			}

			// Act: Convert back
			final, ok2 := mt.ExchangeRate(tt.toCurrency, tt.backToCurrency, intermediate)
			if !ok2 {
				t.Fatalf("Second conversion failed")
			}

			// Assert: Check if we're back to approximately the original value
			delta := final - tt.expectedFinal
			if delta < 0 {
				delta = -delta
			}
			if delta > tt.toleranceDelta {
				t.Errorf("Chain conversion: %f %s -> %f %s -> %f %s; expected ~%f (delta: %f)",
					tt.initialValue, tt.fromCurrency,
					intermediate, tt.toCurrency,
					final, tt.backToCurrency,
					tt.expectedFinal, delta)
			}
		})
	}
}

// TestMoneyTransform_LoadExchange_Reload tests reloading exchange rates
// to ensure rates can be updated dynamically
func TestMoneyTransform_LoadExchange_Reload(t *testing.T) {
	tests := []struct {
		name        string
		initialRate float64
		updatedRate float64
	}{
		{
			name:        "reload updates exchange rates",
			initialRate: 25000,
			updatedRate: 26000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Create with initial rate
			initialLoader := mockExchangeRateLoader{
				rates: map[CurrencyUnit][]ExchangeRate{
					USD: {{Currency: VND, Rate: tt.initialRate}},
				},
			}
			mt := NewMoneyTransform(initialLoader)

			// Verify initial rate
			result1, ok1 := mt.ExchangeRate(USD, VND, 1)
			if !ok1 || result1 != tt.initialRate {
				t.Fatalf("Initial rate not set correctly: %f", result1)
			}

			// Act: Update loader and reload
			mt.loadExchangeRate = mockExchangeRateLoader{
				rates: map[CurrencyUnit][]ExchangeRate{
					USD: {{Currency: VND, Rate: tt.updatedRate}},
				},
			}
			mt.LoadExchange()

			// Assert: Verify updated rate
			result2, ok2 := mt.ExchangeRate(USD, VND, 1)
			if !ok2 {
				t.Error("ExchangeRate() failed after reload")
			}
			if result2 != tt.updatedRate {
				t.Errorf("After reload, rate = %f; expected %f", result2, tt.updatedRate)
			}
		})
	}
}
