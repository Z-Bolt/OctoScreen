
package utils

import (
	"fmt"
	"strings"
)


var translatedTags = [][2]string{{"strong", "b"}}
var disallowedTags = []string{"p"}

func CleanHTML(html string) string {
	for _, tag := range translatedTags {
		html = ReplaceHTMLTag(html, tag[0], tag[1])
	}

	for _, tag := range disallowedTags {
		html = ReplaceHTMLTag(html, tag, " ")
	}

	return html
}

func ReplaceHTMLTag(html, from, to string) string {
	for _, pattern := range []string{"<%s>", "</%s>", "<%s/>"} {
		to := to
		if to != "" && to != " " {
			to = fmt.Sprintf(pattern, to)
		}

		html = strings.Replace(html, fmt.Sprintf(pattern, from), to, -1)
	}

	return html
}

// TODO: Clean up StrEllipsis(), StrEllipsisLen(), and TruncateString() and consolidate
// into a single function.
func StrEllipsis(name string) string {
	l := len(name)
	if l > 32 {
		return name[:12] + "..." + name[l - 17:l]
	}

	return name
}

func StrEllipsisLen(name string, length int) string {
	l := len(name)
	if l > length {
		return name[:(length / 3)] + "..." + name[l - (length / 3):l]
	}

	return name
}

func TruncateString(str string, length int) string {
	if length < 3 {
		length = 3
	}

	truncateString := str
	if len(str) > length {
		if length > 3 {
			length -= 3
		}

		truncateString = str[0:length] + "..."
	}

	return truncateString
}
