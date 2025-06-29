package repository

import (
	"database/sql"
	"fmt"
	"service/models"
)

func need_read(rows *sql.Rows) []models.Need {
	var needs []models.Need

	for rows.Next() {
		var need models.Need
		if err := rows.Scan(&need.Id, &need.Description); err != nil {
			fmt.Print(err)
			continue
		}
		needs = append(needs, need)
	}
	return needs
}

func Need_Get_All() []models.Need {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `Description` FROM Need")
	if err != nil {
		return nil
	}
	defer rows.Close()
	return need_read(rows)

}

func Need_Get_By_Id(id int) *models.Need {
	db := Database_Get()
	rows, err := db.Query("SELECT `Id`, `Description` FROM Need WHERE Id = ?", id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	needs := need_read(rows)
	if len(needs) > 0 {
		return &needs[0]
	}
	return nil
}

func Need_Get_All_By_TechId(id int) []models.Need {
	db := Database_Get()
	rows, err := db.Query("select `Id`, `Description` from Need where Id in (select NeedId from TechNeed where TechId = ?)", id)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()
	return need_read(rows)
}

func Need_Save(need *models.Need) error {
	db := Database_Get()
	if need.Id == 0 {
		result, err := db.Exec("insert into TechNeed () values ()")
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		need.Id = int(id)
	}
	// Do not allow to change tech the need is linked to
	_, err := db.Exec(
		"update TechNeed set `Description` = ? where `Id` = ?",
		need.Description,
		need.Id,
	)
	return err
}
