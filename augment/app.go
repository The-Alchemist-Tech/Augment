package augment

import (
	"database/sql"

	fund "augment/fund"
	investor "augment/investor"

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
		Path("/transaction/create").
		HandlerFunc(cap.CreateTransaction)
}