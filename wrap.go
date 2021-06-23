package cli

import "strings"

func wordWrap(x string, linewidth int) string {
	words := strings.Fields(strings.TrimSpace(x))
	if len(words) == 0 {
		return "" // apparently, there are no words
	}

	wrapped := words[0]
	used := len(wrapped)

	for _, word := range words[1:] {
		if used+1+len(word) > linewidth {
			wrapped += "\n" + word
			used = len(word)
		} else {
			wrapped += " " + word
			used += 1 + len(word)
		}
	}

	return wrapped
}
