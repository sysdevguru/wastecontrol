package controllers

import (
	"strconv"
	. "wastecontrol/models"

	"github.com/gin-gonic/gin"
)

func CreateDisposal(c *gin.Context) {
	var disposal Disposal
	var caller User
	c.BindJSON(&disposal)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()
	disposal.CreatedBy = caller.Id
	disposal.Create()

	c.JSON(200, gin.H{"message": "Company created"})
}

func GetDisposals(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.GetRole()
	switch user.UserType {
	case 1:
		c.JSON(200, GetDisposalsForAdmin())
		return
	case 2:
		c.JSON(200, GetDisposalsForOperator(user.Id))
		return
	case 3:
		c.JSON(200, GetDisposalsForOperator(user.Id))
		return
	case 4:
		c.JSON(200, GetDisposalsForOperator(user.Id))
		return
	}
}

func GetDisposal(c *gin.Context) {
	var disposal Disposal

	disposal.Id, _ = strconv.Atoi(c.Param("Id"))
	disposal.Get()

	c.JSON(200, disposal)
}

func UpdateDisposal(c *gin.Context) {
	var disposal Disposal
	var caller User

	c.BindJSON(&disposal)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	/* if caller.UserType != 1 {
		c.JSON(401, gin.H{
			"message": "Not allowed",
		})
		return
	} */

	disposal.Update()

	c.JSON(200, disposal)
}

func DeleteDisposal(c *gin.Context) {
	var disposal Disposal
	var caller User

	disposal.Id, _ = strconv.Atoi(c.Param("Id"))
	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	disposal.Delete()

	c.JSON(200, "")
}
