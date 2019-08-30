package models

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"wastecontrol/db"
)

type Container struct {
	Id                   int       `json:"id" binding:""`
	SensorId             string    `json:"sensor_id" binding:""`
	MeasurementFrequency int       `json:"measurement_frequency" binding:""`
	WasteDistance        float64   `json:"waste_distance" binding:""`
	LatestMeasurement    *string   `json:"latest_measurement" binding:""`
	Height               int       `json:"height" binding:""`
	Size                 *int      `json:"size" binding:""`
	SizeType             *string   `json:"size_type" binding:""`
	LatestPickup         *string   `json:"latest_pickup" binding:""`
	Note                 *string   `json:"note" binding:""`
	Address              *string   `json:"address" binding:""`
	Latitude             *float64  `json:"lat" binding:""`
	Longitude            *float64  `json:"lng" binding:""`
	CreatedBy            *int      `json:"created_by" binding:""`
	Fraction             *string   `json:"fraction" binding:""`
	FractionId           *int      `json:"fraction_id" binding:""`
	FractionSize         *int      `json:"fraction_size" binding:""`
	ContainerGroupId     *int      `json:"container_group_id" binding:""`
	ContainerTypeId      *int      `json:"container_type_id" binding:""`
	ContainerType        *string   `json:"container_type" binding:""`
	DisposalId           *int      `json:"disposal_id" binding:""`
	DisposalCompany      *string   `json:"disposal_company" binding:""`
	Sensor               Sensor    `json:"sensor" binding:""`
	ErrorRequestSensor   *int      `json:"error_request_sensor" binding:""`
	Viewer               User      `json:"viewer" binding:""`
	Operator             User      `json:"operator" binding:""`
	PickupHistory        []LogData `json:"pickup_history" binding:""`
	LogId                int       `json:"log_id" binding:""`
}

type ContainerType struct {
	Id        int    `json:"id" 					binding:""`
	Type      string `json:"type" 				binding:""`
	SizeLiter *int   `json:"size_liter"	binding:""`
	SizeM3    *int   `json:"size_m3" 		binding:""`
}

type ContainerPickupNote struct {
	Id          int     `json:"id" 					binding:""`
	ContainerId int     `json:"container_id" 				binding:""`
	Weight      *string `json:"weight"	binding:""`
	Note        *string `json:"note" 		binding:""`
}

type ContainerFraction struct {
	Id    int     `json:"id" 														binding:""`
	Type  string  `json:"type" 													binding:""`
	Color *string `json:"color"													binding:""`
	Size  *int    `json:"size_m3" 											binding:""`
	WCF   *int    `json:"weight_calibration_factor" 		binding:""`
}

type ContainerGroup struct {
	Id                int    `json:"id" 									binding:""`
	Name              string `json:"name" 								binding:""`
	ContainerTypeId   int    `json:"container_type_id" 	binding:""`
	DisposalCompanyId int    `json:"disposal_company_id"	binding:""`
	ViewerId          int    `json:"viewer_id" 					binding:""`
	MaxPickupFreq     string `json:"max_pickup_freq" 		binding:""`
	FractionId        int    `json:"fraction_id" 				binding:""`
	Address           string `json:"address" 						binding:""`
	FullContainers    string `json:"full_containers" 		binding:""`
	WasteDistance     string `json:"waste_distance" 			binding:""`
}

type ContainerLog struct {
	ContainerId int        `json:"container_id" binding:""`
	StartDate   string     `json:"start_date" binding:""`
	EndDate     string     `json:"end_date" binding:""`
	Logs        *[]LogData `json:"log_data" binding:""`
	Max100      bool       `json:"max_100" binding:""`
}

type ContainerJSON struct {
	Id              int     `json:"container_id" binding:""`
	Note            *string `json:"note" binding:""`
	LatestPickup    *string `json:"latest_pickup" binding:""`
	Address         *string `json:"address" binding:""`
	ViewerId        *int    `json:"viewer_id" binding:""`
	Fraction        *string `json:"fraction" binding:""`
	DisposalCompany *string `json:"disposal_company" binding:""`
	WasteDistance   int     `json:"waste_distance" binding:""`
	ContainerType   *string `json:"container_type" binding:""`
}

