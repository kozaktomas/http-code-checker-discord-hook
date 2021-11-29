# http-code-checker-discord-hook

[![WTFPL](http://www.wtfpl.net/wp-content/uploads/2012/12/wtfpl-badge-4.png)](http://www.wtfpl.net/)

App periodically creates HTTP request and sends Discord notification when reach expected status code.

```
$ http-code-checker-discord-hook --help
usage: main [<flags>] <url> <code> <discord_hook_url>

Flags:
      --help        Show context-sensitive help (also try --help-long and --help-man).
  -s, --sleep="5m"  Duration between checks (2s, 5m, 10h, 2d). Default 5m.
  -v, --verbose     Verbose mode. Default false.

Args:
  <url>               Requested URL address.
  <code>              Expected HTTP status code.
  <discord_hook_url>  Discord webhook url.
```

## Example:

```bash
http-code-checker-discord-hook --sleep="5m" https://google.com 200 https://discord.com/api/webhooks/xyz/xyz
```
