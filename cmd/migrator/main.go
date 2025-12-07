package main

import (
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/amberdance/url-shortener/internal/app"
)

func main() {
	a, err := app.GetApp()
	if err != nil {
		panic(err)
	}

	u, _ := url.Parse(a.Config().DatabaseDSN)
	db := dbmate.New(u)

	err = db.CreateAndMigrate()
	if err != nil {
		panic(err)
	}
}
