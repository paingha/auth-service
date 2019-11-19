package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/paingha/auth-service/utils"
)

//AuthMiddleware middleware for auth checking
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if the token is found and if it is valid
		switch r.RequestURI {
		case "/api/v1/user":
			next.ServeHTTP(w, r)
		case "/api/v1/user/login":
			next.ServeHTTP(w, r)
		case "/api/v1/user/register":
			next.ServeHTTP(w, r)
		case "/api/v1/user/verify-email":
			next.ServeHTTP(w, r)
		case "/api/v1/user/token-verify":
			next.ServeHTTP(w, r)
		case "/api/v1/user/get-verify-token":
			next.ServeHTTP(w, r)
		default:
			token := r.Header.Get("Authorization")
			verifytoken, _ := utils.VerifyJWT(token)
			if token != "" &&  verifytoken{
				next.ServeHTTP(w, r)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"message":    "Unauthorized",
					"statusCode": 401,
				})
			}
		}

	})
}
