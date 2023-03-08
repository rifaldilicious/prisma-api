package main

import (
	"crypto-project/api"
	"crypto-project/migrations"
)

func main() {
	migrations.Migrate()
	api.StartApi()
}
