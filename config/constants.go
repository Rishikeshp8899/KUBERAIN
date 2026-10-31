package config

import (
	"context"
	"net/http"
)

// ContextKey defines a custom type for context keys to prevent collisions
type ContextKey string

// Application Context Keys
const (
	AdminEmailKey ContextKey = "admin_email"
	AdminRoleKey  ContextKey = "admin_role"
)

// Session & Cookie Configuration
const (
	SessionCookieName = "kuberain_admin_session"
	SessionDuration   = 12 // Session lifetime in hours
)

// System Defaults & Environment
const (
	AppName        = "Kuberain Industry"
	DefaultRole    = "Admin"
	SuperAdminRole = "SuperAdmin"
)

// Helper: GetAdminEmail retrieves the admin email safely from the request context
func GetAdminEmail(r *http.Request) string {
	if email, ok := r.Context().Value(AdminEmailKey).(string); ok {
		return email
	}
	return ""
}

// Helper: SetAdminEmail attaches the admin email to a request context
func SetAdminEmail(r *http.Request, email string) *http.Request {
	ctx := context.WithValue(r.Context(), AdminEmailKey, email)
	return r.WithContext(ctx)
}
