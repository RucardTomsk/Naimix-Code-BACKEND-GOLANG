package helpers

import (
	"errors"
	"strconv"
	"time"
)

/*
ParseUnixTimeStampToString
The function parse Unix TimeStamp takes a string representing a Unix timestamp as input, and returns a formatted string and an error as output.

The input string is first parsed into an integer using the strconv.ParseInt function with a base of 10 and a bit size of 64.
If parsing is successful, the function converts the Unix timestamp into a human-readable format using the time.Unix function and the Format method.
The formatted string is returned along with a nil error.

If there is an error during the parsing of the input string, the function returns an empty string and an error containing the message "parse unix_timestamp error".
*/
func ParseUnixTimeStampToString(stringUnixTimeStamp string) (string, error) {
	unixTimeStamp, err := strconv.ParseInt(
		stringUnixTimeStamp,
		10,
		64)

	if err != nil {
		return "", errors.New("parse unix_timestamp error")
	}
	return time.Unix(unixTimeStamp, 0).Format("2 January 2006 15:04"), nil
}

/*
ParseUnixTimeStampToTime
The function parse Unix TimeStamp takes a string representing a Unix timestamp as input, and returns a formatted time.Time and an error as output.

The input string is first parsed into an integer using the strconv.ParseInt function with a base of 10 and a bit size of 64.
If parsing is successful, the function converts the Unix timestamp into a human-readable format using the time.Unix function and the Format method.
The formatted string is returned along with a nil error.

If there is an error during the parsing of the input string, the function returns an empty string and an error containing the message "parse unix_timestamp error".
*/
func ParseUnixTimeStampToTime(stringUnixTimeStamp string) (*time.Time, error) {
	unixTimeStamp, err := strconv.ParseInt(
		stringUnixTimeStamp,
		10,
		64)

	if err != nil {
		return nil, errors.New("parse unix_timestamp error")
	}
	newTime := time.Unix(unixTimeStamp, 0)
	return &newTime, nil
}
