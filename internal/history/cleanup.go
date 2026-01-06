package history

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Cleanup handles history cleanup operations
type Cleanup struct {
	db      *Database
	storage *Storage
	logger  *zap.Logger
}

// NewCleanup creates a new cleanup handler
func NewCleanup(db *Database, storage *Storage, logger *zap.Logger) *Cleanup {
	return &Cleanup{
		db:      db,
		storage: storage,
		logger:  logger,
	}
}

// CleanupStrategy defines different cleanup strategies
type CleanupStrategy int

const (
	StrategyAgeDays     CleanupStrategy = iota // Delete compositions older than N days
	StrategyCountKeep                          // Keep only N most recent compositions
	StrategyCountDelete                        // Delete N oldest compositions
	StrategyDirectory                          // Delete all compositions in specific directory
	StrategyAll                                // Delete all compositions
)

// CleanupOptions defines options for cleanup operation
type CleanupOptions struct {
	Strategy    CleanupStrategy
	AgeDays     int    // For StrategyAgeDays
	KeepCount   int    // For StrategyCountKeep
	DeleteCount int    // For StrategyCountDelete
	Directory   string // For StrategyDirectory
	DryRun      bool   // If true, don't actually delete, just report
}

// CleanupResult represents the result of a cleanup operation
type CleanupResult struct {
	Strategy      CleanupStrategy
	FilesToDelete []string
	TotalSize     int64
	FileCount     int
	Success       bool
	Error         error
	DryRun        bool
}

// PreviewCleanup previews what would be deleted without actually deleting
func (c *Cleanup) PreviewCleanup(opts CleanupOptions) (*CleanupResult, error) {
	opts.DryRun = true
	return c.ExecuteCleanup(opts)
}

// ExecuteCleanup performs cleanup based on strategy
func (c *Cleanup) ExecuteCleanup(opts CleanupOptions) (*CleanupResult, error) {
	c.logger.Info("Executing cleanup",
		zap.String("strategy", c.strategyToString(opts.Strategy)),
		zap.Bool("dry_run", opts.DryRun),
	)

	result := &CleanupResult{
		Strategy: opts.Strategy,
		DryRun:   opts.DryRun,
	}

	// Get all compositions
	compositions, err := c.db.GetAllCompositions()
	if err != nil {
		result.Error = fmt.Errorf("failed to get compositions: %w", err)
		return result, result.Error
	}

	// Filter compositions based on strategy
	var toDelete []Composition
	switch opts.Strategy {
	case StrategyAgeDays:
		toDelete = c.filterByAge(compositions, opts.AgeDays)
	case StrategyCountKeep:
		toDelete = c.filterByKeepCount(compositions, opts.KeepCount)
	case StrategyCountDelete:
		toDelete = c.filterByDeleteCount(compositions, opts.DeleteCount)
	case StrategyDirectory:
		toDelete = c.filterByDirectory(compositions, opts.Directory)
	case StrategyAll:
		toDelete = compositions
	default:
		result.Error = fmt.Errorf("unknown cleanup strategy: %d", opts.Strategy)
		return result, result.Error
	}

	// Calculate total size
	var totalSize int64
	for _, comp := range toDelete {
		totalSize += int64(comp.CharacterCount)
	}

	result.FilesToDelete = make([]string, len(toDelete))
	result.TotalSize = totalSize
	result.FileCount = len(toDelete)

	for i, comp := range toDelete {
		result.FilesToDelete[i] = comp.FilePath
	}

	// If dry run, just return the preview
	if opts.DryRun {
		c.logger.Info("Cleanup preview",
			zap.Int("files_to_delete", len(toDelete)),
			zap.Int64("total_size", totalSize),
		)
		return result, nil
	}

	// Actually delete files and database entries
	for _, comp := range toDelete {
		if err := c.storage.DeleteHistoryFile(comp.FilePath); err != nil {
			c.logger.Warn("Failed to delete history file",
				zap.String("file_path", comp.FilePath),
				zap.Error(err),
			)
			// Continue with other files
		}

		if err := c.db.DeleteComposition(comp.FilePath); err != nil {
			c.logger.Warn("Failed to delete database entry",
				zap.String("file_path", comp.FilePath),
				zap.Error(err),
			)
			// Continue with other files
		}
	}

	result.Success = true
	c.logger.Info("Cleanup completed",
		zap.Int("files_deleted", len(toDelete)),
		zap.Int64("total_size_freed", totalSize),
	)

	return result, nil
}

