package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
)

const version = "0.1.0"

var (
	inputFlag     = flag.String("input", "", "Input file containing meeting information")
	outputFlag    = flag.String("output", "", "Output markdown file path (optional)")
	helpFlag      = flag.Bool("help", false, "Show help information")
	versionFlag   = flag.Bool("version", false, "Show version information")
	detailsFlag   = flag.Bool("details", false, "Include the meeting details section")
	attendeesFlag = flag.Bool("attendees", false, "Include the attendees section")
	plainFlag     = flag.Bool("plain", false, "Output plain text without markdown formatting")
)

// Meeting represents parsed meeting information
type Meeting struct {
	Title         string
	DateTime      string
	MeetingTime   time.Time
	Duration      string
	Frequency     string
	Location      string
	MeetLink      string
	PhoneInfo     string
	Organizer     string
	Attendees     []Attendee
	Description   string
	Links         []string
}

// Attendee represents meeting participant information
type Attendee struct {
	Name     string
	Status   string
	Location string
	Notes    string
}

func main() {
	flag.Parse()

	if *helpFlag {
		showHelp()
		return
	}

	if *versionFlag {
		fmt.Printf("meetingmate v%s\n", version)
		return
	}

	var input string
	var err error

	if *inputFlag != "" {
		// Read from file
		content, err := os.ReadFile(*inputFlag)
		if err != nil {
			log.Fatalf("Error reading input file: %v", err)
		}
		input = string(content)
	} else {
		// Read from stdin
		input, err = readFromStdin()
		if err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
	}

	if strings.TrimSpace(input) == "" {
		fmt.Println("No input provided. Use --help for usage information.")
		return
	}

	meeting, err := parseMeeting(input)
	if err != nil {
		log.Fatalf("Error parsing meeting information: %v", err)
	}

	var output string
	if *plainFlag {
		output = generatePlainText(meeting, *detailsFlag, *attendeesFlag)
	} else {
		output = generateMarkdown(meeting, *detailsFlag, *attendeesFlag)
	}

	if *outputFlag != "" {
		err := os.WriteFile(*outputFlag, []byte(output), 0644)
		if err != nil {
			log.Fatalf("Error writing to output file: %v", err)
		}
		fmt.Printf("Meeting notes saved to: %s\n", *outputFlag)
	} else {
		// Check if output is being piped (not a terminal)
		isTerminal := isatty.IsTerminal(os.Stdout.Fd())
		
		if *plainFlag || !isTerminal {
			// Plain text output or piped output (like pbcopy)
			if *plainFlag {
				fmt.Print(output)
			} else {
				// Piped markdown without terminal formatting
				fmt.Print(output)
			}
		} else {
			// For terminal display, we need to handle YAML front matter properly
			// Glamour doesn't render YAML front matter well, so we'll extract it
			lines := strings.Split(output, "\n")
			var frontMatter []string
			var markdownContent []string
			inFrontMatter := false
			frontMatterClosed := false
			
			for i, line := range lines {
				if i == 0 && line == "---" {
					inFrontMatter = true
					frontMatter = append(frontMatter, line)
				} else if inFrontMatter && line == "---" {
					frontMatter = append(frontMatter, line)
					inFrontMatter = false
					frontMatterClosed = true
				} else if inFrontMatter {
					frontMatter = append(frontMatter, line)
				} else {
					markdownContent = append(markdownContent, line)
				}
			}
			
			// Display front matter in a styled box
			if frontMatterClosed && len(frontMatter) > 0 {
				frontMatterStyle := lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(lipgloss.Color("#7D56F4")).
					Padding(1).
					Margin(1)
				
				fmt.Print(frontMatterStyle.Render(strings.Join(frontMatter, "\n")))
				fmt.Print("\n")
			}
			
			// Render the rest with glamour
			markdownOnly := strings.Join(markdownContent, "\n")
			if strings.TrimSpace(markdownOnly) != "" {
				renderer, err := glamour.NewTermRenderer(
					glamour.WithAutoStyle(),
					glamour.WithWordWrap(80),
				)
				if err != nil {
					fmt.Print(markdownOnly)
					return
				}

				out, err := renderer.Render(markdownOnly)
				if err != nil {
					fmt.Print(markdownOnly)
					return
				}
				fmt.Print(out)
			}
		}
	}
}

