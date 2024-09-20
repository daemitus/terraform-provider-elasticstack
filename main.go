package main

import (
	"context"
	"flag"
	"log"

	"github.com/daemitus/terraform-provider-elasticstack/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
)

var version string = "dev"

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	serverFactory, err := provider.ProtoV6ProviderServerFactory(context.Background(), version)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt
	if debugMode {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/daemitus/elasticstack",
		serverFactory,
		serveOpts...,
	)
	if err != nil {
		log.Fatal(err)
	}
}
