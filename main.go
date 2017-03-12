package main

import (
	"log"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/fatih/color"
	"github.com/piot/hasty-repl/commander"
	"github.com/piot/hasty-repl/connection"
	"github.com/piot/hasty-repl/repl"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	cacert  = kingpin.Flag("cacert", "A .pem-encoded Certificate Authority for TSL connections. Mostly used for self signed certificates.").String()
	host    = kingpin.Flag("server", "Hasty Server or Load Balancer").Default("localhost:3333").String()
)

func boot() error {
	kingpin.Parse()

	conn, connErr := connection.NewConnection(*host, *cacert)
	if connErr != nil {
		return connErr
	}
	commander := commander.NewCommander(conn)
	repl, replErr := repl.NewRepl(&commander)
	if replErr != nil {
		return replErr
	}

	promptErr := repl.PromptForever()
	return promptErr
}

func main() {
	color.Cyan("HastyRepl v0.1")
	err := boot()
	if err != nil {
		log.Printf("Error: %s", err)
	}
}
