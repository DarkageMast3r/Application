package repository

import (
	"database/sql"
	"fmt"
	"service/models"
)

func category_read(rows *sql.Rows) []models.Category {
	var categories []models.Category

	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.Id, &cat.Naam, &cat.Beschrijving, &cat.GegenereerdOp); err != nil {
			fmt.Print(err)
			continue
		}
		categories = append(categories, cat)
	}
	return categories
}

func Category_Get_All() []models.Category {
	db := Database_Get()
	rows, err := db.Query("SELECT [Id], [Name], [Description], [GeneratedOn] FROM Category")
	if err != nil {
		return nil
	}
	defer rows.Close()
	return category_read(rows)

}

func Category_Get_By_Id(id int) *models.Category {
	db := Database_Get()
	rows, err := db.Query("SELECT [Id], [Name], [Description], [GeneratedOn] FROM Category WHERE Id = ?", id)
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
