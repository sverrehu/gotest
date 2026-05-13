package immutableMath

import (
	"math/big"
	"testing"
)

func TestAdd(t *testing.T) {
	n1 := big.NewFloat(2.0)
	n2 := big.NewFloat(3.0)
	expected := big.NewFloat(5.0)
	result := Add(n1, n2)
	if result.Cmp(expected) != 0 {
		t.Errorf("Add: expected %v, got %v", expected, result)
	}
	assertImmutable(t, n1, 2.0, n2, 3.0)
}

func TestSubtract(t *testing.T) {
	n1 := big.NewFloat(2.0)
	n2 := big.NewFloat(3.0)
	expected := big.NewFloat(-1.0)
	result := Subtract(n1, n2)
	if result.Cmp(expected) != 0 {
		t.Errorf("Subtract: expected %v, got %v", expected, result)
	}
	assertImmutable(t, n1, 2.0, n2, 3.0)
}

func TestMultiply(t *testing.T) {
	n1 := big.NewFloat(2.0)
	n2 := big.NewFloat(3.0)
	expected := big.NewFloat(6.0)
	result := Multiply(n1, n2)
	if result.Cmp(expected) != 0 {
		t.Errorf("Multiply: expected %v, got %v", expected, result)
	}
	assertImmutable(t, n1, 2.0, n2, 3.0)
}

func TestDivide(t *testing.T) {
	n1 := big.NewFloat(10.0)
	n2 := big.NewFloat(2.0)
	expected := big.NewFloat(5.0)
	result := Divide(n1, n2)
	if result.Cmp(expected) != 0 {
		t.Errorf("Divide: expected %v, got %v", expected, result)
	}
	assertImmutable(t, n1, 10.0, n2, 2.0)
}

func assertImmutable(t *testing.T, n1 *big.Float, f1 float64, n2 *big.Float, f2 float64) {
	if big.NewFloat(f1).Cmp(n1) != 0 || big.NewFloat(f2).Cmp(n2) != 0 {
		t.Errorf("Immutability check failed, operator has been modified")
	}
}
