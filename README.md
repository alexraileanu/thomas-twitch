## thomas-the-twitch-bot

### Installation

Installation is fairly straight-forward. You can either download the latest release (linux only) or compilte it yourself.
If you choose to compile it yourself, you only need to clone this repo and run the following command:

`make build`

This will create a thomas binary in the current working directory. Feel free to place it anywhere in your path
(`/usr/bin` or `/usr/local/bin` for instance).

### Set-up

#### Subscribe to an event on twitch.

In order to have twitch send you the proper events, Thomas needs to be exposed to the internet (meaning hosted
somewhere and accessible via either the server's IP or a domain name).

You will also need to create an application on twitch which will give you a `Client ID` and a `Client Secret`.

After this is set up, you can subscribe to an event. To find out more about how this is done, you can read
the docs on twitch [here](https://dev.twitch.tv/docs/eventsub/manage-subscriptions#subscribing-to-events).

Essentially, you post to `https://api.twitch.tv/helix/eventsub/subscriptions` with a JSON body like:
```json
{
  "type": "stream.online",
  "version": "1",
  "condition": {
      "broadcaster_user_id": "1234"
  },
  "transport": {
      "method": "webhook",
      "callback": "https://example.com/callback",
      "secret": "s3cre7"
  }
}
```
A couple of things you need to notice:
- The broadcaster_id you can get by making a request to `https://api.twitch.tv/helix/users?login=$BROADCASTER_NAME` where
`$BROADCASTER_NAME` is the streamer's name on twitch.
- Callback must be the domain/ip of Thomas publicly accessible from the internet. The path **MUST** be `/twitch/alert`.
- The secret value you should generate randomly and store it somewhere safe.


#### Discord setup
Create a bot on discord. More info [here](https://discord.com/developers/applications)

After you created your bot, you need to copy the bot token. Keep in mind, you'll only be able
to see this token once, so make sure you store it somewhere safe!

You will also need a channel where Thomas will send the notifications. You can get its id by right-clicking on
the channel and pressing the `Copy ID` option.

### .env contents

Once you have the appropriate tokens and codes, you can fill the .env file like so:

```dotenv
TWITCH_CLIENT_ID=
TWITCH_CLIENT_SECRET=
STREAM_ID=
SUB_SECRET=
DISCORD_TOKEN=
DISCORD_CHANNEL_ID=
```

These should be somewhat self-explanatory:
- `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET` you get when creating your app on twitch
- `STREAM_ID` you get after making a request to `https://api.twitch.tv/helix/users?login=$BROADCASTER_NAME`
- `SUB_SECRET` is the secret you generated when subscribing to the twitch event
- `DISCORD_TOKEN` is the bot (very important! not the app) id you get when creating your bot on discord
- `DISCORD_CHANNEL_ID` is the channel you want Thomas to send the messages on