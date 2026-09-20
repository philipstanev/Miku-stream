package auth

func getPermissions(endpoint string) []string {
	// The underlying data is safe inside the function scope
	permissions := map[string][]string{
		"POST /register/": {"user", "admin"},
		"POST /login/":    {"user", "admin"},
	}
	return permissions[endpoint]
}

func CheckPermission(endpoint string, role string) bool {
	if len(getPermissions(endpoint)) == 0 {
		return true
	}

	for _, a := range getPermissions(endpoint) {
		if a == role {
			return true
		}
	}

	return false
}
