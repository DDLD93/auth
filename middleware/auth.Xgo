package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/ddld93/auth/utilities"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		_,err := utilities.TokenMaker.VerifyToken(reqToken)
		if err != nil{
			w.WriteHeader(http.StatusUnauthorized)
		}else{

			log.Println(r.RequestURI)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		}

    })
}