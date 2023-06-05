package caches

// BotCountsCache stores bot counts cache values
type BotCountsCache struct {
	Counts       map[int64]int64
	UpdatesTimes map[int64]int64
}

var localBotCountsCache = BotCountsCache{
	Counts:       make(map[int64]int64),
	UpdatesTimes: make(map[int64]int64),
}

// GetBotCount returns the bot count from asked cache value and a boolean for value availability
func GetUpToDateBotCount(identifier int64, nowTime int64) (int64, bool) {
	key := roundIdentifier(identifier)
	value, exists := localBotCountsCache.Counts[key]
	if !exists {
		return value, false
	}
	updateTime, _ := localBotCountsCache.UpdatesTimes[key]
	// value is outdated
	if nowTime-updateTime > CacheRefreshPeriod/2 {
		return value, false
	}
	return value, true
}

// SetBotCount sets the cache value for a given time. Resets cache when too many values in it
func SetBotCount(identifier int64, count int64, updateTime int64) {
	if len(localBotCountsCache.Counts) > maxCacheSize {
		localBotCountsCache.Counts = make(map[int64]int64)
		localBotCountsCache.UpdatesTimes = make(map[int64]int64)
	}
	key := roundIdentifier(identifier)
	localBotCountsCache.Counts[key] = count
	localBotCountsCache.UpdatesTimes[key] = updateTime
}
