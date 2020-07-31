package meander

import "strings"

// ParseCost parses the string cost value to the number value.
func ParseCost(s string) Cost {
	return costStrings[s]
}

// CostRange represents a range between Cost values.
type CostRange struct {
	From Cost
	To   Cost
}

// String parses the CostRange field values to a string representation.
func (r CostRange) String() string {
	return r.From.String() + "..." + r.To.String()
}

// ParseCostRange parses a string cost range to a populated CostRange.
func ParseCostRange(s string) *CostRange {
	segs := strings.Split(s, "...")

	return &CostRange{
		From: ParseCost(segs[0]),
		To:   ParseCost(segs[1]),
	}
}
