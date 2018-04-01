package main

import (
    "encoding/json"
    //"io"
    //"ioutil"
    "os"
    //"os/exec"
    "fmt"
    "strings"
    //"net/http"
    "github.com/jakobleben/slack-request-resource/protocol"
    "github.com/nlopes/slack"
)

func main() {

    var request protocol.CheckRequest

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

    if len(request.Source.AgentId) == 0 {
        fatal1("Missing source field: agent_id.")
    }

    fmt.Fprintf(os.Stderr, "Pattern: %s\n", request.Source.Pattern)

    slack_client := slack.New(request.Source.Token)

    params := slack.NewHistoryParameters()

    if request_version, ok := request.Version["request"]; ok {
        params.Oldest = request_version
        fmt.Fprintf(os.Stderr, "Request version: %s\n", request_version)
    }

    params.Inclusive = true
    params.Count = 100

    var history *slack.History
    history, err = slack_client.GetChannelHistory(request.Source.ChannelId, params)
    if err != nil {
		fatal("getting messages.", err)
	}

    versions := []protocol.Version{}

    for _, msg := range history.Messages {

        version, was_detected := process_message(&msg, request, slack_client)

        if was_detected { break }

        if version != nil {
            versions = append(versions, version)
        }
    }

    response := protocol.CheckResponse{}
    for i := len(versions) - 1; i >= 0; i--  {
        response = append(response, versions[i])
    }

    json.NewEncoder(os.Stdout).Encode(&response)
}

type Channel struct {
    id string
    name string
}

type ChannelsMeta struct {
    next_cursor string
}

type Channels struct {
    ok bool
    channels []Channel
    meta ChannelsMeta
}

// Check if message is a request for us.
// Check if the request was already handled, ignore it if so.
// Extract requested version
// Return extracted version (if any), and a flag whether the request was already handled.

func process_message(message *slack.Message, request protocol.CheckRequest, slack_client *slack.Client) (protocol.Version, bool) {

    is_reply := len(message.Msg.ThreadTimestamp) > 0 &&
        message.Msg.ThreadTimestamp != message.Msg.Timestamp

    if is_reply {
        fmt.Fprintf(os.Stderr, "Message %s is a reply. Skipping.\n", message.Msg.Timestamp)
        return nil, false
    }

    is_by_bot := message.Msg.SubType == "bot_message" || len(message.Msg.User) == 0

    if is_by_bot {
        fmt.Fprintf(os.Stderr, "Message %s is by a bot. Skipping.\n", message.Msg.Timestamp)
        return nil, false
    }

    fmt.Fprintf(os.Stderr, "Message %s: %s \n", message.Msg.Timestamp, message.Msg.Text)

    matches, is_matched, parse_err := protocol.ParseMessage(message.Msg.Text, &request.Source)

    if parse_err != nil {
        fatal("parsing message", parse_err)
    }

    if !is_matched {
        fmt.Fprintf(os.Stderr, "Invalid format.\n")
        return nil, false
    }

    fmt.Fprintf(os.Stderr, "Message parsed: %s\n", strings.Join(matches, ", "))

    if request.Source.IgnoreReplied {
        if has_reply(message, &request, slack_client) {
            fmt.Fprintf(os.Stderr, "Message already replied to.\n")
            return nil, true
        }
    }

    version := protocol.Version{
        "request": message.Msg.Timestamp,
    }

    return version, false
}

func has_reply(message *slack.Message, request *protocol.CheckRequest, slack_client *slack.Client) bool {
    if message.Msg.ReplyCount == 0 {
        return false
    }

    replies, err := slack_client.GetChannelReplies(request.Source.ChannelId, message.Msg.Timestamp)
    if err != nil {
        fatal("getting replies", err)
    }

    for _, reply := range replies {
        has_reply := reply.User == request.Source.AgentId
        if has_reply { return true }
    }

    return false
}

/*
func get_channel_id(request protocol.CheckRequest) {

    var done = false
    var cursor string

    for !done {
        channels := get_channels(cursor)

        for _, channel := range channels.channels {
            fmt.Fprintf(os.Stderr, "Channel: %s %s\n", channel.id, channel.name)
        }

        cursor = channels.meta.next_cursor
        done = len(cursor) == 0
    }
}

func get_channels(cursor string) (Channels) {
    url = "https://slack.com/api/channels.list?" +
        "token=" + request.Source.Token +
        "&exclude_archived=true" +
        "&exclude_members=true"

    if len(cursor) > 0 {
        url += "&cursor=" + cursor
    }

    resp, get_err := http.Get(url)
    if get_err != nil { fatal("getting channels", get_err) }

    body, read_err := ioutil.ReadAll(resp.Body)
    if read_err != nil { fatal("getting channels - reading response body", read_err) }

    var channels Channels
    parse_err := json.Unmarshall(body, &channels)
    if parse_err != nil { fatal("getting channels - parsing response body", parse_err) }

    return channels
}

func get_history(request protocol.CheckRequest, channel_id) {
    url = "https://slack.com/api/channels.history?" +
        "token=" + request.Source.Token
    response, err := http.Get()
}
*/

func fatal(doing string, err error) {
    fmt.Fprintf(os.Stderr, "error " + doing + ": " + err.Error() + "\n")
	os.Exit(1)
}

func fatal1(reason string) {
    fmt.Fprintf(os.Stderr, reason + "\n")
    os.Exit(1)
}