func showHelp() {
	helpStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7D56F4"))

	fmt.Print(helpStyle.Render("MeetingMate") + " - Convert Google Calendar meeting info to Obsidian markdown\n\n")

	fmt.Println("USAGE:")
	fmt.Println("  meetingmate [FLAGS] [< input.txt]")
	fmt.Println("  meetingmate --input meeting.txt --output notes.md")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  --input, -i      Input file containing meeting information")
	fmt.Println("  --output, -o     Output markdown file path")
	fmt.Println("  --details        Include the meeting details section")
	fmt.Println("  --attendees      Include the attendees section")
	fmt.Println("  --plain          Output plain text without markdown formatting")
	fmt.Println("  --help, -h       Show this help message")
	fmt.Println("  --version, -v    Show version information")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Copy clean markdown from clipboard to clipboard (perfect for Obsidian)")
	fmt.Println("  pbpaste | meetingmate | pbcopy")
	fmt.Println()
	fmt.Println("  # Read from file with minimal output (default)")
	fmt.Println("  meetingmate --input meeting.txt")
	fmt.Println()
	fmt.Println("  # Full detailed output with all sections")
	fmt.Println("  meetingmate --input meeting.txt --details --attendees --output notes.md")
	fmt.Println()
	fmt.Println("  # Plain text output for copying to email/chat")
	fmt.Println("  meetingmate --input meeting.txt --plain")
}

func readFromStdin() (string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return strings.Join(lines, "\n"), nil
}

func parseMeeting(input string) (*Meeting, error) {
	lines := strings.Split(input, "\n")
	meeting := &Meeting{
		Attendees: make([]Attendee, 0),
		Links:     make([]string, 0),
	}

	var currentSection string
	var descriptionLines []string
	inDescription := false

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse title (first non-empty line)
		if meeting.Title == "" && line != "" {
			meeting.Title = line
			continue
		}

		// Parse date/time
		if strings.Contains(line, "⋅") && strings.Contains(line, "–") {
			meeting.DateTime = line
			meeting.MeetingTime = parseMeetingTime(line)
			continue
		}

		// Parse frequency
		if strings.HasPrefix(line, "Weekly") || strings.HasPrefix(line, "Daily") || strings.HasPrefix(line, "Monthly") {
			meeting.Frequency = line
			continue
		}

		// Parse meeting links
		if strings.Contains(line, "meet.google.com") || strings.Contains(line, "zoom.us") || strings.Contains(line, "teams.microsoft.com") {
			meeting.MeetLink = line
			meeting.Links = append(meeting.Links, line)
			continue
		}

		// Detect other URLs (http/https)
		if strings.Contains(line, "http://") || strings.Contains(line, "https://") {
			// Extract URLs from the line
			words := strings.Fields(line)
			for _, word := range words {
				if strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
					meeting.Links = append(meeting.Links, word)
				}
			}
			// If this line is just a URL, don't process it as attendee
			if len(strings.Fields(line)) == 1 && (strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://")) {
				continue
			}
		}

		// Parse phone information
		if strings.Contains(line, "Join by phone") {
			currentSection = "phone"
			continue
		}
		if currentSection == "phone" && (strings.Contains(line, "+") || strings.Contains(line, "PIN:") || strings.Contains(line, "ID:")) {
			if meeting.PhoneInfo == "" {
				meeting.PhoneInfo = line
			} else {
				meeting.PhoneInfo += "\n" + line
			}
			continue
		}

		// Parse guest count
		if strings.Contains(line, "guests") || strings.Contains(line, "yes") || strings.Contains(line, "no") || strings.Contains(line, "maybe") {
			currentSection = "guests"
			continue
		}

		// Parse organizer
		if strings.Contains(line, "Organiser") || strings.Contains(line, "Organizer") {
			if i+1 < len(lines) {
				meeting.Organizer = strings.TrimSpace(lines[i+1])
			}
			continue
		}

		// Parse "Created by:" 
		if strings.HasPrefix(line, "Created by:") {
			createdBy := strings.TrimSpace(strings.TrimPrefix(line, "Created by:"))
			if createdBy != "" {
				meeting.Organizer = createdBy
			}
			continue
		}

		// Detect start of description
		if strings.Contains(line, "Hi,") || (len(line) > 50 && strings.Contains(line, " ")) {
			inDescription = true
			descriptionLines = append(descriptionLines, line)
			continue
		}

		if inDescription {
			descriptionLines = append(descriptionLines, line)
			continue
		}

		// Parse attendees
		if line != meeting.Title && !strings.Contains(line, "⋅") && !strings.Contains(line, "meet.") {
			attendee := parseAttendee(line, lines, i)
			if attendee.Name != "" {
				meeting.Attendees = append(meeting.Attendees, attendee)
			}
		}
	}

	if len(descriptionLines) > 0 {
		meeting.Description = strings.Join(descriptionLines, "\n")
	}

	return meeting, nil
}

