package code

import (
	"math/rand/v2"
	"strings"
	"time"
)

const symbols = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func MakeCode(length int) string {
	var builder strings.Builder
	r := rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), uint64(time.Now().UnixMilli())))

	for i := 0; i < length; i++ {
		builder.WriteByte(symbols[r.IntN(len(symbols))])
	}

	return builder.String()
}
