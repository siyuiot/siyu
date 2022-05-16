package qgin

import "github.com/gin-gonic/gin"

type HandlerFunc func(c *gin.Context) (data interface{}, err error)
