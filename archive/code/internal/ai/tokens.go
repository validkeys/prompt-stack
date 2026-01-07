package ai

import (
	"fmt"
	"strings"
)

// TokenBudget manages token budget enforcement for AI context
type TokenBudget struct {
	// Configuration
	modelContextLimit int // Total context limit for current model
	compositionLimit  int // Max tokens for composition (25% of context)
	libraryLimit      int // Max tokens for library prompts (15% of context)
	warningThreshold  int // Warning threshold (15% of context)
	blockThreshold    int // Block threshold (25% of context)
}

// NewTokenBudget creates a new token budget manager
func NewTokenBudget(contextLimit int) *TokenBudget {
	return &TokenBudget{
		modelContextLimit: contextLimit,
		compositionLimit:  int(float64(contextLimit) * 0.25), // 25% for composition
		libraryLimit:      int(float64(contextLimit) * 0.15), // 15% for library
		warningThreshold:  int(float64(contextLimit) * 0.15), // 15% warning
		blockThreshold:    int(float64(contextLimit) * 0.25), // 25% block
	}
}

// EstimateTokensDetailed provides a more detailed token estimation
// Takes into account whitespace, punctuation, and word boundaries
func EstimateTokensDetailed(text string) int {
	if text == "" {
		return 0
	}

	// Count words (split by whitespace)
	words := strings.Fields(text)
	wordCount := len(words)

	// Count characters
	charCount := len(text)

	// Count lines
	lineCount := strings.Count(text, "\n") + 1

	// Weighted estimation:
	// - Words: ~1.3 tokens per word (average)
	// - Characters: ~0.25 tokens per character
	// - Lines: ~0.1 tokens per line (for structure)

	wordTokens := float64(wordCount) * 1.3
	charTokens := float64(charCount) * 0.25
	lineTokens := float64(lineCount) * 0.1

	// Average the estimates
	totalTokens := (wordTokens + charTokens + lineTokens) / 3

	return int(totalTokens)
}

// CheckComposition checks if composition is within token budget
// Returns (withinBudget, atWarning, atBlock, tokenCount)
func (tb *TokenBudget) CheckComposition(content string) (bool, bool, bool, int) {
	tokens := EstimateTokensDetailed(content)

	withinBudget := tokens <= tb.compositionLimit
	atWarning := tokens >= tb.warningThreshold
	atBlock := tokens >= tb.blockThreshold

	return withinBudget, atWarning, atBlock, tokens
}

// CheckLibrary checks if library prompts are within token budget
// Returns (withinBudget, tokenCount)
func (tb *TokenBudget) CheckLibrary(prompts []*IndexedPrompt) (bool, int) {
	totalTokens := 0

	for _, prompt := range prompts {
		totalTokens += EstimateTokensDetailed(prompt.Content)
	}

	withinBudget := totalTokens <= tb.libraryLimit

	return withinBudget, totalTokens
}

// CanAddPrompt checks if adding a prompt would exceed library budget
func (tb *TokenBudget) CanAddPrompt(existingPrompts []*IndexedPrompt, newPrompt *IndexedPrompt) bool {
	existingTokens := 0
	for _, prompt := range existingPrompts {
		existingTokens += EstimateTokensDetailed(prompt.Content)
	}

	newPromptTokens := EstimateTokensDetailed(newPrompt.Content)
	totalTokens := existingTokens + newPromptTokens

	return totalTokens <= tb.libraryLimit
}

// GetCompositionLimit returns the composition token limit
func (tb *TokenBudget) GetCompositionLimit() int {
	return tb.compositionLimit
}

// GetLibraryLimit returns the library token limit
func (tb *TokenBudget) GetLibraryLimit() int {
	return tb.libraryLimit
}

// GetWarningThreshold returns the warning threshold
func (tb *TokenBudget) GetWarningThreshold() int {
	return tb.warningThreshold
}

// GetBlockThreshold returns the block threshold
func (tb *TokenBudget) GetBlockThreshold() int {
	return tb.blockThreshold
}

// GetModelContextLimit returns the model's total context limit
func (tb *TokenBudget) GetModelContextLimit() int {
	return tb.modelContextLimit
}

// FormatTokenCount formats a token count for display
func FormatTokenCount(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d tokens", tokens)
	}

	// Format as K for large numbers
	k := float64(tokens) / 1000.0
	return fmt.Sprintf("%.1fK tokens", k)
}

// FormatTokenPercentage formats a token count as percentage of context limit
func (tb *TokenBudget) FormatTokenPercentage(tokens int) string {
	if tb.modelContextLimit == 0 {
		return "0%"
	}

	percentage := float64(tokens) / float64(tb.modelContextLimit) * 100.0
	return fmt.Sprintf("%.1f%%", percentage)
}

// GetBudgetStatus returns a human-readable status message
func (tb *TokenBudget) GetBudgetStatus(content string) string {
	_, atWarning, atBlock, tokens := tb.CheckComposition(content)

	if atBlock {
		return fmt.Sprintf("⚠️ Composition exceeds token budget (%s / %s limit)",
			FormatTokenCount(tokens),
			FormatTokenCount(tb.compositionLimit))
	}

	if atWarning {
		return fmt.Sprintf("⚠️ Composition approaching token budget (%s / %s limit)",
			FormatTokenCount(tokens),
			FormatTokenCount(tb.compositionLimit))
	}

	return fmt.Sprintf("Composition: %s (%s of context)",
		FormatTokenCount(tokens),
		tb.FormatTokenPercentage(tokens))
}

// GetLibraryBudgetStatus returns a human-readable library budget status
func (tb *TokenBudget) GetLibraryBudgetStatus(prompts []*IndexedPrompt) string {
	withinBudget, tokens := tb.CheckLibrary(prompts)

	if !withinBudget {
		return fmt.Sprintf("⚠️ Library context exceeds budget (%s / %s limit)",
			FormatTokenCount(tokens),
			FormatTokenCount(tb.libraryLimit))
	}

	return fmt.Sprintf("Library context: %s (%s of context)",
		FormatTokenCount(tokens),
		tb.FormatTokenPercentage(tokens))
}
