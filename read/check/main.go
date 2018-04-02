package main

import (
    "encoding/json"
    //"io"
    //"ioutil"
    "os"
    //"os/exec"
    "fmt"
    //"strings"
    //"net/http"
    "github.com/jleben/slack-chat-resource/utils"
    "github.com/nlopes/slack"
)

func main() {

    var request utils.CheckRequest

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

    if request.Source.Filter != nil {
        fmt.Fprintf(os.Stderr, "Filter:\n")
        fmt.Fprintf(os.Stderr, "  - author: %s\n", request.Source.Filter.AuthorId)
        fmt.Fprintf(os.Stderr, "  - pattern: %s\n", request.Source.Filter.TextPattern)
    }

    if request.Source.ReplyFilter != nil {
        fmt.Fprintf(os.Stderr, "Reply Filter:\n")
        fmt.Fprintf(os.Stderr, "  - author: %s\n", request.Source.ReplyFilter.AuthorId)
        fmt.Fprintf(os.Stderr, "  - pattern: %s\n", request.Source.ReplyFilter.TextPattern)
    }

    slack_client := slack.New(request.Source.Token)

    history := get_messages(&request, slack_client)

    versions := []utils.Version{}

    for _, msg := range history.Messages {

        accept, stop := process_message(&msg, &request, slack_client)

        if accept {
            version := utils.Version{ "timestamp": msg.Msg.Timestamp }
            versions = append(versions, version)
        }

        if stop { break }
    }

    response := utils.CheckResponse{}
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

func get_messages(request *utils.CheckRequest, slack_client *slack.Client) *slack.History {

    params := slack.NewHistoryParameters()

    if request_version, ok := request.Version["timestamp"]; ok {
        params.Oldest = request_version
        fmt.Fprintf(os.Stderr, "Request timestamp: %s\n", request_version)
    }

    params.Inclusive = true
    params.Count = 100

    var history *slack.History
    history, err := slack_client.GetChannelHistory(request.Source.ChannelId, params)
    if err != nil {
        fatal("getting messages.", err)
    }

    return history
}

func process_message(message *slack.Message, request *utils.CheckRequest,
                     slack_client *slack.Client) (accept bool, stop bool) {

    is_reply := len(message.Msg.ThreadTimestamp) > 0 &&
        message.Msg.ThreadTimestamp != message.Msg.Timestamp

    if is_reply {
        fmt.Fprintf(os.Stderr, "Message %s is a reply. Skipping.\n", message.Msg.Timestamp)
        return false, false
    }

    fmt.Fprintf(os.Stderr, "- Message %s: %s \n", message.Msg.Timestamp, message.Msg.Text)

    if request.Source.Filter != nil {
        fmt.Fprintf(os.Stderr, "Matching message...\n")
        if !match_message(message, request.Source.Filter) {
            fmt.Fprintf(os.Stderr, "Message did not matched.\n")
            return false, false
        }
    }

    if request.Source.ReplyFilter != nil {
        fmt.Fprintf(os.Stderr, "Matching replies...\n")
        if match_replies(message, request, slack_client) {
            fmt.Fprintf(os.Stderr, "A reply was matched.\n")
            return false, true
        }
    }

    return true, false
}

func match_message(message *slack.Message, filter *utils.MessageFilter) bool {

    author_id := filter.AuthorId
    if len(author_id) > 0 && message.Msg.User != author_id && message.Msg.BotID != author_id {
        fmt.Fprintf(os.Stderr, "Author is not %s.\n", author_id)
        return false
    }

    text_pattern := filter.TextPattern
    if text_pattern != nil && !text_pattern.MatchString(message.Msg.Text) {
        fmt.Fprintf(os.Stderr, "Message text does not match pattern.\n")
        return false
    }

    fmt.Fprintf(os.Stderr, "Message matched.\n")

    return true
}

func match_replies(message *slack.Message, request *utils.CheckRequest, slack_client *slack.Client) bool {

    if message.Msg.ReplyCount == 0 {
        return false
    }

    replies, err := slack_client.GetChannelReplies(request.Source.ChannelId, message.Msg.Timestamp)
    if err != nil {
        fatal("getting replies", err)
    }

    for _, reply := range replies[1:] {
        fmt.Fprintf(os.Stderr, "- A reply: %s\n", reply.Msg.Text)
        if match_message(&reply, request.Source.ReplyFilter) {
            return true
        }
    }

    return false
}

/*
func get_channel_id(request utils.CheckRequest) {

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
*/

func fatal(doing string, err error) {
    fmt.Fprintf(os.Stderr, "error " + doing + ": " + err.Error() + "\n")
	os.Exit(1)
}

func fatal1(reason string) {
    fmt.Fprintf(os.Stderr, reason + "\n")
    os.Exit(1)
}
