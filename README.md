<h1 align="center">whoisusing</h1>

<p align="center">
  Find which process is using a port, and optionally kill it.
</p>

<p align="center">
  <a href="https://github.com/Kishan-Thanki/whoisusing/releases">
    <img src="https://img.shields.io/github/v/release/Kishan-Thanki/whoisusing" alt="Release">
  </a>
  <a href="https://github.com/Kishan-Thanki/whoisusing/actions">
    <img src="https://github.com/Kishan-Thanki/whoisusing/actions/workflows/ci.yaml/badge.svg" alt="CI">
  </a>
  <a href="https://github.com/Kishan-Thanki/whoisusing/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/Kishan-Thanki/whoisusing" alt="License">
  </a>
</p>

## Install

With Go 1.26 or newer:

```bash
go install github.com/Kishan-Thanki/whoisusing@latest
```

Or download a binary for Linux or macOS from the
[Releases](https://github.com/Kishan-Thanki/whoisusing/releases) page.

## Requirements

`lsof` must be installed. It ships with macOS. On minimal Linux images:

```bash
sudo apt install lsof
```

Without sufficient permissions, `lsof` may only show processes you are
allowed to inspect, and `-k` can only terminate processes you are
allowed to signal.

## Usage

```text
whoisusing -p <port> [-k] [-q]
```

| Flag           | Description                                            |
| -------------- | ------------------------------------------------------ |
| `-p <port>`    | Port to inspect (required, 1-65535)                    |
| `-k`           | Gracefully kill (SIGTERM) the processes using the port |
| `-q`           | Quiet mode: print only PIDs, one per line              |
| `-h`, `--help` | Show help                                              |

Short flags can be combined in any order. `-p` takes the next argument
as its value.

```bash
whoisusing -p 8080
whoisusing -kp 8080
whoisusing -kqp 8080
PID=$(whoisusing -q -p 8080)
```

## Exit codes

| Code | Meaning                                                                            |
| ---- | ---------------------------------------------------------------------------------- |
| `0`  | Success                                                                            |
| `1`  | Runtime or operational failure; in quiet mode, also returned when the port is free |
| `2`  | Invalid arguments                                                                  |

In normal mode, an unused port is considered successful.

## License

Licensed under the [MIT License](./LICENSE).
