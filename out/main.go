package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"
    "fmt"
    "github.com/jakobleben/slack-request-resource/protocol"
    "github.com/nlopes/slack"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <source>")
		os.Exit(1)
	}

    source_dir := os.Args[1]

    var request protocol.OutRequest

    request_err := json.NewDecoder(os.Stdin).Decode(&request)
    if request_err != nil {
        fatal("Parsing request.", request_err)
    }

    if len(request.Source.Token) == 0 {
        fatal1("Missing source field: token.")
    }

    if len(request.Source.AgentId) == 0 {
        fatal1("Missing source field: agent_id.")
    }

    if len(request.Source.ChannelId) == 0 {
        fatal1("Missing source field: channel_id.")
    }

    if len(request.Params.ContentFile) == 0 {
        fatal1("Missing params field: content.")
    }

    fmt.Fprintf(os.Stderr, "thread file: %s\n", request.Params.ThreadFile)
    fmt.Fprintf(os.Stderr, "contents file: %s\n", request.Params.ContentFile)

    contents := get_file_contents(filepath.Join(source_dir, request.Params.ContentFile))

    var thread string
    if len(request.Params.ThreadFile) != 0 {
        thread = get_file_contents(filepath.Join(source_dir, request.Params.ThreadFile))
    }

    fmt.Fprintf(os.Stderr, "thread: %s\n", thread)
    fmt.Fprintf(os.Stderr, "contents:\n%s\n", contents)

    slack_client := slack.New(request.Source.Token)

    send(thread, contents, &request, slack_client)

    var response protocol.OutResponse

    response_err := json.NewEncoder(os.Stdout).Encode(&response)
    if response_err != nil {
        fatal("encoding response", response_err)
    }
}

func get_file_contents(path string) string {
    file, open_err := os.Open(path)
    if open_err != nil {
        fatal("opening contents file", open_err)
    }

    data, read_err := ioutil.ReadAll(file)
    if read_err != nil {
        fatal("reading contents file", read_err)
    }

    return string(data)
}

func send(thread string, contents string, request *protocol.OutRequest, slack_client *slack.Client) {

    params := slack.NewPostMessageParameters()
    params.ThreadTimestamp = thread

    _, _, err := slack_client.PostMessage(request.Source.ChannelId, contents, params)
    if err != nil {
        fatal("sending", err)
    }
}

func fatal(doing string, err error) {
    fmt.Fprintf(os.Stderr, "error " + doing + ": " + err.Error() + "\n")
	os.Exit(1)
}

func fatal1(reason string) {
    fmt.Fprintf(os.Stderr, reason + "\n")
    os.Exit(1)
}
