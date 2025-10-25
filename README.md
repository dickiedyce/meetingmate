# MeetingMate

A beautiful CLI tool built with Go and the Charm ecosystem to parse Google Calendar meeting information and convert it into Obsidian-compatible markdown files.

## Features

- 📅 **Parse Google Calendar meeting text** - Extracts all meeting details automatically
- 📝 **Generate Obsidian-compatible markdown** - Complete with YAML front matter
- 🎨 **Beautiful terminal output** - Styled front matter display and Glamour rendering
- 📱 **Multiple input methods** - Files, stdin, or clipboard support
- 🏷️ **Smart metadata extraction** - Auto-generates tags, participants, organizer, and timestamps
- 📋 **Structured output** - Pre-formatted sections for notes, links, and action items
- 🔄 **Intelligent output detection** - Clean markdown for piping, formatted display for terminal
- ⚙️ **Flexible options** - Include/exclude sections as needed
- 📄 **Multiple output formats** - Markdown or plain text

## Installation

### Prerequisites

- Go 1.19 or later

### Build from source

```bash
git clone https://github.com/yourusername/meetingmate.git
cd meetingmate
go mod tidy
go build -o meetingmate
```

Or build directly in VS Code using the "Build MeetingMate" task (Ctrl+Shift+P > Tasks: Run Task)

## Usage

### Basic Usage

```bash
# Read from clipboard and copy clean markdown to clipboard (perfect for Obsidian)
pbpaste | meetingmate | pbcopy

# Read from file and display beautifully formatted output in terminal
meetingmate --input meeting.txt

# Read from file and save to markdown
meetingmate --input meeting.txt --output notes.md

# Plain text output for simple copying
meetingmate --input meeting.txt --plain

# Include optional sections
meetingmate --input meeting.txt --details --attendees
```

### Input Format

MeetingMate expects the text format that you get when copying a Google Calendar event. For example:

```
Digital Lab Architecture weekly
Monday, 27 October⋅14:30 – 15:00
Weekly on Monday

meet.google.com/nzq-jeai-min
Join by phone
‪(GB) +44 20 3873 3170‬ PIN: ‪392 870 041 0894‬#
...
```

### Output Format

The tool generates Obsidian-compatible markdown with comprehensive YAML front matter:

**Front Matter:**
- `tags` - Meeting category and organizer-based tags
- `date` - File creation date
- `meeting` - Parsed meeting timestamp (ISO 8601 format)
- `organiser` - Meeting creator/organizer
- `participants` - Full list of attendee names

**Content Sections:**
- Meeting details (optional with `--details`)
- Attendee list with status and location (optional with `--attendees`)
- Meeting description
- Notes section (placeholder for your notes)
- Links section (automatically extracted URLs)
- Action items (checkbox format)

## Example Output

```markdown
---
tags: [meeting, finja]
date: 2025-10-25
meeting: 2025-10-29T13:00:00Z
organiser: Finja Wrzodek
participants:
  - Finja Wrzodek
  - Lukasz Rakoczy
  - Richard Dyce
  - Michael Gengler
---

# LEAP Technical Meeting

## Description

As we are approaching the phase in the LEAP project where coordination between the systems involved in material management will be needed, the team has decided that we should meet regularly in a technical group.

## Notes

<!-- Add your meeting notes here -->

## Links

- meet.google.com/wxj-xdri-boh

## Action Items

- [ ] 
```

## Command Line Options

- `--input, -i`: Input file containing meeting information
- `--output, -o`: Output markdown file path  
- `--details`: Include the meeting details section (date, time, links, etc.)
- `--attendees`: Include the attendees section with status and locations
- `--plain`: Output plain text without markdown formatting
- `--help, -h`: Show help information
- `--version, -v`: Show version information

## Output Behavior

**Smart Output Detection:**
- **Terminal display**: Beautiful formatted output with styled YAML front matter in a bordered box
- **Piped output**: Clean markdown without ANSI codes (perfect for `pbcopy` → Obsidian)
- **File output**: Raw markdown saved to specified file

**Default Sections:**
- Title, Description, Notes, Links, Action Items (always included)
- Meeting Details and Attendees are **excluded by default** (add `--details` and `--attendees` to include)

## Advanced Usage Examples

```bash
# Minimal output (default - no details/attendees)
pbpaste | meetingmate | pbcopy

# Full detailed output  
meetingmate --input meeting.txt --details --attendees --output full-notes.md

# Plain text for email/chat
meetingmate --input meeting.txt --plain

# Interactive terminal viewing
meetingmate --input meeting.txt --details --attendees

# Pipeline workflow
pbpaste | meetingmate --details | pbcopy
```

## Development

This project uses:

- [Charm Glamour](https://github.com/charmbracelet/glamour) - Beautiful markdown rendering in terminal
- [Charm Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal UI styling and layout
- [go-isatty](https://github.com/mattn/go-isatty) - Smart pipe detection for clean output

### Building

```bash
go mod tidy
go build -o meetingmate
```

### VS Code Development

The project includes VS Code tasks for development:

- **"Build MeetingMate"** - Compiles the binary
- **"Run MeetingMate with Sample"** - Tests with provided sample data
- **"Run MeetingMate Help"** - Shows help documentation

Use `Ctrl+Shift+P` → `Tasks: Run Task` to access these.

### Testing

```bash
# Test with sample data
./meetingmate --input sample_meeting.txt

# Test pipeline workflow  
echo "Meeting content" | ./meetingmate --plain

# Test all features
./meetingmate --input sample_meeting.txt --details --attendees --output test.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.