type LogData struct {
	Id                 int     `json:"id" binding:""`
	DecryptedPlainText string  `json:"decrypted_plain_text" binding:""`
	DevEUI             string  `json:"dev_eui" binding:""`
	ContainerHeight    int     `json:"container_height" binding:""`
	Timestamp          string  `json:"timestamp" binding:""`
	D1                 float64 `json:"d1" binding:""`
	D2                 float64 `json:"d2" binding:""`
	D3                 float64 `json:"d3" binding:""`
	D4                 float64 `json:"d4" binding:""`
	WasteDistance      int     `json:"waste_distance" binding:""`
}

type DecryptedObject struct {
	Battery       int   `json:"BAT" binding:""`
	WasteDistance int64 `json:"WD"	binding:""`
	Error         int   `json:"ERR" binding:""`
	PRD           *int  `json:"PRD" binding:""`
}

type Pickup struct {
	Id                int     `json:"id" binding:""`
	Search            *string `json:"search" binding:""`
	ContainerId       *int    `json:"container_id" binding:""`
	FractionTypeId    *int    `json:"fraction_type_id" binding:""`
	FractionType      *string `json:"fraction_type" binding:""`
	DisposalCompanyId *int    `json:"disposal_company_id" binding:""`
	DisposalCompany   *string `json:"disposal_company" binding:""`
	WasteDistance     *int    `json:"waste_distance" binding:""`
	ContainerTypeId   *int    `json:"container_type_id" binding:""`
	ContainerType     *string `json:"container_type" binding:""`
	ViewerId          *int    `json:"viewer_id" binding:""`
	PickupTime        int     `json:"pickup_time" binding:""`
	UserId            int     `json:"user_id" binding:""`
	UserEmail         *string `json:"user_email" binding:""`
	Format            int     `json:"format" binding:""`
}

type ExternalAPI struct {
	Id                    int      `json:"ContainerNo" binding:""`
	Address               string   `json:"Address" binding:""`
	Latitude              *float64 `json:"Lat" binding:""`
	Longitude             *float64 `json:"Lng" binding:""`
	WasteDistance         int      `json:"fillingDegree" binding:""`
	ContainerHeight       int      `json:"ContainerHeight" binding:""`
	Battery               int      `json:"Battery" binding:""`
	SensorId              string   `json:"volumeUnitId" binding:""`
	LatestMeasurement     string   `json:"LatestMeasurement" binding:""`
	MeasurementFrequency  int      `json:"MeasurementFrequency" binding:""`
	ContainerType         string   `json:"ContainerType" binding:""`
	ContainerFractionType string   `json:"ContainerFractionType" binding:""`
	DisposalCompany       *string  `json:"DisposalCompany" binding:""`
	PickupTime            *int     `json:"PickupTime" binding:""`
	PickupWhenDistance    *string  `json:"Pickup"`
}

/********** Void Functions **********/

/**** Container ****/
func (c *Container) Get() {
	db.Conn.QueryRow("SELECT s.sensor_id, s.measurement_frequency, c.note, c.group_id, c.container_type_id, ct.container_type, c.address, c.size, c.container_height, c.fraction_type_id, f.fraction_type, c.lat, c.lng, c.disposal_company_id, d.company_name, c.size_type, (SELECT user_id FROM container_user WHERE container_id = c.id AND type = 1 ORDER BY id DESC LIMIT 1), (SELECT user_id FROM container_user WHERE container_id = c.id AND type = 3 ORDER BY id DESC LIMIT 1) FROM container c LEFT JOIN container_type ct ON c.container_type_id = ct.id LEFT JOIN container_fraction_type f ON c.fraction_type_id = f.id LEFT JOIN disposal_company d ON c.disposal_company_id = d.id LEFT JOIN sensor s ON c.sensor_id = s.id WHERE c.id = ?", c.Id).Scan(&c.SensorId, &c.MeasurementFrequency, &c.Note, &c.ContainerGroupId, &c.ContainerTypeId, &c.ContainerType, &c.Address, &c.Size, &c.Height, &c.FractionId, &c.Fraction, &c.Latitude, &c.Longitude, &c.DisposalId, &c.DisposalCompany, &c.SizeType, &c.Operator.Id, &c.Viewer.Id)
}

