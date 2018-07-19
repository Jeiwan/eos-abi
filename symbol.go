package eosabi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/gommon/log"
)

type symbol struct {
	Precision uint8
	Token     string
}

func parseSymbol(str string) *symbol {
	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		log.Errorf("parseSymbol: wrong symbol %s", str)
		return nil
	}

	p, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Errorf("parseSymbol: wrong symbol %s", str)
		return nil
	}

	return &symbol{Precision: uint8(p), Token: parts[1]}
}

func (s symbol) String() string {
	return fmt.Sprintf("%d,%s", s.Precision, s.Token)
}
