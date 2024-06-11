package utils

func TrimText(text string, keep int) string {
	if len(text) > keep {
		return text[0:keep] + "..."
	}

	return text
}
