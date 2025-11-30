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
						Value:   5 * time.Minute,
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
					&cli.BoolFlag{
						Name:  "json",
						Usage: "Print the resulting graph to stdout as json",
						Value: false,
					},
				},
				Action: cmd.Discover,
			},
			{
				Name: "simulate",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "graph",
						Aliases:  []string{"g"},
						Usage:    "Path to graph.json",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "do",
						Aliases:  []string{"d"},
						Usage:    "Intervention on a variable (e.g., 'front_calls * 1.2' or 'front_calls = 200')",
						Required: true,
					},
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
						Value:   5 * time.Minute,
					},
					&cli.DurationFlag{
						Name:  "step",
						Usage: "Data resolution (e.g., 1m, 15s)",
						Value: time.Minute,
					},
					&cli.IntFlag{
						Name:  "horizon",
						Usage: "Number of steps for the counterfactual simulation (e.g., if step is 1m and horizon is 60, you will replace the last 60mins with a counterfactual history)",
						Value: 60,
					},
					&cli.StringFlag{
						Name:  "effect",
						Usage: "Variable on which to focus counterfactual predictions (e.g., 'orders_cpu')",
					},
				},
				Action: cmd.Simulate,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
