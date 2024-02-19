package main

import (
	"encoding/json"
	"fmt"
	"os"
	"terraform/internal/plans/planfile"
)

func main() {
	planFile := os.Args[1]
	fmt.Printf("Loading: %s\n", planFile)

	readPlan(planFile)
	patchPlan()
	writePlan()
}

func readPlan(planFile string) {
	r, err := planfile.Open(planFile)
	if err != nil {
		panic(err)
	}

	plan, err := r.ReadPlan()
	if err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(plan, "", "  ")
	fmt.Printf("Plan:\n%s\n", j)
}


func patchPlan() {

}

func writePlan() {

}