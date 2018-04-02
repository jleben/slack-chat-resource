# Slack Chat Resources

This repository provides Concourse resource types to read, act on, and reply to messages on Slack.

There are two resource types:

- `slack-read-resource`: For reading messages.
- `slack-post-resource`: For posting messages.

There are two resource types because a system does not want to respond to messages that it posts itself. Concourse assumes that an output of a resource is also a valid input. Therefore, separate resources are used for reading and posting. Since using a single resource has no benefits over separate resources, reading and posting are split into two resource types.

Docker Store:

- [jakobleben/slack-read-resource](https://store.docker.com/community/images/jakobleben/slack-read-resource)
- [jakobleben/slack-post-resource](https://store.docker.com/community/images/jakobleben/slack-post-resource)

## Version Format

Timestamps of Slack messages are used as resource versions. For example, a message may be represented by a version like this:

    timestamp: 1234567890.123

A timestamp uniquely identifies a message within a channel. See [Slack API](https://api.slack.com/events/message) for details.

## Reading Messages

Usage in a pipeline:

    resource_types:
        - name: slack-read-resource
          type: docker-image
          source:
            repository: jakobleben/slack-read-resource

    resources:
        - name: slack-in
          type: slack-read-resource
          source: ...

### Source Configuration

The `source` field configures the resource for reading messages from a specific channel. It allows filtering messages by their author and text pattern:

- `token`: *Required*. A Slack API token that allows reading all messages on a selected channel.
- `channel_id`: *Required*. The selected channel ID. The resource only reads messages on this channel.
- `matching`: *Optional*. Only report messages matching this filter. See below for details.
- `not_replied_by`: *Optional*. Ignore messages that have a reply matching this filter. See below for details.

The values of `matching` and `not_replied_by` represent message filters. They are maps with the following elements:

- `author`: *Optional*. User ID that must match the author of the message - either the `user` or the `bot_id` field.
  See [Slack API](https://api.slack.com/events/message) regarding authorship.
- `text_pattern`: *Optional*. Regular expression that must match the message text.
  Wrap in single quotes instead of double, to avoid having to escape `\`.
  See [Slack API](https://api.slack.com/docs/message-formatting) for details on text formatting.


The resource only reports messages that begin new threads and not replies to other messages.

When given a message timestamp as the current version, it only reads messages with that timestamp and later. In any case though, it reads at most 100 of the latest messages. Therefore, the resource must be checked often enough to avoid missing messages.

If `source` has a `not_replied_by` filter, and it matches a message that also matches the `matching` filter, then all messages older than the latest such message are also considered obsolete and are not read.

#### Example

    resources:
      - name: slack-in
        type: slack-read-resource
        source:
          token: "xxxx-xxxxxxxxxx-xxxx"
          channel_id: "C11111111"
          matching:
            text_pattern: '<@U22222222>\s+(.+)'
          not_replied_by:
            author: U22222222

This configures a resource reading messages from channel with ID `C11111111`. It reads only messages that begin by mentioning the user with ID `U22222222`. It ignores messages already replied to by that same user.

### `get`: Read a Message

Reads the message with the requested timestamp and produces the following files:

- `timestamp`: The message timestamp.
- `text`: The message text.
- `text_part1`, `text_part2`, etc.: Parts of text parsed using the `text_pattern` parameter described below.

Parameters:

- `text_pattern`: *Optional*. A regular expression to match against the message text.
  The text matched by each [capturing group](https://www.regular-expressions.info/brackets.html)
  is stored into a file `text_part<num>` where `<num>` is the group index starting with 1.
  Wrap in single quotes instead of double, to avoid having to escape `\`.
  See [Slack API](https://api.slack.com/docs/message-formatting) for details on text formatting.

#### Example

    - get: slack-in
      params:
          text_pattern: '([A-Z]+) ([0-9]+)'

When this configuration sees a message with the text `abc 123` and timestamp `111.222`, it will produce the following files and contents:

- `timestamp`: `111.222`
- `text`: `abc 123`
- `text_part1`: `abc`
- `text_part2`: `123`


## Posting Messages

Usage in a pipeline:

    resource_types:
        - name: slack-post-resource
          type: docker-image
          source:
            repository: jakobleben/slack-post-resource

    resources:
        - name: slack-out
          type: slack-post-resource
          source: ...

### Source Configuration

The `source` field configures the resource for posting on a specific channel:

- `token`: *Required*. A Slack API token that allows posting on a selected channel.
- `channel_id`: *Required*. The selected channel ID. The resource only posts messages on this channel.

#### Example

    resources:
      - name: slack-out
        type: slack-post-resource
        source:
          token: "xxxx-xxxxxxxxxx-xxxx"
          channel_id: "C11111111"

This configures the resource to post on the channel with ID `C11111111`.

### `put`: Post a Message

Posts a message to the selected channel.

Parameters:

- `text`: *Required*. The text to post.
- `thread`: *Optional*. The timestamp of the message to reply to. If missing, starts a new thread.

All parameters allow insertion of contents of arbitrary files. Each occurence of the pattern `{{filename}}` is substituted with the contents of the file `filename`.

### Example

Consider a job with the `get: slack-in` step from the example above followed by this step:

    - put: slack-out
      params:
        thread: "{{slack-in/timestamp}}"
        text: "Hi {{slack-in/text_part1}}! I will do {{slack-in/text_part2}} right away!"

This will reply to the message read by the `get` step (since `thread` is the timestamp of the original message), and the reply will read:

    Hi abc! I will do 123 right away!
