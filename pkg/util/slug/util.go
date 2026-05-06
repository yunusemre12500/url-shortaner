package slug

import (
	"crypto/rand"
	"math/big"
)

const CHARSET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Generate(length int) (string, error) {
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(CHARSET))))

		if err != nil {
			return "", err
		}

		result[i] = CHARSET[num.Int64()]
	}

	return string(result), nil
}
