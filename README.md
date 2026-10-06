<h1 align="center">whoisusing</h1>

<p align="center">
  Find which process is using a port, and optionally kill it.
</p>

<p align="center">
  <a href="https://github.com/Kishan-Thanki/whoisusing/releases">
    <img src="https://img.shields.io/github/v/release/Kishan-Thanki/whoisusing" alt="Release">
  </a>
  <a href="https://github.com/Kishan-Thanki/whoisusing/actions/workflows/ci.yaml">
    <img src="https://github.com/Kishan-Thanki/whoisusing/actions/workflows/ci.yaml/badge.svg" alt="CI">
  </a>
  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="MIT License">
  </a>
</p>

## Install

With Go 1.26 or newer:

```bash
go install github.com/kishan-thanki/whoisusing@latest
```

Or download a binary for your platform from the
[Releases](https://github.com/Kishan-Thanki/whoisusing/releases) page.

## Requirements

`lsof` must be installed.

macOS includes `lsof` by default. On Debian/Ubuntu:

```bash
sudo apt install lsof
```

Without sufficient permissions, `lsof` may only show processes you are
allowed to inspect, and `-k` can only terminate processes you are
allowed to signal.

## Supported Platforms

* Linux
* macOS
* FreeBSD
* NetBSD
* OpenBSD
* Solaris
* illumos

## Usage

```text
whoisusing -p <port> [-k] [-q]
```

| Flag           | Description                                  |
| -------------- | -------------------------------------------- |
| `-p <port>`    | Port to inspect (required, 1-65535)          |
| `-k`           | Gracefully kill the processes using the port |
| `-q`           | Quiet mode: print only PIDs, one per line    |
| `-h`, `--help` | Show help                                    |

Short flags can be combined in any order. `-p` takes the next argument
as its value.

```bash
whoisusing -p 8080
whoisusing -kp 8080
whoisusing -kqp 8080
PID=$(whoisusing -q -p 8080)
```

## Exit Codes

| Code | Meaning                                                                            |
| ---- | ---------------------------------------------------------------------------------- |
| `0`  | Success                                                                            |
| `1`  | Runtime or operational failure; in quiet mode, also returned when the port is free |
| `2`  | Invalid arguments                                                                  |

In normal mode, an unused port is considered successful.

## License

Licensed under the [MIT License](./LICENSE).
