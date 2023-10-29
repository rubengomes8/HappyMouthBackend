package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type API struct {
	handler handler
}

func NewAPI(db *gorm.DB) API {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)
	return API{
		handler: h,
	}
}

func (a API) SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/users/register", a.handler.Register)
		v1.POST("/users/login", a.handler.Login)
	}
	return r
}
