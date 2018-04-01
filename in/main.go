package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"
    "fmt"
    "github.com/jleben/trigger-resource/protocol"
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

    if len(request.Source.Context) == 0 {
        fatal1("Missing source field: context.")
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

    reply(request, slack_client)
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

    slack_request := protocol.ParseSlackRequest(message.Msg.Text, &request.Source)

    if slack_request == nil {
        fatal1("Failed to parse message.")
    }

    {
        err := os.MkdirAll(destination, 0755)
        if err != nil {
            fatal("creating destination directory", err)
        }
    }

	{
        err := ioutil.WriteFile(filepath.Join(destination, "contents"), []byte(slack_request.Contents), 0644)
        if err != nil {
            fatal("writing contents file", err)
        }
    }

    var response protocol.InResponse
    response.Version = request.Version
    return response
}

func reply(request protocol.InRequest, slack_client *slack.Client) {

    params := slack.NewPostMessageParameters()
    params.ThreadTimestamp = request.Version["request"]

    text := fmt.Sprintf("Working on it.")

    _, _, err := slack_client.PostMessage(request.Source.ChannelId, text, params)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Warning: Failed to reply to Slack request: %s\n", err.Error())
    }
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}

func fatal1(reason string) {
    fmt.Fprintf(os.Stderr, reason + "\n")
    os.Exit(1)
}
