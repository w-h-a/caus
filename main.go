package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/w-h-a/caus/cmd"
)

func main() {
	app := &cli.App{
		Name:  "caus",
		Usage: "Causal discovery for your metrics and trace aggregates",
		Commands: []*cli.Command{
			{
				Name: "discover",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "vars",
						Aliases:  []string{"v"},
						Usage:    "Path to vars.yml config",
						Required: true,
					},
					&cli.DurationFlag{
						Name:    "start",
						Aliases: []string{"s"},
						Usage:   "How long ago to start (e.g., '2h', '30m')",
						Value:   2 * time.Hour,
					},
					&cli.DurationFlag{
						Name:    "end",
						Aliases: []string{"e"},
						Usage:   "How long ago to end (e.g., '0m' for now)",
						Value:   0,
					},
					&cli.DurationFlag{
						Name:  "step",
						Usage: "Data resolution (e.g., 1m, 15s)",
						Value: time.Minute,
					},
					&cli.IntFlag{
						Name:  "lag",
						Usage: "Max causal lag to check",
						Value: 3,
					},
					&cli.Float64Flag{
						Name:  "alpha",
						Usage: "Significance level (e.g., 0.05)",
						Value: 0.05,
					},
				},
				Action: cmd.Run,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
