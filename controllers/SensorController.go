package controllers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	. "wastecontrol/models"
	"github.com/gin-gonic/gin"
	"wastecontrol/db"
)

type ecb struct {
	b         cipher.Block
	blockSize int
}

type ecbDecrypter ecb
type ecbEncrypter ecb

func CreateSensor(c *gin.Context) {
	var sensor Sensor
	var caller User
	c.BindJSON(&sensor)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()
	if caller.UserType != 1 {
		c.JSON(401, "")
	}
	sensor.Create()

	c.JSON(200, gin.H{"message": "Sensor created"})
}

func CreateNbiotSensor(c *gin.Context) {
	var sensor Sensor
	c.BindJSON(&sensor)

	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])

	sensor.CreateNbiot(reqBody)

	c.JSON(200, sensor)
}

func CreateLorawanSensor(c *gin.Context) {
	var sensor Sensor
	c.BindJSON(&sensor)
	
	buf := make([]byte, 1024)
	num, _ := c.Request.Body.Read(buf)
	reqBody := string(buf[0:num])

	sensor.CreateLorawan(reqBody)

	c.JSON(200, sensor)
}

func GetSensors(c *gin.Context) {
	var user User

	user.Token = GetToken(c.GetHeader("Authorization"))
	user.GetId()
	user.GetRole()
	switch user.UserType {
	case 1:
		c.JSON(200, GetSensorsForAdmin())
		return
	case 2:
		c.JSON(200, GetSensorsForReseller(user))
		return
	case 3:
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	case 4:
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
}

func GetSensor(c *gin.Context) {
	var sensor Sensor

	sensor.SensorId = c.Param("Id")
	sensor.Get()
	resellers := GetResellers()

	c.JSON(200, gin.H{
		"sensor":    sensor,
		"resellers": resellers,
	})
}

func GetSensorLogs(c *gin.Context) {
	c.JSON(200, ReturnSensorLogs())
}

func UpdateSensor(c *gin.Context) {
	var sensor Sensor
	var caller User

	c.BindJSON(&sensor)

	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	sensor.Update()

	c.JSON(200, sensor)
}

func UpdateMassSensor(c *gin.Context) {
	var massSensor MassSensor
	c.BindJSON(&massSensor)

	massSensor.Update()

	c.JSON(200, massSensor)
}

func DeleteSensor(c *gin.Context) {
	var sensor Sensor
	var caller User

	sensor.Id, _ = strconv.Atoi(c.Param("Id"))
	caller.Token = GetToken(c.GetHeader("Authorization"))
	caller.GetId()
	caller.GetRole()

	sensor.Delete()

	c.JSON(200, "")
}

func GetTTNData(c *gin.Context) {
	var ttn TTN

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	ttn.XML = string(bodyBytes)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	c.BindJSON(&ttn)

	if len(ttn.PayloadRaw) < 1 {
		ttn.UploadEmptyLog()
		c.JSON(400, gin.H{"Message": "Empty payload"})
		return
	}

	ttn.IP = c.ClientIP()

	ttn.Upload()

	c.JSON(200, ttn)
}

func GetTeracomData(c *gin.Context) {
	var teracom Teracom

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	teracom.XML = string(bodyBytes)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	c.BindJSON(&teracom)

	if len(teracom.Data) < 1 {
		teracom.UploadEmptyLog()
		c.JSON(200, gin.H{"Message": "Empty payload"})
		return
	}

	teracom.IP = c.ClientIP()
	teracom.Upload()

	c.JSON(200, teracom)
}

func GetNbiotData(c *gin.Context) {
	var nbiot Nbiot

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	nbiot.Payload = bodyBytes
	if len(nbiot.Payload) != 20 {
		nbiot.UploadEmptyLog()
		c.JSON(200, gin.H{"Message": "Wrong payload size"})
		return
	}
	nbiot.ContainerId = int((rune(nbiot.Payload[0]) << 24) + (rune(nbiot.Payload[1]) << 16) + (rune(nbiot.Payload[2]) << 8) + (rune(nbiot.Payload[3])))
	nbiot.GetByContainerId()
	key, _ := hex.DecodeString(nbiot.AppKey)

	var encrypted []byte = nbiot.Payload[4:20]

	nbiot.DecryptedPayload, _ = AesDecrypt(encrypted, key)
	nbiot.WasteDistance = fmt.Sprint(rune(nbiot.DecryptedPayload[8]) + rune(nbiot.DecryptedPayload[9]))
	nbiot.Error = string(rune(nbiot.DecryptedPayload[15]))
	nbiot.IP = c.ClientIP()
	nbiot.Upload()

	//bytes [0..7] -> container ID and UNIX timestamp copied from decrypted payload
	var returnData []byte = nbiot.DecryptedPayload[0:16]
	strMeasurementFreq, _ := strconv.Atoi(nbiot.MeasurementFrequency)
	//bytes [8..9] -> wakeup period
	returnData[8] = byte((strMeasurementFreq >> 8) & 0xff)
	returnData[9] = byte(strMeasurementFreq & 0xff)

	//bytes [10..13] -> IP address of the UDP server (95.179.158.9)
	returnData[10] = 95
	returnData[11] = 179
	returnData[12] = 158
	returnData[13] = 9

	//byte 14 -> state
	returnData[14] = byte(nbiot.State)

	//byte 15 -> not used
	returnData[15] = 0

	encryptedDownlink, _ := AesEncrypt(returnData, key)

	c.String(200, strings.ToUpper(hex.EncodeToString(encryptedDownlink)))
}

func GetNbiotFourObjects(c *gin.Context) {
	var sensorData Trace
	var nbiot Nbiot

	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	nbiot.Payload = bodyBytes
	if len(nbiot.Payload) != 36 {
		nbiot.UploadEmptyLog()
		c.JSON(200, gin.H{"Message": "Wrong payload size"})
		return
	}
	nbiot.ContainerId = int((rune(nbiot.Payload[0]) << 24) + (rune(nbiot.Payload[1]) << 16) + (rune(nbiot.Payload[2]) << 8) + (rune(nbiot.Payload[3])))
	nbiot.GetByContainerId()
	key, _ := hex.DecodeString(nbiot.AppKey)

	response := ""
	response += strconv.Itoa(nbiot.ContainerId) + "\r\n"
	response += "AppKey -> " + hex.EncodeToString(key) + "\r\n"

	var encrypted1 []byte = nbiot.Payload[4:20]
	var encrypted2 []byte = nbiot.Payload[20:36]
	var decrypted1, decrypted2 []byte

	decrypted1, _ = AesDecrypt(encrypted1, key)
	decrypted2, _ = AesDecrypt(encrypted2, key)

	response += "Decrypted1 -> " + hex.EncodeToString(decrypted1) + "\r\n"
	response += "Decrypted2 -> " + hex.EncodeToString(decrypted2) + "\r\n"
	response += "Default preset -> " + strconv.Itoa(nbiot.DefaultPreset) + "\r\n"

	sensorData.Timestamp = int((rune(decrypted1[0]) << 24) + (rune(decrypted1[1]) << 16) + (rune(decrypted1[2]) << 8) + (rune(decrypted1[3])))
	sensorData.ContainerID = int((rune(decrypted1[4]) << 24) + (rune(decrypted1[5]) << 16) + (rune(decrypted1[6]) << 8) + (rune(decrypted1[7])))

	//TODO: check if the sensorData.Timestamp is greater than the previously recorded

	//check if the container number matches
	if nbiot.ContainerId != sensorData.ContainerID {
		nbiot.UploadEmptyLog()
		c.JSON(200, gin.H{"Message": "Container number mismatch"})
		return
	}

	sensorData.D1 = float32((rune(decrypted1[8])<<8)+(rune(decrypted1[9]))) / 10
	sensorData.D2 = float32((rune(decrypted1[10])<<8)+(rune(decrypted1[11]))) / 10
	sensorData.D3 = float32((rune(decrypted1[12])<<8)+(rune(decrypted1[13]))) / 10
	sensorData.D4 = float32((rune(decrypted1[14])<<8)+(rune(decrypted1[15]))) / 10

	sensorData.W1 = int((rune(decrypted2[4]) << 8) + (rune(decrypted2[5])))
	sensorData.W2 = int((rune(decrypted2[6]) << 8) + (rune(decrypted2[7])))
	sensorData.W3 = int((rune(decrypted2[8]) << 8) + (rune(decrypted2[9])))
	sensorData.W4 = int((rune(decrypted2[10]) << 8) + (rune(decrypted2[11])))

	sensorData.A1 = int(rune(decrypted2[12]))
	sensorData.A2 = int(rune(decrypted2[13]))
	sensorData.A3 = int(rune(decrypted2[14]))
	sensorData.A4 = int(rune(decrypted2[15]))

	sensorData.UploadFourObjects()

	//generate response

	//bytes [0..7] -> UNIX timestamp and container ID copied from decrypted payload
	var returnData []byte = decrypted1[0:16]
	strMeasurementFreq, _ := strconv.Atoi(nbiot.MeasurementFrequency)
	//bytes [8..9] -> wakeup period
	returnData[8] = byte((strMeasurementFreq >> 8) & 0xff)
	returnData[9] = byte(strMeasurementFreq & 0xff)

	//bytes [10..13] -> IP address of the UDP server (95.179.158.9)
	returnData[10] = 95
	returnData[11] = 179
	returnData[12] = 158
	returnData[13] = 9

	//byte 14 -> state
	returnData[14] = byte(nbiot.State)

	//byte 15 -> default preset
	returnData[15] = byte(nbiot.DefaultPreset)

	encryptedDownlink, _ := AesEncrypt(returnData, key)

	c.String(200, hex.EncodeToString(encryptedDownlink))
	//c.String(200, strings.ToUpper(hex.EncodeToString(encryptedDownlink)))
	//c.String(200, response)
}

func UploadCSV(c *gin.Context) {
	file, err := c.FormFile("file")
	reseller, _ := strconv.Atoi(c.Param("resellerId"))
	operator, _ := strconv.Atoi(c.Param("operatorId"))
	CheckErr(err)
	newFile, _ := file.Open()
	reader := csv.NewReader(newFile)
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			panic(error)
		}
		var sensor Sensor
		sensor.SensorId = line[0]
		sensor.ResellerId = &reseller
		sensor.OperatorId = &operator
		sensor.CreateFromCSV()
	}
	c.JSON(200, "")
}

