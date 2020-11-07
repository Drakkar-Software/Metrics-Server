package controller

import (
	"log"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var maxRequestsPerDay = uint8(100)

type requestStats struct {
	Day          int
	RequestCount uint8
	mutex        sync.RWMutex
}

func newRequestStats() *requestStats {
	return &requestStats{
		Day:          time.Now().Day(),
		RequestCount: 1,
	}
}

func (req *requestStats) newDay(day int) {
	req.mutex.Lock()
	defer req.mutex.Unlock()
	req.Day = day
	req.RequestCount = 1
}

func (req *requestStats) isMaxedForToday() bool {
	return req.RequestCount >= maxRequestsPerDay
}

func (req *requestStats) incCounter() {
	req.mutex.Lock()
	defer req.mutex.Unlock()
	req.RequestCount++
}

var requestStatsByIP = make(map[string]*requestStats)

// IsIPAllowed returns false if API is getting spammed
func IsIPAllowed(c echo.Context) bool {
	ip := c.RealIP()
	stats, exists := requestStatsByIP[ip]
	if exists {
		// not first request
		now := time.Now().Day()
		if stats.Day != now {
			// new day: reset counter
			stats.newDay(now)
		} else {
			// check counter
			if stats.isMaxedForToday() {
				log.Println("Spam attack")
				return false
			}
			stats.incCounter()
		}
	} else {
		// first request: start stats
		requestStatsByIP[ip] = newRequestStats()
	}
	return true
}
