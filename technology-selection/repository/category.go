package repository

import (
	"database/sql"
	"service/models"
	"service/service"
)

func category_read(rows *sql.Rows) []models.Category {
	var categories []models.Category

	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.Id, &cat.Name, &cat.Description); err != nil {
			service.LogError(err)
			continue
		}
		categories = append(categories, cat)
	}
	return categories
}

func Category_Get_All() []models.Category {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `Name`, `Description` FROM Category")
	if err != nil {
		return nil
	}
	defer rows.Close()
	return category_read(rows)

}

func Category_Get_By_Id(id int) *models.Category {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `Name`, `Description` FROM Category WHERE Id = ?", id)
	if err != nil {
		service.LogError(err)
		return nil
	}
	defer rows.Close()
	categories := category_read(rows)
	if len(categories) > 0 {
		return &categories[0]
	}
	return nil
}

func Category_Save(category *models.Category) error {
	db := Database_Get()
	if category.Id == 0 {
		result, err := db.Exec("insert into Category () values ()")
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		category.Id = int(id)
	}
	_, err := db.Exec(
		"update Category set `Name` = ?, `Description` = ? where `Id` = ? ",
		category.Name,
		category.Description,
		category.Id,
	)
	return err
}

func Category_Delete(id int) error {
	db := Database_Get()
	_, err := db.Exec("delete from Category where `Id` = ?", id)
	return err
}