func AesEncrypt(src, key []byte) ([]byte, error) {
	Block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src) == 0 {
		return nil, errors.New("plaintext empty")
	}
	mode := NewECBEncrypter(Block)
	ciphertext := src
	mode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext, nil
}

func AesDecrypt(src, key []byte) ([]byte, error) {
	Block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(src) == 0 {
		return nil, errors.New("plaintext empty")
	}
	mode := NewECBDecrypter(Block)
	ciphertext := src
	mode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext, nil
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		dst = dst[x.blockSize:]
		src = src[x.blockSize:]
	}
}

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbDecrypter) BlockSize() int {
	return x.blockSize
}

func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Decrypt(dst, src[:x.blockSize])
		dst = dst[x.blockSize:]
		src = src[x.blockSize:]
	}
}

func UCReadParams(c *gin.Context) {
	var sensor Sensor
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	sensor.Id, _ = strconv.Atoi(string(bodyBytes))
	sensor.GetParamsThresholds()

	c.String(200, *sensor.Pga460Params+*sensor.Thresholds)
}

func UCReadTraceData(c *gin.Context) {
	var trace Trace
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	trace.ContainerID, _ = strconv.Atoi(string(bodyBytes))
	trace.GetTraceData()

	c.JSON(200, trace)
}

