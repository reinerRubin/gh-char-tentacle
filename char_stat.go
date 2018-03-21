package tentacle

import (
	"fmt"
	"sort"
)

type (
	// CharStat TBD
	CharStat map[rune]int64
	// Chars TBD
	Chars []rune

	// StatItem TBD
	StatItem struct {
		Char     rune
		Quantity int64
	}

	// SortedStat TBD
	SortedStat []*StatItem
)

func (c Chars) Len() int           { return len(c) }
func (c Chars) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c Chars) Less(i, j int) bool { return c[i] <= c[j] }

func (ss SortedStat) Len() int           { return len(ss) }
func (ss SortedStat) Swap(i, j int)      { ss[i], ss[j] = ss[j], ss[i] }
func (ss SortedStat) Less(i, j int) bool { return ss[i].Quantity >= ss[j].Quantity }

// NewCharStat TBD
func NewCharStat(text string) CharStat {
	stat := make(CharStat)
	for _, char := range text {
		stat[char] = stat[char] + 1
	}

	return stat
}

// String TBD
func (hs CharStat) String() string {
	chars := hs.Chars()
	sort.Sort(chars)

	s := ""
	for i, char := range chars {
		s += fmt.Sprintf("%s: %d", string(char), hs[char])
		if i != len(chars)-1 {
			s += "\n"
		}
	}

	return s
}

// Chars TBD
func (hs CharStat) Chars() Chars {
	chars := make(Chars, 0, len(hs))
	for char := range hs {
		chars = append(chars, char)
	}

	return chars
}

// Merge TBD
func (hs CharStat) Merge(b CharStat) {
	for char, quantity := range b {
		hs[char] = hs[char] + quantity
	}
}

// SortedStat TBD
func (hs CharStat) SortedStat() SortedStat {
	sortedStat := make(SortedStat, 0, len(hs))

	// TODO optimize me
	for char, quantity := range hs {
		sortedStat = append(sortedStat, &StatItem{
			Char:     char,
			Quantity: quantity,
		})
	}
	sort.Sort(sortedStat)

	return sortedStat
}

// TextGraph TBD
func (ss SortedStat) TextGraph(desiredMaxWidth int) string {
	const charLen = 5
	var (
		widthSafyOffset = charLen + 19 // charLen + maxint64.to_s.size
		acc             = ""
		maxSaftyWidth   = desiredMaxWidth - widthSafyOffset
	)
	if maxSaftyWidth <= 0 {
		// at least, we have tried
		maxSaftyWidth = desiredMaxWidth
	}
	if len(ss) == 0 {
		return acc
	}

	biggest := ss[0].Quantity
	count := func(quantity int64) int {
		w := int(quantity * int64(maxSaftyWidth) / biggest)
		if w == 0 {
			w = 1 // for aesthetic
		}
		return w
	}

	for _, statItem := range ss {
		acc += fmt.Sprintf(rightPad(fmt.Sprintf(`%q`, string(statItem.Char)), " ", charLen))
		for i := 0; i < count(statItem.Quantity); i++ {
			acc += "■"
		}

		acc += fmt.Sprintf(" %d\n", statItem.Quantity)
	}

	return acc
}

func times(str string, n int) (out string) {
	for i := 0; i < n; i++ {
		out += str
	}
	return
}

func rightPad(str string, pad string, length int) string {
	return str + times(pad, length-len(str))
}
