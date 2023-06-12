package inputs

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Define the regular expression for semantic versions
var semanticVersRegex = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

type CreateApiKeyInput struct {
	AppVersion     string
	Revoked        bool
	ExpirationTime time.Time
}

func (in *CreateApiKeyInput) Sanitize() {
	in.AppVersion = strings.TrimSpace(in.AppVersion)
	in.AppVersion = strings.ToLower(in.AppVersion)
}

func (in CreateApiKeyInput) Validate() error {
	if !semanticVersRegex.MatchString(in.AppVersion) {
		return fmt.Errorf("app version not valid")
	}
	return nil
}
