package repository

import (
	"database/sql"
	"fmt"
	"service/models"
	"time"
)

func category_read(rows *sql.Rows) []models.Category {
	var categories []models.Category

	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.Id, &cat.Name, &cat.Description, &cat.GeneratedOn); err != nil {
			fmt.Print(err)
			continue
		}
		categories = append(categories, cat)
	}
	return categories
}

func Category_Get_All() []models.Category {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `Name`, `Description`, `GeneratedOn` FROM Category")
	if err != nil {
		return nil
	}
	defer rows.Close()
	return category_read(rows)

}

func Category_Get_By_Id(id int) *models.Category {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `Name`, `Description`, `GeneratedOn` FROM Category WHERE Id = ?", id)
	if err != nil {
		fmt.Println(err)
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
		category.GeneratedOn = time.Now()
		category.Id = int(id)
	}
	_, err := db.Exec(
		"update Category set `Name` = ?, `Description` = ?, `GeneratedOn` = ? where `Id` = ? ",
		category.Name,
		category.Description,
		category.GeneratedOn,
		category.Id,
	)
	return err
}
