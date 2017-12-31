package types

// Name indicates a string that may contain colons.
// Escape its colons before joining with other strings (using colon as a separator).
type Name string

func TranslateEscapedRune(r rune) rune {
	return r
}

func EscapeName(name Name) string {
	// Replace : with \: and \ with \\
	newRunes := []rune{}
	for _, rune := range name {
		switch rune {
		case '\\', ':':
			newRunes = append(newRunes, '\\', rune)
		default:
			newRunes = append(newRunes, rune)
		}
	}

	return string(newRunes)
}

func UnescapeName(name string) Name {
	// Replace \\ with \ and \: with :, but not if the \: was already touched
	newRunes := []rune{}
	isEscaped := false
	for _, rune := range name {
		if isEscaped {
			newRunes = append(newRunes, TranslateEscapedRune(rune))
			isEscaped = false
		} else {
			if rune == '\\' {
				isEscaped = true
			} else {
				newRunes = append(newRunes, rune)
			}
		}
	}

	return Name(newRunes)
}

func SplitAtUnescapedColons(s string) []string {
	segments := []string{}

	segmentStart := 0
	isEscaped := false
	for i, rune := range s {
		if isEscaped {
			isEscaped = false
		} else {
			switch rune {
			case '\\':
				isEscaped = true
			case ':':
				segments = append(segments, s[segmentStart:i])
				segmentStart = i + 1
			}
		}
	}

	if segmentStart <= len(s) {
		segments = append(segments, s[segmentStart:len(s)])
	}

	return segments
}
