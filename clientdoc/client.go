package clientdoc

import "fmt"

func (o *OrgContext) String() string {
	return fmt.Sprintf("%s/%s/%s", o.Env, o.Stack, o.OrgCode)
}
