package handler

import (
	"net"
	"net/http"
)

func TrustSubnetMiddleware(h http.Handler, trustSubnet string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/api/internal/stats" {
			if trustSubnet == "" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			subnet := net.ParseIP(trustSubnet)
			ipStr := r.Header.Get("X-Real-IP")
			if ipStr == "" {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			ip := net.ParseIP(ipStr)
			if ip == nil {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			if !subnet.Equal(ip) {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}
