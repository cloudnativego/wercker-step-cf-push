package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//TODO: accept env from wercker container
//TODO: prefer flags over env
//TODO: validate parameters and supply meaningful output
//TODO: verify cf command exists, if not, install cf

const notSupplied = "<not-supplied>"

var (
	api         string
	usr         string
	pwd         string
	org         string
	spc         string
	appname     string
	errors      []string
	dockerImage string
)

func init() {
	flag.StringVar(&api, "api", notSupplied, "Target CF API URL. Overrides API env variable")
	flag.StringVar(&usr, "user", notSupplied, "CF User. Overrides USER env variable ")
	flag.StringVar(&pwd, "password", notSupplied, "CF Password. Overrides PASSWORD env variable")
	flag.StringVar(&org, "org", notSupplied, "CF Org. Overrides ORG env variable")
	flag.StringVar(&spc, "space", notSupplied, "CF Space. Overrides SPACE env variable")
	flag.StringVar(&appname, "appname", notSupplied, "Name of application to be pushed to Cloud Foundry. Overrides APPNAME env variable")
	flag.StringVar(&dockerImage, "docker-image", notSupplied, "Optional. Path to docker image.  Overrides DOCKER_IMAGE env variable")
}

func main() {
	flag.Parse()
	errors = make([]string, 0)

	api = reconcileWithEnvironment(api, "WERCKER_CF_DEPLOY_API", true)
	usr = reconcileWithEnvironment(usr, "WERCKER_CF_DEPLOY_USER", true)
	pwd = reconcileWithEnvironment(pwd, "WERCKER_CF_DEPLOY_PASSWORD", true)
	org = reconcileWithEnvironment(org, "WERCKER_CF_DEPLOY_ORG", true)
	spc = reconcileWithEnvironment(spc, "WERCKER_CF_DEPLOY_SPACE", true)
	appname = reconcileWithEnvironment(appname, "WERCKER_CF_DEPLOY_APPNAME", true)

	if len(errors) > 0 {
		for _, v := range errors {
			fmt.Println(v)
		}
		os.Exit(1)
	}

	fmt.Println("Downloading and installing CF CLI...")
	if ok := installCF(); !ok {
		fmt.Println("Unable to install CF CLI.")
		os.Exit(1)
	}
	fmt.Println("CF CLI installed.")

	fmt.Println("Generating cf push command...")
	pushCommand := determinePushCommand()
	if len(pushCommand) == 0 {
		fmt.Println("Error: Push command not created.")
		fmt.Printf("API: %s\nUSR: %s\nPWD: %s\nORG: %s\nSPC: %s\n", api, usr, pwd, org, spc)
		os.Exit(1)
	}
	fmt.Printf("Generated Command: cf %v\n", pushCommand)

	fmt.Println("Deploying app...")
	deployCommand := exec.Command("./cf", pushCommand...)
	err := deployCommand.Run()
	if err != nil {
		fmt.Printf("ERROR OCCURRED: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("SUCCESS.\n")
	os.Exit(0)
}

func installCF() bool {
	downloadCommand := exec.Command("wget", "-O", "cf.tgz", "https://cli.run.pivotal.io/stable?release=linux64-binary")
	err := downloadCommand.Run()
	if err != nil {
		fmt.Printf("Error retrieving cf binary: %s", err)
		return false
	}

	unzipCommand := exec.Command("tar", "-zxf", "cf.tgz")
	err = unzipCommand.Run()
	if err != nil {
		fmt.Printf("Error unpacking cf binary: %s", err)
		return false
	}
	return true
}

func determinePushCommand() (cmd []string) {
	//Docker
	dockerImage = reconcileWithEnvironment(dockerImage, "WERCKER_CF_DEPLOY_DOCKER_IMAGE", false)

	if len(dockerImage) > 0 {
		commandString := fmt.Sprintf("push %s -o %s", appname, dockerImage)
		cmd = strings.Split(commandString, " ")
	}
	return
}

func reconcileWithEnvironment(orig string, envName string, required bool) (result string) {
	result = orig
	if orig == notSupplied {
		result = os.Getenv(envName)
	}
	if len(result) == 0 && required {
		errors = append(errors, fmt.Sprintf("%s not supplied via flag or environment.", envName))
	}
	return
}
