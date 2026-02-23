# ECHELON - SS13 Log Analysis & Replay System

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?logo=go)

**ECHELON** is a powerful log analysis tool for Space Station 13 server administrators. Named after the surveillance program, it helps admins efficiently investigate incidents, track player actions, and filter through massive log files with ease.

## Features

### Core Functionality
- ✅ **Multi-file parsing** - Process multiple log files (attack.log, game.log, etc.) simultaneously
- ✅ **Player filtering** - Filter logs by ckey (player key)
- ✅ **Round validation** - Warns if log files are from different rounds
- ✅ **Timestamp filtering** - Filter events relative to round start (e.g., "10-30 minutes into the round")
- ✅ **Multiple output formats** - Export as readable .log, JSON, or color-coded HTML
- ✅ **Clutter Cleanup** - Optional removal of cluttering log entries for cleaner output (eg. Mob ID's)
- ✅ **Drag & Drop** - Easy file input via drag-and-drop interface

### Dual Interface
- **CLI** - Fast command-line tool for quick filtering
- **GUI** - User-friendly graphical interface with drag-and-drop

## Screenshots

<img width="775" height="1067" alt="image" src="https://github.com/user-attachments/assets/c194be5f-eee6-49ae-a659-4988fd0cbde2" />


## Installation

### Prerequisites
- Go 1.21 or higher
- (Windows) MinGW-w64 GCC compiler for GUI build

### Download Pre-built Binary
*Coming soon - check [Releases](https://github.com/Amor4k/Echelon/releases)*

### Build from Source
```bash
# Clone the repository
git clone https://github.com/Amor4k/Echelon.git
cd Echelon

# Install dependencies
go mod tidy

# Build CLI
go build -o echelon-cli.exe cmd/cli/main.go

# Build GUI
go build -o echelon.exe cmd/gui/main.go cmd/gui/output.go

# Optional: Reduce executable size
go build -ldflags="-s -w" -o echelon.exe cmd/gui/*.go
```

## Usage

### GUI Mode

Simply run the executable:
```bash
./echelon.exe
```

1. Enter the player's ckey
2. Drag & drop log files or use "Select Log Files"
3. Configure options (mob ID cleanup, time filtering, output format)
4. Click "Filter Logs"
5. Results saved to the output directory

### CLI Mode
```bash
# Basic usage
echelon-cli --ckey playername attack.log.json game.log.json

# With time filtering (30-45 minutes after round start)
echelon-cli --ckey playername --after-minutes 30 --before-minutes 45 attack.log.json

# Disable mob ID cleaning
echelon-cli --ckey playername --clean-mob-ids=false attack.log.json

# Without banner
echelon-cli --ckey playername --no-banner attack.log.json
```

#### CLI Flags
- `--ckey` - Player ckey to filter (required)
- `--after-minutes` - Filter logs after N minutes from round start
- `--before-minutes` - Filter logs before N minutes from round start
- `--clean-mob-ids` - Remove mob IDs from output (default: true)
- `--no-banner` - Disable ASCII banner

## Log Format Support

ECHELON supports SS13 JSON logs with the following structure:
```json
{
  "ts": "2026-02-02 20:38:29.305",
  "cat": "attack",
  "msg": "playerckey/(Player Name) (mob_123) attacked...",
  "w-state": {...},
  "id": 1234
}
```

Compatible with:
- attack.log.json
- game.log.json
- Any SS13 subsystem logs in JSON format

## Output Formats

### Readable Log (.log)
```
[2026-02-02 20:38:29] [attack] playerckey/(Player Name) attacked target...
[2026-02-02 20:39:15] [game] playerckey/(Player Name) said "Hello world"
```

### JSON (.json)
Structured data for further processing or programmatic analysis.

### HTML (.html)
Color-coded, browser-viewable format with:
- Attack events in red
- Game events in teal
- Timestamps highlighted in blue
- Clean, monospace formatting

## Project Structure
```
Echelon/
├── cmd/
│   ├── cli/           # Command-line interface
│   └── gui/           # Graphical interface
├── internal/
│   └── parser/        # Core parsing logic
├── testdata/          # Sample log files
└── README.md
```

## Roadmap

### Planned Features
- [ ] **Multiple ckey support** - Filter for multiple players simultaneously
- [ ] **Interaction filtering** - Find when two players interact
- [ ] **Statistics view** - Summary of player actions (kills, deaths, chat count)
- [ ] **Settings memory** - Remember last used preferences
- [ ] **Log preview** - Preview results before saving
- [ ] **DM replay system** - Visual replay in BYOND (long-term goal)

### Future Enhancements
- Better false positive filtering (e.g., "ben" vs "ruben")
- Export to Markdown with syntax highlighting
- Real-time log monitoring
- Batch processing for multiple rounds

## Contributing

Contributions are welcome! Feel free to:
- Report bugs
- Suggest features
- Submit pull requests

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Fyne](https://fyne.io/) for the GUI
- Inspired by SS13 admin tooling needs
- Thanks to the SS13 community

## Author

**Amor4k** - [GitHub](https://github.com/Amor4k)

## Support

For issues, questions, or suggestions:
- Open an issue on [GitHub](https://github.com/Amor4k/Echelon/issues)
- Contact server administration

---

**Note:** ECHELON is designed for legitimate administrative use on Space Station 13 servers. Always respect player privacy and follow your server's admin policies.
