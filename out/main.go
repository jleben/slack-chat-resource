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

    if len(request.Source.ChannelId) == 0 {
        fatal1("Missing source field: channel_id.")
    }

    if len(request.Params.TextFile) == 0 && len(request.Params.Text) == 0 {
        fatal1("Missing params field: text or text_file.")
    }

    fmt.Fprintf(os.Stderr, "thread file: %s\n", request.Params.ThreadFile)
    fmt.Fprintf(os.Stderr, "text file: %s\n", request.Params.TextFile)
    fmt.Fprintf(os.Stderr, "input text:\n%s\n", request.Params.Text)

    text := request.Params.Text

    if len(text) == 0 {
        text = get_file_contents(filepath.Join(source_dir, request.Params.TextFile))
    }

    var thread string
    if len(request.Params.ThreadFile) != 0 {
        thread = get_file_contents(filepath.Join(source_dir, request.Params.ThreadFile))
    }

    fmt.Fprintf(os.Stderr, "thread: %s\n", thread)
    fmt.Fprintf(os.Stderr, "output text:\n%s\n", text)

    slack_client := slack.New(request.Source.Token)

    response := send(thread, text, &request, slack_client)

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

func send(thread string, text string, request *protocol.OutRequest, slack_client *slack.Client) protocol.OutResponse {

    params := slack.NewPostMessageParameters()
    params.ThreadTimestamp = thread

    _, timestamp, err := slack_client.PostMessage(request.Source.ChannelId, text, params)
    if err != nil {
        fatal("sending", err)
    }

    var response protocol.OutResponse
    response.Version = protocol.Version { "timestamp": timestamp }
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
