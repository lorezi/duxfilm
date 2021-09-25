// generates a custom JSON type for duration
package data

import (
	"fmt"
	"strconv"
)

type Duration int32

func (d Duration) MarshalJSON() ([]byte, error) {
	// adds mins to the value
	jsonValue := fmt.Sprintf("%d mins", d)

	// format to a valid json string (quoted string)
	quotedJSONValue := strconv.Quote(jsonValue)

	// convert the quoted string value to a byte slice and return it
	return []byte(quotedJSONValue), nil
}
