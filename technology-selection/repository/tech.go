package repository

import (
	"database/sql"
	"fmt"
	"service/models"
)

func tech_read(rows *sql.Rows) []models.Tech {
	var techs []models.Tech

	for rows.Next() {
		var tech models.Tech
		if err := rows.Scan(&tech.Id, &tech.CategoryId, &tech.Name); err != nil {
			fmt.Println(err)
			continue
		}
		tech.Category = Category_Get_By_Id(tech.CategoryId)
		tech.Needs = Need_Get_All_By_TechId(tech.Id)
		techs = append(techs, tech)
	}
	return techs
}

func Tech_Get_All() []models.Tech {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `CategoryId`, `Name` FROM Tech")
	if err != nil {
		return nil
	}
	defer rows.Close()
	return tech_read(rows)

}

func Tech_Get_By_Id(id int) *models.Tech {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `CategoryId`, `Name` FROM Tech WHERE Id = ?", id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	techs := tech_read(rows)
	if len(techs) > 0 {
		return &techs[0]
	}
	return nil
}

func Tech_Save(tech *models.Tech) error {
	db := Database_Get()
	if tech.Id == 0 {
		result, err := db.Exec("insert into Tech () values ()")
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		tech.Id = int(id)
	}

	_, err := db.Exec(
		"update Tech set `Name` = ?, `CategoryId` = ? where `Id` = ?",
		tech.Name,
		tech.CategoryId,
		tech.Id,
	)
	if err != nil {
		return err
	}
	_, err = db.Exec("delete from TechNeed where `TechId` = ?", tech.Id)
	if err != nil {
		return err
	}
	for _, need := range tech.Needs {
		_, err = db.Exec("insert into TechNeed (`TechId`, `NeedId`) values (?, ?)", tech.Id, need.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func Tech_Delete(id int) error {
	db := Database_Get()
	_, err := db.Exec("delete from TechNeed where `TechId` = ?", id)
	if err != nil {
		return err
	}
	_, err = db.Exec("delete from Tech where `Id` = ?", id)
	return err
}
