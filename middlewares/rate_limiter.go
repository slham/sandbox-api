package middlewares

import (
	"log"
	"net/http"

	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
)

var (
	localQuota = throttled.RateQuota{
		MaxRate:  throttled.PerMin(2000),
		MaxBurst: 500,
	}
	prodQuota = throttled.RateQuota{
		MaxRate:  throttled.PerMin(20),
		MaxBurst: 5,
	}
	quotaMap = map[string]throttled.RateQuota{
		"LOCAL": localQuota,
		"PROD":  prodQuota,
	}
)

func RateLimit(env string) func(h http.Handler) http.Handler {
	store, err := memstore.NewCtx(65536)
	if err != nil {
		log.Fatal(err)
	}

	quota := quotaMap[env]
	rateLimiter, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		log.Fatal(err)
	}

	httpRateLimiter := throttled.HTTPRateLimiterCtx{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}

	return httpRateLimiter.RateLimit
}
