package main

import (
	"fmt"
	"math/big"
)

func dump(f float64) {
	b := big.NewFloat(f)
	mant := new(big.Float)
	exp := b.MantExp(mant)
	fmt.Printf("%g: precission: %d, exponent: %d, mantissa: %d\n", f, b.Prec(), exp, mant)
}

func main() {
	dump(3.14)
	dump(3.14e5)
	dump(3.14e-5)
}
