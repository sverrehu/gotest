package immutableMath

import "math/big"

func Add(f0 *big.Float, f1 *big.Float) *big.Float {
	result := new(big.Float)
	result.Copy(f0)
	return result.Add(result, f1)
}

func Subtract(f0 *big.Float, f1 *big.Float) *big.Float {
	result := new(big.Float)
	result.Copy(f0)
	return result.Sub(result, f1)
}

func Divide(f0 *big.Float, f1 *big.Float) *big.Float {
	result := new(big.Float)
	result.Copy(f0)
	return result.Quo(result, f1)
}

func Multiply(f0 *big.Float, f1 *big.Float) *big.Float {
	result := new(big.Float)
	result.Copy(f0)
	return result.Mul(result, f1)
}
