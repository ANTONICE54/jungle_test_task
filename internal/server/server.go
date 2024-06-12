package server

import (
	"app/internal/config"
	"app/internal/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server struct {
	router   *gin.Engine
	Config   config.Config
	Store    *gorm.DB
	CropChan chan int
	wg       *sync.WaitGroup
}

func NewServer() (*Server, error) {

	cfg := config.LoadConfig()

	server := &Server{
		CropChan: make(chan int, 5),
		Config:   cfg,
		wg:       &sync.WaitGroup{},
	}

	db := database.InitDB(&cfg)
	server.Store = db

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {

	router := gin.Default()

	router.POST("/register", server.Register)
	router.POST("/login", server.Login)
	router.POST("/images", server.authMiddleware(), server.PostImage)
	router.GET("/images", server.authMiddleware(), server.GetImages)

	server.router = router
}

func (server *Server) Start() error {
	return server.router.Run(server.Config.ServerAddress)
}

// Listens for shutdown and does all the necessary work before exiting
func (server *Server) ListenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	server.shutdown()
	os.Exit(0)
}

func (server *Server) shutdown() {
	server.wg.Wait()
	close(server.CropChan)

}
