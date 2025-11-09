package main

import (
	"fmt"
	"sync"
)

type CurrencyUnit string

const (
	USD CurrencyUnit = "USD"
	VND CurrencyUnit = "VND"
)

var DefaultLoadExchangeRate = map[CurrencyUnit][]ExchangeRate{
	USD: {
		{
			Currency: USD,
			Rate:     1,
		},
		{
			Currency: VND,
			Rate:     25000,
		},
	},
	VND: {
		{
			Currency: USD,
			Rate:     0.00004,
		},
		{
			Currency: VND,
			Rate:     1,
		},
	},
}

type defaultExchangeRate struct{}

func (d defaultExchangeRate) LoadExchangeRate() (map[CurrencyUnit][]ExchangeRate, error) {
	return DefaultLoadExchangeRate, nil
}

type ILoadExchangeRate interface {
	LoadExchangeRate() (map[CurrencyUnit][]ExchangeRate, error)
}
type ExchangeRate struct {
	Currency CurrencyUnit
	Rate     float64
}
type MoneyTransform struct {
	sync.Map
	loadExchangeRate ILoadExchangeRate
}

func NewMoneyTransform(loadExchangeRate ILoadExchangeRate) *MoneyTransform {
	if loadExchangeRate == nil {
		loadExchangeRate = defaultExchangeRate{}
	}
	var m MoneyTransform
	m.loadExchangeRate = loadExchangeRate
	m.LoadExchange()
	return &m
}

func (m *MoneyTransform) LoadExchange() {
	data, err := m.loadExchangeRate.LoadExchangeRate()
	if err != nil {
		return
	}
	for currency, rates := range data {
		for _, rate := range rates {
			m.Store(fmt.Sprintf("%s_%s", currency, rate.Currency), rate.Rate)
		}
	}
}

func (m *MoneyTransform) SetExchangeRate(currency CurrencyUnit, rate ExchangeRate) {
	m.Store(currency, rate)
}

func (m *MoneyTransform) ExchangeRate(currency CurrencyUnit, currencyTo CurrencyUnit, value float64) (float64, bool) {
	v, ok := m.Load(fmt.Sprintf("%s_%s", currency, currencyTo))
	if !ok {
		return value, false
	}
	exchangeRate, ok := v.(float64)
	if !ok || exchangeRate == 0 {
		return value, false
	}
	return value * exchangeRate, true
}
