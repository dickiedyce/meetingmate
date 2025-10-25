# MeetingMate Project Instructions

This is a Go CLI application using the Charm ecosystem (Glamour, Lipgloss) to parse Google Calendar meeting information and convert it to Obsidian-compatible markdown files.

## Project Status: ✅ COMPLETED

- [x] Verify that the copilot-instructions.md file in the .github directory is created.
- [x] Clarify Project Requirements - Go CLI with Charm ecosystem for meeting parsing
- [x] Scaffold the Project - Go module structure with main.go and dependencies
- [x] Customize the Project - Full CLI implementation with parsing logic
- [x] Install Required Extensions - Go extension installed for VS Code
- [x] Compile the Project - Successfully builds without errors
- [x] Create and Run Task - VS Code tasks for build and run created
- [x] Launch the Project - CLI tested with sample data, help, and version commands
- [x] Ensure Documentation is Complete - README.md updated with usage instructions

## Key Features

- **Smart Google Calendar parsing** - Extracts meeting details, participants, organizer, and timestamps
- **Rich YAML front matter** - Complete metadata for Obsidian integration
- **Intelligent output detection** - Clean markdown for piping, beautiful formatting for terminal
- **Flexible section control** - Include/exclude details and attendees with flags
- **Multiple output formats** - Markdown and plain text options
- **Beautiful terminal UI** - Styled front matter display with Lipgloss, Glamour rendering
- **Pipeline-friendly** - Perfect for clipboard workflows with pbcopy/pbpaste

## Usage

```bash
# Primary workflow: Clipboard to Obsidian
pbpaste | ./meetingmate | pbcopy

# Minimal output (default - no details/attendees)
./meetingmate --input sample_meeting.txt

# Full detailed output
./meetingmate --input sample_meeting.txt --details --attendees --output notes.md

# Plain text for email/chat
./meetingmate --input sample_meeting.txt --plain

# Show help
./meetingmate --help
```

## Current Features Status

✅ **Smart parsing** - Extracts title, datetime, organizer, participants, links, description  
✅ **YAML front matter** - tags, date, meeting timestamp, organiser, participants list  
✅ **Flexible sections** - Optional details and attendees sections  
✅ **Output formats** - Markdown (default) and plain text (--plain)  
✅ **Smart output** - Auto-detects pipes vs terminal for clean clipboard operations  
✅ **Beautiful terminal UI** - Styled front matter box, Glamour markdown rendering  
✅ **VS Code integration** - Build, run, and help tasks available

## Development

Use the VS Code tasks:
- **"Build MeetingMate"** - Compiles the binary (`go build -o meetingmate`)
- **"Run MeetingMate with Sample"** - Tests with sample data (includes --details --attendees)
- **"Run MeetingMate Help"** - Shows help documentation

The tool automatically handles different output contexts and provides the appropriate formatting for each use case.