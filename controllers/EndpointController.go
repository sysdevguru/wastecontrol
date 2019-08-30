package controllers

import (
	"strconv"
	. "wastecontrol/models"

	"github.com/gin-gonic/gin"
)

func CreateEndpoint(c *gin.Context) {
	var endpoint Endpoint
	var caller User
	c.BindJSON(&endpoint)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()
	if caller.UserType != 1 {
		c.JSON(401, "Unauthorized")
		return
	}

	endpoint.Create()
	mail := NewRequest([]string{endpoint.Email}, "New Endpoint", "We've created a new endpoint for you at Wastecontrol. We use basic authentication.<br>The endpoint is: https://wastecontrol.herokuapp.com/api/"+endpoint.Name+"<br>Username is: "+endpoint.Username+"<br>Password is: "+endpoint.Password, nil)
	mail.SendEmail()

	c.JSON(200, "Endpoint created")
}

func GetEndpoint(c *gin.Context) {
	var endpoint Endpoint

	endpoint.Id, _ = strconv.Atoi(c.Param("Id"))
	endpoint.GetSpecific()

	c.JSON(200, endpoint)
}

func GetEndpoints(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.GetRole()
	if user.UserType != 1 {
		c.JSON(401, "Unauthorized")
		return
	}

	c.JSON(200, GetAllEndpoints())
}

func UpdateEndpoint(c *gin.Context) {
	var endpoint Endpoint
	var caller User

	c.BindJSON(&endpoint)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	if caller.UserType != 1 {
		c.JSON(401, "Not allowed")
		return
	}

	endpoint.Update()

	c.JSON(200, endpoint)
}

func DeleteEndpoint(c *gin.Context) {
	var endpoint Endpoint
	var caller User

	endpoint.Id, _ = strconv.Atoi(c.Param("Id"))
	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	if caller.UserType != 1 {
		c.JSON(401, "Not allowed")
		return
	}

	endpoint.Delete()

	c.JSON(200, "")
}
