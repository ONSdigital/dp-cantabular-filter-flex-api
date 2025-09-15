package middleware

import (
	"net/http"

	"github.com/ONSdigital/dp-authorisation/auth"
	dphttp "github.com/ONSdigital/dp-net/v3/http"
	"github.com/ONSdigital/dp-net/v3/request"
	"github.com/ONSdigital/log.go/v2/log"
)

// LogIdentity checks for Service Auth or Florence Token and logs the embedded
// User Identity if present.
func LogIdentity() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if id := request.Caller(ctx); id != "" {
				log.Info(ctx, "caller identity verified", log.Data{
					"caller_identity": id,
					"URI":             r.URL.Path,
				})
				next.ServeHTTP(w, r)
				return
			}

			log.Info(ctx, "caller identity not present", log.Data{
				"URI": r.URL.Path,
			})
			http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		})
	}
}

type authHandler interface {
	Require(auth.Permissions, http.HandlerFunc) http.HandlerFunc
}

// Permissions is the middleware for checking the caller has the required
// permissions (CRUD) for the given route
type Permissions struct {
	handler authHandler
}

// NewPermissions returns a new Permissions middleware struct.
func NewPermissions(zebedeeURL string, enabled bool) *Permissions {
	if !enabled {
		return &Permissions{
			handler: &auth.NopHandler{},
		}
	}

	client := auth.NewPermissionsClient(dphttp.NewClient())
	verifier := auth.DefaultPermissionsVerifier()

	return &Permissions{
		handler: auth.NewHandler(
			auth.NewPermissionsRequestBuilder(zebedeeURL),
			client,
			verifier,
		),
	}
}

// Require is the middleware handler you wrap around each route, providing which
// permissions level is required for the call.
func (p *Permissions) Require(required auth.Permissions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p.handler.Require(required, next.ServeHTTP)(w, r)
		})
	}
}

// RequireRead is a helper function for directly adding a 'Read' requirement to a given route
func (p *Permissions) RequireRead(next http.HandlerFunc) http.HandlerFunc {
	return p.Require(auth.Permissions{Read: true})(next).(http.HandlerFunc)
}

// RequireCreate is a helper function for directly adding a 'Create' requirement to a given route
func (p *Permissions) RequireCreate(next http.HandlerFunc) http.HandlerFunc {
	return p.Require(auth.Permissions{Create: true})(next).(http.HandlerFunc)
}

// RequireUpdate is a helper function for directly adding a 'Update' requirement to a given route
func (p *Permissions) RequireUpdate(next http.HandlerFunc) http.HandlerFunc {
	return p.Require(auth.Permissions{Update: true})(next).(http.HandlerFunc)
}

// RequireDelete is a helper function for directly adding a 'Delete' requirement to a given route
func (p *Permissions) RequireDelete(next http.HandlerFunc) http.HandlerFunc {
	return p.Require(auth.Permissions{Delete: true})(next).(http.HandlerFunc)
}
