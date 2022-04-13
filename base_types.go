package main

import "math"

// 8-bit characters.
type charType = int8
type ucharType = uint8

// 32-bit integer.
type valueType = int32

const maxValue = math.MaxInt32

// 32-bit unsigned integer.
type baseType = uint32

// 32 or 64-bit unsigned integer.
type sizeType = int

type valueComparatorFunc = func(lhs valueType, rhs valueType) bool
