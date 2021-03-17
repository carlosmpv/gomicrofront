package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"plugin"
	"text/template"
	"time"
)

func getAvailablePlugins() ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir("./plugins")
	if err != nil {
		return nil, err
	}

	return files, nil
}

func loadPlugins(mux *http.ServeMux) error {
	files, err := getAvailablePlugins()
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			log.Println("Loaded plugins: ", file.Name())
			pluginPath := path.Join("./plugins", file.Name(), fmt.Sprintf("%s.so", file.Name()))

			plugin, err := plugin.Open(pluginPath)
			if err != nil {
				return err
			}

			pluginEntryPoint, err := plugin.Lookup("OnRequest")
			if err != nil {
				return err
			}

			mux.HandleFunc(fmt.Sprintf("/plugins/%s", file.Name()), pluginEntryPoint.(func(http.ResponseWriter, *http.Request)))
		}
	}

	return nil
}

func serve() {
	mux := http.NewServeMux()

	err := loadPlugins(mux)
	if err != nil {
		log.Fatal(err)
	}

	entry, err := plugin.Open("./server/server.so")
	if err != nil {
		log.Fatal(err)
	}

	initFunc, err := entry.Lookup("Init")
	if err != nil {
		log.Fatal(err)
	}

	routerFunc, err := entry.Lookup("Router")
	if err != nil {
		log.Fatal(err)
	}

	entryPoint, err := entry.Lookup("Entrypoint")
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/", http.FileServer(http.Dir("./public")))
	routerFunc.(func(*http.ServeMux))(mux)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,

		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			code, err := entryPoint.(func(w http.ResponseWriter, r *http.Request) (int, error))(w, r)

			if err != nil {
				http.Error(w, err.Error(), code)
				return
			}

			mux.ServeHTTP(w, r)
		}),
	}

	initFunc.(func())()
	log.Fatal(server.ListenAndServe())
}

func newPlugin(name string) {
	err := os.Mkdir(path.Join("./plugins", name), 0700)
	if err != nil {
		log.Fatal(err)
	}

	tplBd, err := template.ParseFiles("./templates/body.js.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	tplSl, err := template.ParseFiles("./templates/soul.go.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	bdFile, err := os.Create(fmt.Sprintf("./plugins/%s/body.js", name))
	if err != nil {
		log.Fatal(err)
	}

	slFile, err := os.Create(fmt.Sprintf("./plugins/%s/soul.go", name))
	if err != nil {
		log.Fatal(err)
	}

	data := struct{ Name string }{name}

	err = tplBd.Execute(bdFile, data)
	if err != nil {
		log.Fatal(err)
	}

	err = tplSl.Execute(slFile, data)
	if err != nil {
		log.Fatal(err)
	}
}

func buildPlugins() {
	files, err := getAvailablePlugins()
	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range files {
		cmd := exec.Command("go", "build", "-buildmode=plugin")
		cmd.Dir = path.Join("./plugins", folder.Name())

		result, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}

		log.Println(result)
	}
}

func buildServer() {
	cmd := exec.Command("go", "build", "-buildmode=plugin")
	cmd.Dir = "./server"

	result, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}

func createPluginsFolder() {
	err := os.Mkdir("plugins", 0700)
	if err != nil {
		log.Println(err)
	}
}

func createPublicFolder() {
	err := os.Mkdir("public", 0700)
	if err != nil {
		log.Println(err)
	}

	tplIn, err := template.ParseFiles("./templates/index.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	idSv, err := os.Create("./public/index.html")
	if err != nil {
		log.Fatal(err)
	}

	tplCr, err := template.ParseFiles("./templates/core.js.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	crSv, err := os.Create("./public/core.js")
	if err != nil {
		log.Fatal(err)
	}

	err = tplIn.Execute(idSv, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = tplCr.Execute(crSv, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func createServer() {
	err := os.Mkdir("server", 0700)
	if err != nil {
		log.Println(err)
	}

	tplSv, err := template.ParseFiles("./templates/entry.go.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	svFile, err := os.Create("./server/entry.go")
	if err != nil {
		log.Fatal(err)
	}

	err = tplSv.Execute(svFile, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s [serve|new-plugin]", os.Args[0])
	}

	if os.Args[1] == "serve" {
		serve()
		return
	}

	if os.Args[1] == "new-plugin" {
		if len(os.Args) < 3 {
			log.Fatalf("Usage: %s new-plugin <name:string>", os.Args[0])
		}

		newPlugin(os.Args[2])
		return
	}

	if os.Args[1] == "build" {
		buildPlugins()
		buildServer()
		return
	}

	if os.Args[1] == "init" {
		createServer()
		createPluginsFolder()
		createPublicFolder()
		return
	}
}
