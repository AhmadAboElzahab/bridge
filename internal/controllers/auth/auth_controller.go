package auth

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct{}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// NewAuthController creates a new authentication controller
func NewAuthController() *AuthController {
	return &AuthController{}
}

// Signup handles user registration
//
//	@Summary		User Signup
//	@Description	Registers a new user
//	@Tags			Auth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			first_name		formData	string	true	"First Name"
//	@Param			last_name		formData	string	false	"Last Name"
//	@Param			email			formData	string	true	"Email"	format(email)
//	@Param			password		formData	string	true	"Password"
//	@Param			date_of_birth	formData	string	false	"Date of Birth"	format(date)
//	@Param			avatar			formData	file	false	"Avatar Image"
//	@Success		201				{object}	models.User
//	@Failure		400				{object}	utils.ErrorResponse
//	@Failure		500				{object}	utils.ErrorResponse
//	@Router			/api/auth/signup [post]
func (uc *AuthController) Signup(ctx *gin.Context) {
	var body struct {
		FirstName   string `form:"first_name" binding:"required"`
		LastName    string `form:"last_name"`
		Email       string `form:"email" binding:"required,email"`
		Password    string `form:"password" binding:"required"`
		DateOfBirth string `form:"date_of_birth"`
	}

	if err := ctx.ShouldBind(&body); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	db := initializers.DB.WithContext(ctx.Request.Context())

	var user models.User
	if err := db.Where("email = ?", body.Email).First(&user).Error; err == nil {
		utils.ErrorWithCodeJSON(ctx, http.StatusConflict, "EMAIL_EXISTS", "Email already exists")
		return
	}

	var savePath, hash string
	file, err := ctx.FormFile("avatar")
	if err == nil {
		savePath, hash, err = utils.ProcessImageUpload(file, "./storage/users")
		if err != nil {
			utils.ErrorJSON(ctx, http.StatusBadRequest, "Failed to process image", err.Error())
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	newUser := models.User{
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Password:    string(hashedPassword),
		DateOfBirth: body.DateOfBirth,
		Avatar:      savePath,
		Blurhash:    hash,
	}

	if err := db.Create(&newUser).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to create user")
		return
	}

	ctx.JSON(http.StatusCreated, newUser)
}

// Signin handles user authentication
//
//	@Summary		User Signin
//	@Description	Authenticates a user and returns a JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SigninRequest	true	"Signin request body"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	utils.ErrorResponse
//	@Failure		401		{object}	utils.ErrorResponse
//	@Failure		500		{object}	utils.ErrorResponse
//	@Router			/api/auth/signin [post]
func (uc *AuthController) Signin(ctx *gin.Context) {
	var body SigninRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	db := initializers.DB.WithContext(ctx.Request.Context())

	var user models.User
	result := db.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.ErrorJSON(ctx, http.StatusUnauthorized, "Invalid email or password")
		} else {
			utils.ErrorJSON(ctx, http.StatusInternalServerError, "Database error")
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		utils.ErrorJSON(ctx, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	newAccessToken, err := generateJWT(user.ID, 24*time.Hour)
	if err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to generate JWT")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   newAccessToken,
	})
}

func generateJWT(userID uint, duration time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type SigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
