package models

import (
	"encoding/hex"
	"encoding/json"
	"strings"
	"wastecontrol/db"
	"strconv"
	"net/http"
	//"fmt"
)
type Sensor struct {
	Id                   int     `json:"id" binding:""`
	SensorId             string  `json:"sensor_id" binding:""`
	MeasurementFrequency int     `json:"measurement_frequency" binding:""`
	Imsi                 *string `json:"imsi" binding:""`
	ReplaceSensorId      int     `json:"replace_sensor_id" binding""`
	LatestMeasurement    *string `json:"latest_measurement" binding:""`
	DeploymentDate       *string `json:"deployment_date" binding:""`
	Deveui               *string `json:"deveui" binding:""`
	Battery              *int    `json:"battery" binding:""`
	Type                 *string `json:"type" binding:""`
	Ip                   *string `json:"ip" binding:""`
	AppEui               *string `json:"appeui" binding:""`
	AppKey               *string `json:"appkey" binding:""`
	Firmware             *string `json:"firmware_version" binding:""`
	Pga460Params         *string `json:"pga460_params" binding:""`
	Thresholds           *string `json:"thresholds" binding:""`
	Status               string  `json:"status" binding:""`
	ResellerId           *int    `json:"reseller_id" binding:""`
	OperatorId           *int    `json:"operator_id" binding:""`
	ViewerId           	 *int    `json:"viewer_id" binding:""`
}

type TTN struct {
	HardwareSerial string          `json:"hardware_serial" binding:""`
	PayloadRaw     string          `json:"payload_raw" binding:""`
	PayloadFields  DecryptedObject `json:"payload_fields" binding:""`
	XML            string          `json:"raw"`
	IP             string          `json:"ip" binding:""`
}

type Teracom struct {
	Port          int             `json:"port" binding:""`
	EUI           string          `json:"EUI" binding:""`
	Timestamp     int             `json:"ts" binding:""`
	Data          string          `json:"data" binding:""`
	PayloadFields DecryptedObject `json:"payload_fields" binding:""`
	XML           string          `json:"raw"`
	IP            string          `json:"ip" binding:""`
	ContainerId   int             `json:"container_id" binding:""`
	MeasurementFrequency string   `json:"measurement_frequency" binding:""`
}

type Nbiot struct {
	EUI                  string `json:"EUI" binding:""`
	AppKey               string `json:"app_key" binding:""`
	WasteDistance        string `json:"waste_distance" binding:""`
	MeasurementFrequency string `json:"measurement_frequency" binding:""`
	Error                string `json:"waste_distance" binding:""`
	Payload              []byte `json:"payload" binding:""`
	DecryptedPayload     []byte `json:"decrypted_payload" binding:""`
	IP                   string `json:"ip" binding:""`
	ContainerId          int    `json:"container_id" binding:""`
	State                int    `json:"state" binding:""`
	DefaultPreset        int    `json:"default_preset" binding:""`
	Timestamp            int    `json:"timestamp" binding:""`
}

type MassSensor struct {
	Sensors              []Sensor `json:"sensors" binding:""`
	ResellerId           *int     `json:"reseller_id" binding:""`
	OperatorId           *int     `json:"operator_id" binding:""`
	ViewerId 	           *int     `json:"viewer_id" binding:""`
	MeasurementFrequency *int     `json:"measurement_frequency" binding:""`
	Firmware             *string  `json:"firmware_version" binding:""`
	DisposalCompany      *int		  `json:"disposal_company_id" binding:""`
	Address              *string  `json:"address" binding:""`
	Lat		               *string  `json:"lat" binding:""`
	Lng		               *string  `json:"lng" binding:""`
	Height               *int     `json:"height" binding:""`
	ContainerType	       *int     `json:"container_type_id" binding:""`
	Size                 *int     `json:"size" binding:""`
	SizeType             *string  `json:"size_type" binding:""`
	Fraction	           *int     `json:"fraction_id" binding:""`
	Note		             *string  `json:"note" binding:""`
}

type SensorLog struct {
	Timestamp   string `json:"timestamp" binding:""`
	Deveui      string `json:"deveui" binding""`
	Payload     string `json:"payload" binding:""`
	IP          string `json:"ip" binding""`
	ContainerId int    `json:"container_id" binding:""`
}

type Trace struct {
	ContainerID  int     `json:"container_id" binding:""`
	Container    int     `json:"container" binding:""`
	TS           string  `json:"TS" binding:""`
	Timestamp    int     `json:"timestamp" binding""`
	Trace        string  `json:"trace" binding:""`
	Pga460Params string  `json:"pga460_params" binding""`
	Thresholds   string  `json:"thresholds" binding:""`
	Preset       int     `json:"preset" binding:""`
	D1           float32 `json:"d1" binding:""`
	D2           float32 `json:"d2" binding:""`
	D3           float32 `json:"d3" binding:""`
	D4           float32 `json:"d4" binding:""`
	W1           int     `json:"w1" binding:""`
	W2           int     `json:"w2" binding:""`
	W3           int     `json:"w3" binding:""`
	W4           int     `json:"w4" binding:""`
	A1           int     `json:"a1" binding:""`
	A2           int     `json:"a2" binding:""`
	A3           int     `json:"a3" binding:""`
	A4           int     `json:"a4" binding:""`
}

