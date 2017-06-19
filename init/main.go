package main

import (
	db "github.com/davidwalter0/go-persist"
)

func init() {
	db := db.Connect()
	db.Initialize()
}

func main() {
}
