package caches

import (
	"strconv"

	"github.com/Drakkar-Software/Metrics-Server/api/model"
)

// TopCache stores exchanges and trading modes top values
type TopCache struct {
	Exchanges         map[string][]model.Top
	TradingModes      map[string][]model.Top
	Pairs             map[string][]model.Top
	EvaluationConfigs map[string][]model.Top
}

var localTopCache = TopCache{
	Exchanges:         make(map[string][]model.Top),
	TradingModes:      make(map[string][]model.Top),
	Pairs:             make(map[string][]model.Top),
	EvaluationConfigs: make(map[string][]model.Top),
}

func getKey(traderType string, identifier int64) string {
	return traderType + strconv.FormatInt(roundIdentifier(identifier), 10)
}

// GetExchangeTop returns the exchange top from asked cache value and a boolean for value availability
func GetExchangeTop(traderType string, identifier int64) ([]model.Top, bool) {
	value, exists := localTopCache.Exchanges[getKey(traderType, identifier)]
	return value, exists
}

// SetExchangeTop sets the cache value for a given time. Resets cache when too many values in it
func SetExchangeTop(traderType string, identifier int64, top []model.Top) {
	if len(localTopCache.Exchanges) > maxCacheSize {
		localTopCache.Exchanges = make(map[string][]model.Top)
	}
	localTopCache.Exchanges[getKey(traderType, identifier)] = top
}

// GetPairTop returns the exchange top from asked cache value and a boolean for value availability
func GetPairTop(traderType string, identifier int64) ([]model.Top, bool) {
	value, exists := localTopCache.Pairs[getKey(traderType, identifier)]
	return value, exists
}

// SetPairTop sets the cache value for a given time. Resets cache when too many values in it
func SetPairTop(traderType string, identifier int64, top []model.Top) {
	if len(localTopCache.Pairs) > maxCacheSize {
		localTopCache.Pairs = make(map[string][]model.Top)
	}
	localTopCache.Pairs[getKey(traderType, identifier)] = top
}

// GetTradingModeTop returns the trading mode top from asked cache value and a boolean for value availability
func GetTradingModeTop(traderType string, identifier int64) ([]model.Top, bool) {
	value, exists := localTopCache.TradingModes[getKey(traderType, identifier)]
	return value, exists
}

// SetTradingModeTop sets the cache value for a given time. Resets cache when too many values in it
func SetTradingModeTop(traderType string, identifier int64, top []model.Top) {
	if len(localTopCache.TradingModes) > maxCacheSize {
		localTopCache.TradingModes = make(map[string][]model.Top)
	}
	localTopCache.TradingModes[getKey(traderType, identifier)] = top
}

// GetEvaluationConfigTop returns the trading mode top from asked cache value and a boolean for value availability
func GetEvaluationConfigTop(traderType string, identifier int64) ([]model.Top, bool) {
	value, exists := localTopCache.EvaluationConfigs[getKey(traderType, identifier)]
	return value, exists
}

// SetEvaluationConfigTop sets the cache value for a given time. Resets cache when too many values in it
func SetEvaluationConfigTop(traderType string, identifier int64, top []model.Top) {
	if len(localTopCache.EvaluationConfigs) > maxCacheSize {
		localTopCache.EvaluationConfigs = make(map[string][]model.Top)
	}
	localTopCache.EvaluationConfigs[getKey(traderType, identifier)] = top
}
