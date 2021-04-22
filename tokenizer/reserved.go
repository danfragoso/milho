package tokenizer

func isBoolean(buffer string) bool {
	switch buffer {
	case "True", "False":
		return true
	}

	return false
}
