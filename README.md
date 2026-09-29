# whoisusing

Find which process is using a port, and optionally kill it.

## Install

With Go 1.21 or newer:

```bash
go install github.com/kishan-thanki/whoisusing@latest
```

Or download a binary for Linux or macOS from the
[Releases](https://github.com/kishan-thanki/whoisusing/releases) page.

## Requirements

`lsof` must be installed. It ships with macOS. On minimal Linux images:

```bash
sudo apt install lsof
```

Without root, `lsof` only shows your own processes, and `-k` can only
kill your own processes.

## Usage

```
whoisusing -p <port> [-k] [-q]
```

| Flag | Description |
|---|---|
| `-p <port>` | Port to inspect (required, 1-65535) |
| `-k` | Gracefully kill (SIGTERM) the processes using the port |
| `-q` | Quiet mode: print only PIDs, one per line |
| `-h`, `--help` | Show help |

Short flags can be combined in any order. `-p` takes the next argument
as its value.

```bash
whoisusing -p 8080
whoisusing -kp 8080
whoisusing -kqp 8080
PID=$(whoisusing -q -p 8080)
```

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success (normal mode: also when the port is free) |
| 1 | Quiet mode: port is free; or kill failed |
| 2 | Invalid arguments |

## License

MIT
