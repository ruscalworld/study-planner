package main

import (
	"fmt"
	"os"

	"study-planner/internal/server"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "Study Planner Server",
		Usage: "Server for Study Planner providing RESTful HTTP API that is used by browser app",

		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "Starts the server",
				Action: server.RunApp,

				Flags: []cli.Flag{
					server.FlagBindAddress,
					server.FlagDatabaseHost,
					server.FlagDatabaseName,
					server.FlagDatabaseUser,
					server.FlagDatabasePassword,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
}