func Generate_trace_req_downlink(trace Trace, deveui string, flush_queue bool) {
	stmt, _ := db.Conn.Prepare("insert into trace_data (container) values (?)")
	stmt.Exec(&trace.Container)
	
	var trace_id int
	db.Conn.QueryRow("select max(trace_number) from trace_data ").Scan(&trace_id)
	
	payload := []byte {0x00, 0x00, 0x00, 0x00}
	payload[0] = byte((trace_id >> 24) & 0xff)
	payload[1] = byte((trace_id >> 16) & 0xff)
	payload[2] = byte((trace_id >> 8) & 0xff)
	payload[3] = byte(trace_id & 0xff)
	
	data := strings.ToUpper(hex.EncodeToString(payload))	
	port := trace.Preset
	trace.Create_downlink(data, port, flush_queue, deveui)
}

func UCRequest(c *gin.Context) {
	var trace Trace
	c.BindJSON(&trace)

	var deveui, sensor_type string
	db.Conn.QueryRow("select deveui,type from sensor where id = ?", trace.Container).Scan(&deveui, &sensor_type)
	
	if sensor_type[3:7] == "VU-2" {
		trace.Request()
	} else {
		Generate_trace_req_downlink(trace, deveui, true)
	}

	c.String(200, "Container #"+strconv.Itoa(trace.Container)+" trace request active for P"+strconv.Itoa(trace.Preset))
}

