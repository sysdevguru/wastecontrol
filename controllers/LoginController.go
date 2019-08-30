package controllers

import(
	"github.com/gin-gonic/gin"
	. "wastecontrol/models"
)

func Login(c *gin.Context) {
	var user User
	c.BindJSON(&user)

	user.CheckUsernameAndPassword()

	if user.Token != "" {
		c.JSON(200, user)
	} else {
		c.JSON(500, gin.H{
			"data":  	"Wrong username or password",
		})
	}
}
