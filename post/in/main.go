package main

import (
    "os"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "path/filepath"
)

func main() {
    var request map[string]interface{}

    destination := os.Args[1]

    {
        err := json.NewDecoder(os.Stdin).Decode(&request)
        if err != nil {
            fatal("parsing request", err)
        }
    }

    response := make(map[string]interface{})
    response["version"] = request["version"]

    {
        err := ioutil.WriteFile(filepath.Join(destination, "timestamp"), []byte(request["version"].(string)) , 0644)
        if err != nil {
            fatal("writing timestamp file", err)
        }
    }

    {
        err := json.NewEncoder(os.Stdout).Encode(&response)
        if err != nil {
            fatal("serializing response", err)
        }
    }
}

func fatal(doing string, err error) {
    fmt.Fprintf(os.Stderr, "Error " + doing + ": " + err.Error() + "\n")
    os.Exit(1)
}
