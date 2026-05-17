package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/pseudoelement/go-file-downloader/src/utils/common"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

type PremiumClient struct {
	allowedRps int
	allowedRpm int
}

type RateLimiter struct {
	clientsActivity map[string][]int64
	premiumClients  map[string]PremiumClient
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		clientsActivity: make(map[string][]int64),
		premiumClients: map[string]PremiumClient{
			"[::1]":        {allowedRps: 10, allowedRpm: 600},
			"82.146.32.19": {allowedRps: 10, allowedRpm: 600},
		},
	}
}

/**
 * @description check client has less than 3 req/sec and 100 req/min
 */
func (rl *RateLimiter) CreateMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		clientIP := common.GetClientIP(req, true)
		clientActivity, ok := rl.clientsActivity[clientIP]
		if ok {
			now := time.Now().UnixMilli()
			clientActivity = append(clientActivity, now)
			rl.clientsActivity[clientIP] = clientActivity

			var rpsCount, rpmCount int
			for endIdx := len(clientActivity); endIdx > 0; endIdx-- {
				timestamp := clientActivity[endIdx-1]
				if timestamp > now-1_000 {
					rpsCount++
				}
				if timestamp > now-60_000 {
					rpmCount++
				} else {
					// skip full iteration over actibity array cause it can be too large
					break
				}
			}

			allowedRps, allowedRpm := rl.getReqsAllowance(req)
			log.Printf("rps: %d, rpm: %d, allowedRps: %d, allowedRpm: %d\n", rpsCount, rpmCount, allowedRps, allowedRpm)
			if rpsCount > allowedRps || rpmCount > allowedRpm {
				api_module.FailResponse(w, "too many requests", http.StatusTooManyRequests)
				return
			}
		} else {
			slice := make([]int64, 0)
			timestamp := time.Now().UnixMilli()
			slice = append(slice, timestamp)
			rl.clientsActivity[clientIP] = slice
		}

		next.ServeHTTP(w, req)
	})
}

func (rl *RateLimiter) RunCleaner(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(1 * time.Minute)
			now := time.Now().UnixMilli()
			updatedClientsActivity := make(map[string][]int64, len(rl.clientsActivity))
			for clientIP, clientActivity := range rl.clientsActivity {
				updatedActivity := make([]int64, 0)
				for _, timestamp := range clientActivity {
					if timestamp > now-60_000 {
						updatedActivity = append(updatedActivity, timestamp)
					}
				}
				updatedClientsActivity[clientIP] = updatedActivity
			}
			rl.clientsActivity = updatedClientsActivity
		}
	}
}

func (rl *RateLimiter) getReqsAllowance(req *http.Request) (allowedRps int, allowedRpm int) {
	ipAddr := common.GetClientIP(req, false)
	allowedRps = 3
	allowedRpm = 100
	premiumClientAllowance, hasPremium := rl.premiumClients[ipAddr]
	if hasPremium {
		allowedRps = premiumClientAllowance.allowedRps
		allowedRpm = premiumClientAllowance.allowedRpm
	}
	return allowedRps, allowedRpm
}
