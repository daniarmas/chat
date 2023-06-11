package inputs

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Logout   bool   `json:"logout"`
}

func (in *SignInInput) Sanitize() {
	in.Email = strings.TrimSpace(in.Email)
	in.Email = strings.ToLower(in.Email)

}

func (in SignInInput) Validate() error {
	if strings.Contains(in.Email, "@") {
		if !emailRegex.MatchString(in.Email) {
			return fmt.Errorf("email not valid")
		}
	}
	return nil
}