// filterByAge filters compositions older than specified days
func (c *Cleanup) filterByAge(compositions []Composition, ageDays int) []Composition {
	cutoff := time.Now().AddDate(0, 0, -ageDays)
	var filtered []Composition

	for _, comp := range compositions {
		createdAt, err := time.Parse(time.RFC3339, comp.CreatedAt)
		if err != nil {
			c.logger.Warn("Failed to parse creation time",
				zap.String("file_path", comp.FilePath),
				zap.Error(err),
			)
			continue
		}

		if createdAt.Before(cutoff) {
			filtered = append(filtered, comp)
		}
	}

	return filtered
}

// filterByKeepCount keeps only N most recent compositions
func (c *Cleanup) filterByKeepCount(compositions []Composition, keepCount int) []Composition {
	if keepCount >= len(compositions) {
		return []Composition{} // Nothing to delete
	}

	// Sort by creation time (most recent first)
	sorted := make([]Composition, len(compositions))
	copy(sorted, compositions)

	// Sort by created_at descending
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			timeI, _ := time.Parse(time.RFC3339, sorted[i].CreatedAt)
			timeJ, _ := time.Parse(time.RFC3339, sorted[j].CreatedAt)
			if timeI.Before(timeJ) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Delete all but the N most recent
	return sorted[keepCount:]
}

// filterByDeleteCount deletes N oldest compositions
func (c *Cleanup) filterByDeleteCount(compositions []Composition, deleteCount int) []Composition {
	if deleteCount <= 0 {
		return []Composition{} // Nothing to delete
	}

	if deleteCount > len(compositions) {
		deleteCount = len(compositions) // Can't delete more than we have
	}

	// Sort by creation time (oldest first)
	sorted := make([]Composition, len(compositions))
	copy(sorted, compositions)

	// Sort by created_at ascending
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			timeI, _ := time.Parse(time.RFC3339, sorted[i].CreatedAt)
			timeJ, _ := time.Parse(time.RFC3339, sorted[j].CreatedAt)
			if timeI.After(timeJ) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Delete the N oldest
	return sorted[:deleteCount]
}

// filterByDirectory filters compositions in a specific directory
func (c *Cleanup) filterByDirectory(compositions []Composition, directory string) []Composition {
	var filtered []Composition

	for _, comp := range compositions {
		if comp.WorkingDirectory == directory {
			filtered = append(filtered, comp)
		}
	}

	return filtered
}

// strategyToString converts strategy to string for logging
func (c *Cleanup) strategyToString(strategy CleanupStrategy) string {
	switch strategy {
	case StrategyAgeDays:
		return "age_days"
	case StrategyCountKeep:
		return "count_keep"
	case StrategyCountDelete:
		return "count_delete"
	case StrategyDirectory:
		return "directory"
	case StrategyAll:
		return "all"
	default:
		return "unknown"
	}
}

// GetCleanupStatistics returns statistics about history for cleanup decisions
func (c *Cleanup) GetCleanupStatistics() (*CleanupStatistics, error) {
	compositions, err := c.db.GetAllCompositions()
	if err != nil {
		return nil, err
	}

	stats := &CleanupStatistics{
		TotalCompositions: len(compositions),
	}

	// Calculate total size
	for _, comp := range compositions {
		stats.TotalSize += int64(comp.CharacterCount)
	}

	// Find date range
	if len(compositions) > 0 {
		oldest, _ := time.Parse(time.RFC3339, compositions[0].CreatedAt)
		newest, _ := time.Parse(time.RFC3339, compositions[0].CreatedAt)

		for _, comp := range compositions {
			createdAt, _ := time.Parse(time.RFC3339, comp.CreatedAt)
			if createdAt.Before(oldest) {
				oldest = createdAt
			}
			if createdAt.After(newest) {
				newest = createdAt
			}
		}

		stats.OldestComposition = oldest
		stats.NewestComposition = newest
		stats.AgeDays = int(time.Since(oldest).Hours() / 24)
	}

	// Count by directory
	stats.DirectoryCounts = make(map[string]int)
	for _, comp := range compositions {
		stats.DirectoryCounts[comp.WorkingDirectory]++
	}

	return stats, nil
}

// CleanupStatistics represents statistics about history
type CleanupStatistics struct {
	TotalCompositions int
	TotalSize         int64
	OldestComposition time.Time
	NewestComposition time.Time
	AgeDays           int
	DirectoryCounts   map[string]int
}

// FormatSize formats size in human-readable format
func FormatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
