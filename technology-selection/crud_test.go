package main

import (
	"database/sql"
	"io"
	"math/rand"
	"os"
	"service/models"
	"service/repository"
	"testing"
)

func emptyDatabase() *sql.DB {
	db := repository.Database_Get_Test()
	sqlFile, err := os.Open("deploy/TestData.sql")
	if err != nil {
		return db
	}
	defer sqlFile.Close()

	content, err := io.ReadAll(sqlFile)
	if err != nil {
		return db
	}
	db.Exec(string(content))
	return db
}

func createTechnology() (models.Tech, error) {
	categories := repository.Category_Get_All()
	category := categories[rand.Intn(len(categories))]
	needs := repository.Need_Get_All()

	tech := models.Tech{
		CategoryId: category.Id,
		Needs:      needs,
		Name:       "Test_tech",
		Cost:       123,
	}
	err := repository.Tech_Save(&tech)
	return tech, err
}

func TestCreateTechnology(t *testing.T) {
	emptyDatabase()
	categories := repository.Category_Get_All()
	category := categories[rand.Intn(len(categories))]
	needs := repository.Need_Get_All()

	tech := models.Tech{
		CategoryId: category.Id,
		Needs:      needs,
		Name:       "Test_tech",
		Cost:       123.456,
	}
	err := repository.Tech_Save(&tech)
	if err != nil {
		t.Error("Error during saving tech:", err)
	}
	if tech.Id == 0 {
		t.Error("Tech unsuccesfully saved (ID = 0)")
	}

	newTech := repository.Tech_Get_By_Id(tech.Id)
	if newTech.Name != tech.Name {
		t.Errorf("Expected Name = %s, got %s.\n", tech.Name, newTech.Name)
	}
	if newTech.Cost != tech.Cost {
		t.Errorf("Expected Cost = %f, got %f.\n", tech.Cost, newTech.Cost)
	}
	if newTech.CategoryId != tech.CategoryId {
		t.Errorf("Expected CategoryId = %d, got %d.\n", tech.CategoryId, newTech.CategoryId)
	}
	if len(newTech.Needs) != len(tech.Needs) {
		t.Errorf("Expected needs len = %d, got %d.\n", len(tech.Needs), len(newTech.Needs))
	}
}

func TestUpdateTechnology(t *testing.T) {
	emptyDatabase()

	tech, err := createTechnology()
	if err != nil {
		t.Error(err)
	}

	savedTech := repository.Tech_Get_By_Id(tech.Id)
	tech.Name = "Updated name"
	repository.Tech_Save(&tech)
	updatedTech := repository.Tech_Get_By_Id(tech.Id)

	if updatedTech.Name != tech.Name {
		t.Errorf("Expected Name = %s, got %s.\n", tech.Name, updatedTech.Name)
	}
	if updatedTech.Name == savedTech.Name {
		t.Errorf("Expected Name != %s\n", updatedTech.Name)
	}
}

func TestDeleteTechnology(t *testing.T) {
	emptyDatabase()

	tech, err := createTechnology()
	if err != nil {
		t.Error(err)
	}

	repository.Tech_Delete(tech.Id)

	if repository.Tech_Get_By_Id(tech.Id) != nil {
		t.Errorf("Tech (ID = %d) did not get deleted.\n", tech.Id)
	}
}

func createNeed() (models.Need, error) {
	need := models.Need{
		Description: "Test_Description",
	}
	err := repository.Need_Save(&need)
	return need, err
}

func TestCreateNeed(t *testing.T) {
	emptyDatabase()
	need := models.Need{
		Description: "Test_Description",
	}
	err := repository.Need_Save(&need)
	if err != nil {
		t.Error("Error during saving need:", err)
	}
	if need.Id == 0 {
		t.Error("Need unsuccesfully saved (ID = 0)")
	}

	newNeed := repository.Need_Get_By_Id(need.Id)
	if newNeed.Description != need.Description {
		t.Errorf("Expected Description = %s, got %s.\n", need.Description, newNeed.Description)
	}
}

func TestUpdateNeed(t *testing.T) {
	emptyDatabase()

	need, err := createNeed()
	if err != nil {
		t.Error(err)
	}

	savedNeed := repository.Need_Get_By_Id(need.Id)
	need.Description = "Updated description"
	repository.Need_Save(&need)
	updatedNeed := repository.Need_Get_By_Id(need.Id)

	if updatedNeed.Description != need.Description {
		t.Errorf("Expected Description = %s, got %s.\n", need.Description, updatedNeed.Description)
	}
	if updatedNeed.Description == savedNeed.Description {
		t.Errorf("Expected Description != %s\n", updatedNeed.Description)
	}
}

func TestDeleteNeed(t *testing.T) {
	emptyDatabase()

	need, err := createNeed()
	if err != nil {
		t.Error(err)
	}

	repository.Need_Delete(need.Id)

	if repository.Need_Get_By_Id(need.Id) != nil {
		t.Errorf("Need (ID = %d) did not get deleted.\n", need.Id)
	}
}

func createCategory() (models.Category, error) {
	category := models.Category{
		Name:         "Test_Name",
		Description:  "Test_Description",
		Technologies: repository.Tech_Get_All(),
	}
	err := repository.Category_Save(&category)
	return category, err
}

func TestCreateCategory(t *testing.T) {
	emptyDatabase()
	techs := repository.Tech_Get_All()
	category := models.Category{
		Name:         "Test_Name",
		Description:  "Test_Description",
		Technologies: techs,
	}
	err := repository.Category_Save(&category)
	if err != nil {
		t.Error("Error during saving category:", err)
	}
	if category.Id == 0 {
		t.Error("Category unsuccesfully saved (ID = 0)")
	}

	newCategory := repository.Category_Get_By_Id(category.Id)
	if newCategory.Name != category.Name {
		t.Errorf("Expected Name = %s, got %s.\n", category.Name, newCategory.Name)
	}
	if newCategory.Description != category.Description {
		t.Errorf("Expected Description = %s, got %s.\n", category.Description, newCategory.Description)
	}
	if len(newCategory.Technologies) != 0 {
		t.Error("Category technologies got loaded in, should be on-request\n")
	}
}

func TestUpdateCategory(t *testing.T) {
	emptyDatabase()

	category, err := createCategory()
	if err != nil {
		t.Error(err)
	}

	category.Name = "Updated name"
	category.Description = "Updated description"
	category.Technologies = category.Technologies[:1]
	repository.Category_Save(&category)
	updatedCategory := repository.Category_Get_By_Id(category.Id)
	if updatedCategory.Name != category.Name {
		t.Errorf("Expected Name = %s, got %s.\n", category.Name, updatedCategory.Name)
	}
	if updatedCategory.Description != category.Description {
		t.Errorf("Expected Description = %s, got %s.\n", category.Description, updatedCategory.Description)
	}
	if len(updatedCategory.Technologies) != 0 {
		t.Error("Category technologies got loaded in, should be on-request\n")
	}
}

func TestDeleteCategory(t *testing.T) {
	emptyDatabase()

	category, err := createCategory()
	if err != nil {
		t.Error(err)
	}

	repository.Category_Delete(category.Id)

	if repository.Category_Get_By_Id(category.Id) != nil {
		t.Errorf("Category (ID = %d) did not get deleted.\n", category.Id)
	}
}
