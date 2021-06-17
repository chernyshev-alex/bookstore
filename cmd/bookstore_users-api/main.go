package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Args struct {
	ConfigPath string
}

func ProcessArgs(conf conf.Config) Args {
	var args Args

	flags := flag.NewFlagSet("bookstore users api", 1)
	flags.StringVar(&args.ConfigPath, "c", "conf/app.yml", "Path to config file")

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
	args := ProcessArgs(conf.Config{})
	conf, err := conf.LoadCondigFromFile(args.ConfigPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println("starting with config ", args.ConfigPath)

	mysqlConf := users_db.MakeMySQLConfig(*conf)

	app := inject(mysqlConf)
	app.StartApp()
}
