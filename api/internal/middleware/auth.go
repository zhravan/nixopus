package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/raghavyuva/nixopus-api/internal/cache"
	user_storage "github.com/raghavyuva/nixopus-api/internal/features/auth/storage"
	"github.com/raghavyuva/nixopus-api/internal/storage"
	"github.com/raghavyuva/nixopus-api/internal/types"
	"github.com/raghavyuva/nixopus-api/internal/utils"
	"github.com/uptrace/bun"
)

// AuthMiddleware is a middleware that checks if the request has a valid
// authorization token. If the token is valid, it adds both the user and
// the authenticated client to the request context.
func AuthMiddleware(next http.Handler, app *storage.App, cache *cache.Cache) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.SendErrorResponse(w, "No authorization token provided", http.StatusUnauthorized)
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		user, err := verifyToken(token, app.Store.DB, app.Ctx, cache)
		if err != nil {
			log.Printf("Auth error: %v", err)
			utils.SendErrorResponse(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		log.Printf("User authenticated. ID: %s, Phone: %s", user.ID, user.Email)

		organizationID := r.Header.Get("X-Organization-Id")
		if organizationID == "" {
			utils.SendErrorResponse(w, "No organization ID provided", http.StatusBadRequest)
			return
		}

		belongsToOrg, err := cache.GetOrgMembership(r.Context(), user.ID.String(), organizationID)
		if err != nil {
			log.Printf("Cache error for org membership: %v", err)
		}

		if !belongsToOrg {
			userStorage := user_storage.UserStorage{
				DB:  app.Store.DB,
				Ctx: app.Ctx,
			}
			belongsToOrg, err = userStorage.UserBelongsToOrganization(user.ID.String(), organizationID)
			if err != nil {
				log.Printf("Error checking organization membership: %v", err)
				utils.SendErrorResponse(w, "Error verifying organization membership", http.StatusInternalServerError)
				return
			}

			if err := cache.SetOrgMembership(r.Context(), user.ID.String(), organizationID, belongsToOrg); err != nil {
				log.Printf("Error caching org membership: %v", err)
			}
		}

		if !belongsToOrg {
			utils.SendErrorResponse(w, "User does not belong to the specified organization", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), types.UserContextKey, user)
		ctx = context.WithValue(ctx, types.AuthTokenKey, token)
		ctx = context.WithValue(ctx, types.OrganizationIDKey, organizationID)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// verifyToken validates a JWT token and returns the associated user
func verifyToken(tokenString string, db *bun.DB, ctx context.Context, cache *cache.Cache) (*types.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return types.JWTSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token expired")
			}
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}

		user, err := cache.GetUser(ctx, email)
		if err != nil {
			log.Printf("Cache error for user: %v", err)
		}

		if user == nil {
			user_storage := user_storage.UserStorage{
				DB:  db,
				Ctx: ctx,
			}
			user, err = user_storage.FindUserByEmail(email)
			if err != nil {
				return nil, err
			}

			if err := cache.SetUser(ctx, email, user); err != nil {
				log.Printf("Error caching user: %v", err)
			}
		}

		return user, nil
	}

	return nil, fmt.Errorf("invalid token")
}
