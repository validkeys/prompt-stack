package history

import (
	"fmt"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// Search handles full-text search operations with highlighting
type Search struct {
	db     *Database
	logger *zap.Logger
}

// NewSearch creates a new search handler
func NewSearch(db *Database, logger *zap.Logger) *Search {
	return &Search{
		db:     db,
		logger: logger,
	}
}

// SearchResult represents a search result with highlighted matches
type SearchResult struct {
	Composition
	HighlightedContent string
	MatchCount         int
	MatchPositions     []MatchPosition
}

// MatchPosition represents the position of a match in the text
type MatchPosition struct {
	Start int
	End   int
	Text  string
}

// SearchQuery performs a full-text search with query highlighting
func (s *Search) SearchQuery(query string) ([]SearchResult, error) {
	s.logger.Debug("Performing search", zap.String("query", query))

	// Clean and validate query
	cleanQuery := s.cleanQuery(query)
	if cleanQuery == "" {
		return []SearchResult{}, nil
	}

	// Perform FTS5 search
	compositions, err := s.db.SearchCompositions(cleanQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to search compositions: %w", err)
	}

	// Process results with highlighting
	var results []SearchResult
	for _, comp := range compositions {
		highlighted, matches := s.highlightMatches(comp.Content, cleanQuery)

		result := SearchResult{
			Composition:        comp,
			HighlightedContent: highlighted,
			MatchCount:         len(matches),
			MatchPositions:     matches,
		}

		results = append(results, result)
	}

	s.logger.Info("Search completed",
		zap.String("query", query),
		zap.Int("results", len(results)),
	)

	return results, nil
}

// cleanQuery cleans and normalizes the search query
func (s *Search) cleanQuery(query string) string {
	// Trim whitespace
	query = strings.TrimSpace(query)

	// Remove special FTS5 operators that might cause issues
	// Keep basic operators: AND, OR, NOT, quotes
	query = regexp.MustCompile(`[^\w\s"\'\-]`).ReplaceAllString(query, " ")

	// Collapse multiple spaces
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

	return strings.TrimSpace(query)
}

// highlightMatches highlights search terms in the content
func (s *Search) highlightMatches(content, query string) (string, []MatchPosition) {
	// Extract search terms from query
	terms := s.extractSearchTerms(query)
	if len(terms) == 0 {
		return content, []MatchPosition{}
	}

	// Find all matches
	var matches []MatchPosition
	contentLower := strings.ToLower(content)

	for _, term := range terms {
		termLower := strings.ToLower(term)
		start := 0

		for {
			// Find next occurrence
			idx := strings.Index(contentLower[start:], termLower)
			if idx == -1 {
				break
			}

			// Calculate actual position
			pos := start + idx

			// Check if this position overlaps with existing matches
			if !s.overlapsWithExisting(matches, pos, pos+len(term)) {
				matches = append(matches, MatchPosition{
					Start: pos,
					End:   pos + len(term),
					Text:  content[pos : pos+len(term)],
				})
			}

			start = pos + len(term)
		}
	}

	// Sort matches by position
	s.sortMatches(matches)

	// Apply highlighting
	highlighted := s.applyHighlighting(content, matches)

	return highlighted, matches
}

// extractSearchTerms extracts individual search terms from the query
func (s *Search) extractSearchTerms(query string) []string {
	var terms []string

	// Remove quotes and split by operators
	query = regexp.MustCompile(`["']`).ReplaceAllString(query, "")

	// Split by AND, OR, NOT operators
	parts := regexp.MustCompile(`\s+(AND|OR|NOT)\s+`).Split(query, -1)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && !s.isOperator(part) {
			terms = append(terms, part)
		}
	}

	return terms
}

// isOperator checks if a string is a search operator
func (s *Search) isOperator(term string) bool {
	upper := strings.ToUpper(term)
	return upper == "AND" || upper == "OR" || upper == "NOT"
}

// overlapsWithExisting checks if a match overlaps with existing matches
func (s *Search) overlapsWithExisting(matches []MatchPosition, start, end int) bool {
	for _, m := range matches {
		if !(end <= m.Start || start >= m.End) {
			return true
		}
	}
	return false
}

// sortMatches sorts matches by position
func (s *Search) sortMatches(matches []MatchPosition) {
	for i := 0; i < len(matches); i++ {
		for j := i + 1; j < len(matches); j++ {
			if matches[i].Start > matches[j].Start {
				matches[i], matches[j] = matches[j], matches[i]
			}
		}
	}
}

// applyHighlighting applies highlighting markers to the content
func (s *Search) applyHighlighting(content string, matches []MatchPosition) string {
	if len(matches) == 0 {
		return content
	}

	var result strings.Builder
	lastPos := 0

	for _, match := range matches {
		// Add content before match
		result.WriteString(content[lastPos:match.Start])

		// Add highlighted match
		result.WriteString("<<")
		result.WriteString(match.Text)
		result.WriteString(">>")

		lastPos = match.End
	}

	// Add remaining content
	result.WriteString(content[lastPos:])

	return result.String()
}

// SearchByDirectory performs search within a specific working directory
func (s *Search) SearchByDirectory(query, workingDir string) ([]SearchResult, error) {
	s.logger.Debug("Searching by directory",
		zap.String("query", query),
		zap.String("working_directory", workingDir),
	)

	// Get compositions by directory
	compositions, err := s.db.GetCompositionsByDirectory(workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get compositions by directory: %w", err)
	}

	// Filter and highlight matches
	var results []SearchResult
	cleanQuery := s.cleanQuery(query)

	for _, comp := range compositions {
		// Check if composition matches query
		if s.matchesQuery(comp.Content, cleanQuery) {
			highlighted, matches := s.highlightMatches(comp.Content, cleanQuery)

			result := SearchResult{
				Composition:        comp,
				HighlightedContent: highlighted,
				MatchCount:         len(matches),
				MatchPositions:     matches,
			}

			results = append(results, result)
		}
	}

	return results, nil
}

// matchesQuery checks if content matches the search query
func (s *Search) matchesQuery(content, query string) bool {
	terms := s.extractSearchTerms(query)
	if len(terms) == 0 {
		return true
	}

	contentLower := strings.ToLower(content)

	for _, term := range terms {
		if strings.Contains(contentLower, strings.ToLower(term)) {
			return true
		}
	}

	return false
}

// GetMatchContext returns context around a match for preview
func (s *Search) GetMatchContext(content string, match MatchPosition, contextLength int) string {
	start := match.Start - contextLength
	if start < 0 {
		start = 0
	}

	end := match.End + contextLength
	if end > len(content) {
		end = len(content)
	}

	return content[start:end]
}
