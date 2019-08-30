package controllers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	. "wastecontrol/models"

	"github.com/gin-gonic/gin"
)

func PickupCron() {
	t := time.Now()
	h := t.Hour()
	pickups := ContainerPickup(strconv.Itoa(h + 2))

	for _, pickup := range pickups {
		containerQuery := "SELECT c.id, c.lat, c.lng, c.waste_distance, c.container_height, c.address, " + strconv.Itoa(pickup.UserId) + ", cft.fraction_type, ct.container_type, d.company_name, c.note from container c LEFT JOIN container_user cu ON cu.container_id = c.id LEFT JOIN container_type ct ON ct.id = c.container_type_id LEFT JOIN container_fraction_type cft ON cft.id = c.fraction_type_id LEFT JOIN disposal_company d ON d.id = c.disposal_company_id LEFT JOIN pickup_history ph ON ph.container_id = c.id WHERE cu.user_id = " + strconv.Itoa(pickup.UserId) + " AND "

		if *pickup.Search != "" {
			containerQuery += "c.address LIKE " + *pickup.Search + " AND "
		}

		if pickup.ContainerId != nil {
			containerQuery += "c.id = " + strconv.Itoa(*pickup.ContainerId) + " AND "
		}

		if pickup.ViewerId != nil {
			containerQuery += "c.viewer_id = " + strconv.Itoa(*pickup.ViewerId) + " AND "
		}

		if pickup.FractionTypeId != nil {
			containerQuery += "c.fraction_type_id = " + strconv.Itoa(*pickup.FractionTypeId) + " AND "
		}

		if pickup.DisposalCompanyId != nil {
			containerQuery += "c.disposal_company_id = " + strconv.Itoa(*pickup.DisposalCompanyId) + " AND "
		}

		if pickup.ContainerTypeId != nil {
			containerQuery += "c.container_type_id = " + strconv.Itoa(*pickup.ContainerTypeId) + " AND "
		}

		containerQuery = strings.TrimRight(containerQuery, "AND ")
		containerQuery += " GROUP BY c.id"
		containers := pickup.GetContainersForCsv(containerQuery)
		if pickup.Format == 1 {
			file, err := ioutil.TempFile("", "pickups_auto_*.csv")
			CheckErr(err)
			defer os.Remove(file.Name())

			firstLine := [][]string{
				{"Container ID", "Note", "Latest pickup", "Address", "Viewer ID", "Fraction Type", "Customer", "Waste Distance", "Container Type"},
			}
			for _, container := range containers {
				if *pickup.WasteDistance != 0 {
					if container.WasteDistance < float64(*pickup.WasteDistance) {
						continue
					}
				}
				containerNote := ""
				containerAddress := ""
				containerFraction := ""
				containerDisposalCompany := ""
				containerType := ""
				containerLatestPickup := ""
				if container.Note != nil {
					containerNote = *container.Note
				}
				if container.Address != nil {
					containerAddress = *container.Address
				}
				if container.Fraction != nil {
					containerFraction = *container.Fraction
				}
				if container.DisposalCompany != nil {
					containerDisposalCompany = *container.DisposalCompany
				}
				if container.ContainerType != nil {
					containerType = *container.ContainerType
				}
				if container.LatestPickup != nil {
					containerLatestPickup = *container.LatestPickup
				}
				row := []string{strconv.Itoa(container.Id), containerNote, containerLatestPickup, containerAddress, strconv.Itoa(container.Viewer.Id), containerFraction, containerDisposalCompany, fmt.Sprintf("%f", container.WasteDistance), containerType}
				firstLine = append(firstLine, row)
			}

			w := csv.NewWriter(file)
			w.WriteAll(firstLine)

			if err := w.Error(); err != nil {
				panic(err)
			}

			mail := NewFileRequest([]string{*pickup.UserEmail}, "Automatisk pickup download", "Automatisk pickup download. Filen er vedhæftet", file.Name())
			mail.SendEmailWithFile()
		} else {
			file, err := ioutil.TempFile("", "pickups_auto_*.json")
			CheckErr(err)
			defer os.Remove(file.Name())

			_, _ = file.Write([]byte("{\n"))

			for _, container := range containers {
				if *pickup.WasteDistance != 0 {
					wd := CalculateWasteDistance(float64(container.WasteDistance), container.Height, 0, "")
					if wd > *pickup.WasteDistance {
						continue
					}
				}
				var containerJSON ContainerJSON
				containerJSON.Id = container.Id
				containerJSON.Note = container.Note
				containerJSON.LatestPickup = container.LatestPickup
				containerJSON.Address = container.Address
				containerJSON.ViewerId = &container.Viewer.Id
				containerJSON.Fraction = container.Fraction
				containerJSON.DisposalCompany = container.DisposalCompany
				containerJSON.WasteDistance = int(container.WasteDistance)
				containerJSON.ContainerType = container.ContainerType

				res, _ := json.Marshal(containerJSON)
				file.Write(res)
				file.Write([]byte(",\n"))
			}

			file.Write([]byte("}"))

			if err := file.Close(); err != nil {
				fmt.Println(err)
			}

			mail := NewFileRequest([]string{*pickup.UserEmail}, "Automatisk pickup download", "Automatisk pickup download. Filen er vedhæftet", file.Name())
			mail.SendEmailWithFile()
		}
	}
}

func TeliaTest(c *gin.Context) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", "api_wastecontrol")
	data.Set("password", "2gmgwLYa")
	req, err := http.NewRequest("POST", "https://api.teliaiot.com/api/v1/token", strings.NewReader(data.Encode()))
	req.Header.Set("API-Subscription-Key", "b4068f4835334f169a0c2904f23112ba")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	req.Header.Set("Cache-control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var raw map[string]interface{}
	json.Unmarshal(body, &raw)

	// now := time.Now()
	// sec := now.Unix()

	var jsonStr = []byte(`[{"asset_id": "d4fb9e5e-4416-4e7e-9f17-f40764bd92bb", "fill_level": 60, "battery": 100, "register_date": 1557391706000}]`)
	req, err = http.NewRequest("POST", "https://api.teliaiot.com/api/v1/services/wastecontrolteliaoffice/data", bytes.NewBuffer(jsonStr))
	req.Header.Set("API-Subscription-Key", "b4068f4835334f169a0c2904f23112ba")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-control", "no-cache")
	req.Header.Set("Authorization", "Bearer "+raw["access_token"].(string))

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

// Post data
// curl -X POST \
//   https://api.teliaiot.com/api/v1/services/weather/data \
//   -H 'API-Subscription-Key: <API_key> \
//   -H 'Authorization: Bearer <bearer_token>\
//   -H 'Content-Type: application/json' \
//   -H 'cache-control: no-cache' \
//   -d '[{
//     "asset_id": "6134f55e-27e0-40a6-9afe-7c0372f3799b",
//     "fill_level": 50,
//     "battery": 100,
//     "register_date": 1542751865000
//   }]'

// Get Token
// curl -X POST \
// https://api.teliaiot.com/api/v1/services/weather/data \
// -H 'API-Subscription-Key: <API_key> \
// -H 'Authorization: Bearer <bearer_token>\
// -H 'Content-Type: application/json' \
// -H 'cache-control: no-cache' \
// -d '[{
//   "asset_id": "6134f55e-27e0-40a6-9afe-7c0372f3799b",
//   "fill_level": 50,
//   "battery": 100,
//   "register_date": 1542751865000
// }]'
