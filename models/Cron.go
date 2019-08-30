package models

import (
	"wastecontrol/db"
)

func ContainerPickup(time string) []Pickup {
	var pickups []Pickup
	rows, _ := db.Conn.Query("SELECT p.*, (SELECT fraction_type FROM container_fraction_type WHERE id = p.fraction_type_id), (SELECT company_name FROM disposal_company WHERE id = p.disposal_company_id), (SELECT container_type FROM container_type WHERE id = p.container_type_id), (SELECT email FROM user WHERE id = p.user_id) FROM pickup p WHERE p.pickuptime = ?", time)

	for rows.Next() {
		var pickup Pickup
		rows.Scan(&pickup.Id, &pickup.Search, &pickup.ContainerId, &pickup.FractionTypeId, &pickup.DisposalCompanyId, &pickup.WasteDistance, &pickup.ContainerTypeId, &pickup.PickupTime, &pickup.UserId, &pickup.ViewerId, &pickup.Format, &pickup.FractionType, &pickup.DisposalCompany, &pickup.ContainerType, &pickup.UserEmail)
		pickups = append(pickups, pickup)
	}

	return pickups
}

func (p *Pickup) GetContainersForCsv(query string) []Container {
	var containers []Container
	rows, _ := db.Conn.Query(query)

	for rows.Next() {
		var container Container
		rows.Scan(&container.Id, &container.Latitude, &container.Longitude, &container.WasteDistance, &container.Height, &container.Address, &container.Viewer.Id, &container.Fraction, &container.ContainerType, &container.DisposalCompany, &container.Note)
		containers = append(containers, container)
		container.AddPickupLog()
	}

	return containers
}