func UCUpdate(c *gin.Context) {
	var trace Trace
	var deveui, sensor_type string
	var port int
	
	c.BindJSON(&trace)

	db.Conn.QueryRow("select deveui,type from sensor where id = ?", trace.Container).Scan(&deveui, &sensor_type)
	
	if sensor_type[3:7] == "VU-3" {
		port = trace.Preset +2
		trace.Preset = 0
	} else {
		trace.Preset += 2
	}
	
	trace.Update()
	
	if sensor_type[3:7] == "VU-3" {
		trace.Preset = port
		trace.Create_downlink(trace.Pga460Params, port, true, deveui)
		trace.Create_downlink(trace.Thresholds, 5, false, deveui)
		Generate_trace_req_downlink(trace, deveui, false)
	}

	c.String(200, "Container #"+strconv.Itoa(trace.Container)+" parameters updated")
}

func UCAdd(c *gin.Context) {
	var bodyBytes []byte
	bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	bodyStr := string(bodyBytes)
	body, _ := hex.DecodeString(bodyStr)

	if len(body) != 228 {
		c.JSON(200, gin.H{"Message": "Wrong payload size"})
		return
	}

	var trace Trace

	trace.ContainerID = int((rune(body[0]) << 24) + (rune(body[1]) << 16) + (rune(body[2]) << 8) + (rune(body[3])))
	trace.Trace = bodyStr[8:264]
	trace.Pga460Params = bodyStr[264:350]
	trace.Thresholds = bodyStr[350:414]
	trace.Preset = int(rune(body[207]))

	trace.D1 = float32(int((rune(body[208])<<8)+(rune(body[209])))) / 100
	trace.W1 = int((rune(body[210]) << 8) + (rune(body[211])))
	trace.A1 = int(rune(body[212]))

	trace.D2 = float32(int((rune(body[213])<<8)+(rune(body[214])))) / 100
	trace.W2 = int((rune(body[215]) << 8) + (rune(body[216])))
	trace.A2 = int(rune(body[217]))

	trace.D3 = float32(int((rune(body[218])<<8)+(rune(body[219])))) / 100
	trace.W3 = int((rune(body[220]) << 8) + (rune(body[221])))
	trace.A3 = int(rune(body[222]))

	trace.D4 = float32(int((rune(body[223])<<8)+(rune(body[224])))) / 100
	trace.W4 = int((rune(body[225]) << 8) + (rune(body[226])))
	trace.A4 = int(rune(body[227]))
	/*
		fmt.Println("Container     ->", trace.ContainerID)
		fmt.Println("Trace data    ->", trace.Trace)
		fmt.Println("PGA460 params ->", trace.Pga460Params)
		fmt.Println("Thresholds    ->", trace.Thresholds)
		fmt.Println("Preset        ->", trace.Preset)
		fmt.Println("Object1       -> D =", trace.D1, "\tW =", trace.W1, "\tA =", trace.A1)
		fmt.Println("Object2       -> D =", trace.D2, "\tW =", trace.W2, "\tA =", trace.A2)
		fmt.Println("Object3       -> D =", trace.D3, "\tW =", trace.W3, "\tA =", trace.A3)
		fmt.Println("Object4       -> D =", trace.D4, "\tW =", trace.W4, "\tA =", trace.A4)
	*/
	trace.Add()

	c.String(200, "New post entered successfully")
}
