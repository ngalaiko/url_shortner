package helpers

import (
	"bytes"
	"strconv"
)

// Uint64sToString concats ids to string with ',' delimmer
func Uint64sToString(ids []uint64) string {
	b := bytes.Buffer{}
	for i, id := range ids {
		if i > 0 {
			b.WriteRune(',')
		}

		b.WriteString(strconv.Itoa(int(id)))
	}

	return b.String()
}