func (c *Container) GetPickupHistory() {
	pRows, _ := db.Conn.Query("SELECT timestamp FROM pickup_history WHERE container_id = '" + strconv.Itoa(c.Id) + "'")
	defer pRows.Close()
	for pRows.Next() {
		var log LogData
		pRows.Scan(&log.Timestamp)
		c.PickupHistory = append(c.PickupHistory, log)
	}

}

func (c *Container) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE container SET lat = ?, lng = ?, container_type_id = ?, address = ?, size = ?, size_type = ?, container_height = ?, fraction_type_id = ?, disposal_company_id = ?, group_id = ?, note = ? WHERE id = ?")
	stmt.Exec(&c.Latitude, &c.Longitude, &c.ContainerTypeId, &c.Address, &c.Size, &c.SizeType, &c.Height, &c.FractionId, &c.DisposalId, &c.ContainerGroupId, &c.Note, &c.Id)

	stmt, _ = db.Conn.Prepare("UPDATE sensor SET measurement_frequency = ? WHERE id = ?")
	stmt.Exec(&c.MeasurementFrequency, &c.Id)

	if c.Operator.Id != 0 {
		stmt, _ = db.Conn.Prepare("DELETE FROM container_user WHERE container_id = ? AND type = 1")
		stmt.Exec(&c.Id)
		stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 1)")
		stmt.Exec(&c.Operator.Id, &c.Id)
	}

	if c.Viewer.Id != 0 {
		stmt, _ = db.Conn.Prepare("DELETE FROM container_user WHERE container_id = ? AND type = 3")
		stmt.Exec(&c.Id)
		stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 3)")
		stmt.Exec(&c.Viewer.Id, &c.Id)
	}

	defer stmt.Close()
}

func (c *Container) UpdateLatLng() {
	stmt, _ := db.Conn.Prepare("UPDATE container SET lat = ?, lng = ? WHERE id = ?")
	stmt.Exec(&c.Latitude, &c.Longitude, &c.Id)

	defer stmt.Close()
}

func (c *Container) Delete() {
	stmt, _ := db.Conn.Prepare("UPDATE container SET deleted = 1 WHERE id = ?")
	stmt.Exec(&c.Id)
	defer stmt.Close()
}

func (c *Container) AddPickupLog() {
	stmt, _ := db.Conn.Prepare("INSERT INTO pickup_history (container_id) VALUES (?)")
	stmt.Exec(&c.Id)

	defer stmt.Close()
}

func (c *Container) GetContainerUser() User {
	var u User
	db.Conn.QueryRow("u.email, u.first_name, u.last_name FROM container_user cu JOIN user u ON u.id = cu.user_id WHERE cu.container_id = ?", c.Id).Scan(&u.Email, &u.FirstName, &u.LastName)

	return u
}

/**** Container fraction ****/
func (cf *ContainerFraction) Create(u *User) {
	stmt, _ := db.Conn.Prepare("INSERT INTO container_fraction_type (fraction_type, color) VALUES (?, ?)")
	res, _ := stmt.Exec(&cf.Type, &cf.Color)
	id, _ := res.LastInsertId()

	stmt, _ = db.Conn.Prepare("INSERT INTO container_fraction_type_user (user_id, container_fraction_type_id) VALUES (?, ?)")
	res, _ = stmt.Exec(&u.Id, &id)

	defer stmt.Close()
}

func (cf *ContainerFraction) Get() {
	db.Conn.QueryRow("SELECT fraction_type, color FROM container_fraction_type WHERE id = ?", cf.Id).Scan(&cf.Type, &cf.Color)
}

