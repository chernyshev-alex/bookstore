package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chernyshev-alex/bookstore/cmd/bookstore_items-api/config"
	"github.com/ilyakaznacheev/cleanenv"
)

type Args struct {
	ConfigPath string
}

func ProcessArgs(conf config.Config) Args {

	var args Args

	flags := flag.NewFlagSet("bookstore users api", 1)
	flags.StringVar(&args.ConfigPath, "c", "config/config.yml", "Path to config file")

	fu := flags.Usage
	flags.Usage = func() {
		fu()
		help, _ := cleanenv.GetDescription(conf, nil)
		fmt.Fprintln(flags.Output())
		fmt.Fprintln(flags.Output(), help)
	}
	flags.Parse(os.Args[1:])
	return args
}

func main() {

	args := ProcessArgs(config.Config{})
	conf, err := config.LoadCondigFromFile(args.ConfigPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("starting with config ", args.ConfigPath)

	app := inject(conf)
	app.StartApp()
}
