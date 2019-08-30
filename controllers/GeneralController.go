package controllers

import (
	"unicode"
	. "wastecontrol/models"

	"github.com/gin-gonic/gin"
)

/**
 * Responding with an error
 */
func response(code int, message string, c *gin.Context) {
	resp := map[string]string{"error": message}

	c.JSON(code, resp)
	c.Abort()
}

/**
 * Check the auth endpoints
 * @type {[type]}
 */
func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		if token == "" {
			response(401, "Bearer token required", c)
			return
		}

		if !CheckToken(token) {
			response(402, "Bearer Token expired", c)
			return
		}

		c.Next()
	}
}

/**
 * Check basicauth for external API
 */
func CheckBasicAuth() gin.HandlerFunc {
	logins := GetAPILogins()
	return gin.BasicAuth(logins)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*, authorization, content-type")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

/**
 * Checks if input is an integer
 * @type {int}
 */
func IsInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}

	return true
}