func (cf *ContainerFraction) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE container_fraction_type SET fraction_type = ?, color = ? WHERE id = ?")
	stmt.Exec(&cf.Type, &cf.Color, &cf.Id)

	defer stmt.Close()
}

func (cf *ContainerFraction) Delete() {
	stmt, _ := db.Conn.Prepare("UPDATE container_fraction_type SET deleted = 1 WHERE id = ?")
	stmt.Exec(&cf.Id)

	defer stmt.Close()
}

/**** Container types ****/
func (ct *ContainerType) Create(u *User) {
	stmt, _ := db.Conn.Prepare("INSERT INTO container_type (container_type) VALUES (?)")
	res, _ := stmt.Exec(&ct.Type)
	id, _ := res.LastInsertId()

	stmt, _ = db.Conn.Prepare("INSERT INTO container_type_user (user_id, container_type_id) VALUES (?, ?)")
	res, _ = stmt.Exec(&u.Id, &id)

	defer stmt.Close()
}

func (ct *ContainerType) Delete() {
	stmt, _ := db.Conn.Prepare("UPDATE container_type SET deleted = 1 WHERE id = ?")
	stmt.Exec(&ct.Id)

	defer stmt.Close()
}

/**** Container Group ****/
func (cg *ContainerGroup) Create(u *User) {
	stmt, _ := db.Conn.Prepare("INSERT INTO container_group (name, user_id, container_type_id, disposal_company_id, viewer_id, max_pickup_freq, fraction_id, address, full_containers, waste_distance) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	stmt.Exec(&cg.Name, &u.Id, &cg.ContainerTypeId, &cg.DisposalCompanyId, &cg.ViewerId, &cg.MaxPickupFreq, &cg.FractionId, &cg.Address, &cg.FullContainers, &cg.WasteDistance)

	defer stmt.Close()
}

func (cg *ContainerGroup) Get() {
	db.Conn.QueryRow("SELECT name, container_type_id, disposal_company_id, viewer_id, max_pickup_freq, fraction_id, address, full_containers, waste_distance FROM container_group WHERE id = ?", cg.Id).Scan(&cg.Name, &cg.ContainerTypeId, &cg.DisposalCompanyId, &cg.ViewerId, &cg.MaxPickupFreq, &cg.FractionId, &cg.Address, &cg.FullContainers, &cg.WasteDistance)
}

func (cg *ContainerGroup) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE container_group SET name = ?, container_type_id = ?, disposal_company_id = ?, viewer_id = ?, max_pickup_freq = ?, fraction_id = ?, address = ?, full_containers = ?, waste_distance = ? WHERE id = ?")
	stmt.Exec(&cg.Name, &cg.ContainerTypeId, &cg.DisposalCompanyId, &cg.ViewerId, &cg.MaxPickupFreq, &cg.FractionId, &cg.Address, &cg.FullContainers, &cg.WasteDistance, &cg.Id)

	defer stmt.Close()
}

func (cg *ContainerGroup) Delete() {
	stmt, _ := db.Conn.Prepare("UPDATE container_group SET deleted = 1 WHERE id = ?")
	stmt.Exec(&cg.Id)

	defer stmt.Close()
}

