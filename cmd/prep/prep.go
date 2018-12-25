package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/alexvanboxel/reactor/pkg/chem"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

var (
	Project string
)

type Info struct {
	Module    string
	Component string
	Project   string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func exists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}
	return false
}

func writeDeploy(deployDir string, tmpl *template.Template, deploy Info) {
	f, err := os.Create(fmt.Sprintf("%s/%s-%s.yaml", deployDir, deploy.Module, deploy.Component))
	check(err)
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	err = tmpl.Execute(f, deploy)
	if err != nil {
		panic(err)
	}

}

func writeService(deployDir string, tmpl *template.Template, deploy Info) {
	f, err := os.Create(fmt.Sprintf("%s/%s-%s.yaml", deployDir, deploy.Module, deploy.Component))
	check(err)
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	err = tmpl.Execute(f, deploy)
	if err != nil {
		panic(err)
	}

}

func sa() {
	deployDir := "work/" + Project + "/misc"
	_ = os.Mkdir(deployDir, os.ModePerm)

	content, err := ioutil.ReadFile("k8s/sa-template.yaml")
	check(err)
	tmpl, err := template.New("sa").Parse(string(content))
	check(err)

	sa, err := ioutil.ReadFile("work/" + Project + "/service-account.json")
	check(err)
	sab64 := base64.StdEncoding.EncodeToString(sa)

	type Info struct {
		ServiceAccount string
	}

	f, err := os.Create(fmt.Sprintf("%s/reactor-service-account.yaml", deployDir))
	check(err)
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	deploy := Info{
		ServiceAccount: sab64,
	}
	err = tmpl.Execute(f, deploy)
	check(err)

}

func deploy(atoms *chem.Atoms) {

	deployDir := "work/" + Project + "/deploy"
	_ = os.Mkdir(deployDir, os.ModePerm)
	serviceDir := "work/" + Project + "/service"
	_ = os.Mkdir(serviceDir, os.ModePerm)

	content, err := ioutil.ReadFile("k8s/deploy-template.yaml")
	check(err)
	tmpl, err := template.New("deploy").Parse(string(content))
	check(err)

	content, err = ioutil.ReadFile("k8s/service-template.yaml")
	check(err)
	serviceTemplate, err := template.New("service").Parse(string(content))
	check(err)

	for k, v := range atoms.Symbols {
		if v.Period > 0 && v.Period <= 2 {
			deploy := Info{
				Module:    "atom",
				Component: strings.ToLower(k),
				Project:   Project,
			}
			writeDeploy(deployDir, tmpl, deploy)
			writeService(serviceDir, serviceTemplate, deploy)
		}
	}

	deploy := Info{
		Module:    "orbit",
		Component: "inf",
		Project:   Project,
	}
	writeDeploy(deployDir, tmpl, deploy)
	writeService(serviceDir, serviceTemplate, deploy)
	for i := 1; i <= 5; i++ {
		deploy := Info{
			Module:    "orbit",
			Component: strconv.Itoa(i),
			Project:   Project,
		}
		writeDeploy(deployDir, tmpl, deploy)
		writeService(serviceDir, serviceTemplate, deploy)
	}

	deploy = Info{
		Module:    "reactor",
		Component: "api",
		Project:   Project,
	}
	writeDeploy(deployDir, tmpl, deploy)
	writeService(serviceDir, serviceTemplate, deploy)
}

func main() {
	fmt.Println("prep started.")
	flag.StringVar(&Project, "project", "", "df")
	flag.Parse()

	if !exists("work/" + Project) {
		log.Fatal("work/" + Project + " doesn't exist")
	}

	atoms := chem.NewAtoms()

	deploy(atoms)
	sa()
}