func parseAttendee(line string, lines []string, index int) Attendee {
	attendee := Attendee{}

	// Skip obvious non-attendee lines
	if strings.Contains(line, "guests") || strings.Contains(line, "yes") || 
		strings.Contains(line, "awaiting") || strings.Contains(line, "Edit") ||
		strings.Contains(line, "More joining") || strings.Contains(line, "http") ||
		strings.Contains(line, "@") || strings.Contains(line, "event_busy") ||
		strings.Contains(line, "Out of office") || strings.Contains(line, "bedtime") ||
		strings.Contains(line, "Outside working hours") || line == "–" ||
		line == "Home" || line == "Office" || line == "Scotland" || 
		strings.Contains(line, "Declined because") || len(line) < 3 ||
		strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
		return attendee
	}

	// Check if this looks like a person's name (contains letters and potentially spaces)
	if len(line) > 0 && !strings.Contains(line, "meet.") && !strings.Contains(line, ".com") {
		// Simple heuristic: if it's 2-4 words and doesn't contain numbers/symbols, likely a name
		words := strings.Fields(line)
		if len(words) >= 1 && len(words) <= 4 {
			hasNumbers := false
			for _, word := range words {
				if strings.ContainsAny(word, "0123456789@.") {
					hasNumbers = true
					break
				}
			}
			
			if !hasNumbers {
				attendee.Name = line
				
				// Look ahead for status and location info
				if index+1 < len(lines) {
					nextLine := strings.TrimSpace(lines[index+1])
					if strings.Contains(nextLine, "event_busy") || strings.Contains(nextLine, "Out of office") {
						attendee.Status = "Declined"
					} else if strings.Contains(nextLine, "Optional") {
						attendee.Status = "Optional"
					} else if nextLine == "Home" || nextLine == "Office" || nextLine == "Scotland" {
						attendee.Location = nextLine
					}
				}
			}
		}
	}

	return attendee
}

