package controllers

import (
	"strconv"
	"strings"
	. "wastecontrol/models"

	"github.com/gin-gonic/gin"
)

func GetContainerLog(c *gin.Context) {
	var containerLog ContainerLog

	c.BindJSON(&containerLog)

	containerLog.GetLogs()

	c.JSON(200, containerLog)
}

func GetContainers(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.GetRole()
	switch user.UserType {
	case 1:
		c.JSON(200, GetContainersForAdmin())
		return
	case 2:
		c.JSON(200, GetContainersForUser(user.Id))
		return
	case 3:
		c.JSON(200, GetContainersForUser(user.Id))
		return
	case 4:
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
}

func GetContainersFrontpage(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.GetRole()

	c.JSON(200, GetContainersForFrontpage(user.UserType, user.Id))
}

func GetContainer(c *gin.Context) {
	var container Container

	container.Id, _ = strconv.Atoi(c.Param("Id"))
	container.Get()

	container.Sensor.SensorId = container.SensorId
	container.Sensor.Get()

	c.JSON(200, container)
}

func UpdateContainer(c *gin.Context) {
	var container Container

	c.BindJSON(&container)

	container.Update()

	c.JSON(200, container)
}

func UpdateContainerLatLng(c *gin.Context) {
	var container Container

	c.BindJSON(&container)

	container.UpdateLatLng()

	c.JSON(200, container)
}

func DeleteContainer(c *gin.Context) {
	var container Container
	container.Id, _ = strconv.Atoi(c.Param("Id"))

	container.Delete()

	c.JSON(200, "")
}

func CreateContainerPickup(c *gin.Context) {
	var pickup Pickup
	var user User

	c.BindJSON(&pickup)
	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()

	pickup.UserId = user.Id
	pickup.Create()

	c.JSON(200, "")
}

func UpdateContainerPickup(c *gin.Context) {
	var pickup Pickup

	c.BindJSON(&pickup)
	pickup.Update()

	c.JSON(200, "")
}

func DeleteContainerPickup(c *gin.Context) {
	var pickup Pickup
	pickup.Id, _ = strconv.Atoi(c.Param("Id"))

	pickup.Delete()

	c.JSON(200, "")
}

func GetContainerPickups(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()

	c.JSON(200, GetUsersContainerPickups(user))
}

func CreatePickupNote(c *gin.Context) {
	var note ContainerPickupNote
	var user User

	c.BindJSON(&note)

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()

	note.Create(&user)

	c.JSON(200, "")
}

func ContainersReadyForPickup(c *gin.Context) {
	var user User
	var pickup Pickup
	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	containerQuery := "SELECT c.id, c.lat, c.lng, c.waste_distance, c.container_height, c.address, " + strconv.Itoa(user.Id) + ", cft.fraction_type, ct.container_type, d.company_name, c.note from container c LEFT JOIN container_user cu ON cu.container_id = c.id LEFT JOIN container_type ct ON ct.id = c.container_type_id LEFT JOIN container_fraction_type cft ON cft.id = c.fraction_type_id LEFT JOIN disposal_company d ON d.id = c.disposal_company_id LEFT JOIN pickup_history ph ON ph.container_id = c.id WHERE cu.user_id = " + strconv.Itoa(user.Id) + " GROUP BY c.id"

	c.JSON(200, pickup.GetContainersForCsv(containerQuery))
}

/**** Container types ****/

func CreateContainerType(c *gin.Context) {
	var containerType ContainerType
	var user User

	c.BindJSON(&containerType)

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()

	containerType.Create(&user)

	c.JSON(200, "")
}

func GetContainerTypes(c *gin.Context) {
	c.JSON(200, ReturnContainerTypes())
}

func DeleteContainerType(c *gin.Context) {
	var containerType ContainerType
	var caller User

	containerType.Id, _ = strconv.Atoi(c.Param("Id"))
	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	containerType.Delete()

	c.JSON(200, "")
}

/**** Container fractions ****/

func CreateContainerFraction(c *gin.Context) {
	var containerFraction ContainerFraction
	var user User

	c.BindJSON(&containerFraction)

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()

	containerFraction.Create(&user)

	c.JSON(200, "")
}

func GetContainerFractions(c *gin.Context) {
	c.JSON(200, ReturnContainerFractions())
}

func GetContainerFraction(c *gin.Context) {
	var containerFraction ContainerFraction

	containerFraction.Id, _ = strconv.Atoi(c.Param("Id"))
	containerFraction.Get()

	c.JSON(200, containerFraction)
}

func UpdateContainerFraction(c *gin.Context) {
	var containerFraction ContainerFraction
	var caller User

	c.BindJSON(&containerFraction)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	containerFraction.Update()

	c.JSON(200, containerFraction)
}

func DeleteContainerFraction(c *gin.Context) {
	var containerFraction ContainerFraction
	var caller User

	containerFraction.Id, _ = strconv.Atoi(c.Param("Id"))
	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	containerFraction.Delete()

	c.JSON(200, "")
}

/**** Container Groups ****/

func CreateContainerGroup(c *gin.Context) {
	var containerGroup ContainerGroup
	var user User

	c.BindJSON(&containerGroup)

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()

	containerGroup.Create(&user)

	c.JSON(200, "")
}

func GetContainerGroups(c *gin.Context) {
	var user User
	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	c.JSON(200, ReturnContainerGroups(user.Id))
}

func GetContainerGroup(c *gin.Context) {
	var containerGroup ContainerGroup

	containerGroup.Id, _ = strconv.Atoi(c.Param("Id"))
	containerGroup.Get()

	c.JSON(200, containerGroup)
}

func UpdateContainerGroup(c *gin.Context) {
	var containerGroup ContainerGroup
	var caller User

	c.BindJSON(&containerGroup)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	containerGroup.Update()

	c.JSON(200, containerGroup)
}

func DeleteContainerGroup(c *gin.Context) {
	var containerGroup ContainerGroup
	var caller User

	containerGroup.Id, _ = strconv.Atoi(c.Param("Id"))
	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	containerGroup.Delete()

	c.JSON(200, "")
}

func ReturnContainersEndpoint(c *gin.Context) {
	name := c.Param("Name")
	sensors := GetSensorNamesForExternalEndpoint(name)
	if len(sensors) == 0 {
		c.JSON(400, "No sensors for this")
		return
	}
	containerQuery := `SELECT
    	c.id,
    	c.address,
    	c.lat,
    	c.lng,
    	c.waste_distance,
    	c.container_height,
    	s.battery_level,
    	s.sensor_id,
    	s.latest_measurement,
    	s.measurement_frequency,
    	ct.container_type,
    	cft.fraction_type,
    	dc.company_name,
    	p.pickuptime,
    	p.waste_distance AS pickup_waste_distance
    FROM
    	container c
    	JOIN container_user cu ON cu.container_id = c.id
    	JOIN sensor s ON s.id = c.id
    	LEFT JOIN pickup p ON p.container_id = c.id
    	LEFT JOIN container_type ct ON ct.id = c.container_type_id
		LEFT JOIN container_fraction_type cft ON cft.id = c.fraction_type_id
		LEFT JOIN disposal_company dc ON dc.id = c.disposal_company_id
	WHERE `
	for _, name := range sensors {
		containerQuery += `s.sensor_id = "` + name + `" OR `
	}
	containerQuery = strings.TrimRight(containerQuery, "OR ")
	containerQuery += " GROUP BY s.sensor_id"

	c.JSON(200, GetContainersForExternalEndpoint(containerQuery))
}
