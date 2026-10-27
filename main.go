package main

import (
	"os"
	"time"

	"github.com/kishan-thanki/whoisusing/v2/internal/app"
	"github.com/kishan-thanki/whoisusing/v2/internal/lookup"
	"github.com/kishan-thanki/whoisusing/v2/internal/process"
)

func main() {
	runner := app.New(
		lookup.NewLsof(),
		process.NewManager(),
		time.Sleep,
	)

	os.Exit(runner.Run(os.Args[1:], os.Stdout, os.Stderr))
}
