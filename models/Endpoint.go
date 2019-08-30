package models

import(
	"wastecontrol/db"
)

type Endpoint struct {
	Id 				int 		`json:"id" binding:""`
	Name			string	`json:"name" binding:""`
	From			int			`json:"from" binding:""`
	Email			string	`json:"email" binding:""`
	To				int			`json:"to" binding:""`
	Username	string	`json:"username" binding:""`
	Password	string	`json:"password" binding:""`
}

/********** Void Functions **********/

func (e *Endpoint) GetSpecific() {
	db.Conn.QueryRow("SELECT username, password, name FROM endpoint WHERE id = ?", e.Id).Scan(&e.Username, &e.Password, &e.Name)
	db.Conn.QueryRow("SELECT MIN(sensor_id), MAX(sensor_id) FROM endpoint_sensor WHERE endpoint_id = ?", e.Id).Scan(&e.From, &e.To)
}

func (e *Endpoint) Create() {
	stmt, _ := db.Conn.Prepare("INSERT INTO endpoint (username, password, name, email) VALUES (?, ?, ?, ?)")
	res, _ := stmt.Exec(&e.Username, &e.Password, &e.Name, &e.Email)
	id, _ := res.LastInsertId()
	for i := e.From; i <= e.To; i++ {
		stmt, _ = db.Conn.Prepare("INSERT INTO endpoint_sensor (endpoint_id, sensor_id) VALUES (?, ?)")
		stmt.Exec(&id, &i)
	}

	defer stmt.Close()
}

func (e *Endpoint) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE endpoint SET name = ?, username = ?, password = ?, email = ? WHERE id = ?")
	stmt.Exec(&e.Name, &e.Username, &e.Password, &e.Email, &e.Id)

	stmt, _ = db.Conn.Prepare("DELETE FROM endpoint_sensor WHERE endpoint_id = ?")
	stmt.Exec(&e.Id)

	for i := e.From; i <= e.To; i++ {
		stmt, _ = db.Conn.Prepare("INSERT INTO endpoint_sensor (endpoint_id, sensor_id) VALUES (?, ?)")
		stmt.Exec(&e.Id, &i)
	}

	defer stmt.Close()
}

func (e *Endpoint) Delete() {
	stmt, _ := db.Conn.Prepare("UPDATE endpoint SET deleted = 1 WHERE id = ?")
	stmt.Exec(&e.Id)

	stmt, _ = db.Conn.Prepare("DELETE FROM endpoint_sensor WHERE endpoint_id = ?")
	stmt.Exec(&e.Id)

  defer stmt.Close()
}

/********** Return type Functions **********/

func GetAllEndpoints() []Endpoint {
	var endpoints []Endpoint
	rows, _ := db.Conn.Query("SELECT id, name FROM endpoint WHERE deleted != 1")
	defer rows.Close()

	for rows.Next() {
		var endpoint Endpoint
		rows.Scan(&endpoint.Id, &endpoint.Name)
		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}
