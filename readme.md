# Discord Bot
Simple discord music recommendation bot created using golang, leveraging [discordgo](https://github.com/bwmarrin/discordgo), [LangChain Go](https://pkg.go.dev/github.com/tmc/langchaingo), and [Spotify](https://pkg.go.dev/github.com/zmb3/spotify@v1.3.0).

## Running
Navigate to file structure in Command Prompt or Terminal and execute the following command:
```
go run main.go
```

This will activate the bot in the Discord channel until the command shell is cancelled (e.g. ctrl + C)

## Commands
|Command|Input|Description|
|-|-|-|
|**-r**|{artist/song/album}|Returns a ten song list based on input with Spotify links|
|**-m**|{song} {artist}|*not yet implemented*|