package ai

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// ContextSelector handles intelligent selection of library prompts for AI context
type ContextSelector struct {
	// Configuration
	maxContextPrompts int // Maximum number of prompts to include in context
}

// NewContextSelector creates a new context selector
func NewContextSelector() *ContextSelector {
	return &ContextSelector{
		maxContextPrompts: 5, // Default: top 5 prompts
	}
}

// IndexedPrompt represents a prompt in the library index
type IndexedPrompt struct {
	PromptID      string
	Title         string
	Description   string
	Tags          []string
	Category      string
	WordFrequency map[string]int
	LastUsed      time.Time
	UseCount      int
	Content       string
}

// PromptScore represents a prompt with its relevance score
type PromptScore struct {
	Prompt    *IndexedPrompt
	Score     int
	Reasoning []string // Human-readable explanation of score
}

// ScorePrompts scores library prompts based on composition keywords and metadata
func (cs *ContextSelector) ScorePrompts(
	prompts []*IndexedPrompt,
	compositionKeywords map[string]int,
	compositionTags []string,
	compositionCategory string,
) []*PromptScore {
	var scored []*PromptScore

	for _, prompt := range prompts {
		score, reasoning := cs.calculatePromptScore(
			prompt,
			compositionKeywords,
			compositionTags,
			compositionCategory,
		)

		scored = append(scored, &PromptScore{
			Prompt:    prompt,
			Score:     score,
			Reasoning: reasoning,
		})
	}

	// Sort by score (descending)
	cs.sortByScore(scored)

	return scored
}

// calculatePromptScore calculates a single prompt's relevance score
func (cs *ContextSelector) calculatePromptScore(
	prompt *IndexedPrompt,
	compositionKeywords map[string]int,
	compositionTags []string,
	compositionCategory string,
) (int, []string) {
	var score int
	var reasoning []string

	// 1. Tag matches: +10 per matching tag
	tagMatches := cs.countTagMatches(prompt.Tags, compositionTags)
	if tagMatches > 0 {
		tagScore := tagMatches * 10
		score += tagScore
		reasoning = append(reasoning, fmt.Sprintf("Tag matches: +%d (%d tags)", tagScore, tagMatches))
	}

	// 2. Category bonus: +5 if same category
	if compositionCategory != "" && prompt.Category == compositionCategory {
		score += 5
		reasoning = append(reasoning, "Category match: +5")
	}

	// 3. Keyword overlap: +1 per matching word (weighted by frequency)
	keywordScore := cs.calculateKeywordScore(prompt.WordFrequency, compositionKeywords)
	if keywordScore > 0 {
		score += keywordScore
		reasoning = append(reasoning, fmt.Sprintf("Keyword overlap: +%d", keywordScore))
	}

	// 4. Recently used bonus: +3 if used in last session
	if cs.isRecentlyUsed(prompt.LastUsed) {
		score += 3
		reasoning = append(reasoning, "Recently used: +3")
	}

	// 5. Frequently used bonus: +use_count
	if prompt.UseCount > 0 {
		score += prompt.UseCount
		reasoning = append(reasoning, fmt.Sprintf("Frequently used: +%d", prompt.UseCount))
	}

	return score, reasoning
}

// countTagMatches counts how many tags match between prompt and composition
func (cs *ContextSelector) countTagMatches(promptTags, compositionTags []string) int {
	tagSet := make(map[string]bool)
	for _, tag := range compositionTags {
		tagSet[strings.ToLower(tag)] = true
	}

	matches := 0
	for _, tag := range promptTags {
		if tagSet[strings.ToLower(tag)] {
			matches++
		}
	}

	return matches
}

// calculateKeywordScore calculates keyword overlap score weighted by frequency
func (cs *ContextSelector) calculateKeywordScore(
	promptKeywords, compositionKeywords map[string]int,
) int {
	score := 0

	for keyword, promptFreq := range promptKeywords {
		if compFreq, exists := compositionKeywords[keyword]; exists {
			// Weight by frequency in both prompt and composition
			// Higher frequency = more important keyword
			weight := (promptFreq + compFreq) / 2
			score += weight
		}
	}

	return score
}