/********** Void Functions **********/

func (s *Sensor) Get() {
	db.Conn.QueryRow("SELECT s.id, s.latest_measurement, s.measurement_frequency, s.battery_level, s.firmware_version, s.appkey, (SELECT user_id FROM container_user WHERE container_id = s.id AND type = 1), (SELECT user_id FROM container_user WHERE container_id = s.id AND type = 2) FROM sensor s WHERE s.sensor_id = ?", s.SensorId).Scan(&s.Id, &s.LatestMeasurement, &s.MeasurementFrequency, &s.Battery, &s.Firmware, &s.AppKey, &s.OperatorId, &s.ResellerId)
}

func (s *Sensor) Create() {
	if s.ReplaceSensorId == 0 {
		stmt, _ := db.Conn.Prepare("INSERT INTO sensor (sensor_id, measurement_frequency, imsi) VALUES (?, ?, ?)")
		res, _ := stmt.Exec(&s.SensorId, &s.MeasurementFrequency, &s.Imsi)
		id, _ := res.LastInsertId()

		stmt, _ = db.Conn.Prepare("INSERT INTO container (sensor_id, container_height) VALUES (?, 100)")
		stmt.Exec(&id)

		if s.ResellerId != nil {
			stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 2)")
			stmt.Exec(&s.ResellerId, &id)
		}

		if s.OperatorId != nil {
			stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 1)")
			stmt.Exec(&s.OperatorId, &id)
		}

		defer stmt.Close()
	} else {
		stmt, _ := db.Conn.Prepare("UPDATE sensor SET sensor_id = ?, measurement_frequency = ?, imsi = ? WHERE id = ?")
		stmt.Exec(&s.SensorId, &s.MeasurementFrequency, &s.Imsi, &s.ReplaceSensorId)

		defer stmt.Close()
	}
}

func (s *Sensor) GetParamsThresholds() {
	db.Conn.QueryRow("SELECT pga460_params, thresholds FROM sensor WHERE id = ?", s.Id).Scan(&s.Pga460Params, &s.Thresholds)
}

func (s *Sensor) CreateNbiot(payload string) {
	checkSensor, _ := db.Conn.Query("SELECT id FROM sensor WHERE sensor_id = ?", s.Imsi)

	defer checkSensor.Close()

	if CheckCount(checkSensor) > 0 {
		stmt, _ := db.Conn.Prepare("UPDATE sensor SET measurement_frequency = ?, appkey = ?, firmware_version = ?, type = ?, imsi = ? WHERE sensor_id = ?")
		stmt.Exec(&s.MeasurementFrequency, &s.AppKey, &s.Firmware, &s.Type, &s.Imsi, &s.Imsi)
		s.Status = "Updated"
		
		db.Conn.QueryRow("SELECT id FROM sensor WHERE sensor_id = ?", s.Imsi).Scan(&s.Id)

		stmt, _ = db.Conn.Prepare("INSERT INTO container_log (id, payload, ip, httpstatus, actiondone) VALUES (?, ?, ?, ?, 'updated')")
		stmt.Exec(&s.Imsi, &payload, &s.Ip, 200)
	} else {
		stmt, _ := db.Conn.Prepare("INSERT INTO sensor (sensor_id, measurement_frequency, appkey, firmware_version, type, imsi) VALUES (?, ?, ?, ?, ?, ?)")
		res, _ := stmt.Exec(&s.Imsi, &s.MeasurementFrequency, &s.AppKey, &s.Firmware, &s.Type, &s.Imsi)
		id, _ := res.LastInsertId()
		s.Id = int(id)
		s.Status = "Created"

		stmt, _ = db.Conn.Prepare("INSERT INTO container (sensor_id) VALUES (?)")
		stmt.Exec(&id)

		stmt, _ = db.Conn.Prepare("INSERT INTO container_log (id, payload, ip, httpstatus, actiondone) VALUES (?, ?, ?, ?, 'created')")
		stmt.Exec(&id, &payload, &s.Ip, 200)

		defer stmt.Close()
	}
}

