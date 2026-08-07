package handler

import (
	"net"
	"net/http"
)

func TrustSubnetMiddleware(h http.Handler, trustSubnet string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/api/internal/stats" {
			if trustSubnet == "" {
				http.Error(w, "", http.StatusForbidden)
				return
			}
			_, subnet, err := net.ParseCIDR(trustSubnet)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			ipStr := r.Header.Get("X-Real-IP")
			if ipStr == "" {
				http.Error(w, "", http.StatusForbidden)
				return
			}
			ip := net.ParseIP(ipStr)
			if ip == nil {
				http.Error(w, "", http.StatusForbidden)
				return
			}
			if !subnet.Contains(ip) {
				http.Error(w, "", http.StatusForbidden)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
