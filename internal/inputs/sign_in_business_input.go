package inputs

import (
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type SignInBusinessInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Logout   bool   `json:"logout"`
}

func (in *SignInBusinessInput) Sanitize() {
	in.Email = strings.TrimSpace(in.Email)
	in.Email = strings.ToLower(in.Email)

}

func (in SignInBusinessInput) Validate() error {
	if strings.Contains(in.Email, "@") {
		if !emailRegex.MatchString(in.Email) {
			return fmt.Errorf("email not valid")
		}
	}
	return nil
}
