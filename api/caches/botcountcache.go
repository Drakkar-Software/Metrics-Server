package caches

// BotCountsCache stores bot counts cache values
type BotCountsCache struct {
	Counts         map[int64]int64
}

var localBotCountsCache = BotCountsCache{
	Counts: make(map[int64]int64),
}

// GetBotCount returns the bot count from asked cache value and a boolean for value availability
func GetBotCount(identifier int64) (int64, bool) {
	value, exists := localBotCountsCache.Counts[roundIdentifier(identifier)]
	return value, exists
}

// SetBotCount sets the cache value for a given time. Resets cache when too many values in it
func SetBotCount(identifier int64, count int64) {
	if len(localBotCountsCache.Counts) > maxCacheSize {
		localBotCountsCache.Counts = make(map[int64]int64)
	}
	localBotCountsCache.Counts[roundIdentifier(identifier)] = count
}
