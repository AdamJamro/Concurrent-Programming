/////////////////////////////////////////////////////////////////////
// provide a lists all viable operations and their implementations //
/////////////////////////////////////////////////////////////////////

package configuration

import (
	common "factorySimulator/commonModels"
)

// Operations is a list of all available operations in boss-worker-client circulation
var Operations = [4]common.Operation{
	{SUM, '+'},
	{DIFFERENCE, '-'},
	{PRODUCT, '*'},
	{QUOTIENT, '/'},
}

// operations' definitions
var (
	SUM        = func(args ...any) any { return fold(AnyToIntWrapper(addImpl), 0, args) }
	DIFFERENCE = func(args ...any) any {
		return fold(AnyToIntWrapper(subImpl), args[0], args[1:])
	}
	PRODUCT = func(args ...any) any {
		return fold(AnyToIntWrapper(mulImpl), 1, args)
	}
	QUOTIENT = func(args ...any) any {
		if len(args) == 0 {
			return 0
		}
		return fold(AnyToIntWrapper(divImpl), args[0], args[1:])
	}
)

// addImpl returns the integer sum.
func addImpl(a int, b int) int {
	return a + b
}

// subImpl returns the integer difference.
func subImpl(a int, b int) int {
	return a - b
}

// mulImpl returns the integer product.
func mulImpl(a int, b int) int {
	return a * b
}

// divImpl = If divider is not 0 it returns result of division, otherwise it returns 0.
func divImpl(a int, b int) int {
	if b != 0 {
		return a / b
	}
	return 0
}

// fold folds from left to right a slice with an accumulator starting at initial value using binary operation
func fold(f func(any, any) any, initial any, slice []any) any {
	accumulator := initial
	for _, elem := range slice {
		accumulator = f(accumulator, elem)
	}
	return accumulator
}

// AnyToIntWrapper is a helper wrapper function which tries to cast and change the type signature of the fn function
func AnyToIntWrapper(fn func(int, int) int) func(any, any) any {
	wrapped := func(a any, b any) any {
		valA, okA := a.(int)
		valB, okB := b.(int)
		if !okA {
			println("Error: ", valA, ":", a, " is not an integer")
			return nil
		}
		if !okB {
			println("Error: ", valB, ":", b, " is not an integer")
			return nil
		}
		return fn(valA, valB)
	}
	return wrapped
}