func (s *Sensor) CreateLorawan(payload string) {
	checkSensor, _ := db.Conn.Query("SELECT id FROM sensor WHERE sensor_id = ?", s.Deveui)
	
	defer checkSensor.Close()

	if CheckCount(checkSensor) > 0 {
		stmt, _ := db.Conn.Prepare("UPDATE sensor SET measurement_frequency = ?, deveui = ?, appeui = ?, appkey = ?, firmware_version = ?, type = ? WHERE sensor_id = ?")
		stmt.Exec(&s.MeasurementFrequency, &s.Deveui, &s.AppEui, &s.AppKey, &s.Firmware, &s.Type, &s.Deveui)
		s.Status = "Updated"

		db.Conn.QueryRow("SELECT id FROM sensor WHERE sensor_id = ?", s.Deveui).Scan(&s.Id)

		stmt, _ = db.Conn.Prepare("INSERT INTO container_log (id, payload, ip, httpstatus, actiondone) VALUES (?, ?, ?, ?, 'updated')")
		stmt.Exec(&s.Deveui, &payload, &s.Ip, 200)
	} else {
		stmt, _ := db.Conn.Prepare("INSERT INTO sensor (sensor_id, measurement_frequency, deveui, appeui, appkey, firmware_version, type) VALUES (?, ?, ?, ?, ?, ?, ?)")
		res, _ := stmt.Exec(&s.Deveui, &s.MeasurementFrequency, &s.Deveui, &s.AppEui, &s.AppKey, &s.Firmware, &s.Type)
		id, _ := res.LastInsertId()
		s.Id = int(id)
		s.Status = "Created"

		stmt, _ = db.Conn.Prepare("INSERT INTO container (sensor_id, container_height) VALUES (?, 100)")
		stmt.Exec(&id)

		stmt, _ = db.Conn.Prepare("INSERT INTO container_log (id, payload, ip, httpstatus, actiondone) VALUES (?, ?, ?, ?, 'created')")
		stmt.Exec(&id, &payload, &s.Ip, 200)

		defer stmt.Close()
	}
}

func (s *Sensor) CreateFromCSV() {
	stmt, _ := db.Conn.Prepare("INSERT INTO sensor (sensor_id) VALUES (?)")
	res, _ := stmt.Exec(&s.SensorId)
	id, _ := res.LastInsertId()

	stmt, _ = db.Conn.Prepare("INSERT INTO container (sensor_id) VALUES (?)")
	stmt.Exec(&id)

	if s.ResellerId != nil {
		stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 1)")
		stmt.Exec(&s.ResellerId, &id)
	}

	if s.OperatorId != nil {
		stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 2)")
		stmt.Exec(&s.OperatorId, &id)
	}
}

func (s *Sensor) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE sensor SET sensor_id = ?, measurement_frequency = ?, firmware_version = ? WHERE id = ?")
	stmt.Exec(&s.SensorId, &s.MeasurementFrequency, &s.Firmware, &s.Id)
	defer stmt.Close()
	if s.ResellerId != nil {
		stmt, _ = db.Conn.Prepare("DELETE FROM container_user WHERE container_id = ? AND type = 2")
		stmt.Exec(&s.Id)
		stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 2)")
		stmt.Exec(&s.ResellerId, &s.Id)
		defer stmt.Close()
	}

}

func (s *Sensor) Delete() {
	stmt, _ := db.Conn.Prepare("UPDATE sensor SET deleted = 1 WHERE id = ?")
	stmt.Exec(&s.Id)

	defer stmt.Close()
}

func (t *Trace) GetTraceData() {
	db.Conn.QueryRow("SELECT timestamp, trace, pga460_params, thresholds, preset, D1, D2, D3, D4, W1, W2, W3, W4, A1, A2, A3, A4 FROM trace_data WHERE container = ? ORDER BY timestamp DESC LIMIT 1", t.ContainerID).Scan(&t.TS, &t.Trace, &t.Pga460Params, &t.Thresholds, &t.Preset, &t.D1, &t.D2, &t.D3, &t.D4, &t.W1, &t.W2, &t.W3, &t.W4, &t.A1, &t.A2, &t.A3, &t.A4)
}

func (t *Trace) Request() {
	stmt, _ := db.Conn.Prepare("UPDATE sensor SET state = ? WHERE id = ?")
	stmt.Exec(&t.Preset, &t.Container)

	defer stmt.Close()
}

func (t *Trace) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE sensor SET pga460_params = ?, thresholds = ?, state = ? WHERE id = ?")
	stmt.Exec(&t.Pga460Params, &t.Thresholds, &t.Preset, &t.Container)

	defer stmt.Close()
}

