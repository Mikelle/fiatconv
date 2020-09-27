package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/fiatconv/commands"
)

func main() {
	app := cli.NewApp()
	app.Name = "fiatconv"
	app.Commands = []cli.Command{
		commands.NewFiatconvCmd(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

