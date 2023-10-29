package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type API struct {
	Handler Handler
}

func NewAPI(db *gorm.DB) API {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewAuthHandler(svc)
	return API{
		Handler: h,
	}

}
