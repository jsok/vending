package coins

import (
	"fmt"
	"strings"
)

type Denomination int
type DenominationSlice []int

func (d DenominationSlice) String() string {
	s := make([]string, len(d))
	for i := range d {
		s[i] = fmt.Sprintf("%dc", d[i])
	}
	return fmt.Sprintf("DenominationSlice[%v]", strings.Join(s, " "))
}
