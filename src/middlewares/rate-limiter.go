package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pseudoelement/go-file-downloader/src/utils/common"
	api_module "github.com/pseudoelement/golang-utils/src/api"
)

type PremiumClient struct {
	allowedRps int64
	allowedRpm int64
	allowedRph int64
}

type RateLimiter struct {
	clientsActivity map[string][]int64
	premiumClients  map[string]PremiumClient
	allowedOrigins  []string
}

func NewRateLimiter(allowedOrigins []string) *RateLimiter {
	return &RateLimiter{
		allowedOrigins:  allowedOrigins,
		clientsActivity: make(map[string][]int64),
		premiumClients: map[string]PremiumClient{
			"[::1]": {allowedRps: 10, allowedRpm: 600, allowedRph: 36_000},
			// doesn't work for clients, cause every client has own IP, not an IP of server, where client app deployed
			"82.146.32.19": {allowedRps: 10, allowedRpm: 600, allowedRph: 36_000},
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
			rpsCount, rpmCount, rphCount := int64(1), int64(1), int64(1)
			for endIdx := len(clientActivity); endIdx > 0; endIdx-- {
				timestamp := clientActivity[endIdx-1]
				if timestamp > now-1_000 {
					rpsCount++
				}
				if timestamp > now-60_000 {
					rpmCount++
				}
				if timestamp > now-60*60_000 {
					rphCount++
				} else {
					// skip full iteration over actibity array cause it can be too large
					break
				}
			}

			allowedRps, allowedRpm, allowedRph := rl.getReqsAllowance(req)
			log.Printf(
				"rps: %v, rpm: %v, rph: %v, allowedRps: %d, allowedRpm: %d, allowedRph: %d\n",
				rpsCount,
				rpmCount,
				rphCount,
				allowedRps,
				allowedRpm,
				allowedRph,
			)
			// prevents memore overflow of spammers, clientActivity could be with billions of items
			if rphCount < allowedRph {
				clientActivity = append(clientActivity, now)
				rl.clientsActivity[clientIP] = clientActivity
			}
			if rpsCount > allowedRps || rpmCount > allowedRpm || rphCount > allowedRph {
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
			for clientIP, clientActivity := range rl.clientsActivity {
				updatedActivity := make([]int64, 0)
				for _, timestamp := range clientActivity {
					oneHourAgoTimestamp := now - 60*60_000
					if timestamp > oneHourAgoTimestamp {
						updatedActivity = append(updatedActivity, timestamp)
					}
				}
				rl.clientsActivity[clientIP] = updatedActivity
			}
		}
	}
}

func (rl *RateLimiter) getReqsAllowance(req *http.Request) (allowedRps int64, allowedRpm int64, allowedRph int64) {
	ipAddr := common.GetClientIP(req, false)
	allowedRps, allowedRpm, allowedRph = 3, 100, 6_000
	premiumClientAllowance, hasPremium := rl.premiumClients[ipAddr]
	if hasPremium {
		allowedRps = premiumClientAllowance.allowedRps
		allowedRpm = premiumClientAllowance.allowedRpm
		allowedRph = premiumClientAllowance.allowedRph
	}

	// for browser requests same as origin but has / in the end
	referer := req.Header.Get("Referer")
	origin := req.Header.Get("Origin")
	for _, allowedOrigin := range rl.allowedOrigins {
		if strings.Contains(origin, allowedOrigin) && strings.Contains(referer, allowedOrigin) {
			allowedRps, allowedRpm, allowedRph = 10, 500, 30_000
		}
	}

	return allowedRps, allowedRpm, allowedRph
}
