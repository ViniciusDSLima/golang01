package auth

import (
	"context"
	"fmt"
	"github.com/ViniciusDSLima/golang01/config"
	"github.com/ViniciusDSLima/golang01/types"
	"github.com/ViniciusDSLima/golang01/utils"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

type contextKey string

const UserKey contextKey = "userId"

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := utils.GetTokenFromRequest(r)

		token, err := ValidateJWT(tokenString)

		if err != nil {
			log.Printf("Error while validating token: %v", err)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("Invalid token")
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		userId := claims["userId"].(string)

		u, err := store.GetUserById(userId)

		if err != nil {
			log.Printf("Error while getting user: %v", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.Id)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userId string) (string, error) {
	expiration := time.Second * time.Duration(config.Env.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userId,
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Env.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) string {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return ""
	}

	return userId
}
