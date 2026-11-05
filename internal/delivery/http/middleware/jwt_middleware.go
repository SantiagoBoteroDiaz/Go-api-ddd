package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc{
	secret := os.Getenv("JWT_SECRET")

	return func(c *gin.Context) {
		if secret == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"Error" : "JWT no configurado", 
			})
		}

		authHeader := c.GetHeader("Authorization") 

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Token no proporcionado o formato inválido",
            })
            return
        }

		token, err := jwt.Parse(
			parts[1],
			func(token *jwt.Token) (any, error){
				if token.Method != jwt.SigningMethodHS256 {
					return nil, errors.New("Algoritmo jwt no permitido ")
				}

				return []byte(secret), nil
			},
			
			jwt.WithValidMethods([]string{"HS256"}),
		)

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error" : "Token invalido o expirado",
			})
			return 
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Error" : "Claims invalidos", 
			})
			return 
		}

		c.Set("userID", claims["sub"])

		c.Next()
	}
}