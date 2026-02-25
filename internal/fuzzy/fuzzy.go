// Package fuzzy cung cấp thuật toán khớp chuỗi mờ (fzf-style) tối ưu cho
// danh sách anime/manga tiếng Việt và tiếng Anh.
//
// Thuật toán scoring:
//   - +16  mỗi ký tự liên tiếp khớp (consecutive bonus)
//   - +8   ký tự đầu từ khớp (word boundary bonus, sau space/dash/dot/bracket)
//   - +4   ký tự đầu chuỗi khớp
//   - -1   mỗi ký tự bị bỏ qua giữa hai lần khớp (gap penalty)
//   - Score âm hoặc không có subsequence → không khớp
package fuzzy

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Result chứa thông tin kết quả sau khi score.
type Result struct {
	Index int    // vị trí trong slice gốc
	Score int    // điểm số (cao hơn = khớp hơn)
	Str   string // chuỗi gốc
}

const (
	scoreMatch      = 16 // ký tự liên tiếp khớp
	scoreBoundary   = 8  // đầu từ
	scoreStart      = 4  // đầu chuỗi
	scoreCamelUpper = 2  // chữ hoa (CamelCase)
	penaltyGap      = -1 // mỗi ký tự bị skip
	penaltyGapStart = -3 // bắt đầu một khoảng skip mới
)

// isWordBoundary trả về true nếu ký tự trước là khoảng cách / dấu phân cách.
func isWordBoundary(r rune) bool {
	return r == ' ' || r == '-' || r == '_' || r == '.' ||
		r == '(' || r == ')' || r == '[' || r == ']' ||
		r == '/' || r == ':' || r == ','
}

// normalize chuyển về chữ thường, giữ nguyên Unicode (hỗ trợ tiếng Việt).
func normalize(s string) string {
	return strings.ToLower(s)
}

// Score tính điểm khớp của pattern trong text theo thuật toán fzf-style.
// Trả về (score, matched):
//   - matched=false nếu pattern không phải subsequence của text.
//   - score: điểm số (lớn hơn = tốt hơn).
func Score(text, pattern string) (int, bool) {
	if pattern == "" {
		return 0, true
	}

	normText := normalize(text)
	normPat := normalize(pattern)

	textRunes := []rune(normText)
	origRunes := []rune(text) // dùng để detect word boundary và camelCase
	patRunes := []rune(normPat)

	n := len(textRunes)
	m := len(patRunes)

	if m > n {
		return 0, false
	}

	// Kiểm tra subsequence (forward scan)
	firstMatch := make([]int, 0, m)
	pi := 0
	for ti := 0; ti < n && pi < m; ti++ {
		if textRunes[ti] == patRunes[pi] {
			firstMatch = append(firstMatch, ti)
			pi++
		}
	}
	if pi < m {
		return 0, false // không phải subsequence
	}

	// Tính score với dynamic scoring
	score := 0
	inGap := false
	prevMatchIdx := -1

	for rank, ti := range firstMatch {
		_ = rank

		// Gap penalty
		if prevMatchIdx >= 0 {
			gap := ti - prevMatchIdx - 1
			if gap > 0 {
				if !inGap {
					score += penaltyGapStart
					inGap = true
				}
				score += gap * penaltyGap
			} else {
				inGap = false
			}
		}

		// Match bonus
		bonus := scoreMatch

		// Consecutive bonus
		if prevMatchIdx == ti-1 {
			bonus += scoreMatch // double cho consecutive
		}

		// Word boundary bonus
		if ti == 0 {
			bonus += scoreStart
		} else {
			prev := origRunes[ti-1]
			if isWordBoundary(prev) {
				bonus += scoreBoundary
			} else if unicode.IsUpper(origRunes[ti]) && unicode.IsLower(prev) {
				bonus += scoreCamelUpper // CamelCase
			}
		}

		score += bonus
		prevMatchIdx = ti
	}

	// Bonus cho exact prefix match
	if strings.HasPrefix(normText, normPat) {
		score += len(patRunes) * 2
	}

	return score, true
}

// FilterAndRank nhận một slice chuỗi và pattern, trả về các Result đã được
// filter (chỉ những cái khớp) và sắp xếp theo score giảm dần.
// Hàm này chạy synchronously và rất nhanh cho < 1000 phần tử.
func FilterAndRank(items []string, pattern string) []Result {
	if pattern == "" {
		results := make([]Result, len(items))
		for i, s := range items {
			results[i] = Result{Index: i, Score: 0, Str: s}
		}
		return results
	}

	results := make([]Result, 0, len(items))
	for i, s := range items {
		if score, ok := Score(s, pattern); ok {
			results = append(results, Result{Index: i, Score: score, Str: s})
		}
	}

	// Sắp xếp ổn định: score giảm dần, cùng score thì theo index gốc
	sort.SliceStable(results, func(a, b int) bool {
		if results[a].Score != results[b].Score {
			return results[a].Score > results[b].Score
		}
		return results[a].Index < results[b].Index
	})

	return results
}

// ContainsQuery kiểm tra nhanh xem chuỗi có phải subsequence của pattern không.
// Dùng để filter sơ bộ trước khi gọi Score.
func ContainsQuery(text, pattern string) bool {
	if pattern == "" {
		return true
	}
	t := normalize(text)
	p := normalize(pattern)

	tRunes := []rune(t)
	pRunes := []rune(p)

	pi := 0
	for _, r := range tRunes {
		if pi >= len(pRunes) {
			break
		}
		if r == pRunes[pi] {
			pi++
		}
	}
	return pi == len(pRunes)
}

// ByteLen trả về số byte UTF-8 của pattern (dùng để detect empty).
func ByteLen(s string) int {
	return utf8.RuneCountInString(s)
}
