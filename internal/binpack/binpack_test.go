package binpack_test

import (
	"strings"
	"testing"

	"github.com/Lumin-Digital/terraform-aws-iam-policy-splitter/internal/binpack"
)

func TestBestFit_BasicDistribution(t *testing.T) {
	t.Parallel()

	items := []string{"aaaa", "bb", "ccc", "d", "ee"}
	bins, err := binpack.BestFit(items, 2, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bins) != 2 {
		t.Fatalf("expected 2 bins, got %d", len(bins))
	}
	// Verify all items are distributed.
	totalItems := 0
	for _, b := range bins {
		totalItems += len(b.Items)
		if b.TotalLength() > 10 {
			t.Errorf("bin exceeds size limit: total length %d > 10", b.TotalLength())
		}
	}
	if totalItems != 5 {
		t.Errorf("expected 5 total items, got %d", totalItems)
	}
}

func TestBestFit_SingleBin(t *testing.T) {
	t.Parallel()

	items := []string{"ab", "cd", "ef"}
	bins, err := binpack.BestFit(items, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bins[0].Items) != 3 {
		t.Errorf("expected 3 items in single bin, got %d", len(bins[0].Items))
	}
	if bins[0].TotalLength() != 6 {
		t.Errorf("expected total length 6, got %d", bins[0].TotalLength())
	}
}

func TestBestFit_ItemExceedsBinSize(t *testing.T) {
	t.Parallel()

	items := []string{"toolongstring"}
	_, err := binpack.BestFit(items, 2, 5)
	if err == nil {
		t.Fatal("expected error for item exceeding bin size")
	}
}

func TestBestFit_CannotFitAll(t *testing.T) {
	t.Parallel()

	items := []string{"aaaa", "bbbb", "cccc"}
	_, err := binpack.BestFit(items, 2, 5)
	if err == nil {
		t.Fatal("expected error when items cannot fit")
	}
}

func TestBestFit_EmptyInput(t *testing.T) {
	t.Parallel()

	bins, err := binpack.BestFit([]string{}, 3, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, b := range bins {
		if len(b.Items) != 0 {
			t.Errorf("expected empty bin, got %d items", len(b.Items))
		}
	}
}

func TestBestFit_InvalidNumBins(t *testing.T) {
	t.Parallel()

	_, err := binpack.BestFit([]string{"a"}, 0, 10)
	if err == nil {
		t.Fatal("expected error for zero bins")
	}
}

func TestBestFit_InvalidBinSize(t *testing.T) {
	t.Parallel()

	_, err := binpack.BestFit([]string{"a"}, 1, 0)
	if err == nil {
		t.Fatal("expected error for zero bin size")
	}
}

func TestBestFit_BestFitBehavior(t *testing.T) {
	t.Parallel()

	// With best-fit, items should go into the tightest-fitting bin.
	// Bins of size 10, items: 7, 5, 3, 3, 2.
	// Sorted descending: 7, 5, 3, 3, 2.
	// Bin 1: 7 (rem 3), Bin 2: 5 (rem 5).
	// 3 -> Bin 1 (rem 3, tighter fit), Bin 1: 7,3 (rem 0).
	// 3 -> Bin 2 (rem 5), Bin 2: 5,3 (rem 2).
	// 2 -> Bin 2 (rem 2, exact fit), Bin 2: 5,3,2 (rem 0).
	items := []string{"1234567", "12345", "123", "abc", "xy"}
	bins, err := binpack.BestFit(items, 2, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Both bins should be fully packed.
	for i, b := range bins {
		if b.RemainingSize != 0 {
			t.Errorf("bin %d has %d remaining, expected 0", i, b.RemainingSize)
		}
	}
}

func TestBestFit_ErrorMessageTruncation(t *testing.T) {
	t.Parallel()

	// Create a long item that exceeds bin size — error should not contain the full string.
	longItem := make([]byte, 200)
	for i := range longItem {
		longItem[i] = 'x'
	}
	_, err := binpack.BestFit([]string{string(longItem)}, 2, 5)
	if err == nil {
		t.Fatal("expected error for oversized item")
	}
	// The error should contain a truncated prefix, not the full 200-char string.
	errMsg := err.Error()
	if len(errMsg) > 150 {
		t.Errorf("error message is too long (%d chars), may be leaking full item content", len(errMsg))
	}
	if !strings.Contains(errMsg, "prefix:") {
		t.Error("error message should contain 'prefix:' indicating truncation")
	}
}

func TestBestFitWithSizes_UsesProvidedSizes(t *testing.T) {
	t.Parallel()

	// Items have whitespace, but sizes reflect stripped length.
	// "ab cd" is 5 chars raw, but effective size is 4 (no space).
	// "efgh" is 4 chars raw, effective size is 4.
	// With binSize=5, both fit in one bin if we use stripped sizes (4+4=8 > 5, needs 2 bins).
	// With binSize=8, both fit in one bin using stripped sizes.
	items := []string{"ab cd", "efgh"}
	sizes := []int{4, 4} // stripped sizes
	bins, err := binpack.BestFitWithSizes(items, sizes, 1, 8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Both items should be in the single bin.
	if len(bins[0].Items) != 2 {
		t.Errorf("expected 2 items in bin, got %d", len(bins[0].Items))
	}
	// Original items (with whitespace) should be preserved.
	found := false
	for _, item := range bins[0].Items {
		if item == "ab cd" {
			found = true
		}
	}
	if !found {
		t.Error("original item 'ab cd' with whitespace should be preserved in bin")
	}
}

func TestBestFitWithSizes_MismatchedLengths(t *testing.T) {
	t.Parallel()

	_, err := binpack.BestFitWithSizes([]string{"a", "b"}, []int{1}, 1, 10)
	if err == nil {
		t.Fatal("expected error for mismatched items/sizes lengths")
	}
}

func TestBestFitWithSizes_WhitespacePolicyFits(t *testing.T) {
	t.Parallel()

	// A policy with lots of whitespace that would exceed bin size by raw length,
	// but fits when whitespace is excluded from sizing.
	policy := `{  "Version":  "2012-10-17",  "Statement":  []  }` // 50 chars raw
	strippedSize := 38                                            // without whitespace

	bins, err := binpack.BestFitWithSizes([]string{policy}, []int{strippedSize}, 1, 40)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bins[0].Items[0] != policy {
		t.Errorf("expected original policy preserved, got %q", bins[0].Items[0])
	}
}

func TestBin_TotalLength(t *testing.T) {
	t.Parallel()

	b := binpack.Bin{Items: []string{"abc", "de", "f"}}
	if b.TotalLength() != 6 {
		t.Errorf("expected 6, got %d", b.TotalLength())
	}
}
