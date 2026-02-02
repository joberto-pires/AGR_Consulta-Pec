package middleware

import (
	"net/http"
	"strings"
)

// CacheMiddleware adiciona headers de cache para otimização
func CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remover headers de compressão problemáticos
		w.Header().Del("Content-Encoding")
		w.Header().Del("Transfer-Encoding")
		
		// Cache estático por 1 hora
		if strings.HasPrefix(r.URL.Path, "/static/") {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		
		next.ServeHTTP(w, r)
	})
}

// NoCompressionMiddleware - Middleware para garantir sem compressão
func NoCompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remover qualquer encoding problemático
		w.Header().Del("Content-Encoding")
		w.Header().Del("Transfer-Encoding")
		
		// Forçar charset UTF-8 para HTML
		if strings.Contains(r.Header.Get("Accept"), "text/html") {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		
		next.ServeHTTP(w, r)
	})
}
