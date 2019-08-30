package models

import (
	"wastecontrol/db"
)

type Disposal struct {
	Id          int             `json:"id" binding:""`
	CompanyName string          `json:"company_name" binding:""`
	Address     string          `json:"address" binding:""`
	Zip         string          `json:"zip" binding:""`
	Contact     *string         `json:"contact" binding:""`
	City        string          `json:"city" binding:""`
	Phone       *string         `json:"phone" binding:""`
	Emails      []DisposalEmail `json:"emails" binding:""`
	CreatedBy   int             `json:"created_by" binding:""`
}

type DisposalEmail struct {
	Email string `json:"email" binding:""`
}

/********** Void Functions **********/

func (d *Disposal) Get() {
	db.Conn.QueryRow("SELECT company_name, address, postal, city, phone FROM disposal_company WHERE id = ?", d.Id).Scan(&d.CompanyName, &d.Address, &d.Zip, &d.City, &d.Phone)
	d.GetDisposalEmails()
}

func (d *Disposal) Create() {
	stmt, _ := db.Conn.Prepare("INSERT INTO disposal_company (company_name, address, postal, city, phone, contact_person) VALUES (?, ?, ?, ?, ?, ?)")
	res, _ := stmt.Exec(&d.CompanyName, &d.Address, &d.Zip, &d.City, &d.Phone, &d.Contact)
	id, _ := res.LastInsertId()
	d.Id = int(id)

	stmt, _ = db.Conn.Prepare("INSERT INTO disposal_company_user (disposal_company_id, user_id) VALUES (?, ?)")
	stmt.Exec(&d.Id, &d.CreatedBy)

	d.CreateEmails()

	defer stmt.Close()
}

func (d *Disposal) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE disposal_company SET company_name = ?, address = ?, postal = ?, city = ?, phone = ?, contact_person = ? WHERE id = ?")
	stmt.Exec(&d.CompanyName, &d.Address, &d.Zip, &d.City, &d.Phone, &d.Contact, &d.Id)

	d.DeleteEmails()
	d.CreateEmails()

	defer stmt.Close()
}

func (d *Disposal) Delete() {
	stmt, _ := db.Conn.Prepare("UPDATE disposal_company SET deleted = 1 WHERE id = ?")
	stmt.Exec(&d.Id)
	d.DeleteEmails()

	defer stmt.Close()
}

func (d *Disposal) GetDisposalEmails() {
	rows, _ := db.Conn.Query("SELECT email FROM disposal_company_email WHERE disposal_company_id = ?", d.Id)
	for rows.Next() {
		var email DisposalEmail
		rows.Scan(&email.Email)
		d.Emails = append(d.Emails, email)
	}

	defer rows.Close()
}

func (d *Disposal) DeleteEmails() {
	stmt, _ := db.Conn.Prepare("DELETE FROM disposal_company_email WHERE disposal_company_id = ?")
	stmt.Exec(&d.Id)

	defer stmt.Close()
}

func (d *Disposal) CreateEmails() {
	stmt, _ := db.Conn.Prepare("INSERT INTO disposal_company_email (disposal_company_id, email) VALUES (?, ?)")
	for _, email := range d.Emails {
		stmt.Exec(&d.Id, &email.Email)
	}

	defer stmt.Close()
}

/********** Return type Functions **********/

func GetDisposalsForAdmin() []Disposal {
	var disposals []Disposal
	rows, _ := db.Conn.Query("SELECT id, company_name FROM disposal_company WHERE deleted != 1")
	defer rows.Close()
	for rows.Next() {
		var disposal Disposal
		rows.Scan(&disposal.Id, &disposal.CompanyName)
		disposals = append(disposals, disposal)
	}

	return disposals
}

func GetDisposalsForOperator(userId int) []Disposal {
	var disposals []Disposal
	rows, _ := db.Conn.Query("SELECT d.id, d.company_name FROM disposal_company d WHERE deleted != 1 AND d.id IN (SELECT disposal_company_id FROM disposal_company_user WHERE user_id = ?)", userId)
	defer rows.Close()
	for rows.Next() {
		var disposal Disposal
		rows.Scan(&disposal.Id, &disposal.CompanyName)
		disposals = append(disposals, disposal)
	}

	return disposals
}
