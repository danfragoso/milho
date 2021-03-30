package tokenizer

func isReserved(buffer string) bool {
	switch buffer {
	case "True", "False", "defn", "def":
		return true
	}

	return false
}
