# Slack Request Resource

A Concourse resource to read and reply to messages on Slack.

Docker Store: [jakobleben/slack-request-resource](https://store.docker.com/community/images/jakobleben/slack-request-resource)

## Source Configuration

- `token`: *Required*. A Slack API token that allows reading all messages on a selected channel and posting on it.
- `channel_id`: *Required*. The selected channel ID. The resource will only read messages and post on this channel.
- `matching`: *Optional*. Only read messages matching this filter. See below for details.
- `not_replied_by`: *Optional*. Ignore messages that have a reply matching this filter. See below for details.

The value of `matching` and `not_replied_by` represents a message filter. It is a map with the following elements:

- `author`: User ID representing the author of the message.
- `text_pattern`: Regular expression to match against the message text.

## Behavior

### `check`: List messages

When the resource is checked, it reads messages on the selected channel and reports their timestamps as versions, for example:

    { "timestamp": "<message timestamp>" }

It always reads at most 100 of the latest messages

Moreover, if given a message timestamp as the current `version`, it only reads messages with that timestamp and later.

It only reports messages that start threads and not replies to other messages.

If `source` has a `matching` filter, only messages that match the filter are reported.

If `source` has a `not_replied_by` filter, then a message that would otherwise be reported is ignored if it has a reply that matches the filter. In addition, only message later than the latest such message are reported.


### `in`: Read message

Reads the message with the requested timestamp and produces the following files:

- `timestamp`: The message timestamp.
- `text`: The message text.
- `text_part1`, `text_part2`, etc.: Parts of text parsed using the `text_pattern` parameter described below.

Parameters:

- `text_pattern`: *Optional*. A regular expression to match against the message text. The text matched by each capturing group is stored into a file `text_part<num>` where `<num>` is the group index starting with 1. This index is an integer representing the relative order of the beginnings of groups from left to right.

For example:

    - get: slack
      params:
          text_pattern: '([A-Z]+) ([0-9]+)'

When this configuration sees a message with the text `abc 123` and timestamp `111.222`, it will produce the following files and contents:

- `timestamp`: `111.222`
- `text`: `abc 123`
- `text_part1`: `abc`
- `text_part2`: `123`


### `out`: Post message

Posts a message to the selected channel.

Parameters:

- `text`: *Required*. The text to post.
- `thread`: *Optional*. The timestamp of the message to reply to. If missing, starts a new thread.

All parameters allow insertion of contents of arbitrary files. Each occurence of the pattern `{{filename}}` is substituted with the contents of the file `filename`.

For example, consider a job with the `get` example above followed by this:

    - put: slack
      params:
        thread: "{{slack/timestamp}}"
        text: "Hi {{slack/text_part1}}! I will do {{slack/text_part2}} right away!"

This will reply to the message read by the `put` step (since `thread` is the timestamp of the original message), and the reply will read:

    Hi abc! I will do 123 right away!
