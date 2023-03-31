package s3util

import (
	"regexp"
	"strings"

	"github.com/grokify/mogo/type/stringsutil"
)

var rxSafeChars = regexp.MustCompile(`^[0-9a-zA-Z!-_\.*'\(\)]+$`)

// ObjectKeyNameIsSafe checks to see if all chars are within the safe characgter set.
// See more here: https://docs.aws.amazon.com/AmazonS3/latest/userguide/object-keys.html
func ObjectKeyNameIsSafe(name string) bool {
	if name == "" {
		return false
	} else if strings.Index(name, ".") == 0 {
		return false
	} else if stringsutil.ReverseIndex(name, ".") == 0 {
		return false
	} else if !rxSafeChars.MatchString(name) {
		return false
	}
	return true
}

/*
0-9

a-z

A-Z

Special characters
Exclamation point (!)

Hyphen (-)

Underscore (_)

Period (.)

Asterisk (*)

Single quote (')

Open parenthesis (()

Close parenthesis ())
*/
