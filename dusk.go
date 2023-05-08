package main
import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "io/fs"
  "net/http"
  "os"
  "strings"
  "errors"

  // haha chinese package go brrrrrrr
  "github.com/gookit/color"
)

// safety go brrrrrr
func safe(e error) {
  if e != nil {
    panic(e)
  }
}

// long boi function chain
func vparse(ver interface{}) string {
  body := strings.Fields(fmt.Sprint(ver))

  ver_k := strings.Join(body[:3], ".")
  if len(body) > 3 {
    ver_k += "-" + body[3];
  }
  return strings.Trim(ver_k, "[]");
}


func main() {
  meta := make(map[string]map[string]interface{}) 
  // not paying attention to error, because am smort
  os.Mkdir("./pkgs", 0755)

  // haha colours go brrrrrr
  indigo := color.C256(63, false)
  lime := color.C256(155, false)
  pink := color.C256(164, false)
  orange := color.C256(202, false)
  bold := color.Style{color.Bold}
  rederr := color.Style{color.FgRed}

  // if not exist
   if _, err := os.Stat("./pkgs/metadata.json"); !errors.Is(err, fs.ErrNotExist) {
    // readfile
    dat, err := ioutil.ReadFile("./pkgs/metadata.json")
    safe(err)
    // put contents of file in meta, parsed as json
    err = json.Unmarshal([]byte(dat), &meta)
    safe(err)
  }
   
  if len(os.Args) == 1 || (len(os.Args) <= 2 && os.Args[1] != "ls") {
    // usage
    fmt.Fprintln(os.Stderr, "usage: dusk command package [more packages]")
    os.Exit(1)
  } else if os.Args[1] == "ls" {
    // listing packages
    if len(meta) == 0 {
      indigo.Println("Unfortunately, no packages are installed.")
    } else {
      indigo.Println("Installed packages are:")
      for ky := range meta {
        bold.Print(pink.Sprint("Package "))
        orange.Printf("%s ", ky)
        lime.Printf("v%s\n", vparse(meta[ky]["version"]))
      }
    }
    os.Exit(0)
  }
  for _, pkg := range os.Args[2:] { // iterates over all the packages passed
    if os.Args[1] == "del" {
      err := os.Remove("./pkgs/" + pkg + ".night") // delete night file
      safe(err)
      bold.Print(pink.Sprint("Uninstalled package "))
      orange.Printf("%s ", pkg)
      lime.Printf("v%s\n", vparse(meta[pkg]["version"]))

      delete(meta, pkg)
    } else {
      client := &http.Client{} // creates client

      // res, err := client.Get("https://duskcdn.firefish.repl.co/cdn/" + pkg)
      var method string
      if method = "GET"; os.Args[1] == "inf" {
        method = "POST"
      }
      req, err := http.NewRequest(method, fmt.Sprintf(
        "https://duskcdn.firefish.repl.co/cdn/%s", pkg,
      ), nil) // initialise request
      safe(err)

      req.Header.Add("X-Requested-With", "night-dusk-pkg")
      res, err := client.Do(req) // send request
      safe(err)
      if res.StatusCode != 200 {
        panic(fmt.Sprintf("Status Code %d", res.StatusCode))
        os.Exit(1)
      }

      if os.Args[1] != "inf" {
        body, err := ioutil.ReadAll(res.Body) // read body of response
        res.Body.Close()
        safe(err)
        switch os.Args[1] {
        case "upd":
          if _, err = os.Stat("./pkgs/" + pkg + ".night"); !os.IsNotExist(err) {

            err = ioutil.WriteFile("./pkgs/"+pkg+".night", body, 0666)
            safe(err) // write to night file

            if vparse(res.Header["X-Package-Version"]) == vparse(meta[pkg]["version"]) {
              fmt.Fprintln(os.Stderr, rederr.Sprintf("Package %s already at latest version.", pkg))
              os.Exit(0)
            }

            bold.Print(pink.Sprint("Updated package "))
            orange.Printf("%s ", pkg)
            lime.Printf("v%s to v%s\n",
              vparse(meta[pkg]["version"]),
              vparse(res.Header["X-Package-Version"]),
            )

            meta[pkg]["version"] = res.Header["X-Package-Version"]
            break
          }
          fmt.Fprintln(os.Stderr, rederr.Sprint("Warning: destination file doesn't exist, installing package instead"))
          fallthrough
        case "add":
          if _, err = os.Stat("./pkgs/" + pkg + ".night"); !os.IsNotExist(err) {
            panic(fmt.Sprintf("File already exists, please use %s",
              lime.Sprint("dusk upd **[packages]**"),
            ))
          }

          err = ioutil.WriteFile("./pkgs/"+pkg+".night", body, 0666)
          safe(err) // write to night file

          bold.Print(pink.Sprint("Installed package "))
          orange.Printf("%s ", pkg)
          lime.Printf("v%s\n", vparse(res.Header["X-Package-Version"]))

          meta[pkg] = make(map[string]interface{})
          meta[pkg]["version"] = res.Header["X-Package-Version"]
        }
      } else {
        bold.Print(pink.Sprint("Package "))
        orange.Printf("%s ", pkg)
        lime.Printf("v%s\n", vparse(res.Header["X-Package-Version"]))
      }
    }
  }
  dat, err := json.Marshal(meta)
  safe(err)
  err = ioutil.WriteFile("./pkgs/metadata.json", dat, 0666)
  safe(err)
}
