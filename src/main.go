package main

import (
	"pf"
	"fmt"
)

func main() {
	sim := pf.NewSim()
	/*sim.AddUnit(pf.P(0, 200), pf.P(200, 200))
	sim.AddUnit(pf.P(200, 200), pf.P(0, 210))*/
	
	black := pf.NewColor(0,0,0,1)
	//blue := pf.NewColor(0,0,1,1)
	//red := pf.NewColor(1,0,0,1)
	
	sim.AddBuilding(20, 50, 50, 100)
	//sim.AddBuilding(50, 150, 100, 50)
	//sim.AddBuilding(50, 25, 100, 25)
	
	sim.AddUnit(pf.P(200, 100), pf.P(0, 100), black)
	//sim.AddUnit(pf.P(200, 110), pf.P(0, 120), blue)
	
	for i:=0; i<10; i++ {
		sim.AddUnit(pf.P(200, 80+i*5), pf.P(0, 100), black)
	}
	
	sim.Init()
	
	fmt.Println("Starting sim")
	//go sim.Update()
	
	sim.Run()
}