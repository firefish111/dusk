package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func safe(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	os.Mkdir("./dusk_pkgs", 0755)
	// metadata things
	if _, err := os.Stat("./dusk_pkgs/metadata.json"); os.IsNotExist(err) {
		ioutil.WriteFile("./dusk_pkgs/metadata.json", []byte("{}"), 0666)
	}
	meta := make(map[string]interface{})
	dat, err := ioutil.ReadFile("./dusk_pkgs/metadata.json")
	safe(err)
	err = json.Unmarshal([]byte(dat), &meta)
	safe(err)
  // help menu
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "dusk usage: %s command package [more packages]\n", os.Args[0]) // incorrectly formatted
		os.Exit(1)
	} else if os.Args[1] == "help" {
		fmt.Println("Usage: %s command package [more packages]\nExample:\n\n\tdusk add myPkg\nThe above installs myPkg.\n\n\tdusk del myPkg\nThe above deletes myPkg.\n\n\tdusk upd myPkg\nThe above updates myPkg.") // prints help
		os.Exit(0)
	}
	for _, pkg := range os.Args[2:] { // iterates over all the packages passed
		if os.Args[1] != "del" {
			client := &http.Client{} // creates client

			res, err := client.Get("https://duskcdn.firefish.repl.co/cdn/" + pkg)
			req, err := http.NewRequest("GET", "https://duskcdn.firefish.repl.co/cdn/"+pkg, nil) // initialise request
			safe(err)

			req.Header.Add("X-Requested-With", "night-dusk-pkg")
			res, err = client.Do(req) // send request
			safe(err)
			if res.StatusCode != 200 {
				panic(fmt.Sprintf("Status Code %d", res.StatusCode))
				os.Exit(1)
			}

			body, err := ioutil.ReadAll(res.Body) // read body of response
			res.Body.Close()
			safe(err)

			if os.Args[1] == "add" {
				if _, err = os.Stat("./dusk_pkgs/" + pkg + ".night"); !os.IsNotExist(err) {
					panic("File already exists, please use \x1b[38;5;155mdusk upd **[packages]**\x1b[0m")
				}

				err = ioutil.WriteFile("./dusk_pkgs/"+pkg+".night", body, 0666)
				safe(err) // write to night file

				fmt.Printf("\x1b[1m\x1b[38;5;164mInstalled package \x1b[38;5;202m%s \x1b[38;5;155mv%s\x1b[0m\n", pkg, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(res.Header["X-Package-Version"])), "."), "[]"))

			} else if os.Args[1] == "upd" {
				if _, err = os.Stat("./dusk_pkgs/" + pkg + ".night"); os.IsNotExist(err) {
					fmt.Fprintln(os.Stderr, "\x1b[1m\x1b[38;5;9mWarning: destination file doesn't exist, installing package standalone")
				}

				err = ioutil.WriteFile("./dusk_pkgs/"+pkg+".night", body, 0666)
				safe(err) // write to night file

				fmt.Printf("\x1b[1m\x1b[38;5;164mUpdated package \x1b[38;5;202m%s \x1b[38;5;155mto v%s\x1b[0m\n", pkg, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(res.Header["X-Package-Version"])), "."), "[]"))

			}
		} else {
			err := os.Remove("./dusk_pkgs/" + pkg + ".night") // delete night file
			safe(err)
			fmt.Printf("\x1b[1m\x1b[38;5;164mUninstalled package \x1b[38;5;202m%s \x1b[38;5;155m\n\x1b[0m", pkg)
		}
	}
}
