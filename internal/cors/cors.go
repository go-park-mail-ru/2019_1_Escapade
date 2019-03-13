package cors

func IsAllowed(get string, origins []string) (allowed bool) {
	allowed = false
	for _, str := range origins {
		if str == get {
			allowed = true
			break
		}
	}

	return
}
