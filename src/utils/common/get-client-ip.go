package common

import (
	"log"
	"net"
	"net/http"
	"strings"
)

/**
 * for WebSocket connection ip is empty string "" every time, need to check ip from conn.RemoteAddr().String()
 */
func GetClientIP(req *http.Request, withLogs bool) string {
	xff := req.Header.Get("X-Forwarded-For")
	xri := req.Header.Get("X-Real-Ip")
	if withLogs {
		log.Println("X-Forwarded-For:", xff)
		log.Println("X-Real-Ip:", xri)
	}

	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	if xri != "" {
		return xri
	}

	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}
	return host
}
