package gowebdav

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
)

// User represents a user entry from .netrc
type User struct {
	Login    string
	Password string
	Machine  string
}

func parseLine(s string) (login, pass string) {
	fields := strings.Fields(s)
	for i, f := range fields {
		if f == "login" {
			login = fields[i+1]
		}
		if f == "password" {
			pass = fields[i+1]
		}
	}
	return login, pass
}

// ReadConfig reads a configuration from ~/.netrc
func ReadConfig(uri string) (*User, error) {
	var usr User
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	usr.Machine = u.Host

	curU, err := user.Current()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path.Join(curU.HomeDir, ".netrc"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		re := fmt.Sprintf(`^.*machine %s.*$`, usr.Machine)
		s := scanner.Text()
		matched, err := regexp.MatchString(re, s)
		if err != nil {
			return nil, err
		}
		if matched {
			usr.Login, usr.Password = parseLine(s)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &usr, nil
}
