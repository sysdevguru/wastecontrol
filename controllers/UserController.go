package controllers

import (
	"strconv"
	. "wastecontrol/models"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user User
	var caller User

	c.BindJSON(&user)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()
	if (user.UserType == 1 && caller.UserType == 1) || (user.UserType == 2 && caller.UserType == 1) || (user.UserType == 3 && caller.UserType <= 2) || (user.UserType == 4 && caller.UserType <= 3) {
		caller.Create(&user)
		mail := NewRequest([]string{user.Email}, "Ny bruger", "Hej "+user.FirstName+",<br><br>Dit brugernavn er: "+user.Username+"<br>Dit kodeord er: "+user.Password, nil)
		mail.SendEmail()
		c.JSON(200, "")
	} else {
		c.JSON(401, "")
	}

}

func GetUser(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.GetSpecific()

	c.JSON(200, user)
}

func AcceptGdpr(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.UpdateGdpr()

	c.JSON(200, user)
}

func GetUsers(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.GetRole()
	switch user.UserType {
	case 1:
		c.JSON(200, GetUsersForAdmin())
		return
	case 2:
		c.JSON(200, GetUsersCreatedByUser(&user))
		return
	case 3:
		c.JSON(200, GetUsersCreatedByUser(&user))
		return
	case 4:
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
}

func GetUserResellers(c *gin.Context) {
	c.JSON(200, GetResellers())
}

func GetUserOperators(c *gin.Context) {
	c.JSON(200, gin.H{
		"operators":    GetOperators(),
		"nextSensorId": GetNextSensorId(),
	})
}

func GetUserViewers(c *gin.Context) {
	var user User
	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.GetRole()

	c.JSON(200, GetViewers(user.UserType, user.Id))
}

func GetSpecificUser(c *gin.Context) {
	var user User

	user.Id, _ = strconv.Atoi(c.Param("Id"))
	user.GetSpecific()

	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	var user User
	var caller User

	c.BindJSON(&user)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	if caller.Id != user.Id && caller.UserType == 4 {
		c.JSON(401, gin.H{
			"message": "Not allowed",
		})
		return
	}

	user.Update()

	c.JSON(200, user)
}

func UpdatePassword(c *gin.Context) {
	var user User

	c.BindJSON(&user)

	user.Token = GetToken(c.GetHeader("Authorization"))
	if user.Id < 1 {
		user.GetId()
	}

	user.UpdatePass()

	c.JSON(200, gin.H{"message": "Password updated"})
}

func DeleteUser(c *gin.Context) {
	var user User
	var caller User

	user.Id, _ = strconv.Atoi(c.Param("Id"))
	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	user.Delete()

	c.JSON(200, "")
}
