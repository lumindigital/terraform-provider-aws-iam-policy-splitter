// Package binpack implements a best-fit decreasing bin packing algorithm
// for distributing strings across a fixed number of bins with size constraints.
package binpack

import (
	"fmt"
	"math"
	"sort"
)

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
	if numBins <= 0 {
		return nil, fmt.Errorf("numBins must be positive, got %d", numBins)
	}
	if binSize <= 0 {
		return nil, fmt.Errorf("binSize must be positive, got %d", binSize)
	}

	// Validate all items fit in a single bin.
	for _, item := range items {
		if len(item) > binSize {
			return nil, fmt.Errorf("item %q (length %d) exceeds bin size %d", item, len(item), binSize)
		}
	}

	// Sort items by length descending (best-fit decreasing).
	sorted := make([]string, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) > len(sorted[j])
	})

	// Initialize bins.
	bins := make([]Bin, numBins)
	for i := range bins {
		bins[i].RemainingSize = binSize
	}

	// Place each item in the best-fit bin.
	for _, item := range sorted {
		bestIdx := findBestFitBin(bins, len(item))

		if bestIdx < 0 {
			return nil, fmt.Errorf(
				"cannot fit item %q (length %d) into %d bins of size %d",
				item, len(item), numBins, binSize,
			)
		}

		bins[bestIdx].Items = append(bins[bestIdx].Items, item)
		bins[bestIdx].RemainingSize -= len(item)
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
