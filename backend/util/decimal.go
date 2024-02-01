package util

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Converts Decimal128 to Float64 and returns the float
*/
func Decimal128ToFloat64(dec primitive.Decimal128) float64 {
	f, _ := strconv.ParseFloat(dec.String(), 64)
	return f
}

/*
Converts a Float64 into a primitive.Decimal128 with 2 decimal places
and returns the Decimal128
*/
func Float64ToDecimal128(f float64) primitive.Decimal128 {
	str := fmt.Sprintf("%.2f", f)
	dec, _ := primitive.ParseDecimal128(str)
	return dec
}
