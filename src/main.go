package main

import (
	"pf"
	"fmt"
	"geo"
)

func main() {
	sim := pf.NewSim()
	/*sim.AddUnit(pf.P(0, 200), pf.P(200, 200))
	sim.AddUnit(pf.P(200, 200), pf.P(0, 210))*/
	
	black := pf.NewColor(0,0,0,1)
	//blue := pf.NewColor(0,0,1,1)
	//red := pf.NewColor(1,0,0,1)
	
	sim.AddBuilding(100, 50, 50, 80)
	//sim.AddBuilding(50, 150, 100, 50)
	//sim.AddBuilding(50, 25, 100, 25)
	
	sim.Init()
	
	// TODO if the target point is on the corner of a poly, an error occurrs
	sim.AddUnit(geo.Vec{200, 100}, geo.Vec{1, 100}, black)
	
	for i:=0; i<10; i++ {
		sim.AddUnit(geo.Vec{200, 80.0+float64(i)*10.0}, geo.Vec{1, 100}, black)
	}
	
	fmt.Println("Starting sim")
	//go sim.Update()
	
	sim.Run()
}