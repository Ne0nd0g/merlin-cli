package main

import (
	// Standard
	"flag"
	"fmt"

	// Internal
	"github.com/Ne0nd0g/merlin-cli/services/cli"
	"github.com/Ne0nd0g/merlin-cli/version"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:50051", "The address of the Merlin server to connect to")
	password := flag.String("password", "merlin", "the password to connect to the Merlin server")
	secure := flag.Bool("secure", false, "Require server TLS certificate verification")
	tlsKey := flag.String("tlsKey", "", "TLS private key file path")
	tlsCert := flag.String("tlsCert", "", "TLS certificate file path")
	tlsCA := flag.String("tlsCA", "", "TLS Certificate Authority file path")
	v := flag.Bool("version", false, "Print the version number and exit")
	flag.Parse()

	if *v {
		fmt.Printf("Merlin Version: %s, Build: %s\n", version.Version, version.Build)
		return
	}

	// Start Merlin Command Line Interface
	cliService := cli.NewCLIService(*password, *secure, *tlsKey, *tlsCert, *tlsCA)
	cliService.Run(*addr)
}
