package mem

import (
	"fmt"
	"regexp"
)

// This section has not been implemented, so skip it

type StrengthDuration struct{}

func (section StrengthDuration) Header() string {
	return "STRENGTH-DURATION DATA"
}

func (sd StrengthDuration) String() string {
	return fmt.Sprintf("StrengthDuration{Import not implemented}")
}

func (sd StrengthDuration) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^SD\.\d+.*`)
}

func (section *StrengthDuration) ParseLine(result []string) error {
	return nil
}
