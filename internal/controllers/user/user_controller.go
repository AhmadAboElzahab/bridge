package user

import (
	"net/http"

	"github.com/AhmadAboElzahab/bridge/internal/controllers/base"
	"github.com/AhmadAboElzahab/bridge/internal/initializers"
	"github.com/AhmadAboElzahab/bridge/internal/models"
	"github.com/AhmadAboElzahab/bridge/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	base.BaseController
}

func NewUserController() *UserController {
	return &UserController{
		BaseController: base.BaseController{
			Model: &models.User{},
		},
	}
}

// UserInput is the request body for creating or updating a user.
type UserInput struct {
	FirstName string `json:"first_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
}

// Index godoc
// @Summary     List users with filters
// @Description Query users with advanced filtering, search, and pagination
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body  services.TabPayload  true  "Filter and pagination payload"
// @Success     200   {object}  object{data=[]models.User,meta=object{totalRowCount=int64}}
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/users/index [post]
func (uc *UserController) Index(ctx *gin.Context) {
	uc.BaseController.Index(ctx)
}

// Store godoc
// @Summary     Create a user
// @Description Create a new user with first name and email
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body  UserInput  true  "User data"
// @Success     201   {object}  models.User
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     409   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/users [post]
func (uc *UserController) Store(ctx *gin.Context) {
	var input UserInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ErrorJSON(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	var existing models.User
	if err := initializers.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		utils.ErrorWithCodeJSON(ctx, http.StatusConflict, "EMAIL_EXISTS", "Email already in use")
		return
	}

	user := models.User{FirstName: input.FirstName, Email: input.Email}
	if err := initializers.DB.Create(&user).Error; err != nil {
		utils.ErrorJSON(ctx, http.StatusInternalServerError, "Failed to create user", err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

// Show godoc
// @Summary     Get a user by ID
// @Description Get a single user by their ID
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Param       id    path    int  true  "User ID"
// @Success     200   {object}  models.User
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     404   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/users/{id} [get]
func (uc *UserController) Show(ctx *gin.Context) {
	uc.BaseController.Show(ctx)
}

// Update godoc
// @Summary     Update a user
// @Description Update user fields by ID
// @Tags        users
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path    int         true  "User ID"
// @Param       body  body    UserInput   true  "Fields to update"
// @Success     200   {object}  models.User
// @Failure     400   {object}  utils.ErrorResponse
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     404   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/users/{id} [put]
func (uc *UserController) Update(ctx *gin.Context) {
	uc.BaseController.Update(ctx)
}

// Delete godoc
// @Summary     Delete a user
// @Description Delete a user by ID (soft delete)
// @Tags        users
// @Produce     json
// @Security    BearerAuth
// @Param       id    path    int  true  "User ID"
// @Success     200   {object}  object{message=string}
// @Failure     401   {object}  utils.ErrorResponse
// @Failure     404   {object}  utils.ErrorResponse
// @Failure     500   {object}  utils.ErrorResponse
// @Router      /api/users/{id} [delete]
func (uc *UserController) Delete(ctx *gin.Context) {
	uc.BaseController.Delete(ctx)
}
