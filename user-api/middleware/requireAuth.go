package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"user-reservation-api/initializers"
	"user-reservation-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Middleware para autenticación de usuarios
func RequireAuth(c *gin.Context) {
	// Obtener el token desde la cookie
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Token no encontrado"})
		c.Abort()
		return
	}

	// Decodificar el token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Token inválido"})
		c.Abort()
		return
	}

	// Extraer los claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["sub"] == nil || claims["role"] == nil || claims["exp"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Claims inválidos o incompletos"})
		c.Abort()
		return
	}

	// Validar expiración del token
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Token expirado"})
		c.Abort()
		return
	}

	// Obtener el ID del usuario y el rol
	userID := uint(claims["sub"].(float64))
	role := claims["role"].(string)

	// Buscar al usuario en la base de datos
	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil || user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Usuario no encontrado"})
		c.Abort()
		return
	}

	// Guardar los datos del usuario en el contexto
	c.Set("user", map[string]interface{}{
		"id":   userID,
		"role": role,
	})

	// Continuar con la solicitud
	c.Next()
}

// Middleware para restringir acceso solo a administradores
func RequireAdmin(c *gin.Context) {
	// Verificar si el usuario está autenticado
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Usuario no autenticado"})
		c.Abort()
		return
	}

	// Obtener el rol del usuario desde el contexto
	user, ok := userData.(map[string]interface{})
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No autorizado: Datos de usuario inválidos"})
		c.Abort()
		return
	}

	// Verificar si el rol es 'admin'
	if user["role"] != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Acceso denegado: Se requieren permisos de administrador"})
		c.Abort()
		return
	}

	// Continuar con la solicitud si el rol es adecuado
	c.Next()
}
