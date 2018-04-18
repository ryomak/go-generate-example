package main

import (
	"github.com/BurntSushi/toml"
	"os"
	"html/template"
	"strings"
)

type Config struct {
	API []APIConfig
	NotFound NotFoundConfig
}

type APIConfig struct {
	Name string
	Method string
	EndPoint string
	Params []string
}

type NotFoundConfig struct {
	Name string
}

const(
	input = "api.toml"
	output = "../api_route.go"
)

//go:generate go run api.go
func main(){
	var config Config
	_, err := toml.DecodeFile(input, &config)
	if err != nil {
		panic(err)
	}
	file ,err := os.OpenFile(output,os.O_WRONLY|os.O_CREATE,0666)
	if err != nil{
		panic(err)
	}
	t := template.New(output)
	_,err = t.Parse(RouteSoure)
	if err != nil{
		panic(err)
	}
	err = t.Execute(file,config)
	if err != nil{
		panic(err)
	}

	//controller
	for _,api := range config.API{
		filePath := "../../controller/"+strings.ToLower(api.Name)+"_controller.go"
		file ,err := os.OpenFile(filePath,os.O_WRONLY|os.O_CREATE,0666)
		if err != nil{
			panic(err)
		}
		t,err = template.New(filePath).Parse(controllerSource)
		if err != nil{
			panic(err)
		}
		err = t.Execute(file,api)
		if err != nil{
			panic(err)
		}
	}

	//notFound
	filePath := "../../controller/"+strings.ToLower(config.NotFound.Name)+"_controller.go"
	file ,err = os.OpenFile(filePath,os.O_WRONLY|os.O_CREATE,0666)
	if err != nil{
		panic(err)
	}
	t,err = template.New(filePath).Parse(controllerSource)
	if err != nil{
		panic(err)
	}
	err = t.Execute(file,config.NotFound)
	if err != nil{
		panic(err)
	}

}

const RouteSoure = `
package api

import (
	"github.com/gorilla/mux"
	"github.com/ryomak/go-generate-example/src/controller"
	"github.com/urfave/negroni"
	"log"
	"net/http"
)

type Router struct {
	*mux.Router
}

func New() *Router {
	return &Router{
		mux.NewRouter(),
	}
}

func (r *Router) Route(addr string) {
	a := r.PathPrefix("/api").Subrouter()
	r.NotFoundHandler = http.HandlerFunc(controller.NotFound)
	{{- range .API}}
	{{- if .Params}}
	a.Path("{{.EndPoint}}").Queries(
		{{- range .Params}}
			"{{. -}}","{{.}}",
		{{- end}}
	).HandlerFunc(controller.{{.Name}}).Methods("{{.Method}}")
	{{- else}}
	a.Path("{{.EndPoint}}").HandlerFunc(controller.{{.Name}}).Methods("{{.Method}}")
	{{- end}}
	{{- end}}
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(r)
	log.Println("server start :" + addr)
	http.ListenAndServe(":"+addr, n)
}
`

const controllerSource = `
package controller
import (
	"net/http"
)

func {{.Name}}(w http.ResponseWriter,r *http.Request){
	
}
`