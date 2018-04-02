package main

import (
    "os"
    "fmt"
    "encoding/json"
)

func main() {
    var request map[string]interface{}

    {
        err := json.NewDecoder(os.Stdin).Decode(&request)
        if err != nil {
            fatal("parsing request", err)
        }
    }

    response := []interface{}{}

    if request["version"] != nil {
        response = append(response, request["version"])
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
