package protocol

import (
    "strings"
    "regexp"
    "errors"
)

type Source struct {
    ChannelId string `json:"channel_id"`
    AgentId string `json:"agent_id"`
    Token string `json:"token"`
    Pattern string `json:"pattern"`
    IgnoreReplied bool `json:"ignore_replied"`
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

func ParseMessage(text string, source *Source) ([]string, bool, error) {

    text = strings.TrimLeft(text, " ")

    bot_mention := "<@" + source.AgentId + ">"
    if !strings.HasPrefix(text, bot_mention) { return nil, false, nil }

    text = text[len(bot_mention):]
    text = strings.TrimLeft(text, " ")

    if len(source.Pattern) == 0 { return []string{text}, true, nil }

    regexp, regexp_err := regexp.Compile("^" + source.Pattern + "$")
    if regexp_err != nil {
        return nil, false, errors.New("Invalid pattern")
    }

    matches := regexp.FindStringSubmatch(text)

    matched := len(matches) > 0

    return matches, matched, nil
}

