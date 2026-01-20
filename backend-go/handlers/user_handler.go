package handlers

import (
	"backend/models"
	"backend/repository"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

// GetAllUsers handles GET /api/users
// @Summary      Get All Users
// @Description  Retrieve a list of all registered users (Admin Dashboard)
// @Tags         Users
// @Produce      json
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.Repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser handles POST /api/users
// @Summary      Create New User
// @Description  Admin creates a new user (Lecturer/Staff/Admin)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User Data"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	user.Password = string(hashedPassword)

	// Save to DB
	if err := h.Repo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// UpdateUser handles PUT /api/users/:id
// @Summary      Update User
// @Description  Update user details. If password is empty, it remains unchanged.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path int true "User ID"
// @Param        user body models.User true "Updated User Data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	// 1. Get Existing User
	existingUser, err := h.Repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 2. Bind New Data
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// 3. Update Fields
	existingUser.Name = input.Name
	existingUser.Email = input.Email
	existingUser.Role = input.Role

	// 4. Handle Password Update (Only if provided)
	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		existingUser.Password = string(hashedPassword)
	}

	// 5. Save Changes
	if err := h.Repo.Update(existingUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser handles DELETE /api/users/:id
// @Summary      Delete User
// @Description  Permanently remove a user from the system
// @Tags         Users
// @Produce      json
// @Param        id   path int true "User ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}
	if err := h.Repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}