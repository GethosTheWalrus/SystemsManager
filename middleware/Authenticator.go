package middleware

import (
	"ThePooReview/models"
	"fmt"
	"net/http"
)

func AuthenticatorMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionToken := r.Header.Get("Authorization")
		// currentTime := time.Now()

		session := models.Session{}
		Db.Where("expiration_time > DATE(NOW())").First(&session, "id=?", sessionToken)

		if len(session.Id) > 0 {

			next.ServeHTTP(w, r)

		} else {

			w.WriteHeader(401)
			fmt.Fprintln(w, "unauthorized")
			return

		}

	})

}
