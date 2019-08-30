package models

import(
	"wastecontrol/db"
	"strings"
)

func CheckToken(token string) bool {
	token = GetToken(token)
	checkToken, _ := db.Conn.Query("SELECT expiration FROM auth WHERE token = ? AND expiration > NOW()", token)

	defer checkToken.Close()

	if CheckCount(checkToken) > 0 {
		return true
	}

	return false

}

func GetToken(token string) string {
	token = strings.Replace(token, "Bearer ", "", 1)
	return token
}
