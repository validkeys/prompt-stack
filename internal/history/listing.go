package history

import (
	"sort"
	"time"

	"go.uber.org/zap"
)

// Listing handles history composition listing with various sorting options
type Listing struct {
	db     *Database
	logger *zap.Logger
}

// NewListing creates a new listing handler
func NewListing(db *Database, logger *zap.Logger) *Listing {
	return &Listing{
		db:     db,
		logger: logger,
	}
}

// SortBy defines sorting options for history listings
type SortBy int

const (
	SortByRecent    SortBy = iota // Sort by most recent (default)
	SortByOldest                  // Sort by oldest first
	SortByDirectory               // Sort by working directory
	SortBySize                    // Sort by character count
	SortByLines                   // Sort by line count
)

// ListOptions defines options for listing compositions
type ListOptions struct {
	SortBy     SortBy
	WorkingDir string // Filter by working directory (empty = all)
	Limit      int    // Maximum number of results (0 = unlimited)
	Offset     int    // Skip first N results (for pagination)
}

// ListCompositions returns compositions based on listing options
func (l *Listing) ListCompositions(opts ListOptions) ([]Composition, error) {
	l.logger.Debug("Listing compositions",
		zap.String("sort_by", l.sortByToString(opts.SortBy)),
		zap.String("working_dir", opts.WorkingDir),
		zap.Int("limit", opts.Limit),
		zap.Int("offset", opts.Offset),
	)

	var compositions []Composition
	var err error

	// Get compositions based on filter
	if opts.WorkingDir != "" {
		compositions, err = l.db.GetCompositionsByDirectory(opts.WorkingDir)
	} else {
		compositions, err = l.db.GetAllCompositions()
	}

	if err != nil {
		return nil, err
	}

	// Sort compositions
	l.sortCompositions(compositions, opts.SortBy)

	// Apply offset and limit
	if opts.Offset > 0 {
		if opts.Offset >= len(compositions) {
			return []Composition{}, nil
		}
		compositions = compositions[opts.Offset:]
	}

	if opts.Limit > 0 && opts.Limit < len(compositions) {
		compositions = compositions[:opts.Limit]
	}

	l.logger.Debug("Listed compositions", zap.Int("count", len(compositions)))
	return compositions, nil
}

// sortCompositions sorts compositions based on sort option
func (l *Listing) sortCompositions(compositions []Composition, sortBy SortBy) {
	switch sortBy {
	case SortByRecent:
		// Sort by most recent (descending created_at)
		sort.Slice(compositions, func(i, j int) bool {
			timeI, _ := time.Parse(time.RFC3339, compositions[i].CreatedAt)
			timeJ, _ := time.Parse(time.RFC3339, compositions[j].CreatedAt)
			return timeI.After(timeJ)
		})

	case SortByOldest:
		// Sort by oldest (ascending created_at)
		sort.Slice(compositions, func(i, j int) bool {
			timeI, _ := time.Parse(time.RFC3339, compositions[i].CreatedAt)
			timeJ, _ := time.Parse(time.RFC3339, compositions[j].CreatedAt)
			return timeI.Before(timeJ)
		})

	case SortByDirectory:
		// Sort by working directory (alphabetical)
		sort.Slice(compositions, func(i, j int) bool {
			return compositions[i].WorkingDirectory < compositions[j].WorkingDirectory
		})

	case SortBySize:
		// Sort by character count (descending)
		sort.Slice(compositions, func(i, j int) bool {
			return compositions[i].CharacterCount > compositions[j].CharacterCount
		})

	case SortByLines:
		// Sort by line count (descending)
		sort.Slice(compositions, func(i, j int) bool {
			return compositions[i].LineCount > compositions[j].LineCount
		})
	}
}

// sortByToString converts SortBy to string for logging
func (l *Listing) sortByToString(sortBy SortBy) string {
	switch sortBy {
	case SortByRecent:
		return "recent"
	case SortByOldest:
		return "oldest"
	case SortByDirectory:
		return "directory"
	case SortBySize:
		return "size"
	case SortByLines:
		return "lines"
	default:
		return "unknown"
	}
}

// GetRecentCompositions returns most recent compositions
func (l *Listing) GetRecentCompositions(limit int) ([]Composition, error) {
	return l.ListCompositions(ListOptions{
		SortBy: SortByRecent,
		Limit:  limit,
	})
}

// GetOldestCompositions returns oldest compositions
func (l *Listing) GetOldestCompositions(limit int) ([]Composition, error) {
	return l.ListCompositions(ListOptions{
		SortBy: SortByOldest,
		Limit:  limit,
	})
}

// GetCompositionsByDirectory returns compositions for a directory, sorted by recent
func (l *Listing) GetCompositionsByDirectory(workingDir string, limit int) ([]Composition, error) {
	return l.ListCompositions(ListOptions{
		SortBy:     SortByRecent,
		WorkingDir: workingDir,
		Limit:      limit,
	})
}

// GetLargestCompositions returns compositions with most characters
func (l *Listing) GetLargestCompositions(limit int) ([]Composition, error) {
	return l.ListCompositions(ListOptions{
		SortBy: SortBySize,
		Limit:  limit,
	})
}

// GetLongestCompositions returns compositions with most lines
func (l *Listing) GetLongestCompositions(limit int) ([]Composition, error) {
	return l.ListCompositions(ListOptions{
		SortBy: SortByLines,
		Limit:  limit,
	})
}

// GetUniqueDirectories returns list of unique working directories
func (l *Listing) GetUniqueDirectories() ([]string, error) {
	compositions, err := l.db.GetAllCompositions()
	if err != nil {
		return nil, err
	}

	// Use map to track unique directories
	dirMap := make(map[string]bool)
	for _, comp := range compositions {
		dirMap[comp.WorkingDirectory] = true
	}

	// Convert map to slice
	directories := make([]string, 0, len(dirMap))
	for dir := range dirMap {
		directories = append(directories, dir)
	}

	// Sort alphabetically
	sort.Strings(directories)

	return directories, nil
}

// GetCompositionCount returns total number of compositions
func (l *Listing) GetCompositionCount() (int, error) {
	compositions, err := l.db.GetAllCompositions()
	if err != nil {
		return 0, err
	}
	return len(compositions), nil
}

// GetCompositionCountByDirectory returns number of compositions in a directory
func (l *Listing) GetCompositionCountByDirectory(workingDir string) (int, error) {
	compositions, err := l.db.GetCompositionsByDirectory(workingDir)
	if err != nil {
		return 0, err
	}
	return len(compositions), nil
}

// GetTotalSize returns total size of all compositions in characters
func (l *Listing) GetTotalSize() (int, error) {
	compositions, err := l.db.GetAllCompositions()
	if err != nil {
		return 0, err
	}

	total := 0
	for _, comp := range compositions {
		total += comp.CharacterCount
	}

	return total, nil
}

// GetDateRange returns the date range of all compositions
func (l *Listing) GetDateRange() (time.Time, time.Time, error) {
	compositions, err := l.db.GetAllCompositions()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if len(compositions) == 0 {
		return time.Time{}, time.Time{}, nil
	}

	// Find earliest and latest dates
	var earliest, latest time.Time
	for i, comp := range compositions {
		compTime, err := time.Parse(time.RFC3339, comp.CreatedAt)
		if err != nil {
			continue
		}

		if i == 0 {
			earliest = compTime
			latest = compTime
		} else {
			if compTime.Before(earliest) {
				earliest = compTime
			}
			if compTime.After(latest) {
				latest = compTime
			}
		}
	}

	return earliest, latest, nil
}
