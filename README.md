# snowflake
discord bot that tracks snowfall accumulation (US) 

create a `config.json` in the project root
```js
{
  "token": "",
  "guild_id": ""
}
```
`guild_id` is optional as it's the id of the guild to create application commands during development.

### Commands
- `/snowfall` `day` | display top snowfall accumilation withing the last 24hr
- `/snowfall` `season` | display top snowfall accumilation during the current season.
