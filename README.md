# jp-prev-bizday

Get the previous business day considering Japanese holidays.

> 日本の祝日に対応した直近の営業日を取得するツール

## Overview

`jp-prev-bizday` is a command-line tool that finds the most recent business day (excluding weekends and Japanese holidays) before a specified date. It uses the Japanese Holiday API to accurately determine holidays in Japan.

## Requirements

Accurate Japanese holiday detection using [jp-holiday.net API](https://jp-holiday.net/)

## Limitations

- Searches up to 30 days back for a business day
- Requires `internet connection` for holiday checking
- Only supports Japanese holidays

## Installation

```bash
$ go install github.com/yourusername/jp-prev-bizday@latest
```

## Usage

### Basic usage

Get the previous business day from today:

```bash
$ jp-prev-bizday
```

Output:
```
2025-07-23
```

### Specify a base date

Get the previous business day from a specific date:

```bash
jp-prev-bizday -date 2025-07-24
```

### Verbose output

Display detailed information including the base date:

```bash
jp-prev-bizday -verbose
```

Output:
```
基準日: 2025-07-24 (木)
直前の営業日: 2025-07-23 (水)
```

## License

This project is licensed under the [MIT License](./LICENSE).
