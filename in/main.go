package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"
    "fmt"
    //"strings"
    "github.com/jakobleben/slack-request-resource/protocol"
    "github.com/nlopes/slack"
)


// FIMXE: Pass params to target resource

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <destination>")
		os.Exit(1)
	}

    destination := os.Args[1]

    var request protocol.InRequest

    var err error

	err = json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("Parsing request.", err)
	}

	if len(request.Source.Token) == 0 {
        fatal1("Missing source field: token.")
    }

    if len(request.Source.ChannelId) == 0 {
        fatal1("Missing source field: channel_id.")
    }

	if _,ok := request.Version["request"]; !ok {
        fatal1("Missing version field: request")
    }

    fmt.Fprintf(os.Stderr,"request version: %v\n", request.Version["request"])

    slack_client := slack.New(request.Source.Token)

    response := get(&request, destination, slack_client)

    err = json.NewEncoder(os.Stdout).Encode(&response)
    if err != nil {
        fatal("encoding response", err)
    }
}

func get(request *protocol.InRequest, destination string, slack_client *slack.Client) protocol.InResponse {

    params := slack.NewHistoryParameters()
    params.Latest = request.Version["request"]
    params.Inclusive = true
    params.Count = 1

    history, history_err := slack_client.GetChannelHistory(request.Source.ChannelId, params)
    if history_err != nil {
		fatal("getting message", history_err)
	}

	if len(history.Messages) < 1 {
        fatal1("Message could not be found.")
    }

    message := history.Messages[0]

    fmt.Fprintf(os.Stderr, "Text: %s\n", message.Msg.Text)

    {
        err := os.MkdirAll(destination, 0755)
        if err != nil {
            fatal("creating destination directory", err)
        }
    }

    parts := []string {}

    if request.Params.TextPattern != nil {
        fmt.Fprintf(os.Stderr, "Pattern: %s\n", request.Params.TextPattern)
        parts = request.Params.TextPattern.FindStringSubmatch(message.Msg.Text)
    }

    {
        err := ioutil.WriteFile(filepath.Join(destination, "text"), []byte(message.Msg.Text), 0644)
        if err != nil {
            fatal("writing text file", err)
        }
    }

    for i := 1; i < len(parts); i++ {
        part := parts[i]
        fmt.Fprintf(os.Stderr, "Part: %s\n", part)
        filename := fmt.Sprintf("text_part%d", i)
        err := ioutil.WriteFile(filepath.Join(destination, filename), []byte(part), 0644)
        if err != nil {
            fatal("writing text part file", err)
        }
    }

    var response protocol.InResponse
    response.Version = request.Version
    return response
}

func fatal(doing string, err error) {
    fmt.Fprintf(os.Stderr, "error " + doing + ": " + err.Error() + "\n")
    os.Exit(1)
}

func fatal1(reason string) {
    fmt.Fprintf(os.Stderr, reason + "\n")
    os.Exit(1)
}
