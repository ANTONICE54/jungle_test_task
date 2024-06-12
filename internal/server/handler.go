package server

import (
	"app/internal/models"
	"app/internal/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"image"
	"net/http"
)

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func (server *Server) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
	}
	server.Store.Create(&user)

	tokenString, err := util.CreateToken(user.ID, []byte(server.Config.SecretKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	rsp := tokenResponse{
		Token: tokenString,
	}
	ctx.JSON(http.StatusOK, rsp)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (server *Server) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	user := models.User{}
	server.Store.Where("username = ?", req.Username).First(&user)

	err := util.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := util.CreateToken(user.ID, []byte(server.Config.SecretKey))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}

	rsp := tokenResponse{
		Token: tokenString,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) PostImage(ctx *gin.Context) {
	userID := ctx.Keys["user_id"]

	//Getting file from the request
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer file.Close()

	imageName := util.GetImageName(header.Filename)

	//converting file to an image object
	img, format, err := image.Decode(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for size := 30; size <= 80; size++ {
		server.wg.Add(1)
		server.CropChan <- size
		go util.CropAndSaveImage(server.wg, server.CropChan, img, uint(userID.(float64)), format, imageName, server.Store)
	}

	ctx.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", imageName))
}

func (server *Server) GetImages(ctx *gin.Context) {
	rsp := make(map[string]string)

	userID := ctx.Keys["user_id"]

	URLList := make([]models.ImageURLs, 0)
	server.Store.Find(&URLList, "user_id = ?", userID)

	for _, URLObject := range URLList {
		base64Image, err := util.ConvertToBase64(URLObject.ImageURL)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		imageName := util.GetImageName(URLObject.ImageURL)
		rsp[imageName] = base64Image
	}

	ctx.JSON(http.StatusOK, rsp)
}
