//go:generate go-enum -f=$GOFILE

package stockindex

// Index is an enumeration of indices that are allowed.
// ENUM(
// TWSE
// DJI
// NASDAQ
// SP500
// )
type Index int32
