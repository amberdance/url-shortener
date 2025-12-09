package main

import (
	"log"
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

	u, err := url.Parse(a.Config().DatabaseDSN)
	if err != nil {
		log.Fatalln(err)
	}

	db := dbmate.New(u)

	err = db.CreateAndMigrate()
	if err != nil {
		panic(err)
	}
}
