package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/philipstanev/Miku-stream/internal/sessions"
)

func getPermissions(endpoint string) []string {
	// The underlying data is safe inside the function scope
	permissions := map[string][]string{
		"POST /register/": {"user", "admin"},
		"POST /login/":    {"user", "admin"},
	}
	return permissions[endpoint]
}

func checkPermissionTable(endpoint string, role string) error {
	if len(getPermissions(endpoint)) == 0 {
		return nil
	}

	for _, a := range getPermissions(endpoint) {
		if a == role {
			return nil
		}
	}

	return errors.New("Unauthorized")
}

func getSessionIdFromHeader(r *http.Request) (string, error) {
	authField := r.Header.Get("Authorization")
	if authField == "" {
		return "", errors.New("Missing session ID")
	}

	split := strings.SplitN(authField, " ", 2)
	if split[0] != "Bearer" {
		return "", errors.New("Authentication scheme missing or incorrect")
	}

	return split[1], nil
}
func CheckPermission(w http.ResponseWriter, r *http.Request, s sessions.Session) error {
	sessionID, err := getSessionIdFromHeader(r)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return err
	}
	sessionInfo, err := s.GetInfo(sessionID)
	if err != nil {
		http.Error(w, "Unauthenticated", 401)
		return err
	}
	err = checkPermissionTable(r.Pattern, sessionInfo.Role)
	if err != nil {
		http.Error(w, err.Error(), 403)
		return err
	}

	return nil

}
