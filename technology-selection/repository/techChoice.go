package repository

import (
	"database/sql"
	"fmt"
	"service/models"
)

func techChoice_read(rows *sql.Rows) []models.TechChoice {
	var techChoices []models.TechChoice

	for rows.Next() {
		var techChoice models.TechChoice
		if err := rows.Scan(&techChoice.Id, &techChoice.TechId, &techChoice.ClientId); err != nil {
			fmt.Print(err)
			continue
		}
		techChoice.Tech = Tech_Get_By_Id(techChoice.TechId)
		techChoices = append(techChoices, techChoice)
	}
	return techChoices
}

func TechChoice_Get_All() []models.TechChoice {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `TechId`, `ClientId`, `Status` FROM TechChoice")
	if err != nil {
		return nil
	}
	defer rows.Close()
	return techChoice_read(rows)

}

func TechChoice_Get_By_Id(id int) *models.TechChoice {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `TechId`, `ClientId`, `Status` FROM TechChoice WHERE Id = ?", id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	techChoicec := techChoice_read(rows)
	if len(techChoicec) > 0 {
		return &techChoicec[0]
	}
	return nil
}

func TechChoice_Save(techChoice *models.TechChoice) error {
	db := Database_Get()
	if techChoice.Id == 0 {
		result, err := db.Exec("insert into TechChoice (TechId, ClientId) values (?, ?)", techChoice.TechId, techChoice.ClientId)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		techChoice.Id = int(id)
	}
	_, err := db.Exec(
		"update TechChoice set `Status` = ? where `Id` = ?",
		techChoice.Status,
		techChoice.Id,
	)
	return err
}
