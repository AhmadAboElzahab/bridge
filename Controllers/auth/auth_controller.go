package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/AhmadAboElzahab/bridge/initializers"
	"github.com/AhmadAboElzahab/bridge/models"
	"github.com/AhmadAboElzahab/bridge/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct{}
type Claims struct {
	UserID uint `json:"user_id"` // Change `user_id` to `UserID` (uppercase)
	jwt.StandardClaims
}

// Constructor function for creating a new AuthController
func NewAuthController() *AuthController {

	return &AuthController{}
}

// Store (Create a new auth)
func (uc *AuthController) Signup(ctx *gin.Context) {
	var body struct {
		First_Name    string `json:"first_name"`
		Last_Name     string `json:"last_name"`
		Email         string `json:"email"`
		Password      string `json:"password"`
		Date_of_Birth string `json:"Date_of_Birth"`
	}
	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	// Process the image using the reusable function
	savePath, hash, err := utils.ProcessImageUpload(file, "./storage/users/")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user = models.User{}
	if err := initializers.DB.Where("email = ?", body.Email).First(&user).Error; err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	new_user := models.User{
		First_Name:    body.First_Name,
		Last_Name:     body.Last_Name,
		Email:         body.Email,
		Password:      string(hashedPassword),
		Date_of_Birth: body.Date_of_Birth,
		Avatar:        savePath,
		Blurhash:      hash,
	}

	if err := initializers.DB.Create(&new_user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create auth"})
		return
	}

	ctx.JSON(http.StatusCreated, new_user)
}
func (uc *AuthController) Signin(ctx *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Check password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	newAccessToken, err := generateJWT(user.ID, 15*time.Minute)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   newAccessToken,
	})
}

// Generate JWT Token
func generateJWT(user_id uint, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID: user_id, // Change `user_id` to `UserID`
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
