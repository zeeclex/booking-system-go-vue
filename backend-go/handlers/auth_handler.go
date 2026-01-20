package handlers

import (
	"backend/models"
	"backend/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
}

// Login handles user authentication
// @Summary      Login User
// @Description  Authenticate using Email and Password to get JWT Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Email and Password"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest

	// Validate Input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input data"})
		return
	}

	// 1. Check Email
	user, err := h.UserRepo.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email not registered"})
		return
	}

	// 2. Check Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid password"})
		return
	}

	// 3. Success
	// [TODO] For Production: Replace this hardcoded string with real JWT implementation (github.com/golang-jwt/jwt)
	token := "iti-shield-secret-token-123"
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"role":  user.Role,
		"token": token,
	})
}