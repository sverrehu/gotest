package main

import (
	"fmt"
	"math"
	"math/big"
)

func dump(f float64) {
	b := big.NewFloat(f)
	mant := new(big.Float)
	exp := b.MantExp(mant)
	fmt.Printf("%s: precission: %d, exponent: %d, mantissa: %d\n", b.String(), b.Prec(), exp, mant)
}

func dump2(s string) {
	b := new(big.Float)
	precisionBits := uint(0.5 + float64(len(s))*math.Log2(10))
	b.SetPrec(precisionBits)
	_, _, err := b.Parse(s, 10)
	if err != nil {
		panic(err)
	}
	mant := new(big.Float)
	exp := b.MantExp(mant)
	fmt.Printf("%s: precission: %d, exponent: %d, mantissa: %d\n", b.Text('g', -2), b.Prec(), exp, mant)
}

func main() {
	dump(3.14)
	dump(3.14e5)
	dump(3.14e-5)
	dump2("3.1415926535897932384626433832795028841971693993751058209749445923078164062862089986280348253421170679821480865132823066e-1000")
}
