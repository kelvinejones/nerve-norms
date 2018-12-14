package mem

import "fmt"

// This section has not been implemented, so skip it

type StrengthDuration struct{}

func (section StrengthDuration) Header() string {
	return "STRENGTH-DURATION DATA"
}

func (sd StrengthDuration) String() string {
	return fmt.Sprintf("StrengthDuration{Import not implemented}")
}

func (section *StrengthDuration) Parse(result []string) error {
	return nil
}