func (t *Trace) Add() {
	flags := 0x7f
	stmt, _ := db.Conn.Prepare("INSERT INTO trace_data (container, trace, PGA460_params, thresholds, preset, D1, W1, A1, D2, W2, A2, D3, W3, A3, D4, W4, A4, flags) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	stmt.Exec(&t.ContainerID, &t.Trace, &t.Pga460Params, &t.Thresholds, &t.Preset, &t.D1, &t.W1, &t.A1, &t.D2, &t.W2, &t.A2, &t.D3, &t.W3, &t.A3, &t.D4, &t.W4, &t.A4, &flags)

	stmt, _ = db.Conn.Prepare("UPDATE sensor SET state = 0 WHERE id = ?")
	stmt.Exec(&t.ContainerID)

	defer stmt.Close()
}

func (t *Trace) UploadFourObjects() {
	var count int
	db.Conn.QueryRow("select count(*) from pga460_measurements where container = ? and node_TS >= ?", t.ContainerID, t.Timestamp).Scan(&count)
	
	if count == 0 {
		stmt, _ := db.Conn.Prepare("INSERT INTO pga460_measurements (container, node_TS, preset, D1, W1, A1, D2, W2, A2, D3, W3, A3, D4, W4, A4) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
		stmt.Exec(&t.ContainerID, &t.Timestamp, &t.Preset, &t.D1, &t.W1, &t.A1, &t.D2, &t.W2, &t.A2, &t.D3, &t.W3, &t.A3, &t.D4, &t.W4, &t.A4)
	
		stmt, _ = db.Conn.Prepare("UPDATE container c, sensor s SET c.waste_distance = ?, c.latest_measurement = NOW(), s.latest_measurement = NOW() WHERE s.id = ? AND c.id = s.id")
		stmt.Exec(&t.D1, &t.ContainerID)
	
		defer stmt.Close()	
	}
}

func (n *Nbiot) GetByContainerId() {
	db.Conn.QueryRow("SELECT appkey, sensor_id, measurement_frequency, state, default_preset FROM sensor WHERE id = ?", n.ContainerId).Scan(&n.AppKey, &n.EUI, &n.MeasurementFrequency, &n.State, &n.DefaultPreset)
}

func (n *Nbiot) UploadEmptyLog() {
	stmt, _ := db.Conn.Prepare("INSERT INTO initial_api_log (dev_eui, xml_post, ip_address) VALUES (?, ?, ?)")
	stmt.Exec(&n.EUI, &n.Payload, &n.IP)

	defer stmt.Close()
}

func (n *Nbiot) Upload() {
	stmt, _ := db.Conn.Prepare("INSERT INTO container_api_log (dev_eui, encrypted_plaintext, decrypted_plaintext, ip_address, file_name) VALUES (?, ?, ?, ?, ?)")
	stmt.Exec(&n.EUI, strings.ToUpper(hex.EncodeToString(n.Payload)), "{\"WD\":"+n.WasteDistance+"}", &n.IP, "nbiot")

	stmt, _ = db.Conn.Prepare("UPDATE container c, sensor s SET c.waste_distance = ?, c.latest_measurement = NOW(), s.latest_measurement = NOW() WHERE s.deveui = ? AND c.sensor_id = s.id")
	stmt.Exec(&n.WasteDistance, &n.EUI)

	defer stmt.Close()
}

func (t *TTN) UploadEmptyLog() {
	stmt, _ := db.Conn.Prepare("INSERT INTO initial_api_log (dev_eui, xml_post, ip_address) VALUES (?, ?, ?)")
	stmt.Exec(&t.HardwareSerial, &t.XML, &t.IP)

	defer stmt.Close()
}

func (t *TTN) Upload() {
	decrypt, _ := json.Marshal(t.PayloadFields)
	stmt, _ := db.Conn.Prepare("INSERT INTO container_api_log (dev_eui, encrypted_plaintext, decrypted_plaintext, xml_post, ip_address, file_name) VALUES (?, ?, ?, ?, ?, ?)")
	stmt.Exec(&t.HardwareSerial, &t.PayloadRaw, string(decrypt), &t.XML, &t.IP, "ttn")

	stmt, _ = db.Conn.Prepare("UPDATE container c, sensor s SET c.waste_distance = ?, c.latest_measurement = NOW(), s.battery_level = ?, s.latest_measurement = NOW() WHERE s.deveui = ? AND c.sensor_id = s.id")
	stmt.Exec(&t.PayloadFields.WasteDistance, &t.PayloadFields.Battery, &t.HardwareSerial)

	defer stmt.Close()
}

func (t *Teracom) GetByDevEui() {
	db.Conn.QueryRow("SELECT id, measurement_frequency FROM sensor WHERE deveui = ?", t.EUI).Scan(&t.ContainerId, &t.MeasurementFrequency)
}

func (t *Teracom) UploadEmptyLog() {
	stmt, _ := db.Conn.Prepare("INSERT INTO initial_api_log (dev_eui, xml_post, ip_address) VALUES (?, ?, ?)")
	stmt.Exec(&t.EUI, &t.XML, &t.IP)

	defer stmt.Close()
}

func (t *Teracom) Upload() {
	var payload string
	
	if t.Port == 88 {	//PORT 88 -> JSON both in up- and downlink
		tmp, _ := hex.DecodeString(t.Data)
		payload = string(tmp)
		
		var payloadFields DecryptedObject
		json.Unmarshal(tmp, &payloadFields)
		t.PayloadFields = payloadFields
		
		stmt, _ := db.Conn.Prepare("INSERT INTO container_api_log (dev_eui, encrypted_plaintext, decrypted_plaintext, xml_post, ip_address, file_name) VALUES (?, ?, ?, ?, ?, ?)")
		stmt.Exec(&t.EUI, &t.Data, payload, &t.XML, &t.IP, "teracom")
	
		stmt, _ = db.Conn.Prepare("UPDATE container c, sensor s SET c.waste_distance = ?, c.latest_measurement = NOW(), s.battery_level = ?, s.latest_measurement = NOW() WHERE s.deveui = ? AND c.sensor_id = s.id")
		stmt.Exec(&t.PayloadFields.WasteDistance, &t.PayloadFields.Battery, &t.EUI)
	
		//check wakeup period and generate downlink on mismatch
		t.GetByDevEui()	//read container number from DB
		prd, _ := strconv.ParseInt(t.MeasurementFrequency, 10, 64)
		if int(prd) != *t.PayloadFields.PRD {
			
			var appeui, app_id, api_key, dl_token string
			db.Conn.QueryRow("select appeui from sensor where deveui = ?", t.EUI).Scan(&appeui)
			db.Conn.QueryRow("select APP_ID, API_key, Downlink_token from teracom_account where AppEUI = ?", appeui).Scan(&app_id, &api_key, &dl_token)
			
			type downlink_payload struct {
				PRD	int
			}
			
			pl := downlink_payload {
				PRD:	int(prd),
			}
	
			jsonPayload, _ := json.Marshal(pl)	
			
			type teracom_downlink struct {
				Appid		string	`json:"appid"`
				Cmd			string	`json:"cmd"`
				EUI			string	`json:"EUI"`
				Port		int		`json:"port"`
				Confirmed	bool	`json:"confirmed"`
				Data		string	`json:"data"`
			}
			
			dl := teracom_downlink {
				Appid:		app_id,
				Cmd:		"tx",
				EUI:		t.EUI,
				Port:		88,
				Confirmed:	false,
				Data:		strings.ToUpper(hex.EncodeToString(jsonPayload)),
			}
			var jsonData []byte
			jsonData, _ = json.Marshal(dl)
			
			//create and execute request to TERACOM API
			client := &http.Client{}
			
			flushURL := "https://iotnet.teracom.dk/1/nwk/app/" + app_id + "/device/" + t.EUI + "/flush_queue"
			flushReq, _ := http.NewRequest("POST", flushURL, strings.NewReader(string("")))
			flushReq.Header.Set("Authorization", "Bearer " + api_key)
		    flushReq.Header.Set("Content-Type", "application/json")
		    client.Do(flushReq)
	
			req, _ := http.NewRequest("POST", "https://iotnet.teracom.dk/1/rest", strings.NewReader(string(jsonData)))
		    req.Header.Set("Authorization", dl_token)
		    req.Header.Set("Content-Type", "application/json")
		    client.Do(req)
		}

		defer stmt.Close()
	}
	
	if t.Port == 42 {	//PORT 42 -> binary data (4 objects)
		payload, _ := hex.DecodeString(t.Data)
		
		stmt, _ := db.Conn.Prepare("INSERT INTO container_api_log (dev_eui, encrypted_plaintext, decrypted_plaintext, xml_post, ip_address, file_name) VALUES (?, ?, ?, ?, ?, ?)")
		stmt.Exec(&t.EUI, &t.Data, payload, &t.XML, &t.IP, "teracom")
		
		var sensorData Trace
		db.Conn.QueryRow("select id from sensor where deveui = ?", t.EUI).Scan(&sensorData.ContainerID)
		sensorData.Timestamp =  t.Timestamp / 1000
		
		sensorData.D1 = float32((rune(payload[0]) << 8) + (rune(payload[1]))) / 10
		sensorData.D2 = float32((rune(payload[2]) << 8) + (rune(payload[3]))) / 10
		sensorData.D3 = float32((rune(payload[4]) << 8) + (rune(payload[5]))) / 10
		sensorData.D4 = float32((rune(payload[6]) << 8) + (rune(payload[7]))) / 10
	
		sensorData.W1 = int((rune(payload[8]) << 8) + (rune(payload[9])))
		sensorData.W2 = int((rune(payload[10]) << 8) + (rune(payload[11])))
		sensorData.W3 = int((rune(payload[12]) << 8) + (rune(payload[13])))
		sensorData.W4 = int((rune(payload[14]) << 8) + (rune(payload[15])))
	
		sensorData.A1 = int(rune(payload[16]))
		sensorData.A2 = int(rune(payload[17]))
		sensorData.A3 = int(rune(payload[18]))
		sensorData.A4 = int(rune(payload[19]))
		
		sensorData.Preset = int(rune(payload[20]))
		
		PRD := int((rune(payload[21]) << 8) + (rune(payload[22])))
		var PRD_DB, preset_DB, state int
		db.Conn.QueryRow("select measurement_frequency, default_preset, state from sensor where deveui = ?", t.EUI).Scan(&PRD_DB, &preset_DB, &state)
		
		if (PRD != PRD_DB) || (sensorData.Preset != preset_DB) {
			//reset state to 0
			stmt, _ = db.Conn.Prepare("update sensor set state = 0 where deveui = ?")
			stmt.Exec(&t.EUI)
			
			//generate downlink
			var data string
			var port int
			
			payload := []byte {0x00, 0x00, 0x00}
			payload[0] = byte((PRD_DB >> 8) & 0xff)
			payload[1] = byte(PRD_DB & 0xff)
			payload[2] = byte(preset_DB & 0xff)

			data = strings.ToUpper(hex.EncodeToString(payload))				
			port = 42
			sensorData.Create_downlink(data, port, false, t.EUI)
		}
	
		sensorData.UploadFourObjects()
	}
	if t.Port == 101 || t.Port == 102 || t.Port == 103 || t.Port == 104  {	//PORT 101..104 -> trace data (message #1..4)
		payload, _ := hex.DecodeString(t.Data)
		
		stmt, _ := db.Conn.Prepare("INSERT INTO container_api_log (dev_eui, encrypted_plaintext, decrypted_plaintext, xml_post, ip_address, file_name) VALUES (?, ?, ?, ?, ?, ?)")
		stmt.Exec(&t.EUI, &t.Data, payload, &t.XML, &t.IP, "teracom")
		
		trace_number := (rune(payload[0]) << 24) + (rune(payload[1]) << 16) + (rune(payload[2]) << 8) + (rune(payload[3]))
		
		var trace string
		var flags int
		var new_trace string
		
		db.Conn.QueryRow("select trace, flags from trace_data where trace_number = ?", trace_number).Scan(&trace, &flags)
		
		switch t.Port {
			case 101:
				new_trace = t.Data[8:] + trace[64:]
				flags |= 0x01
			case 102:
				new_trace = trace[:64] + t.Data[8:] + trace[128:]
				flags |= 0x02
			case 103:
				new_trace = trace[:128] + t.Data[8:] + trace[192:]
				flags |= 0x04
			case 104:
				new_trace = trace[:192] + t.Data[8:]
				flags |= 0x08
		}
		stmt, _ = db.Conn.Prepare("update trace_data set trace = ?, flags = ? where trace_number = ?")
		stmt.Exec(&new_trace, &flags, &trace_number)
	}
	if t.Port == 105 || t.Port == 106 {	//PORT 105 -> pga460 params (message #5); PORT 106 -> thresholds (message #6)
		payload, _ := hex.DecodeString(t.Data)
		
		stmt, _ := db.Conn.Prepare("INSERT INTO container_api_log (dev_eui, encrypted_plaintext, decrypted_plaintext, xml_post, ip_address, file_name) VALUES (?, ?, ?, ?, ?, ?)")
		stmt.Exec(&t.EUI, &t.Data, payload, &t.XML, &t.IP, "teracom")
		
		trace_number := (rune(payload[0]) << 24) + (rune(payload[1]) << 16) + (rune(payload[2]) << 8) + (rune(payload[3]))
		
		var flags int
		
		db.Conn.QueryRow("select flags from trace_data where trace_number = ?", trace_number).Scan(&flags)
		
		data := t.Data[8:]
		
		switch t.Port {
			case 105:
				flags |= 0x10
				stmt, _ = db.Conn.Prepare("update trace_data set pga460_params = ?, flags = ? where trace_number = ?")
			case 106:
				flags |= 0x20
				stmt, _ = db.Conn.Prepare("update trace_data set thresholds = ?, flags = ? where trace_number = ?")
				
		}
		stmt.Exec(&data, &flags, &trace_number)
	}
	if t.Port == 107 {	//PORT 107 -> distance matrix (message #7)
		payload, _ := hex.DecodeString(t.Data)
		
		stmt, _ := db.Conn.Prepare("INSERT INTO container_api_log (dev_eui, encrypted_plaintext, decrypted_plaintext, xml_post, ip_address, file_name) VALUES (?, ?, ?, ?, ?, ?)")
		stmt.Exec(&t.EUI, &t.Data, payload, &t.XML, &t.IP, "teracom")
		
		trace_number := (rune(payload[0]) << 24) + (rune(payload[1]) << 16) + (rune(payload[2]) << 8) + (rune(payload[3]))
		
		var flags int
		db.Conn.QueryRow("select flags from trace_data where trace_number = ?", trace_number).Scan(&flags)
		flags |= 0x40;
		
		Preset := int(rune(payload[4]))

		D1 := float32(int((rune(payload[5]) << 8) + (rune(payload[6])))) / 100
		W1 := int((rune(payload[7]) << 8) + (rune(payload[8])))
		A1 := int(rune(payload[9]))
	
		D2 := float32(int((rune(payload[10]) << 8) + (rune(payload[11])))) / 100
		W2 := int((rune(payload[12]) << 8) + (rune(payload[13])))
		A2 := int(rune(payload[14]))
	
		D3 := float32(int((rune(payload[15]) << 8) + (rune(payload[16])))) / 100
		W3 := int((rune(payload[17]) << 8) + (rune(payload[18])))
		A3 := int(rune(payload[19]))
	
		D4 := float32(int((rune(payload[20]) << 8) + (rune(payload[21])))) / 100
		W4 := int((rune(payload[22]) << 8) + (rune(payload[23])))
		A4 := int(rune(payload[24]))

		stmt, _ = db.Conn.Prepare("update trace_data set preset = ?, D1 = ?, W1 = ?, A1 = ?, D2 = ?, W2 = ?, A2 = ?, D3 = ?, W3 = ?, A3 = ?, D4 = ?, W4 = ?, A4 = ?, flags = ? where trace_number = ?")
		stmt.Exec(&Preset, &D1, &W1, &A1, &D2, &W2, &A2, &D3, &W3, &A3, &D4, &W4, &A4, &flags, &trace_number)
		
		var PRD_DB, preset_DB, state int
		db.Conn.QueryRow("select measurement_frequency, default_preset, state from sensor where deveui = ?", t.EUI).Scan(&PRD_DB, &preset_DB, &state)
		
		//generate downlink
		var data string
		var port int
		
		payload = []byte {0x00, 0x00, 0x00}
		payload[0] = byte((PRD_DB >> 8) & 0xff)
		payload[1] = byte(PRD_DB & 0xff)
		payload[2] = byte(preset_DB & 0xff)

		data = strings.ToUpper(hex.EncodeToString(payload))				
		port = 42
		var sensorData Trace
		sensorData.Create_downlink(data, port, false, t.EUI)
	}
}

func (t *Trace) Create_downlink(data string, port int, flush_queue bool, deveui string) {  
	
	var appeui, app_id, api_key, dl_token string
	db.Conn.QueryRow("select appeui from sensor where deveui = ?", deveui).Scan(&appeui)
	db.Conn.QueryRow("select APP_ID, API_key, Downlink_token from teracom_account where AppEUI = ?", appeui).Scan(&app_id, &api_key, &dl_token)
	
	type teracom_downlink struct {
		Appid		string	`json:"appid"`
		Cmd			string	`json:"cmd"`
		EUI			string	`json:"EUI"`
		Port		int		`json:"port"`
		Confirmed	bool	`json:"confirmed"`
		Data		string	`json:"data"`
	}
	
	dl := teracom_downlink {
		Appid:		app_id,
		Cmd:		"tx",
		EUI:		deveui,
		Port:		port,
		Confirmed:	false,
		Data:		data,
	}
	
	var jsonData []byte
	jsonData, _ = json.Marshal(dl)
	
	//create and execute request to TERACOM API
	client := &http.Client{}
	
	if flush_queue {
		flushURL := "https://iotnet.teracom.dk/1/nwk/app/" + app_id + "/device/" + deveui + "/flush_queue"
		flushReq, _ := http.NewRequest("POST", flushURL, strings.NewReader(string("")))
		flushReq.Header.Set("Authorization", "Bearer " + api_key)
	    flushReq.Header.Set("Content-Type", "application/json")
	    client.Do(flushReq)	
	}
    
	req, _ := http.NewRequest("POST", "https://iotnet.teracom.dk/1/rest", strings.NewReader(string(jsonData)))
    req.Header.Set("Authorization", dl_token)
    req.Header.Set("Content-Type", "application/json")
    client.Do(req)
}

func (m *MassSensor) Update() {

	var sensorValues []interface{}
	var containerValues []interface{}
	sensorQuery := "UPDATE sensor SET"
	containerQuery := "UPDATE container SET"

	if m.MeasurementFrequency != nil {
		sensorQuery += " measurement_frequency = ?,"
		sensorValues = append(sensorValues, m.MeasurementFrequency)
	}
	if m.Firmware != nil {
		sensorQuery += " firmware_version = ?,"
		sensorValues = append(sensorValues, m.Firmware)
	}
	if m.DisposalCompany != nil {
		containerQuery += " disposal_company_id = ?,"
		containerValues = append(containerValues, m.DisposalCompany)
	}
	if m.ContainerType != nil {
		containerQuery += " container_type_id = ?,"
		containerValues = append(containerValues, m.ContainerType)
	}
	if m.Address != nil {
		containerQuery += " address = ?,"
		containerValues = append(containerValues, m.Address)
	}
	if m.Lat != nil {
		containerQuery += " lat = ?,"
		containerValues = append(containerValues, m.Lat)
	}
	if m.Lng != nil {
		containerQuery += " lng = ?,"
		containerValues = append(containerValues, m.Lng)
	}
	if m.Height != nil {
		containerQuery += " container_height = ?,"
		containerValues = append(containerValues, m.Height)
	}
	if m.Size != nil {
		containerQuery += " size = ?,"
		containerValues = append(containerValues, m.Size)
	}
	if m.SizeType != nil {
		containerQuery += " size_type = ?,"
		containerValues = append(containerValues, m.SizeType)
	}
	if m.Fraction != nil {
		containerQuery += " fraction_type_id = ?,"
		containerValues = append(containerValues, m.Fraction)
	}
	if m.Note != nil {
		containerQuery += " note = ?,"
		containerValues = append(containerValues, m.Note)
	}

	sensorQuery = strings.TrimRight(sensorQuery, ",")
	sensorQuery += " WHERE id = ?"
	containerQuery = strings.TrimRight(containerQuery, ",")
	containerQuery += " WHERE id = ?"

	for _, sensor := range m.Sensors {
		if len(sensorValues) > 0 {
			sensorValues = append(sensorValues, sensor.Id)
			stmt, _ := db.Conn.Prepare(sensorQuery)
			stmt.Exec(sensorValues...)
			defer stmt.Close()
			sensorValues = sensorValues[:len(sensorValues)-1]
		}
		if len(containerValues) > 0 {
			containerValues = append(containerValues, sensor.Id)
			stmt, _ := db.Conn.Prepare(containerQuery)
			stmt.Exec(containerValues...)
			defer stmt.Close()
			containerValues = containerValues[:len(containerValues)-1]
		}

		if m.OperatorId != nil {
			stmt, _ := db.Conn.Prepare("DELETE FROM container_user WHERE container_id = ? AND type = 1")
			stmt.Exec(&sensor.Id)
			stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 1)")
			stmt.Exec(&m.OperatorId, &sensor.Id)
			defer stmt.Close()
		}
		if m.ResellerId != nil {
			stmt, _ := db.Conn.Prepare("DELETE FROM container_user WHERE container_id = ? AND type = 2")
			stmt.Exec(&sensor.Id)
			stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 2)")
			stmt.Exec(&m.ResellerId, &sensor.Id)
			defer stmt.Close()
		}
		if m.ViewerId != nil {
			stmt, _ := db.Conn.Prepare("DELETE FROM container_user WHERE container_id = ? AND type = 3")
			stmt.Exec(&sensor.Id)
			stmt, _ = db.Conn.Prepare("INSERT INTO container_user (user_id, container_id, type) VALUES (?, ?, 3)")
			stmt.Exec(&m.ViewerId, &sensor.Id)
			defer stmt.Close()
		}
	}
}