func (cl *ContainerLog) GetLogs() {
	var logDataList []LogData

	// New NbIoT measurements
	rows, _ := db.Conn.Query("SELECT p.D1, p.D2, p.D3, p.D4, p.timestamp, c.container_height FROM pga460_measurements p JOIN container c ON c.id = p.container WHERE p.container = ? AND p.timestamp >= ? AND p.timestamp <= ? ORDER BY p.id DESC", &cl.ContainerId, &cl.StartDate, &cl.EndDate)
	defer rows.Close()
	if cl.Max100 {
		rows, _ = db.Conn.Query("SELECT p.D1, p.D2, p.D3, p.D4, p.timestamp, c.container_height FROM pga460_measurements p JOIN container c ON c.id = p.container WHERE p.container = ? ORDER BY p.id DESC LIMIT 100", &cl.ContainerId)
		defer rows.Close()
	}
	for rows.Next() {
		var logData LogData
		rows.Scan(&logData.D1, &logData.D2, &logData.D3, &logData.D4, &logData.Timestamp, &logData.ContainerHeight)
		logData.WasteDistance = CalculateWasteDistance(CalculateWasteDistanceFourDimensional(logData.D1, logData.D2, logData.D3, logData.D4, logData.ContainerHeight), logData.ContainerHeight, 0, "")
		logDataList = append(logDataList, logData)
	}

	rows, _ = db.Conn.Query("SELECT cal.id, cal.dev_eui, cal.decrypted_plaintext, c.container_height, cal.timestamp FROM container_api_log cal JOIN sensor s ON s.deveui = cal.dev_eui JOIN container c ON c.sensor_id = s.id AND c.id = ? WHERE cal.decrypted_plaintext != '' AND cal.dev_eui != '' AND cal.timestamp >= ? AND cal.timestamp <= ? ORDER BY cal.id DESC", &cl.ContainerId, &cl.StartDate, &cl.EndDate)
	defer rows.Close()
	if cl.Max100 {
		rows, _ = db.Conn.Query("SELECT cal.id, cal.dev_eui, cal.decrypted_plaintext, c.container_height, cal.timestamp FROM container_api_log cal JOIN sensor s ON s.deveui = cal.dev_eui JOIN container c ON c.sensor_id = s.id AND c.id = ? WHERE cal.decrypted_plaintext != '' AND cal.dev_eui != '' ORDER BY cal.id DESC LIMIT 100", &cl.ContainerId)
		defer rows.Close()
	}
	for rows.Next() {
		var logData LogData
		rows.Scan(&logData.Id, &logData.DevEUI, &logData.DecryptedPlainText, &logData.ContainerHeight, &logData.Timestamp)
		var wasteDistance float64 = 0.0
		if len(logData.DecryptedPlainText) < 6 {
			wasteDistance, _ = strconv.ParseFloat(logData.DecryptedPlainText, 64)
		} else if strings.Contains(logData.DecryptedPlainText, "{") {
			var decryptedObject DecryptedObject
			json.Unmarshal([]byte(logData.DecryptedPlainText), &decryptedObject)
			wasteDistance = float64(decryptedObject.WasteDistance)
		}
		logData.WasteDistance = CalculateWasteDistance(wasteDistance, logData.ContainerHeight, logData.Id, logData.DevEUI)
		logDataList = append(logDataList, logData)
	}

	cl.Logs = &logDataList
}

func (p *Pickup) Create() {
	stmt, _ := db.Conn.Prepare("INSERT INTO pickup (search, container_id, fraction_type_id, disposal_company_id, waste_distance, container_type_id, pickuptime, user_id, viewer_id, format) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	stmt.Exec(&p.Search, &p.ContainerId, &p.FractionTypeId, &p.DisposalCompanyId, &p.WasteDistance, &p.ContainerTypeId, &p.PickupTime, &p.UserId, &p.ViewerId, &p.Format)
	defer stmt.Close()
}

func (p *Pickup) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE pickup SET search = ?, container_id = ?, fraction_type_id = ?, disposal_company_id = ?, waste_distance = ?, container_type_id = ?, pickuptime = ?, viewer_id = ?, format = ? WHERE id = ?")
	stmt.Exec(&p.Search, &p.ContainerId, &p.FractionTypeId, &p.DisposalCompanyId, &p.WasteDistance, &p.ContainerTypeId, &p.PickupTime, &p.ViewerId, &p.Format, &p.Id)

	defer stmt.Close()
}

func (p *Pickup) Delete() {
	stmt, _ := db.Conn.Prepare("DELETE FROM pickup WHERE id = ?")
	stmt.Exec(&p.Id)

	defer stmt.Close()
}

