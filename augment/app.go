package augment

import (
	"database/sql"

	fund "augment/fund"
	investor "augment/investor"
	cap "augment/cap"

	"github.com/gorilla/mux"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	app.Router.
		Methods("POST").
		Path("/fund/create").
		HandlerFunc(fund.CreateFund)
	app.Router.
		Methods("POST").
		Path("/investor/create").
		HandlerFunc(investor.CreateInvestor)
	app.Router.
		Methods("POST").
		Path("/transfer/create").
		HandlerFunc(cap.CreateTransfer)
	app.Router.
		Methods("GET").
		Path("/cap/fund"). // Use query param "id" with an interger for the fund ID (EG "/cap/fund?id=1")
		HandlerFunc(cap.GetFundCap)
	app.Router.
		Methods("GET").
		Path("/cap/fund/history"). // Same query param as above - "/cap/fund/history?id=1"
		HandlerFunc(cap.GetFundCapHistory)
}