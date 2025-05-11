package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {

	os.Setenv("TEST_ENV", "true")

	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