// isRecentlyUsed checks if prompt was used in the last session (last 24 hours)
func (cs *ContextSelector) isRecentlyUsed(lastUsed time.Time) bool {
	if lastUsed.IsZero() {
		return false
	}

	// Consider "recently used" as within the last 24 hours
	recentThreshold := time.Now().Add(-24 * time.Hour)
	return lastUsed.After(recentThreshold)
}

// sortByScore sorts prompts by score in descending order
func (cs *ContextSelector) sortByScore(scored []*PromptScore) {
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].Score > scored[i].Score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}
}

// SelectTopPrompts selects top N prompts that fit within token budget
func (cs *ContextSelector) SelectTopPrompts(
	scored []*PromptScore,
	tokenBudget int,
) []*IndexedPrompt {
	var selected []*IndexedPrompt
	totalTokens := 0

	for _, scoredPrompt := range scored {
		// Estimate tokens for this prompt
		promptTokens := cs.estimateTokens(scoredPrompt.Prompt.Content)

		// Check if adding this prompt would exceed budget
		if totalTokens+promptTokens > tokenBudget {
			break
		}

		selected = append(selected, scoredPrompt.Prompt)
		totalTokens += promptTokens

		// Stop if we've reached max prompts
		if len(selected) >= cs.maxContextPrompts {
			break
		}
	}

	return selected
}

// estimateTokens provides a rough token estimation (~4 chars = 1 token)
func (cs *ContextSelector) estimateTokens(text string) int {
	// Rough estimation: ~4 characters per token
	return len(text) / 4
}

// KeywordExtraction extracts keywords from composition content using word frequency analysis
func (cs *ContextSelector) KeywordExtraction(content string) map[string]int {
	// Normalize content: lowercase and remove special characters
	normalized := strings.ToLower(content)

	// Remove markdown syntax and code blocks
	normalized = removeMarkdownSyntax(normalized)

	// Tokenize into words
	words := tokenizeWords(normalized)

	// Filter out stopwords and short words
	filtered := filterStopwords(words)

	// Calculate word frequency
	frequency := calculateWordFrequency(filtered)

	return frequency
}

// removeMarkdownSyntax removes markdown formatting and code blocks
func removeMarkdownSyntax(content string) string {
	// Remove code blocks (```...```)
	codeBlockRegex := regexp.MustCompile("```[\\s\\S]*?```")
	content = codeBlockRegex.ReplaceAllString(content, " ")

	// Remove inline code (`...`)
	inlineCodeRegex := regexp.MustCompile("`[^`]+`")
	content = inlineCodeRegex.ReplaceAllString(content, " ")

	// Remove markdown headers (#...)
	headerRegex := regexp.MustCompile("(?m)^#+\\s+")
	content = headerRegex.ReplaceAllString(content, " ")

	// Remove markdown links ([text](url))
	linkRegex := regexp.MustCompile(`\[[^\]]+\]\([^)]+\)`)
	content = linkRegex.ReplaceAllString(content, " ")

	// Remove markdown emphasis (*, _, **, __)
	emphasisRegex := regexp.MustCompile(`\*+|_+`)
	content = emphasisRegex.ReplaceAllString(content, " ")

	// Remove markdown lists (-, *, +, 1., 2., etc.)
	listRegex := regexp.MustCompile(`(?m)^[\s]*[-*+]|\d+\.\s+`)
	content = listRegex.ReplaceAllString(content, " ")

	// Remove special characters but keep alphanumeric and spaces
	specialCharRegex := regexp.MustCompile(`[^a-z0-9\s]`)
	content = specialCharRegex.ReplaceAllString(content, " ")

	// Collapse multiple spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	content = spaceRegex.ReplaceAllString(content, " ")

	return strings.TrimSpace(content)
}

// tokenizeWords splits content into individual words
func tokenizeWords(content string) []string {
	return strings.Fields(content)
}