/********** Return type Functions **********/

func GetSensorsForAdmin() []Sensor {
	var sensors []Sensor
	rows, _ := db.Conn.Query("SELECT id, sensor_id, latest_measurement, deployment_date, imsi FROM sensor WHERE deleted != 1")
	for rows.Next() {
		var sensor Sensor
		rows.Scan(&sensor.Id, &sensor.SensorId, &sensor.LatestMeasurement, &sensor.DeploymentDate, &sensor.Imsi)
		sensors = append(sensors, sensor)
	}

	defer rows.Close()

	return sensors
}

func GetSensorsForReseller(u User) []Sensor {
	var sensors []Sensor
	rows, _ := db.Conn.Query("SELECT id, sensor_id, latest_measurement, deployment_date, imsi FROM sensor WHERE user_id = ? AND deleted != 1", u.Id)
	for rows.Next() {
		var sensor Sensor
		rows.Scan(&sensor.Id, &sensor.SensorId, &sensor.LatestMeasurement, &sensor.DeploymentDate, &sensor.Imsi)
		sensors = append(sensors, sensor)
	}

	defer rows.Close()

	return sensors
}

func GetNextSensorId() int {
	var next int
	db.Conn.QueryRow("SELECT MAX(id) FROM sensor").Scan(&next)
	return next
}