func (cpn *ContainerPickupNote) Create(u *User) {
	stmt, _ := db.Conn.Prepare("INSERT INTO pickup_note (container_id, weight, note, user_id) VALUES (?, ?, ?, ?)")
	stmt.Exec(&cpn.ContainerId, &cpn.Weight, &cpn.Note, &u.Id)

	defer stmt.Close()
}

/********** Return type Functions **********/

func GetContainersForAdmin() []Container {
	var containers []Container
	rows, _ := db.Conn.Query("SELECT c.id, c.address, c.lat, c.lng, c.waste_distance, f.fraction_type FROM container c LEFT JOIN container_fraction_type f ON c.fraction_type_id = f.id WHERE c.deleted != 1")
	for rows.Next() {
		var container Container
		rows.Scan(&container.Id, &container.Address, &container.Latitude, &container.Longitude, &container.WasteDistance, &container.Fraction)
		containers = append(containers, container)
	}

	defer rows.Close()

	return containers
}

func GetContainersForUser(userID int) []Container {
	var containers []Container
	rows, _ := db.Conn.Query(`SELECT
		c.id,
		c.address,
		c.lat,
		c.lng,
		c.waste_distance,
		f.fraction_type
	FROM
		container c
		LEFT JOIN container_fraction_type f ON c.fraction_type_id = f.id
		LEFT JOIN container_user cu ON cu.container_id = c.id
	WHERE
		cu.user_id = ? AND
		c.deleted != 1
	ORDER BY c.id`, userID)
	for rows.Next() {
		var container Container
		rows.Scan(&container.Id, &container.Address, &container.Latitude, &container.Longitude, &container.WasteDistance, &container.Fraction)
		containers = append(containers, container)
	}

	defer rows.Close()

	return containers
}

func GetContainersForFrontpage(userType int, userID int) []Container {
	var containers []Container
	rows, _ := db.Conn.Query(`SELECT c.id, s.type, s.battery_level, ct.container_type, c.address, f.fraction_type, c.lat, c.lng, c.waste_distance, c.container_height, d.company_name
		FROM container c
		LEFT JOIN container_type ct ON c.container_type_id = ct.id
		LEFT JOIN container_fraction_type f ON c.fraction_type_id = f.id
		LEFT JOIN disposal_company d ON c.disposal_company_id = d.id
		LEFT JOIN sensor s ON c.sensor_id = s.id
		WHERE c.deleted != 1`)
	defer rows.Close()
	if userType != 1 {
		rows, _ = db.Conn.Query(`SELECT c.id, s.type, s.battery_level, ct.container_type, c.address, f.fraction_type, c.lat, c.lng, c.waste_distance, c.container_height, d.company_name
		FROM container c
		LEFT JOIN container_type ct ON c.container_type_id = ct.id
		LEFT JOIN container_fraction_type f ON c.fraction_type_id = f.id
		LEFT JOIN disposal_company d ON c.disposal_company_id = d.id
		LEFT JOIN sensor s ON c.sensor_id = s.id
		LEFT JOIN container_user cu ON cu.container_id = c.id
		WHERE cu.user_id = ? AND c.deleted != 1`, userID)
		defer rows.Close()
	}
	for rows.Next() {
		var container Container
		rows.Scan(&container.Id, &container.Sensor.Type, &container.Sensor.Battery, &container.ContainerType, &container.Address, &container.Fraction, &container.Latitude, &container.Longitude, &container.WasteDistance, &container.Height, &container.DisposalCompany)
		container.WasteDistance = float64(CalculateWasteDistance(container.WasteDistance, container.Height, 0, ""))
		container.GetPickupHistory()
		containers = append(containers, container)
	}

	return containers
}

func ReturnContainerTypes() []ContainerType {
	var containerTypes []ContainerType

	rows, _ := db.Conn.Query("SELECT id, container_type, container_size_liter, container_sizem3 FROM container_type WHERE deleted != 1")
	defer rows.Close()
	for rows.Next() {
		var containerType ContainerType
		rows.Scan(&containerType.Id, &containerType.Type, &containerType.SizeLiter, &containerType.SizeM3)
		containerTypes = append(containerTypes, containerType)
	}

	return containerTypes
}

