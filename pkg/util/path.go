package util

import (
	"log"
	u "os/user"
	p "path"
	"strings"
)

// Abs expands "~/example" to "/User/$user/example" and returns a
// clean path.
func Abs(path string) string {
	var user *u.User
	var err error

	if user, err = u.Current(); err != nil {
		log.Fatalf("Unable to identify current user")
	}

	if strings.HasPrefix(path, "~/") {
		return p.Join(user.HomeDir, strings.TrimPrefix(path, "~/"))
	}

	return p.Clean(path)
}
