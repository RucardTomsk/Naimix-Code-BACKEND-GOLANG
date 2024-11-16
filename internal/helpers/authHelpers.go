package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
)

func GetDeviceID(c *gin.Context) string {
	userAgent := c.Request.Header.Get("User-Agent")
	// Можно хэшировать User-Agent для создания уникального DeviceID
	hasher := sha256.New()
	hasher.Write([]byte(userAgent))
	return hex.EncodeToString(hasher.Sum(nil))
}
