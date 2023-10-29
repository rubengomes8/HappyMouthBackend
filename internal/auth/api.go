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

func NewAPI(db *gorm.DB) *gin.Engine {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)
	api := API{
		handler: h,
	}

	return api.SetupRouter()
}

func (a API) SetupRouter() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	{
		v1.POST("/auth/register", a.handler.Register)
		v1.POST("/auth/login", a.handler.Login)
	}
	return r
}
