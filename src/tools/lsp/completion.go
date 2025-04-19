package lsp

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"go.lsp.dev/protocol"
)

// FuzzySearchResult represents a search result with the original string and its match score
type FuzzySearchResult struct {
	Keyword Keyword
	Score   int
}

// LevenshteinDistance calculates the minimum number of single-character edits
// required to change one string into another
func LevenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	// Create a matrix of size (len(s1)+1) x (len(s2)+1)
	rows := len(s1) + 1
	cols := len(s2) + 1

	// Initialize the matrix
	matrix := make([][]int, rows)
	for i := 0; i < rows; i++ {
		matrix[i] = make([]int, cols)
		matrix[i][0] = i
	}

	for j := 0; j < cols; j++ {
		matrix[0][j] = j
	}

	// Fill the matrix
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[rows-1][cols-1]
}

func min(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}

	m := nums[0]
	for _, num := range nums {
		if num < m {
			m = num
		}
	}
	return m
}

// FuzzySearch performs a fuzzy search on a list of strings
// and returns results sorted by relevance (lower score = better match)
// limit determines the maximum number of results to return (0 = no limit)
func fuzzySearchKeywords(keywords []Keyword, searchTerm string, limit int) []FuzzySearchResult {
	var results []FuzzySearchResult

	for _, item := range keywords {
		// Calculate distance between search term and item
		distance := LevenshteinDistance(searchTerm, item.Text)

		// Add to results
		results = append(results, FuzzySearchResult{
			Keyword: item,
			Score:   distance,
		})
	}

	// Sort results by score (lower is better)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	// Apply limit if specified
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}

	return results
}

func (se *Server) Completion(ctx context.Context, params *protocol.CompletionParams) (result *protocol.CompletionList, err error) {
	logger.Sugar().Debug("Completion", params.TextDocument.URI, params.Position)

	// Get document content
	document, err := se.getDocument(params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	// Extract the current line up to the cursor position
	line := document.LineAt(params.Position.Line)
	prefix := line[:params.Position.Character]

	// Extract the current word being typed
	// This regex finds the last word boundary before cursor
	re := regexp.MustCompile(`\b[\w\d_]*$`)
	match := re.FindString(prefix)

	// If no match found, return empty completion list
	if match == "" {
		return &protocol.CompletionList{
			IsIncomplete: false,
			Items:        []protocol.CompletionItem{},
		}, nil
	}

	// Get list of keywords to search against
	keywords := se.getKeywords()

	// Perform fuzzy search on keywords
	matches := fuzzySearchKeywords(keywords, match, 10)

	// Convert matches to CompletionItems
	items := make([]protocol.CompletionItem, 0, len(matches))
	for _, m := range matches {
		items = append(items, protocol.CompletionItem{
			Label:  m.Keyword.Text,
			Kind:   m.Keyword.Kind,
			Detail: m.Keyword.Detail,
			Documentation: protocol.MarkupContent{
				Kind:  protocol.PlainText,
				Value: m.Keyword.Description,
			},
			FilterText: m.Keyword.Text,
			SortText:   fmt.Sprintf("%05d", m.Score), // Format score for proper sorting
			InsertText: m.Keyword.Text,
		})
	}

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}
func (se *Server) CompletionResolve(ctx context.Context, params *protocol.CompletionItem) (result *protocol.CompletionItem, err error) {
	return nil, nil
}