func ReturnContainerFractions() []ContainerFraction {
	var containerFractions []ContainerFraction

	rows, _ := db.Conn.Query("SELECT id, fraction_type, color, size, weight_calibration_factor FROM container_fraction_type WHERE deleted != 1")
	defer rows.Close()
	for rows.Next() {
		var containerFraction ContainerFraction
		rows.Scan(&containerFraction.Id, &containerFraction.Type, &containerFraction.Color, &containerFraction.Size, &containerFraction.WCF)
		containerFractions = append(containerFractions, containerFraction)
	}

	return containerFractions
}

func ReturnContainerGroups(userId int) []ContainerGroup {
	var containerGroups []ContainerGroup

	rows, _ := db.Conn.Query("SELECT id, name FROM container_group WHERE user_id = ? AND deleted != 1", userId)
	defer rows.Close()
	for rows.Next() {
		var containerGroup ContainerGroup
		rows.Scan(&containerGroup.Id, &containerGroup.Name)
		containerGroups = append(containerGroups, containerGroup)
	}

	return containerGroups
}

func CalculateWasteDistance(distance float64, height int, logId int, deveui string) int {
	if height == 0 {
		return 0
	}
	calcFilled := math.Round((1 - distance/float64(height)) * 100)

	if logId == 1 {
		rows, _ := db.Conn.Query("SELECT id FROM container_api_log WHERE dev_eui = ? ORDER BY id DESC LIMIT 1", deveui)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&logId)
		}
	}

	if logId != 0 && logId != 1 {
		rows, _ := db.Conn.Query("SELECT decrypted_plaintext FROM container_api_log WHERE dev_eui = ? AND id < ? ORDER BY id DESC LIMIT 1", &deveui, logId)
		defer rows.Close()
		for rows.Next() {
			var logData LogData
			rows.Scan(&logData.DecryptedPlainText)
			var wasteDistance float64 = 0.0

			if len(logData.DecryptedPlainText) < 6 {
				wasteDistance, _ = strconv.ParseFloat(logData.DecryptedPlainText, 64)
				if wasteDistance*2 < distance {
					calcFilled = calcFilled / 2
				}

			} else if strings.Contains(logData.DecryptedPlainText, "{") {
				var decryptedObject DecryptedObject
				json.Unmarshal([]byte(logData.DecryptedPlainText), &decryptedObject)
				wasteDistance = float64(decryptedObject.WasteDistance)
				if wasteDistance*2 < distance {
					calcFilled = calcFilled / 2
				}
			}
		}
	}

	if distance <= 12.0 && logId != 0 && logId != 1 {
		rows, _ := db.Conn.Query("SELECT id, decrypted_plaintext FROM container_api_log WHERE dev_eui = ? AND id < ? ORDER BY id DESC LIMIT 3", &deveui, logId)
		defer rows.Close()
		for rows.Next() {
			var logData LogData
			rows.Scan(&logData.Id, &logData.DecryptedPlainText)
			var wasteDistance float64 = 0.0

			if len(logData.DecryptedPlainText) < 6 {
				wasteDistance, _ = strconv.ParseFloat(logData.DecryptedPlainText, 64)
				if wasteDistance > 12.0 {
					distance = wasteDistance
					break
				}

			} else if strings.Contains(logData.DecryptedPlainText, "{") {
				var decryptedObject DecryptedObject
				json.Unmarshal([]byte(logData.DecryptedPlainText), &decryptedObject)
				wasteDistance = float64(decryptedObject.WasteDistance)
				if wasteDistance > 12.0 {
					distance = wasteDistance
					break
				}
			}
		}
		if distance <= 12.0 {
			return 100
		} else {
			calcFilled = math.Round((1 - distance/float64(height)) * 100)
		}
	}

	if int(distance) > height {
		calcFilled = 0
	} else if calcFilled >= 100 || distance > float64(height) {
		calcFilled = 100
	} else if distance == 0 {
		calcFilled = 100
	} else if calcFilled < 0 {
		calcFilled = 0
	}

	if distance == 1123 && height <= 100 {
		calcFilled = 100
	} else if distance == 1123 && height > 100 {
		calcFilled = 0
	}

	return int(calcFilled)
}

