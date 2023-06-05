package caches

// Refresh cache every 4 hours
var CacheRefreshPeriod = int64(3600 * 4)

var maxCacheSize = 10000

func roundIdentifier(identifier int64) int64 {
	return identifier - (identifier % CacheRefreshPeriod)
}
