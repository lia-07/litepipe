# LitePipe

LitePipe is a simple CI/CD tool written in Go using only the standard library. It's goal is to be lightweight and reliable. In a nutshell, it listens for webhook POST requests, verifies them, checks if some user-defined criteria were met, and if so runs some user-defined commands. It currently exclusively supports GitHub pull request webhooks, but in the future I aim to add more flexibility.

## Setup

After cloning this repo, (assuming you have Go installed) simply run `go build litepipe.go`.
Then, you must create a configuration file.

### The configuration file

The configuration file should looks something like this:

```json
{
  "port": 3001,
  "webhookSecret": "secret",
  "triggerDirectories": [
    "directory/*",
    "directory2/*.html",
    "directory2/*/*.css"
  ],
  "tasks": ["echo Hello, World"],
  "tasksDirectory": "directory/another_directory"
}
```

Here are descriptions of the different arguments:

- **`port`:** Sets the localhost port LitePipe should listen on. Optional, and defaults to `3001`.
- **`webhookSecret`:** Where you put the webhook secret.
- **`triggerDirectories`:** If a commit affects one of the directories listed here, the tasks will be executed. Optional, and defaults to `*`.
- **`tasks`:** The tasks (commands) to be executed.
- **`tasksDirectory`:** The directory the tasks will be executed in. Optional, and defaults to the directory LitePipe is running in.

By default, LitePipe looks for `config.json` in the same directory it is in. The configuration file should looks something like this:
If necessary (for example, if LitePipe is running in a `systemctl` environment) you can specify a custom path for the configuration file with the `-config` flag. For example, on a UNIX/UNIX-like system you could run `./litepipe -config "/home/user/litepipe/config.json"`.

After creating the configuration file, you simply need to execute the LitePipe binary. Personally, I set it up behind an Nginx reverse proxy, and haven't tested other configurations.
