## MeetingMate v0.1.0 - Initial Release

**MeetingMate** is a beautiful CLI tool built with Go and the Charm ecosystem that transforms Google Calendar meeting information into Obsidian-compatible markdown files with rich metadata.

### **What MeetingMate Does**

Simply copy a Google Calendar meeting invite, pipe it through MeetingMate, and get perfectly formatted meeting notes ready for your Obsidian vault or any markdown system.

```bash
# The magic one-liner
pbpaste | meetingmate | pbcopy
```

### **Key Features**

- **Smart Parsing** - Automatically extracts meeting title, date/time, organizer, participants, links, and descriptions from Google Calendar text
- **Rich YAML Front Matter** - Complete metadata including tags, timestamps, organizer, and participant lists for Obsidian integration
- **Beautiful Output** - Intelligent display that provides clean markdown for piping and gorgeous terminal formatting for interactive use
- **Pipeline Perfect** - Seamlessly integrates with clipboard workflows (`pbcopy`/`pbpaste`) without hidden characters or formatting issues
- **Flexible Control** - Optional sections (`--details`, `--attendees`) and output formats (`--plain`) for different use cases
- **Charm Ecosystem** - Built with Glamour for markdown rendering and Lipgloss for styled terminal UI

### **Input Support**
- Copy/paste from Google Calendar invites
- File input (`--input meeting.txt`)
- Stdin piping from clipboard tools
- Interactive terminal input

### **Output Formats**
- **Markdown** (default) - Full Obsidian-compatible format with YAML front matter
- **Plain text** (`--plain`) - Clean text suitable for emails and chat
- **File output** (`--output notes.md`) - Save directly to files
- **Smart terminal display** - Styled front matter boxes and formatted content

### **Perfect For**
- **Obsidian users** who want structured meeting notes with searchable metadata
- **Knowledge workers** who attend lots of meetings and need consistent documentation
- **Teams** who want standardized meeting note formats
- **Anyone** who copies meeting info between tools frequently

### **Usage Examples**

```bash
# Primary workflow: Clean meeting notes to Obsidian
pbpaste | meetingmate | pbcopy

# Minimal output (just essentials)
meetingmate --input meeting.txt

# Full detailed output 
meetingmate --input meeting.txt --details --attendees --output notes.md

# Plain text for sharing
meetingmate --input meeting.txt --plain
```

### **Technical Highlights**
- Built in **Go** for speed and reliability
- **Charm Glamour** for beautiful markdown terminal rendering  
- **Charm Lipgloss** for styled terminal UI components
- **Smart pipe detection** automatically provides clean output for clipboard operations
- **Robust parsing** handles complex meeting formats, attendee statuses, and edge cases
- **Zero dependencies** for end users - single binary deployment

### **What Gets Parsed**
- Meeting titles (including complex formatting)
- Date/time with timezone handling  
- Organizer and "Created by" information
- Participant lists with status (attending, declined, optional, etc.)
- Meeting locations and participant locations
- Phone and video conference details
- Meeting links (Google Meet, Zoom, Teams, etc.)
- All URLs mentioned in meeting details
- Rich meeting descriptions and agendas

### **Output Structure**
```yaml
---
tags: [meeting, organizer-name]
date: 2025-10-25
meeting: 2025-11-04T09:00:00Z
organiser: Sarah Johnson
participants:
  - Sarah Johnson
  - Alex Chen  
  - Maria Garcia
---

# Meeting Title

## Description
[Parsed meeting description]

## Notes  
<!-- Ready for your notes -->

## Links
- [All extracted URLs]

## Action Items
- [ ] [Checkbox format for tasks]
```

### **User Experience**
- **Clean by default** - Minimal output focuses on essentials (no details/attendees unless requested)
- **Beautiful terminal** - Styled YAML front matter in bordered boxes, syntax-highlighted markdown
- **Clipboard friendly** - Automatically detects piped output and removes all formatting characters
- **Consistent formatting** - Standardized structure makes notes searchable and processable

### **Workflow Integration**
MeetingMate is designed for the modern knowledge worker's clipboard-to-vault workflow. No file management, no complex setup - just copy, pipe, and paste into your notes system.

This initial release establishes the core functionality for transforming meeting chaos into structured, searchable knowledge that fits seamlessly into your existing note-taking workflow.

---
