package stonaldoc

import (
	"fmt"
	"strings"
)

var (
	// ErrFailedToSplitPath is returned when a path cannot be split properly.
	ErrFailedToSplitPath = fmt.Errorf("failed to split path")
)

func (o *OrgContext) String() string {
	return fmt.Sprintf("%s/%s/%s", o.Env, o.Stack, o.OrgCode)
}

// ParseOrgContext parses an organization context from a path.
func ParseOrgContext(s string) (OrgContext, string, error) {
	const expectedParts = 4
	spl := strings.SplitN(s, "/", expectedParts)
	if len(spl) < expectedParts {
		return OrgContext{}, "", fmt.Errorf("%w: %s", ErrFailedToSplitPath, s)
	}
	return OrgContext{
		Env:     spl[0],
		Stack:   spl[1],
		OrgCode: spl[2],
	}, spl[3], nil
}