# LitePipe

LitePipe is a simple CI/CD tool written in Go using only the standard library. Currently, it will listen for GitHub push webhooks, verify them, and then perform certain actions based on which path the change occurred in. Support for different types of webhooks will be added as needed (feel free to contribute).

Backwards compatibility is not yet guaranteed.

## Setup

After cloning this repo, (assuming you have Go installed) simply run `go build litepipe.go`.
Then, you must create a configuration file.

### The configuration file

The configuration file should looks something like this:

```json
{
  "port": 3001,
  "webhookSecret": "secret",
  "triggerPaths": ["directory/*", "directory2/*.html", "directory2/*/*.css"],
  "tasks": ["echo Hello, World"],
  "tasksDirectory": "directory/another_directory"
}
```

Here are descriptions of the different arguments:

- **`port`:** Sets the localhost port LitePipe should listen on. Optional, and defaults to `3001`.
- **`webhookSecret`:** Where you put the webhook secret.
- **`triggerPaths`:** If a commit affects one of the directories or files listed here, the tasks will be executed. Optional, and defaults to `*`.
- **`tasks`:** The tasks (commands) to be executed.
- **`tasksDirectory`:** The directory the tasks will be executed in. Optional, and defaults to the directory LitePipe is running in.

By default, LitePipe looks for `config.json` in the same directory it is in.
If necessary (for example, if LitePipe is running in a `systemctl` environment) you can specify a custom path for the configuration file with the `-config` flag. For example, on a UNIX/UNIX-like system you could run `./litepipe -config "/home/user/litepipe/config.json"`.

After creating the configuration file, you simply need to execute the LitePipe binary. I have it set up with `systemctl` on my Ubuntu VPS, behind an Nginx reverse proxy. 
