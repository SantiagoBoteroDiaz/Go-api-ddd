package handlers

import (
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

type AunthHandler struct {}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthHandler() *AunthHandler {
	return &AunthHandler{}
}

func(h *AunthHandler) GenerateToken(c *gin.Context){
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error" : "Credenciales mal formadas", 
		})
		return
	}

	if req.Username != os.Getenv("AUTH_USER") || 
	req.Password != os.Getenv("AUTH_PASSWORD") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error" : "Credenciales incorrectas", 
		})
		return
	}

	claims := jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(time.Hour).Unix(), 
		"iat" : time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")), 
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error" : "No se pudo generar el token", 
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": signedToken,
        "token_type":   "Bearer",
        "expires_in":   3600,
	})
}