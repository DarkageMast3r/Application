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
		if err := rows.Scan(&techChoice.Id, &techChoice.TechId, &techChoice.CaseId, &techChoice.Status, &techChoice.Reasoning); err != nil {
			fmt.Println(err)
			continue
		}
		techChoice.Tech = Tech_Get_By_Id(techChoice.TechId)
		techChoices = append(techChoices, techChoice)
	}
	return techChoices
}

func TechChoice_Get_All() []models.TechChoice {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `TechId`, `CaseId`, `Status`, `Reasoning` FROM TechChoice")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	return techChoice_read(rows)
}

func TechChoice_Get_All_By_Case(caseId int) []models.TechChoice {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `TechId`, `CaseId`, `Status`, `Reasoning` FROM TechChoice where CaseId = ?", caseId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	return techChoice_read(rows)

}

func TechChoice_Get_By_Id(id int) *models.TechChoice {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `TechId`, `CaseId`, `Status`, `Reasoning` FROM TechChoice WHERE Id = ?", id)
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
		result, err := db.Exec("insert into TechChoice (TechId, CaseId) values (?, ?)", techChoice.TechId, techChoice.CaseId)
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
		"update TechChoice set `Status` = ?, `Reasoning` = ? where `Id` = ?",
		techChoice.Status,
		techChoice.Reasoning.String,
		techChoice.Id,
	)
	return err
}
