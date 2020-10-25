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
	os.Mkdir("./pkgs", 0755)
	// metadata things
	if _, err := os.Stat("./pkgs/metadata.json"); os.IsNotExist(err) {
		ioutil.WriteFile("./pkgs/metadata.json", []byte("{}"), 0666)
	}
	meta := make(map[string]map[string]interface{})
	dat, err := ioutil.ReadFile("./pkgs/metadata.json")
	safe(err)
	err = json.Unmarshal([]byte(dat), &meta)
	safe(err)
	// usage thingy
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "usage: dusk command package [more packages]")
		os.Exit(1)
	} else if os.Args[1] == "ls" {
		if len(meta) == 0 {
			fmt.Println("\033[38;5;63mUnfortunately, no packages are installed.\033[0m")
		} else {
			for ky := range meta {
				fmt.Println("\033[38;5;63mInstalled packages are:\033[0m")
				fmt.Printf("\033[1m\033[38;5;164mPackage \033[38;5;202m%s \033[38;5;155mv%s\033[0m\n", ky, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(meta[ky]["version"])), "."), "[]"))
			}
		}
	} else if len(os.Args) <= 2 {
		fmt.Fprintln(os.Stderr, "usage: dusk command package [more packages]")
		os.Exit(1)
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
			switch os.Args[1] {
			case "upd":
				if _, err = os.Stat("./pkgs/" + pkg + ".night"); !os.IsNotExist(err) {

					err = ioutil.WriteFile("./pkgs/"+pkg+".night", body, 0666)
					safe(err) // write to night file

					fmt.Printf("\033[1m\033[38;5;164mUpdated package \033[38;5;202m%s \033[38;5;155mv%s to v%s\033[0m\n", pkg, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(meta[pkg]["version"])), "."), "[]"), strings.Trim(strings.Join(strings.Fields(fmt.Sprint(res.Header["X-Package-Version"])), "."), "[]"))
					meta[pkg]["version"] = res.Header["X-Package-Version"]
					break
				}
				fmt.Fprintln(os.Stderr, "\033[38;5;9mWarning: destination file doesn't exist, installing package standalone")
				fallthrough
			case "add":
				if _, err = os.Stat("./pkgs/" + pkg + ".night"); !os.IsNotExist(err) {
					panic("File already exists, please use \033[38;5;155mdusk upd **[packages]**\033[0m")
				}

				err = ioutil.WriteFile("./pkgs/"+pkg+".night", body, 0666)
				safe(err) // write to night file

				fmt.Printf("\033[1m\033[38;5;164mInstalled package \033[38;5;202m%s \033[38;5;155mv%s\033[0m\n", pkg, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(res.Header["X-Package-Version"])), "."), "[]"))

				meta[pkg] = make(map[string]interface{})
				meta[pkg]["version"] = res.Header["X-Package-Version"]
			case "inf":
				fmt.Printf("\033[1m\033[38;5;164mPackage \033[38;5;202m%s \033[38;5;155mv%s\033[0m\n", pkg, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(res.Header["X-Package-Version"])), "."), "[]"))
			}
		} else {
			err := os.Remove("./pkgs/" + pkg + ".night") // delete night file
			safe(err)
			fmt.Printf("\033[1m\033[38;5;164mUninstalled package \033[38;5;202m%s \033[38;5;155mv%s\033[0m\n", pkg, strings.Trim(strings.Join(strings.Fields(fmt.Sprint(meta[pkg]["version"])), "."), "[]"))
			delete(meta, pkg)
		}
	}
	dat, err = json.Marshal(meta)
	safe(err)
	err = ioutil.WriteFile("./pkgs/metadata.json", dat, 0666)
	safe(err)
}
