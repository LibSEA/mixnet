package dht

import "errors"

var ErrValueTooLarge = errors.New("value too large")
var ErrValueTooSmall = errors.New("value too small")
var ErrKeyWrongSize = errors.New("key wrong size")
var ErrSignatureInvalid = errors.New("signature invalid")
