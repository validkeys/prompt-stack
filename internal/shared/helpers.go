package shared

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileExists reports whether the given path exists (file or dir).
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// ReadJSONFile reads a JSON file and unmarshals into map[string]interface{}.
func ReadJSONFile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out map[string]interface{}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ReadFileString returns file contents as a string.
func ReadFileString(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Contains is a thin wrapper around strings.Contains to mirror prior helpers.
func Contains(s, substr string) bool { return strings.Contains(s, substr) }

// FindLineWithPrefix returns the first trimmed line that begins with prefix.
func FindLineWithPrefix(content, prefix string) string {
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

// CountOccurrences returns the number of non-overlapping occurrences of substr in s.
func CountOccurrences(s, substr string) int {
	if substr == "" {
		return 0
	}
	count := 0
	idx := 0
	for {
		i := strings.Index(s[idx:], substr)
		if i == -1 {
			break
		}
		count++
		idx += i + len(substr)
	}
	return count
}

// SplitLines wraps strings.Split but guarantees consistent behavior.
func SplitLines(s string) []string { return strings.Split(s, "\n") }

// JoinLines wraps strings.Join for string slices.
func JoinLines(lines []string) string { return strings.Join(lines, "\n") }

// TrimWhitespace trims common whitespace characters from both ends.
func TrimWhitespace(s string) string { return strings.TrimSpace(s) }

// ExtractSectionBetween returns the substring from startMarker up to (but not including) endMarker.
// If endMarker is not found, returns everything from startMarker.
func ExtractSectionBetween(content, startMarker, endMarker string) string {
	startIdx := strings.Index(content, startMarker)
	if startIdx < 0 {
		return ""
	}
	remaining := content[startIdx:]
	endIdx := strings.Index(remaining, endMarker)
	if endIdx < 0 {
		return remaining
	}
	return remaining[:endIdx]
}

// CountFileEntries heuristically counts quoted file path lines inside a section.
func CountFileEntries(section string) int {
	count := 0
	for _, line := range SplitLines(section) {
		trimmed := TrimWhitespace(line)
		// consider a quoted string on its own line a file entry
		if len(trimmed) > 0 && strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"") {
			count++
			continue
		}
		// or a dash followed by a quoted path: - "path"
		if strings.HasPrefix(trimmed, "- \"") && strings.HasSuffix(trimmed, "\"") {
			count++
		}
	}
	return count
}

var idRE = regexp.MustCompile(`-?\s*id:\s*\"?([A-Za-z0-9\-_.]+)\"?`)

// ExtractID attempts to pull an id token from a line like `- id: "m0-001"`.
func ExtractID(line string) string {
	m := idRE.FindStringSubmatch(line)
	if len(m) >= 2 {
		return m[1]
	}
	// fallback: scan for quoted token
	q := ""
	for _, part := range strings.Fields(line) {
		if strings.HasPrefix(part, "\"") && strings.HasSuffix(part, "\"") {
			q = strings.Trim(part, "\"")
			break
		}
	}
	return q
}

// ExtractTaskSection extracts a task block starting at `- id: "<taskID>"` until the next `- id:` at same indent.
func ExtractTaskSection(content, taskID string) string {
	taskStart := fmt.Sprintf("- id: \"%s\"", taskID)
	startIdx := strings.Index(content, taskStart)
	if startIdx < 0 {
		return ""
	}
	lines := SplitLines(content[startIdx:])
	var section []string
	for _, line := range lines {
		if strings.Contains(line, "- id:") && !strings.Contains(line, taskID) {
			break
		}
		section = append(section, line)
	}
	return JoinLines(section)
}

// FindRepoRoot walks up from cwd to find a directory containing go.mod.
func FindRepoRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// ReadDirFiles returns a list of files under a directory (non-recursive) that match pattern.
func ReadDirFiles(dir, pattern string) ([]fs.DirEntry, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if pattern == "" {
		return ents, nil
	}
	var out []fs.DirEntry
	for _, e := range ents {
		if matched, _ := filepath.Match(pattern, e.Name()); matched {
			out = append(out, e)
		}
	}
	return out, nil
}
