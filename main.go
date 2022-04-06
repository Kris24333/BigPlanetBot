package main

import "country_bot/src"

func main() {
	src.Connection.InitializeDB()
	src.TB.InitializeTB()
	src.TB.Start()
}
