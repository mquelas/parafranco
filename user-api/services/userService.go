package services

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"time"
	"user-reservation-api/dtos"
	"user-reservation-api/initializers"
	"user-reservation-api/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// SignUp registra un nuevo usuario
func SignUp(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body", "details": err.Error()})
		return
	}

	// Validar el formato del email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(body.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Verificar si el email ya existe en la base de datos
	var existingUser models.User
	if err := initializers.DB.First(&existingUser, "email = ?", body.Email).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	// Validar la longitud mínima del password
	if len(body.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		return
	}

	// Hashear el password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
		return
	}

	// Crear el usuario en la base de datos
	user := models.User{Email: body.Email, Password: string(hash), Role: body.Role}
	if err := initializers.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Login maneja la autenticación del usuario
func Login(dto dtos.LoginUserDTO, c *gin.Context) (*models.User, error) {
	var user models.User
	if err := initializers.DB.First(&user, "email = ?", dto.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return nil, errors.New("invalid email or password")
	}

	// Comparar el password con el hash almacenado
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return nil, errors.New("invalid email or password")
	}

	// Generar el token JWT
	tokenString, err := GenerateJWT(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return nil, err
	}

	// Guardar el token en una cookie segura
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "/", "", false, true)

	return &user, nil
}

// GenerateJWT genera un token JWT para autenticación
func GenerateJWT(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role, // 🔥 Se agrega el rol del usuario al token
		"exp":  time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", errors.New("JWT secret is not set")
	}

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Validate obtiene el usuario autenticado desde el contexto
func Validate(c *gin.Context) (*models.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return nil, errors.New("Unauthorized")
	}

	// Convertir la interfaz a *models.User
	userModel, ok := user.(*models.User)
	if !ok {
		return nil, errors.New("Invalid user data")
	}

	return userModel, nil
}
