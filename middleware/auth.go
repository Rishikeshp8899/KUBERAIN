package middleware

import (
	"context"
	"go-mvc-app/config"
	"net/http"
	"time"
)

const SessionCookieName = "admin_session"

// RequireAdminAuth restricts routes to authenticated sessions
func RequireAdminAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(SessionCookieName)
		if err != nil || cookie.Value == "" {
			http.Redirect(w, r, "/admin/login?error=unauthorized", http.StatusSeeOther)
			return
		}

		// Inject authenticated email into request context
		ctx := context.WithValue(r.Context(), config.AdminEmailKey, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// SetAuthCookie sets an HTTP-only session cookie
func SetAuthCookie(w http.ResponseWriter, email string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    email,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})
}

// ClearAuthCookie removes the session cookie on logout
func ClearAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Now().Add(-1 * time.Hour),
	})
}
