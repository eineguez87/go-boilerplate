package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"go-boilerplate/internal/auth"
)

type Handler struct {
	DB        *pgxpool.Pool
	JWTSecret string
}

func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	hash, _ := auth.HashPassword(req.Password)

	_, err := h.DB.Exec(c, `
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
`, req.Email, hash)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	var id, hash string

	err := h.DB.QueryRow(c,
		`SELECT id, password_hash FROM users WHERE email=$1`,
		req.Email,
	).Scan(&id, &hash)

	if err != nil || auth.CheckPassword(hash, req.Password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, _ := auth.GenerateToken(id, h.JWTSecret)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetString("userID")

	var email string
	h.DB.QueryRow(c,
		`SELECT email FROM users WHERE id=$1`,
		userID,
	).Scan(&email)

	c.JSON(http.StatusOK, User{ID: userID, Email: email})
}