func parseMeetingTime(dateTimeStr string) time.Time {
	// Parse formats like "Monday, 27 October⋅14:30 – 15:00"
	
	// Split by the bullet point (⋅)
	parts := strings.Split(dateTimeStr, "⋅")
	if len(parts) != 2 {
		return time.Time{} // Return zero time if parsing fails
	}
	
	datePart := strings.TrimSpace(parts[0])
	timePart := strings.TrimSpace(parts[1])
	
	// Extract start time (before the dash)
	timeRange := strings.Split(timePart, "–")
	if len(timeRange) == 0 {
		return time.Time{}
	}
	startTime := strings.TrimSpace(timeRange[0])
	
	// Parse the date part - try to extract day and month
	// Format: "Monday, 27 October"
	dateFields := strings.Fields(datePart)
	if len(dateFields) < 3 {
		return time.Time{}
	}
	
	dayStr := strings.TrimSuffix(dateFields[1], ",")
	monthStr := dateFields[2]
	
	// Convert month name to number
	monthMap := map[string]int{
		"January": 1, "February": 2, "March": 3, "April": 4,
		"May": 5, "June": 6, "July": 7, "August": 8,
		"September": 9, "October": 10, "November": 11, "December": 12,
	}
	
	month, exists := monthMap[monthStr]
	if !exists {
		return time.Time{}
	}
	
	// Parse day
	day := 1
	if d, err := time.Parse("2", dayStr); err == nil {
		day = d.Day()
	}
	
	// Parse time (format: "14:30")
	timeParts := strings.Split(startTime, ":")
	if len(timeParts) != 2 {
		return time.Time{}
	}
	
	hour := 0
	minute := 0
	if h, err := time.Parse("15", timeParts[0]); err == nil {
		hour = h.Hour()
	}
	if m, err := time.Parse("04", timeParts[1]); err == nil {
		minute = m.Minute()
	}
	
	// Use current year as default
	year := time.Now().Year()
	
	// Create the time
	meetingTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	return meetingTime
}

func generateMarkdown(meeting *Meeting, includeDetails bool, includeAttendees bool) string {
	var md strings.Builder

	// Front matter at the top
	md.WriteString("---\n")
	md.WriteString("tags: [meeting")
	if meeting.Organizer != "" {
		// Simple tag generation from organizer name
		orgTag := strings.ToLower(strings.ReplaceAll(strings.Split(meeting.Organizer, " ")[0], " ", "-"))
		md.WriteString(fmt.Sprintf(", %s", orgTag))
	}
	md.WriteString("]\n")
	md.WriteString(fmt.Sprintf("date: %s\n", time.Now().Format("2006-01-02")))
	
	// Add meeting timestamp if parsed successfully
	if !meeting.MeetingTime.IsZero() {
		md.WriteString(fmt.Sprintf("meeting: %s\n", meeting.MeetingTime.Format("2006-01-02T15:04:05Z07:00")))
	}
	
	// Add organiser if present
	if meeting.Organizer != "" {
		md.WriteString(fmt.Sprintf("organiser: %s\n", meeting.Organizer))
	}
	
	// Add participants list
	if len(meeting.Attendees) > 0 || meeting.Organizer != "" {
		md.WriteString("participants:\n")
		
		// Add organizer first if present
		if meeting.Organizer != "" {
			md.WriteString(fmt.Sprintf("  - %s\n", meeting.Organizer))
		}
		
		// Add attendees
		for _, attendee := range meeting.Attendees {
			if attendee.Name != "" && attendee.Name != meeting.Organizer {
				md.WriteString(fmt.Sprintf("  - %s\n", attendee.Name))
			}
		}
	}
	
	md.WriteString("---\n")

	// Title
	md.WriteString(fmt.Sprintf("\n# %s\n", meeting.Title))

	// Meeting details (optional)
	if includeDetails {
		md.WriteString("\n## Meeting Details\n")
		
		if meeting.DateTime != "" {
			md.WriteString(fmt.Sprintf("**Date & Time:** %s\n", meeting.DateTime))
		}
		
		if meeting.Frequency != "" {
			md.WriteString(fmt.Sprintf("**Frequency:** %s\n", meeting.Frequency))
		}

		if meeting.MeetLink != "" {
			md.WriteString(fmt.Sprintf("**Meeting Link:** %s\n", meeting.MeetLink))
		}

		if meeting.PhoneInfo != "" {
			md.WriteString("**Phone Information:**\n```\n")
			md.WriteString(meeting.PhoneInfo)
			md.WriteString("\n```\n")
		}

		if meeting.Organizer != "" {
			md.WriteString(fmt.Sprintf("**Organizer:** %s\n", meeting.Organizer))
		}
	}

	// Attendees (optional)
	if includeAttendees && len(meeting.Attendees) > 0 {
		md.WriteString("\n## Attendees\n")
		for _, attendee := range meeting.Attendees {
			md.WriteString(fmt.Sprintf("- **%s**", attendee.Name))
			if attendee.Status != "" {
				md.WriteString(fmt.Sprintf(" (%s)", attendee.Status))
			}
			if attendee.Location != "" {
				md.WriteString(fmt.Sprintf(" - %s", attendee.Location))
			}
			md.WriteString("\n")
		}
	}

	// Description
	if meeting.Description != "" {
		md.WriteString("\n## Description\n")
		md.WriteString(meeting.Description)
		md.WriteString("\n")
	}

	// Notes section
	md.WriteString("\n## Notes\n")
	md.WriteString("<!-- Add your meeting notes here -->\n")

	// Links
	if len(meeting.Links) > 0 {
		md.WriteString("\n## Links\n")
		for _, link := range meeting.Links {
			md.WriteString(fmt.Sprintf("- %s\n", link))
		}
	}

	// Action items
	md.WriteString("\n## Action Items\n")
	md.WriteString("- [ ] \n")

	return md.String()
}

