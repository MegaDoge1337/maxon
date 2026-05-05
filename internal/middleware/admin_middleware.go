package middleware

import (
	"net/http"

	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/pkg/response"
)

func AdminMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			role := ctx.Value(RoleKey)

			if role != domain.RolesRegistry["admin"] {
				response.WriteError(w, http.StatusUnauthorized, "permissions denied")
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