// filterStopwords removes common stopwords and short words
func filterStopwords(words []string) []string {
	stopwords := map[string]bool{
		// Common English stopwords
		"a": true, "an": true, "and": true, "are": true, "as": true,
		"at": true, "be": true, "by": true, "for": true, "from": true,
		"has": true, "have": true, "he": true, "in": true, "is": true,
		"it": true, "its": true, "of": true, "on": true, "that": true,
		"the": true, "to": true, "was": true, "were": true, "will": true,
		"with": true, "this": true, "but": true, "they": true,
		"i": true, "you": true, "we": true, "your": true, "my": true,
		"our": true, "their": true, "his": true, "her": true, "him": true,
		"me": true, "us": true, "them": true, "what": true, "which": true,
		"who": true, "whom": true, "when": true, "where": true, "why": true,
		"how": true, "all": true, "each": true, "every": true, "both": true,
		"few": true, "more": true, "most": true, "other": true, "some": true,
		"such": true, "no": true, "nor": true, "not": true, "only": true,
		"own": true, "same": true, "so": true, "than": true, "too": true,
		"very": true, "can": true, "just": true, "should": true,
		"now": true, "do": true, "does": true, "did": true, "don": true,
		"doesn": true, "didn": true, "couldn": true, "wouldn": true,
		// Programming-related stopwords
		"code": true, "function": true, "method": true, "class": true,
		"variable": true, "string": true, "number": true, "boolean": true,
		"array": true, "object": true, "null": true, "undefined": true,
		"true": true, "false": true, "return": true, "if": true, "else": true,
		"while": true, "case": true, "break": true, "continue": true,
		"import": true, "export": true, "const": true, "let": true,
		"var": true, "new": true, "super": true, "extends": true,
		"static": true, "public": true, "private": true, "protected": true,
		"async": true, "await": true, "try": true, "catch": true, "finally": true,
		"throw": true, "throws": true, "interface": true, "type": true, "enum": true,
		"implements": true, "abstract": true, "final": true, "override": true,
		"package": true, "struct": true, "func": true, "go": true,
		"defer": true, "select": true, "range": true, "chan": true, "map": true,
		"make": true, "append": true, "len": true, "cap": true, "copy": true,
		"delete": true, "close": true, "print": true, "println": true, "fmt": true,
		"log": true, "error": true, "errors": true, "panic": true, "recover": true,
	}

	var filtered []string
	for _, word := range words {
		// Skip short words (less than 3 characters)
		if len(word) < 3 {
			continue
		}

		// Skip stopwords
		if stopwords[word] {
			continue
		}

		// Skip words that are mostly numbers
		if isMostlyNumeric(word) {
			continue
		}

		filtered = append(filtered, word)
	}

	return filtered
}

// isMostlyNumeric checks if a word is mostly numeric characters
func isMostlyNumeric(word string) bool {
	numericCount := 0
	for _, r := range word {
		if unicode.IsDigit(r) {
			numericCount++
		}
	}
	return numericCount > len(word)/2
}

// calculateWordFrequency calculates the frequency of each word
func calculateWordFrequency(words []string) map[string]int {
	frequency := make(map[string]int)

	for _, word := range words {
		frequency[word]++
	}

	return frequency
}

// GetTopKeywords returns the top N keywords by frequency
func (cs *ContextSelector) GetTopKeywords(frequency map[string]int, n int) []string {
	// Convert map to slice of keyword-frequency pairs
	type keywordFreq struct {
		keyword   string
		frequency int
	}

	var pairs []keywordFreq
	for keyword, freq := range frequency {
		pairs = append(pairs, keywordFreq{keyword, freq})
	}

	// Sort by frequency (descending)
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].frequency > pairs[i].frequency {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}

	// Return top N keywords
	topN := n
	if topN > len(pairs) {
		topN = len(pairs)
	}

	result := make([]string, topN)
	for i := 0; i < topN; i++ {
		result[i] = pairs[i].keyword
	}

	return result
}
