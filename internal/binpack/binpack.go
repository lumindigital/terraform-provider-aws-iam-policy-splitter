// Package binpack implements a best-fit decreasing bin packing algorithm
// for distributing strings across a fixed number of bins with size constraints.
package binpack

import (
	"fmt"
	"math"
	"sort"
)

// maxTruncateLen is the maximum number of characters shown in error messages
// to avoid leaking sensitive policy data.
const maxTruncateLen = 50

// Bin represents a single bin containing packed strings.
type Bin struct {
	Items         []string
	RemainingSize int
}

// TotalLength returns the sum of string lengths in this bin.
func (b *Bin) TotalLength() int {
	total := 0
	for _, s := range b.Items {
		total += len(s)
	}
	return total
}

// BestFit distributes the given strings across numBins bins, where each bin
// has a maximum capacity of binSize (measured by total string length).
// It uses a best-fit decreasing strategy: strings are sorted by length
// (largest first), then each string is placed in the bin with the least
// remaining space that still fits it.
//
// Returns an error if any single string exceeds binSize or if the strings
// cannot fit within the given number of bins.
func BestFit(items []string, numBins int, binSize int) ([]Bin, error) {
	sizes := make([]int, len(items))
	for i, item := range items {
		sizes[i] = len(item)
	}
	return BestFitWithSizes(items, sizes, numBins, binSize)
}

// BestFitWithSizes distributes items across bins using explicitly provided sizes
// for capacity calculations. The actual items stored in bins are the original strings,
// but the packing algorithm uses the corresponding sizes[] values to determine fit.
// This is useful when the effective size of an item differs from its string length
// (e.g., IAM policies where whitespace doesn't count toward the size limit).
func BestFitWithSizes(items []string, sizes []int, numBins int, binSize int) ([]Bin, error) {
	if len(items) != len(sizes) {
		return nil, fmt.Errorf(
			"items and sizes must have the same length: got %d items and %d sizes",
			len(items),
			len(sizes),
		)
	}
	if numBins <= 0 {
		return nil, fmt.Errorf("numBins must be positive, got %d", numBins)
	}
	if binSize <= 0 {
		return nil, fmt.Errorf("binSize must be positive, got %d", binSize)
	}

	// Validate all items fit in a single bin.
	for i, item := range items {
		if sizes[i] > binSize {
			return nil, fmt.Errorf(
				"item of size %d exceeds bin size %d (prefix: %q)",
				sizes[i],
				binSize,
				truncate(item, maxTruncateLen),
			)
		}
	}

	// Build index pairs and sort by size descending (best-fit decreasing).
	type indexedItem struct {
		item string
		size int
	}
	sortedItems := make([]indexedItem, len(items))
	for i := range items {
		sortedItems[i] = indexedItem{item: items[i], size: sizes[i]}
	}
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].size > sortedItems[j].size
	})

	// Initialize bins.
	bins := make([]Bin, numBins)
	for i := range bins {
		bins[i].RemainingSize = binSize
	}

	// Place each item in the best-fit bin.
	for _, si := range sortedItems {
		bestIdx := findBestFitBin(bins, si.size)

		if bestIdx < 0 {
			return nil, fmt.Errorf(
				"cannot fit item of size %d into %d bins of size %d (prefix: %q)",
				si.size, numBins, binSize, truncate(si.item, maxTruncateLen),
			)
		}

		bins[bestIdx].Items = append(bins[bestIdx].Items, si.item)
		bins[bestIdx].RemainingSize -= si.size
	}

	return bins, nil
}

// findBestFitBin returns the index of the bin with the smallest remaining
// capacity that can still fit an item of the given size. Returns -1 if no
// bin can fit the item.
func findBestFitBin(bins []Bin, itemSize int) int {
	bestIdx := -1
	bestRemaining := math.MaxInt

	for i := range bins {
		if bins[i].RemainingSize >= itemSize && bins[i].RemainingSize < bestRemaining {
			bestIdx = i
			bestRemaining = bins[i].RemainingSize
		}
	}

	return bestIdx
}

// truncate returns the first maxLen characters of s, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
