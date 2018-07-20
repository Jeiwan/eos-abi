package eosabi

import (
	"fmt"
	"math"
	"math/big"
)

type asset struct {
	Amount uint64
	Symbol *symbol
}

func (a asset) String() string {
	precision := int(a.Symbol.Precision)
	p := big.NewFloat(math.Pow(10.0, float64(precision)))

	bigA := big.NewFloat(float64(a.Amount))
	bigA.Quo(bigA, p)

	amm := bigA.Text('f', precision)

	return fmt.Sprintf("%s %s", amm, a.Symbol.Token)
}
