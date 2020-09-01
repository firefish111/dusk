
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
		fmt.Fprintf(os.Stderr, "dusk usage: %s command package [more packages]\n", os.Args[0])
		os.Exit(1)
	} else if os.Args[1] == "help" {
		fmt.Printf("Usage: %s command package [more packages]\nExample:\n\n\t%s add myPkg\nThe above installs myPkg.\n\n\t%s del myPkg\nThe above deletes myPkg.\n\n\t%s inf myPkg\nThe above shows information about myPkg.\n", os.Args[0], os.Args[0], os.Args[0], os.Args[0])
		os.Exit(0)
	}
	for _, pkg := range os.Args[2:] {
		client := &http.Client{}

    res, err := client.Get("https://duskcdn.firefish.repl.co/cdn/"+pkg)
    req, err := http.NewRequest("GET", "https://duskcdn.firefish.repl.co/cdn/"+pkg, nil)
    safe(err)

    req.Header.Add("X-Requested-With", "night-dusk-pkg")
    res, err = client.Do(req)
    safe(err)

    body, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    safe(err)
    f, err := os.Create("./dusk_modules/"+pkg+".night")
    safe(err)
    f.Close()
    err = ioutil.WriteFile("./dusk_modules/"+pkg+".night", body, 0644)
    safe(err)

    fmt.Printf("%s", body)
	}
}
