package main

import (
	"fmt"
	"io/ioutil"
  "net/http"
	"os"
)

func safe(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "dusk usage: %s command package [more packages]\n", os.Args[0]) // incorrectly formatted
		os.Exit(1)
	} else if os.Args[1] == "help" {
		fmt.Printf("Usage: %s command package [more packages]\nExample:\n\n\t%s add myPkg\nThe above installs myPkg.\n\n\t%s del myPkg\nThe above deletes myPkg.\n\n\t%s inf myPkg\nThe above shows information about myPkg.\n", os.Args[0], os.Args[0], os.Args[0], os.Args[0]) // prints help
		os.Exit(0)
	}
	for _, pkg := range os.Args[2:] { // iterates over all the packkages passed
		client := &http.Client{} // creates client

		res, err := client.Get("https://duskcdn.firefish.repl.co/cdn/" + pkg)
		req, err := http.NewRequest("GET", "https://duskcdn.firefish.repl.co/cdn/"+pkg, nil) // initialise request
		safe(err)

		req.Header.Add("X-Requested-With", "night-dusk-pkg")
		res, err = client.Do(req) // send request
		safe(err)

		body, err := ioutil.ReadAll(res.Body) // read body of response
		res.Body.Close()
		safe(err)

    fmt.Printf("\x1b[1m\x1b[38;5;164mInstalling \x1b[38;5;202m%s\x1b[38;5;155mv", pkg)
    for indx, err := range res.Header["X-Package-Version"] {
      fmt.Printf("%s", err)
      if indx != len(res.Header["X-Package-Version"])-1 {
        fmt.Print(".")
      }
    }
    fmt.Print("\x1b[0m\n")

		os.Mkdir("./dusk_modules", 0755) 
		err = ioutil.WriteFile("./dusk_modules/"+pkg+".night", body, 0666)
		safe(err) // create night file

		fmt.Printf("%s", body) // echo night file to console: debug feature
	}
}
