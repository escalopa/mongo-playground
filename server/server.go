package server

// use gin
import (
	"context"
	_ "embed"
	"errors"
	"net/http"

	"github.com/escalopa/mongo-playground/domain"
	"github.com/gin-gonic/gin"
)

var (
	//go:embed static/index.html
	indexHTML string
)

type storage interface {
	CreateUser(ctx context.Context, user domain.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, updatedUser domain.User) error
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]domain.User, error)
}

type Server struct {
	srv     *gin.Engine
	storage storage
}

func New(storage storage) *Server {
	s := &Server{
		srv:     gin.Default(),
		storage: storage,
	}

	s.setup()

	return s
}

func (s *Server) setup() {
	const (
		prefix = "/api/v1"
		users  = prefix + "/users"
	)

	s.srv.GET("/", s.Home)

	s.srv.GET(users, s.ListUsers)
	s.srv.POST(users, s.CreateUser)
	s.srv.PUT(users+"/:id", s.UpdateUser)
	s.srv.GET(users+"/:id", s.GetUserByID)
	s.srv.DELETE(users+"/:id", s.DeleteUser)
}

func (s *Server) Run(addr string) error {
	return s.srv.Run(addr)
}

func (s *Server) Home(c *gin.Context) {
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, indexHTML)
}

func (s *Server) CreateUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := s.storage.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (s *Server) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := s.storage.GetUserByID(c.Request.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var updatedUser domain.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.storage.UpdateUser(c.Request.Context(), id, updatedUser)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (s *Server) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	err := s.storage.DeleteUser(c.Request.Context(), id)
	if errors.Is(err, domain.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrNotFound})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (s *Server) ListUsers(c *gin.Context) {
	users, err := s.storage.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
