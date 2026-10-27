# whoisusing

<p align="center">
  Find which process is listening on a TCP port, and optionally terminate it.
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

Or download a binary for your platform from the [Releases](https://github.com/Kishan-Thanki/whoisusing/releases) page.

## Requirements

`lsof` must be installed.

macOS includes `lsof` by default. On Debian/Ubuntu:

```bash
sudo apt install lsof
```

Without sufficient permissions, `lsof` may only show processes you are allowed to inspect, and `-k` can only terminate processes you are allowed to signal.

## Supported Platforms

- Linux
- macOS
- FreeBSD
- NetBSD
- OpenBSD
- Solaris
- illumos

## Usage

```text
whoisusing -p <port> [-k] [-q]
```

| Flag | Description |
|---|---|
| `-p <port>` | TCP port to inspect (required, range: 1–65535) |
| `-k` | Gracefully terminate processes listening on the specified TCP port |
| `-q` | Quiet mode: print only PIDs, one per line |
| `-h`, `--help` | Display help |

Use each short flag separately. The `-p` flag takes the next argument as its value.

```bash
whoisusing -p 8080
whoisusing -k -p 8080
whoisusing -q -p 8080
whoisusing -k -q -p 8080
PID=$(whoisusing -q -p 8080)
```

## Exit Codes

| Code | Meaning |
|---|---|
| `0` | Success; a free port is also considered successful in normal mode |
| `1` | Runtime or operational failure; in quiet mode, also returned when no TCP listener is found |
| `2` | Invalid arguments |

## License

Licensed under the [MIT License](./LICENSE).
