package main

import (
	"flag"

	"github.com/MattWindsor91/act-tester/internal/pkg/director"
)

const usagePlanFile = "Read from this plan `file` instead of stdin."

// direct is the Director being built and run by this command.
var direct director.Director

func init() {
	flag.StringVar(&direct.PlanFile, "i", "", usagePlanFile)
}

func main() {
	flag.Parse()

}
