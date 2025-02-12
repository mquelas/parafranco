package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// RequireAuth verifica que el usuario esté autenticado
func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Token no encontrado"})
		c.Abort()
		return
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error del servidor: JWT secret no configurado"})
		c.Abort()
		return
	}

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Token inválido"})
		c.Abort()
		return
	}

	// Verifica que el token no haya expirado
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Token expirado"})
		c.Abort()
		return
	}

	// Obtiene datos del usuario
	user := map[string]interface{}{
		"id":   claims["sub"],
		"role": claims["role"], // Asegura que el token tenga esta información
	}

	fmt.Println("Usuario autenticado en RequireAuth:", user) // Depuración

	// Guarda el usuario en el contexto
	c.Set("user", user)

	c.Next()
}

func RequireAdmin(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Usuario no autenticado"})
		c.Abort()
		return
	}

	// Verifica que el usuario sea del tipo esperado
	user, ok := userData.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Datos de usuario inválidos"})
		c.Abort()
		return
	}

	// Verifica si el campo "role" existe y es un string
	role, ok := user["role"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: No se encontró el rol en los datos del usuario"})
		c.Abort()
		return
	}

	fmt.Println("Usuario en RequireAdmin (hotel-api):", user) // Depuración

	// Verifica si el usuario es admin
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: Se requieren permisos de administrador"})
		c.Abort()
		return
	}

	c.Next()
}
