package main

import (
    "os"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "path/filepath"
    "github.com/jleben/slack-chat-resource/utils"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
    }
    
    destination := os.Args[1]

    var request utils.InRequest

    err := json.NewDecoder(os.Stdin).Decode(&request)
    if err != nil {
        fatal("parsing request", err)
    }

    var response utils.InResponse
    response.Version = request.Version

    {
        fmt.Fprintf(os.Stderr, "The Latest versions is:\n")
        fmt.Fprintf(os.Stderr, "%s\n", request.Version["timestamp"])
    }

    {
        err := ioutil.WriteFile(filepath.Join(destination, "timestamp"), []byte(request.Version["timestamp"]) , 0644)
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
