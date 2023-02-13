package main

import (
	"fmt"
	"os"

	"github.com/mgjules/batimag/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "Batimag"
	app.HelpName = "batimag"
	app.Usage = "Batch Image Processor!"
	app.Description = "Batimag applies a set of image processing functions to images in a given directory recursively."
	app.Authors = []*cli.Author{
		{
			Name:  "Michaël Giovanni Jules",
			Email: "julesmichaelgiovanni@gmail.com",
		},
	}
	app.Copyright = "(c) 2023 Michaël Giovanni Jules"
	app.Commands = cmd.Commands

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("failed to execute cmd: %v\n", err)
		os.Exit(1)
	}
}
