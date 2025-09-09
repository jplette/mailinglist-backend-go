package requestValidator

import (
	"fmt"
	"mailinglist-backend-go/services/configReader"
	"mailinglist-backend-go/services/jwtValidator"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Name     string
	LastName string
	Email    string
	Admin    bool
}

func ValidateRequest(r *http.Request) (jwt.MapClaims, error) {
	publicKey := configReader.Value("KEYCLOAK_PUBLIC_KEY")
	publicKeyComplete := "-----BEGIN PUBLIC KEY-----\n" + publicKey + "\n-----END PUBLIC KEY-----"
	bearerToken := r.Header.Get("Authorization")
	token := strings.Split(bearerToken, "Bearer ")

	if len(token) < 2 {
		return nil, fmt.Errorf("No token found in header")
	}

	return jwtValidator.ValidateToken(token[1], publicKeyComplete)
}

func isAdmin(claims jwt.MapClaims) bool {
	var groups []interface{}
	groups = claims["groups"].([]interface{})

	for _, group := range groups {
		if group == "Admin" {
			return true
		}
	}
	return false
}

func CurrentUser(claims jwt.MapClaims) User {
	return User{
		Name:     claims["given_name"].(string),
		LastName: claims["family_name"].(string),
		Email:    claims["email"].(string),
		Admin:    isAdmin(claims),
	}
}
