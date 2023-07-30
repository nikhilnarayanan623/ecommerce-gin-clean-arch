package handler

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	Engine *gin.Engine
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
