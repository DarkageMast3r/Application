package repository

import (
	"database/sql"
	"service/models"
	"service/service"
)

func case_read(rows *sql.Rows) []models.Case {
	var cases []models.Case

	for rows.Next() {
		var clientCase models.Case
		if err := rows.Scan(
			&clientCase.Id,
			&clientCase.Name,
			&clientCase.Description,
			&clientCase.ClientId,
			&clientCase.IsClosed,
		); err != nil {
			service.LogError(err)
			continue
		}
		cases = append(cases, clientCase)
	}
	return cases
}

func Case_Get_All() []models.Case {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `Name`, `Description`, `ClientId`, `IsClosed` FROM `Case` WHERE `IsClosed` = 0")
	if err != nil {
		return nil
	}
	defer rows.Close()
	return case_read(rows)

}

func Case_Get_By_Id(id int) *models.Case {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `Name`, `Description`, `ClientId`, `IsClosed` FROM `Case` WHERE Id = ?", id)
	if err != nil {
		service.LogError(err)
		return nil
	}
	defer rows.Close()
	categories := case_read(rows)
	if len(categories) > 0 {
		return &categories[0]
	}
	return nil
}

func Case_Save(clientCase *models.Case) error {
	db := Database_Get()
	if clientCase.Id == 0 {
		result, err := db.Exec("insert into `Case` (`ClientId`) values (?)", clientCase.ClientId)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		clientCase.Id = int(id)
	}
	_, err := db.Exec(
		"update `Case` set `Name` = ?, `Description` = ?, `IsClosed` = ? where `Id` = ? ",
		clientCase.Name,
		clientCase.Description,
		clientCase.IsClosed,
		clientCase.Id,
	)
	return err
}
