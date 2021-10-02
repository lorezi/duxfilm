// generates a custom JSON type for duration
package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidDurationFormat = errors.New("invalid duration format")

type Duration int32

func (d *Duration) MarshalJSON() ([]byte, error) {
	// adds mins to the value
	jsonValue := fmt.Sprintf("%d mins", d)

	// format to a valid json string (quoted string)
	quotedJSONValue := strconv.Quote(jsonValue)

	// convert the quoted string value to a byte slice and return it
	return []byte(quotedJSONValue), nil
}

func (d *Duration) UnmarshalJSON(jsonValue []byte) error {
	var parts []string
	// remove the double quotes for the JSON value
	unQuotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidDurationFormat
	}

	// added validation logic
	if !strings.Contains(unQuotedJSONValue, " ") {
		return ErrInvalidDurationFormat
	}

	parts = strings.Split(unQuotedJSONValue, " ")

	// sanity check
	if len(parts) != 2 || parts[1] != "mins" {

		return ErrInvalidDurationFormat
	}

	// convert the string to int32
	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {

		return ErrInvalidDurationFormat
	}

	// convert the i to duration type
	*d = Duration(i)

	return nil
}
