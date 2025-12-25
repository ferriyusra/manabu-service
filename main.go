package main

//go:generate swag init -g cmd/main.go -o docs

import "manabu-service/cmd"

func main() {
	cmd.Run()
}
