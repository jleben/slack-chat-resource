# Slack Chat Resource

A Concourse resource to read, act on, and reply to messages on Slack.

Docker Store: [jakobleben/slack-chat-resource](https://store.docker.com/community/images/jakobleben/slack-chat-resource)

## Usage

    resource_types:
        - name: slack-chat-resource
          type: docker-image
          source:
            repository: jakobleben/slack-chat-resource

    resources:
        - name: slack
          type: slack-chat-resource
          source: ...

## Source Configuration

When the resource is checked, it reads messages according to the `source` configuration. For each message that matches the configuration, it outputs a resource version in the form:

    timestamp: 1234567890.123

A timestamp uniquely identifies a message within a channel (see [Slack API](https://api.slack.com/events/message) for details).

The `source` field may have the following elements:

- `token`: *Required*. A Slack API token that allows reading all messages on a selected channel and posting on it.
- `channel_id`: *Required*. The selected channel ID. The resource only reads and posts messages on this channel.
- `matching`: *Optional*. Only reports messages matching this filter. See below for details.
- `not_replied_by`: *Optional*. Ignores messages that have a reply matching this filter. See below for details.

The values of `matching` and `not_replied_by` represent message filters. They are maps with the following elements:

- `author`: *Optional*. User ID that must match the author of the message - either the `user` or the `bot_id` field.
  See [Slack API](https://api.slack.com/events/message) regarding authorship.
- `text_pattern`: *Optional*. Regular expression that must match the message text.
  See [Slack API](https://api.slack.com/docs/message-formatting) for details on text formatting.

The resource only reports messages that begin new threads and not replies to other messages.

When given a message timestamp as the current version, it only reads messages with that timestamp and later. In any case though, it reads at most 100 of the latest messages. Therefore, the resource must be checked often enough to avoid missing messages.

If `source` has a `not_replied_by` filter, and it matches a message that also matches the `matching` filter, then all messages older than the latest such message are also considered obsolete and are not read.

### Example

    resources:
      - name: slack
        type: slack-chat-resource
        source:
          token: "xxxx-xxxxxxxxxx-xxxx"
          channel_id: "C11111111"
          matching:
            text_pattern: '<@U22222222>\\s+(.+)'
          not_replied_by:
            author: U22222222

This configures a resource reading messages from channel with ID `C11111111`. It reads only messages the begin by mentioning the user with ID `U22222222`. It ignores messages already replied to by that same user.

## `get`: Read message

Reads the message with the requested timestamp and produces the following files:

- `timestamp`: The message timestamp.
- `text`: The message text.
- `text_part1`, `text_part2`, etc.: Parts of text parsed using the `text_pattern` parameter described below.

Parameters:

- `text_pattern`: *Optional*. A regular expression to match against the message text. The text matched by each capturing group is stored into a file `text_part<num>` where `<num>` is the group index starting with 1. This index is an integer representing the relative order of the beginnings of groups from left to right.
  See [Slack API](https://api.slack.com/docs/message-formatting) for details on text formatting.

### Example

    - get: slack
      params:
          text_pattern: '([A-Z]+) ([0-9]+)'

When this configuration sees a message with the text `abc 123` and timestamp `111.222`, it will produce the following files and contents:

- `timestamp`: `111.222`
- `text`: `abc 123`
- `text_part1`: `abc`
- `text_part2`: `123`


## `put`: Post message

Posts a message to the selected channel.

Parameters:

- `text`: *Required*. The text to post.
- `thread`: *Optional*. The timestamp of the message to reply to. If missing, starts a new thread.

All parameters allow insertion of contents of arbitrary files. Each occurence of the pattern `{{filename}}` is substituted with the contents of the file `filename`.

### Example

Consider a job with the `get` example above followed by this:

    - put: slack
      params:
        thread: "{{slack/timestamp}}"
        text: "Hi {{slack/text_part1}}! I will do {{slack/text_part2}} right away!"

This will reply to the message read by the `get` step (since `thread` is the timestamp of the original message), and the reply will read:

    Hi abc! I will do 123 right away!
