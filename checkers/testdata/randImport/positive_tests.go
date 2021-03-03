package checker_test

import (
	/*! ouch */
	cryptoRand "crypto/rand"
	"math/rand"
)

func init() {
	_ = rand.Rand{}
	_ = cryptoRand.Rand{}
}
