package main

import (
	"time"
	"wastecontrol/controllers"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		for {
			time.Sleep(time.Duration(60) * time.Minute)
			go controllers.PickupCron()
		}
	}()

	raven.SetDSN("https://25a144ee016841168bf39efad57f3560:7ed9ba80d9624dfe9fce4630ec87d35a@sentry.io/1471983")
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(controllers.CORSMiddleware())

	/* ---------------------------  Public routes  --------------------------- */

	/**** Public ****/
	r.POST("/login", controllers.Login)
	r.GET("/telia-test", controllers.TeliaTest)

	/**** UC endpoints ****/
	r.POST("/uc/add", controllers.UCAdd)
	r.POST("/uc/read-params", controllers.UCReadParams)
	r.POST("/uc/read-trace-data", controllers.UCReadTraceData)
	r.POST("/uc/request", controllers.UCRequest)
	r.POST("/uc/update", controllers.UCUpdate)

	/**** Sensor endpoints ****/
	r.POST("/sensor/ttn", controllers.GetTTNData)
	r.POST("/sensor/teracom", controllers.GetTeracomData)
	r.POST("/sensor/nbiot", controllers.GetNbiotData)
	r.POST("/sensor/nbiot/create", controllers.CreateNbiotSensor)
	r.POST("/sensor/lorawan/create", controllers.CreateLorawanSensor)
	r.POST("/sensor/nbiot/fourobjects", controllers.GetNbiotFourObjects)

	/**** External endpoints ****/
	externalAuth := r.Group("/", controllers.CheckBasicAuth())
	externalAuth.POST("/api/:Name", controllers.ReturnContainersEndpoint)
	/* ---------------------------  Public routes  --------------------------- */

	/* ---------------------------  Auth routes  --------------------------- */

	r.Use(controllers.CheckAuth())

	/**** User ****/
	r.GET("/user", controllers.GetUser)
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/resellers", controllers.GetUserResellers)
	r.GET("/users/operators", controllers.GetUserOperators)
	r.GET("/users/viewers", controllers.GetUserViewers)
	r.GET("/user/:Id", controllers.GetSpecificUser)
	r.PUT("/user", controllers.UpdateUser)
	r.PUT("/accept", controllers.AcceptGdpr)
	r.PUT("/user/password", controllers.UpdatePassword)
	r.POST("/user", controllers.CreateUser)
	r.DELETE("/user/:Id", controllers.DeleteUser)

	/**** Sensor ****/
	r.GET("/sensors", controllers.GetSensors)
	r.GET("/sensors/:Id", controllers.GetSensor)
	r.GET("/sensor-logs", controllers.GetSensorLogs)
	r.PUT("/sensor", controllers.UpdateSensor)
	r.PUT("/sensor/mass", controllers.UpdateMassSensor)
	r.POST("/sensor", controllers.CreateSensor)
	r.POST("/sensors/uploadCSV/:resellerId/:operatorId", controllers.UploadCSV)
	r.DELETE("/sensor/:Id", controllers.DeleteSensor)

	/**** Disposal companies ****/
	r.GET("/disposals", controllers.GetDisposals)
	r.GET("/disposals/:Id", controllers.GetDisposal)
	r.PUT("/disposal", controllers.UpdateDisposal)
	r.POST("/disposal", controllers.CreateDisposal)
	r.DELETE("/disposal/:Id", controllers.DeleteDisposal)

	/**** Containers ****/
	r.GET("/containers", controllers.GetContainers)
	r.GET("/containers-pickup", controllers.ContainersReadyForPickup)
	r.GET("/containers-frontpage", controllers.GetContainersFrontpage)
	r.GET("/containers/:Id", controllers.GetContainer)
	r.GET("/container-pickups", controllers.GetContainerPickups)
	r.PUT("/container", controllers.UpdateContainer)
	r.PUT("/container-operator", controllers.UpdateContainerLatLng)
	r.PUT("/container-pickups", controllers.UpdateContainerPickup)
	r.POST("/container/log", controllers.GetContainerLog)
	r.POST("/container/pickup-note", controllers.CreatePickupNote)
	r.POST("/container-pickups", controllers.CreateContainerPickup)
	r.DELETE("/container-pickup/:Id", controllers.DeleteContainerPickup)
	r.DELETE("/container/:Id", controllers.DeleteContainer)

	/**** Container types ****/
	r.GET("/container-types", controllers.GetContainerTypes)
	r.POST("/container-type", controllers.CreateContainerType)
	r.DELETE("/container-type/:Id", controllers.DeleteContainerType)

	/**** Container fractions ****/
	r.GET("/container-fractions", controllers.GetContainerFractions)
	r.GET("/container-fraction/:Id", controllers.GetContainerFraction)
	r.PUT("/container-fraction", controllers.UpdateContainerFraction)
	r.POST("/container-fraction", controllers.CreateContainerFraction)
	r.DELETE("/container-fraction/:Id", controllers.DeleteContainerFraction)

	/**** Container Groups ****/
	r.GET("/container-groups", controllers.GetContainerGroups)
	r.GET("/container-group/:Id", controllers.GetContainerGroup)
	r.PUT("/container-group", controllers.UpdateContainerGroup)
	r.POST("/container-group", controllers.CreateContainerGroup)
	r.DELETE("/container-group/:Id", controllers.DeleteContainerGroup)

	/**** Endpoints ****/
	r.GET("/endpoints", controllers.GetEndpoints)
	r.GET("/endpoint/:Id", controllers.GetEndpoint)
	r.PUT("/endpoint", controllers.UpdateEndpoint)
	r.POST("/endpoint", controllers.CreateEndpoint)
	r.DELETE("/endpoint/:Id", controllers.DeleteEndpoint)

	/* ---------------------------  Auth routes  --------------------------- */

	r.Run()
}