func generatePlainText(meeting *Meeting, includeDetails bool, includeAttendees bool) string {
	var text strings.Builder

	// Title
	text.WriteString(fmt.Sprintf("%s\n", meeting.Title))
	text.WriteString(strings.Repeat("=", len(meeting.Title)) + "\n")

	// Meeting details (optional)
	if includeDetails {
		text.WriteString("\nMeeting Details:\n")
		
		if meeting.DateTime != "" {
			text.WriteString(fmt.Sprintf("Date & Time: %s\n", meeting.DateTime))
		}
		
		if meeting.Frequency != "" {
			text.WriteString(fmt.Sprintf("Frequency: %s\n", meeting.Frequency))
		}

		if meeting.MeetLink != "" {
			text.WriteString(fmt.Sprintf("Meeting Link: %s\n", meeting.MeetLink))
		}

		if meeting.PhoneInfo != "" {
			text.WriteString("Phone Information:\n")
			text.WriteString(meeting.PhoneInfo)
			text.WriteString("\n")
		}

		if meeting.Organizer != "" {
			text.WriteString(fmt.Sprintf("Organizer: %s\n", meeting.Organizer))
		}
	}

	// Attendees (optional)
	if includeAttendees && len(meeting.Attendees) > 0 {
		text.WriteString("\nAttendees:\n")
		for _, attendee := range meeting.Attendees {
			text.WriteString(fmt.Sprintf("- %s", attendee.Name))
			if attendee.Status != "" {
				text.WriteString(fmt.Sprintf(" (%s)", attendee.Status))
			}
			if attendee.Location != "" {
				text.WriteString(fmt.Sprintf(" - %s", attendee.Location))
			}
			text.WriteString("\n")
		}
	}

	// Description
	if meeting.Description != "" {
		text.WriteString("\nDescription:\n")
		text.WriteString(meeting.Description)
		text.WriteString("\n")
	}

	// Notes section
	text.WriteString("\nNotes:\n")
	text.WriteString("(Add your meeting notes here)\n")

	// Links
	if len(meeting.Links) > 0 {
		text.WriteString("\nLinks:\n")
		for _, link := range meeting.Links {
			text.WriteString(fmt.Sprintf("- %s\n", link))
		}
	}

	// Action items
	text.WriteString("\nAction Items:\n")
	text.WriteString("- [ ] \n")

	return text.String()
}