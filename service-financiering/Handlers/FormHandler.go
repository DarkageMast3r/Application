package handlers

import (
	// m "Financiering/Models"
	r "Financiering/Repositories"
	"fmt"
	"net/http"
	// "strconv"
)

func AddDossier(wr http.ResponseWriter, rq *http.Request) {
	db := r.Database_Get()
	// var budget m.Budget
	// var nieuwDossier m.FinancieringsDossier

	rq.ParseForm()

	clientid := rq.Form.Get("ClientID")
	zorgtechid := rq.Form.Get("ZorgTechID")
	// formBedrag := rq.Form.Get("MaxBedrag")
	// maxBedrag, err := strconv.Atoi(formBedrag)
	// if err != nil {
	// 	fmt.Println("AddDossier1: ", err)
	// 	return
	// }
	// nieuwDossier.VraagBudgetAan(float64(maxBedrag))

	stmt, err := db.Prepare("INSERT INTO financieringsdossier(ClientID, ZorgTechID) VALUES(?,?)")
	if err != nil {
		fmt.Println("AddDossier2: ", err)
		return
	}
	_, err = stmt.Exec(clientid, zorgtechid)
	if err != nil {
		fmt.Println("AddDossier3: ", err)
		return
	}
	HomePageHandler(wr, rq)
}
