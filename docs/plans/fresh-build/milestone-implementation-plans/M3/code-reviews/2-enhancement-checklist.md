# Enhancement Checklist - M3 Code Review

**Created:** 2026-01-08  
**Source:** Code Review #2 (Uncommitted Changes)

---

## Code Quality Improvements

- [x] Remove dead code in validation block (`internal/prompt/frontmatter.go:61-64`)
- [x] Implement proper UTF-8 validation using `utf8.Valid(data)` (`internal/prompt/storage.go:94-98`)
- [x] Replace hardcoded newlines with theme spacing constants (`ui/app/model.go:97-102`)
- [x] Replace custom string functions with `strings.Contains` from stdlib (`internal/prompt/storage_test.go:716-726`)
- [x] Expand package comment to list style categories (`ui/theme/theme.go:1-4`)

---

## Future Enhancements

- [x] Add integration tests for end-to-end file I/O workflows
- [x] Add benchmark tests for theme rendering performance
- [x] Add file extension validation (.md files only)
- [x] Add validation for frontmatter fields (required fields, max lengths, etc.)
- [x] Add theme examples and demo mode
