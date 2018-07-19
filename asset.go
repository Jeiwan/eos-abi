package eosabi

import (
	"fmt"
	"strconv"
)

type asset struct {
	Amount uint64
	Symbol *symbol
}

func (a asset) String() string {

	aStr := strconv.Itoa(int(a.Amount))
	p := int(a.Symbol.Precision)

	return fmt.Sprintf("%s.%s %s", aStr[0:len(aStr)-p], aStr[len(aStr)-p:], a.Symbol.Token)
}