func CalculateWasteDistanceFourDimensional(d1 float64, d2 float64, d3 float64, d4 float64, height int) float64 {
	var measurements []float64
	if int(d1) > height {
		d1 = 100
	} else if d1 > 0 && d1 <= 12 {
		d1 = 0
	} else if d1 > 12 && d1 <= 400 {
		measurements = append(measurements, d1)
	} else {
		d1 = 100
	}

	if int(d2) > height {
		d2 = 100
	} else if d2 > 0 && d2 <= 12 {
		d2 = 0
	} else if d2 > 12 && d2 <= 400 {
		measurements = append(measurements, d2)
	} else {
		d2 = 100
	}

	if int(d3) > height {
		d3 = 100
	} else if d3 > 0 && d3 <= 12 {
		d3 = 0
	} else if d3 > 12 && d3 <= 400 {
		measurements = append(measurements, d3)
	} else {
		d3 = 100
	}

	if int(d4) > height {
		d4 = 100
	} else if d4 > 0 && d4 <= 12 {
		d4 = 0
	} else if d4 > 12 && d4 <= 400 {
		measurements = append(measurements, d4)
	} else {
		d4 = 100
	}

	if d1 == 100 && d2 == 100 && d3 == 100 && d4 == 100 {
		return 0.0
	}

	if len(measurements) == 4 {
		return (measurements[0] * 0.80) + (measurements[1] * 0.10) + (measurements[2] * 0.05) + (measurements[3] * 0.05)
	} else if len(measurements) == 3 {
		return (measurements[0] * 0.85) + (measurements[1] * 0.10) + (measurements[2] * 0.05)
	} else if len(measurements) == 2 {
		return (measurements[0] * 0.80) + (measurements[1] * 0.2)
	} else if len(measurements) == 1 {
		return measurements[0]
	}

	return d1
}

func GetUsersContainerPickups(u User) []Pickup {
	var pickups []Pickup
	rows, _ := db.Conn.Query("SELECT p.*, (SELECT fraction_type FROM container_fraction_type WHERE id = p.fraction_type_id), (SELECT company_name FROM disposal_company WHERE id = p.disposal_company_id), (SELECT container_type FROM container_type WHERE id = p.container_type_id) FROM pickup p WHERE p.user_id = ?", u.Id)
	defer rows.Close()

	for rows.Next() {
		var pickup Pickup
		rows.Scan(&pickup.Id, &pickup.Search, &pickup.ContainerId, &pickup.FractionTypeId, &pickup.DisposalCompanyId, &pickup.WasteDistance, &pickup.ContainerTypeId, &pickup.PickupTime, &pickup.UserId, &pickup.ViewerId, &pickup.Format, &pickup.FractionType, &pickup.DisposalCompany, &pickup.ContainerType)
		pickups = append(pickups, pickup)
	}

	return pickups
}

func GetContainersForExternalEndpoint(query string) []ExternalAPI {
	var containers []ExternalAPI
	rows, _ := db.Conn.Query(query)
	defer rows.Close()

	for rows.Next() {
		var container ExternalAPI
		rows.Scan(&container.Id, &container.Address, &container.Latitude, &container.Longitude, &container.WasteDistance, &container.ContainerHeight, &container.Battery, &container.SensorId, &container.LatestMeasurement, &container.MeasurementFrequency, &container.ContainerType, &container.ContainerFractionType, &container.DisposalCompany, &container.PickupTime, &container.PickupWhenDistance)
		container.WasteDistance = CalculateWasteDistance(float64(container.WasteDistance), container.ContainerHeight, 0, "")
		containers = append(containers, container)
	}

	return containers
}
