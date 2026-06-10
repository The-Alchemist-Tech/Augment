package augment

import (
	"database/sql"
	"log"
	"net/http"

	fund "augment/fund"

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
}