func GetSensorInfoFromContainerId(id int) Sensor {
	var sensor Sensor

	db.Conn.QueryRow("SELECT appkey, measurement_frequency, deveui FROM sensor WHERE id = ?", id).Scan(&sensor.AppKey, &sensor.MeasurementFrequency, &sensor.Deveui)

	return sensor
}

func ReturnSensorLogs() []SensorLog {
	var sensorLogs []SensorLog
	rows, _ := db.Conn.Query("SELECT api.timestamp, api.dev_eui, api.decrypted_plaintext, api.ip_address, (SELECT id FROM container WHERE sensor_id = (SELECT id FROM sensor WHERE deveui = api.dev_eui)) FROM container_api_log api LIMIT 100")
	for rows.Next() {
		var sensorLog SensorLog
		rows.Scan(&sensorLog.Timestamp, &sensorLog.Deveui, &sensorLog.Payload, &sensorLog.IP, &sensorLog.ContainerId)
		sensorLogs = append(sensorLogs, sensorLog)
	}

	defer rows.Close()

	return sensorLogs
}

func GetSensorNamesForExternalEndpoint(name string) []string {
	var names []string
	rows, _ := db.Conn.Query("SELECT s.sensor_id FROM endpoint e LEFT JOIN endpoint_sensor es ON e.id = es.endpoint_id LEFT JOIN sensor s ON es.sensor_id = s.id WHERE e.name = ? AND e.deleted != 1", name)
	for rows.Next() {
		var name string
		rows.Scan(&name)
		names = append(names, name)
	}

	defer rows.Close()

	return names
}
