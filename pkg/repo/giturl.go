package repo

import (
	"fmt"
	"net/url"
)

// AddGitURLUserPassword combines the optional username and password to make a git url for cloning git
func AddGitURLUserPassword(rawurl, username, password string) (string, error) {
	if username == "" && password == "" {
		return rawurl, nil
	}

	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl, fmt.Errorf("failed to parse git URL %s: %w", rawurl, err)
	}

	user := u.User
	if user != nil {
		if username == "" {
			username = user.Username()
		}
		if password == "" {
			password, _ = user.Password()
		}
	}
	u.User = url.UserPassword(username, password)
	return u.String(), nil
}
