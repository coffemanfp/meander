package meander

type Cost int8

const (
	_ Cost = iota
	Cost1
	Cost2
	Cost3
	Cost4
	Cost5
)

var costStrings = map[string]Cost{
	"$":     Cost1,
	"$$":    Cost2,
	"$$$":   Cost3,
	"$$$$":  Cost4,
	"$$$$$": Cost5,
}

// String returns the string representation of a Cost.
func (l Cost) String() string {
	for s, v := range costStrings {
		if l == v {
			return s
		}
	}

	return "invalid"
}
