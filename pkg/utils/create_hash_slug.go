package utils

import (
	"crypto/rand"
	"math/big"

	wraperrors "github.com/jhonVitor-rs/url-shortener/pkg/wrap_errors"
)

var lettersRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func CreateHashSlug() (string, error) {
	b := make([]rune, 10)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(lettersRune))))
		if err != nil {
			return "", wraperrors.InternalErr("Failed to create a slug", err)
		}

		b[i] = lettersRune[num.Int64()]
	}

	return string(b), nil
}
