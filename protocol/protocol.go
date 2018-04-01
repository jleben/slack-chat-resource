package protocol

import (
    //"strings"
    "regexp"
    //"errors"
    "encoding/json"
)

type Regexp struct { regexp.Regexp }

type MessageFilter struct {
    AuthorId string `json:"author"`
    TextPattern *Regexp `json:"text_pattern"`
}

type Source struct {
    Token string `json:"token"`
    ChannelId string `json:"channel_id"`
    Filter *MessageFilter `json:"matching"`
    ReplyFilter *MessageFilter `json:"not_replied_by"`
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
    Params InParams `json:"params"`
}

type InResponse struct {
    Version  Version  `json:"version"`
    Metadata Metadata `json:"metadata"`
}

type InParams struct {
    TextPattern *Regexp `json:"text_pattern"`
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

func (r *Regexp) UnmarshalJSON(payload []byte) error {
    var pattern string
    err := json.Unmarshal(payload, &pattern)
    if err != nil { return err }

    regexp, regexp_err := regexp.Compile(pattern)
    if regexp_err != nil { return regexp_err }

    *r = Regexp{*regexp}

    return nil
}
