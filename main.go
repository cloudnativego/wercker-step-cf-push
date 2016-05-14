package main

import (
	"flag"
	"fmt"
	"os"
)

//TODO: accept env from wercker container
//TODO: prefer flags over env
//TODO: validate parameters and supply meaningful output
//TODO: verify cf command exists, if not, install cf

const notSupplied = "<not-supplied>"

var (
	api     string
	usr     string
	pwd     string
	org     string
	spc     string
	appname string
	errors  []string
)

func init() {
	flag.StringVar(&api, "api", notSupplied, "Target CF API URL. Override with CF_API")
}

func main() {
	flag.Parse()
	errors = make([]string, 0)

	api = reconcileWithEnvironment(api, "CF_API")

	if len(errors) > 0 {
		for _, v := range errors {
			fmt.Println(v)
		}
		os.Exit(1)
	}

	fmt.Printf("API: %s\nUSR: %s\nPWD: %s\nORG: %s\nSPC: %s\n", api, usr, pwd, org, spc)

	//	api = os.Getenv("WERCKER_CF_API")
	//	usr = os.Getenv("WERCKER_CF_USERNAME")
	//	pwd = os.Getenv("WERCKER_CF_PASSWORD")
	//	org = os.Getenv("WERCKER_CF_ORG")
	//	spc = os.Getenv("WERCKER_CF_SPACE")
	//	appname = os.Getenv("WERCKER_CF_APP_NAME")
}

func reconcileWithEnvironment(orig string, envName string) (result string) {
	result = orig
	if orig == notSupplied {
		result = os.Getenv(envName)
	}
	if len(result) == 0 {
		errors = append(errors, fmt.Sprintf("%s not supplied via flag or environment.", envName))
	}
	return
}
