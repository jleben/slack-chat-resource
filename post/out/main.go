package main

import (
    "encoding/json"
    "io/ioutil"
    "os"
    "path/filepath"
    "fmt"
    "github.com/jleben/slack-chat-resource/utils"
    "github.com/nlopes/slack"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <source>")
		os.Exit(1)
	}

    source_dir := os.Args[1]

    var request utils.OutRequest

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

    if len(request.Params.Text) == 0 {
        fatal1("Missing params field: text.")
    }

    fmt.Fprintf(os.Stderr, "thread input: %s\n", request.Params.Thread)
    var thread string
    if len(request.Params.Thread) != 0 {
        thread = interpolate(request.Params.Thread, source_dir)
    }
    fmt.Fprintf(os.Stderr, "thread output: %s\n", thread)

    fmt.Fprintf(os.Stderr, "text input:\n%s\n", request.Params.Text)
    text := interpolate(request.Params.Text, source_dir)
    fmt.Fprintf(os.Stderr, "text output:\n%s\n", text)

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
        fatal("opening file", open_err)
    }

    data, read_err := ioutil.ReadAll(file)
    if read_err != nil {
        fatal("reading file", read_err)
    }

    return string(data)
}

func interpolate(text string, source_dir string) string {

    var out_text string

    start_var := 0
    end_var := 0
    inside_var := false
    c0 := '_'

    for pos, c1 := range text {
        if inside_var {
            if c0 == '}' && c1 == '}' {
                inside_var = false
                end_var = pos + 1
                var_name := text[start_var+2:end_var-2]
                value := get_file_contents(filepath.Join(source_dir, var_name))
                out_text += value
            }
        } else {
            if c0 == '{' && c1 == '{' {
                inside_var = true
                start_var = pos - 1
                out_text += text[end_var:start_var]
            }
        }
        c0 = c1
    }

    out_text += text[end_var:]

    return out_text
}

func send(thread string, text string, request *utils.OutRequest, slack_client *slack.Client) utils.OutResponse {

    params := slack.NewPostMessageParameters()
    params.ThreadTimestamp = thread

    _, timestamp, err := slack_client.PostMessage(request.Source.ChannelId, text, params)
    if err != nil {
        fatal("sending", err)
    }

    var response utils.OutResponse
    response.Version = utils.Version { "timestamp": timestamp }
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
