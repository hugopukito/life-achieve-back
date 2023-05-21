package cors

import (
	"net/http"
)

func EnableCors(w *http.ResponseWriter, r *http.Request) string {
	header := (*w).Header()

	// Ajouter le dns du serveur et l'adresse du front en local

	allowList := map[string]bool{
		"http://localhost:8082": true,
	}

	if origin := r.Header.Get("Origin"); allowList[origin] {
		header.Add("Access-Control-Allow-Origin", origin)
	}

	header.Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
	header.Add("Access-Control-Allow-Methods", "GET, PUT, PATCH, POST, DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		(*w).Header().Add("Access-Control-Max-Age", "3600")
		(*w).WriteHeader(http.StatusOK)
		return "options"
	}
	return ""
}
