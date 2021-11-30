package caches

import (
	"github.com/Drakkar-Software/Metrics-Server/api/model"
)

// TopCache stores exchanges and trading modes top values
type TopCache struct {
	Exchanges         map[int64][]model.Top
	TradingModes      map[int64][]model.Top
}

var localTopCache = TopCache{
	Exchanges: make(map[int64][]model.Top),
	TradingModes: make(map[int64][]model.Top),
}

// GetExchangeTop returns the exchange top from asked cache value and a boolean for value availability
func GetExchangeTop(identifier int64) ([]model.Top, bool){
	value, exists := localTopCache.Exchanges[roundIdentifier(identifier)]
	return value, exists
}

// SetExchangeTop sets the cache value for a given time. Resets cache when too many values in it
func SetExchangeTop(identifier int64, top []model.Top){
	if len(localTopCache.Exchanges) > maxCacheSize {
		localTopCache.Exchanges = make(map[int64][]model.Top)
	}
	localTopCache.Exchanges[roundIdentifier(identifier)] = top
}

// GetTradingModeTop returns the trading mode top from asked cache value and a boolean for value availability
func GetTradingModeTop(identifier int64) ([]model.Top, bool){
	value, exists := localTopCache.TradingModes[roundIdentifier(identifier)]
	return value, exists
}

// SetTradingModeTop sets the cache value for a given time. Resets cache when too many values in it
func SetTradingModeTop(identifier int64, top []model.Top){
	if len(localTopCache.TradingModes) > maxCacheSize {
		localTopCache.TradingModes = make(map[int64][]model.Top)
	}
	localTopCache.TradingModes[roundIdentifier(identifier)] = top
}
