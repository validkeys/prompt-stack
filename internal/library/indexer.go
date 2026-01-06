package library

import (
	"time"

	"github.com/kyledavis/prompt-stack/internal/prompt"
	"go.uber.org/zap"
)

// BuildIndex creates an in-memory index of all prompts
func (l *Library) BuildIndex() error {
	l.logger.Info("Building library index")

	// Clear existing index
	l.Index.Prompts = make(map[string]*prompt.IndexedPrompt)

	// Build index for each prompt
	for filePath, p := range l.Prompts {
		indexed := &prompt.IndexedPrompt{
			PromptID:      p.ID,
			Title:         p.Title,
			Tags:          p.Tags,
			Category:      p.Category,
			WordFrequency: prompt.ExtractKeywords(p.Content),
			LastUsed:      p.UsageStats.LastUsed,
			UseCount:      p.UsageStats.UseCount,
		}

		l.Index.Prompts[filePath] = indexed
	}

	l.Index.LastBuilt = time.Now()
	l.logger.Info("Index built successfully", zap.Int("prompt_count", len(l.Index.Prompts)))
	return nil
}

// ScorePrompts scores prompts based on composition keywords and tags
func (l *Library) ScorePrompts(compositionKeywords map[string]int, compositionTags []string) []ScoredPrompt {
	var scored []ScoredPrompt

	for filePath, indexed := range l.Index.Prompts {
		score := l.calculateScore(indexed, compositionKeywords, compositionTags)

		scored = append(scored, ScoredPrompt{
			FilePath: filePath,
			Prompt:   indexed,
			Score:    score,
		})
	}

	// Sort by score (highest first)
	sortScoredPrompts(scored)

	return scored
}

// ScoredPrompt represents a prompt with its relevance score
type ScoredPrompt struct {
	FilePath string
	Prompt   *prompt.IndexedPrompt
	Score    int
}

// calculateScore calculates relevance score for a prompt
func (l *Library) calculateScore(indexed *prompt.IndexedPrompt, compositionKeywords map[string]int, compositionTags []string) int {
	score := 0

	// Tag matches: +10 per matching tag
	for _, tag := range compositionTags {
		for _, promptTag := range indexed.Tags {
			if tag == promptTag {
				score += 10
				break
			}
		}
	}

	// Keyword matches: +1 per matching word (weighted by frequency)
	for keyword, freq := range compositionKeywords {
		if _, exists := indexed.WordFrequency[keyword]; exists {
			// Weight by frequency in composition
			score += freq
		}
	}

	// Recently used: +3 if used in last session
	if !indexed.LastUsed.IsZero() {
		if time.Since(indexed.LastUsed) < 24*time.Hour {
			score += 3
		}
	}

	// Frequently used: +use_count
	score += indexed.UseCount

	return score
}

// sortScoredPrompts sorts scored prompts by score (highest first)
func sortScoredPrompts(scored []ScoredPrompt) {
	// Simple bubble sort for now
	// Could use more efficient algorithm for large libraries
	n := len(scored)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if scored[j].Score < scored[j+1].Score {
				scored[j], scored[j+1] = scored[j+1], scored[j]
			}
		}
	}
}

// GetTopPrompts returns top N prompts by score
func (l *Library) GetTopPrompts(compositionKeywords map[string]int, compositionTags []string, n int) []*prompt.IndexedPrompt {
	scored := l.ScorePrompts(compositionKeywords, compositionTags)

	// Return top N
	if len(scored) < n {
		n = len(scored)
	}

	var result []*prompt.IndexedPrompt
	for i := 0; i < n; i++ {
		result = append(result, scored[i].Prompt)
	}

	return result
}
