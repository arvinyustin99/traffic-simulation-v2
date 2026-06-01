package main

import "github.com/arvinyustin99/traffic-simulation-v2/model"

func main() {
	model.InitTerminal()
	newSimulation := model.NewSimulation()
	newSimulation.Run()
}
