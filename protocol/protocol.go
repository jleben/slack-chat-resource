package protocol

import (
    "strings"
)

type Source struct {
    ChannelId string `json:"channel_id"`
    AgentId string `json:"agent_id"`
    Token string `json:"token"`
    Context string `json:"context"`
}

type Version map[string]string

type Metadata []MetadataField

type MetadataField struct {
    Name  string `json:"name"`
    Value string `json:"value"`
}

type InRequest struct {
    Source  Source  `json:"source"`
    Version Version `json:"version"`
}

type InResponse struct {
    Version  Version  `json:"version"`
    Metadata Metadata `json:"metadata"`
}

type OutParams struct {
    ThreadFile string `json:"thread"`
    ContentFile string `json:"content"`
}

type OutRequest struct {
    Source  Source  `json:"source"`
    Params OutParams `json:"params"`
}

type OutResponse struct {
    Version  Version  `json:"version"`
    Metadata Metadata `json:"metadata"`
}

type CheckRequest struct {
    Source  Source  `json:"source"`
    Version Version `json:"version"`
}

type CheckResponse []Version

type SlackRequest struct {
    Contents string
}

func ParseSlackRequest(text string, source *Source) *SlackRequest {

    text = strings.TrimLeft(text, " ")

    bot_mention := "<@" + source.AgentId + ">"
    if !strings.HasPrefix(text, bot_mention) { return nil }

    text = text[len(bot_mention):]
    text = strings.TrimLeft(text, " ")

    parts := strings.SplitN(text, " ", 2)
    if len(parts) < 2 { return nil }

    context := parts[0]
    if len(context) == 0 { return nil }
    if context != source.Context { return nil }

    text = strings.Trim(parts[1], " ")

    request := new(SlackRequest)
    request.Contents = text

    return request
}

