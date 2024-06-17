package middleware

import (
	"net"
	"net/http"
	"strings"
)

// WithCheckSubnet - Check if a certain ip in a cidr range.
func WithCheckSubnet(trusted *net.IPNet) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ipStr := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(ipStr)
			if ip == nil {
				ips := r.Header.Get("X-Forwarded-For")
				ipStrs := strings.Split(ips, ",")
				ipStr = ipStrs[0]
				ip = net.ParseIP(ipStr)
			}
			if ip == nil {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("failed parse ip from http header"))
				return
			}
			if !trusted.Contains(ip) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("IP address is not allowed"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
