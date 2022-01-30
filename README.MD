# Files To Tar Bot

A Telegram bot to convert a list of files to a single raw tar archive.

## Running

You can download one of the binaries from release file. Then set these global variables before running the application:

```bash
export APP_ID=1234
export APP_HASH=abcd
export BOT_TOKEN=zyxw
export ALLOWED_USERS=1,2,3 # Optional. Can be omitted to run the bot as public bot
```

You can also build the application by running `go build`

## Bot

Bot currently supports 4 actions:

* `/start`: On this message a welcome message is shown
* `/reset`: Delete the list of files which will be placed in archive
* `/done`: Tar all files which the user has been sent to bot and upload it for them. You can also specify a filename
  after that to get the tar file with that name. For example `/done my-files` will send a file with name `my-files.tar`
* On media received: Bot keeps a list of files which each user has sent to it. When the user sends a media to bot, it
  will be added to this list.