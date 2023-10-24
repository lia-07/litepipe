# LitePipe

LitePipe is a simple CI/CD program written in Go using only the standard library. It's goal is to be lightweight and reliable. Put simply, it's job is to receive a webhook POST request, verify it, and then if some user-defined parameters are met, run arbitrary commands. It currently exclusively supports GitHub pull request webhooks, but in the future I aim to add more flexibility.

## Setup

After cloning this repo, (assuming you have Go installed) simply run `go build litepipe.go`. Then, you must create a configuration file.

### The configuration file

The configuration file should looks something like this:

```{
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

The `port` argument sets the localhost port LitePipe should listen on. The

The `port`, `triggerDirectories`, and `tasksDirectory` are optional, and default to `3001`, `["*"]`, and `""` (current directory) if unspecified.

By default, LitePipe looks for `config.json` in the same directory it is in. The configuration file should looks something like this:
If necessary (for example, if LitePipe is running in a `systemctl` environment) you can specify a custom path for the configuration file with the `-config` flag. For example, on a UNIX/UNIX-like system you could run `./litepipe -config "/home/user/litepipe/config.json"`.

After creating the configuration file, you simply need to execute the LitePipe binary. By default
