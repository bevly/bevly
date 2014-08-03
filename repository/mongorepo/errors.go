package mongorepo

import (
	"fmt"
	"strings"
)

type compositeError struct {
	Errors []error
}

func (c *compositeError) Add(err error) {
	if err != nil {
		c.Errors = append(c.Errors, err)
	}
}

func (c *compositeError) IsError() bool {
	return len(c.Errors) > 0
}

func (c *compositeError) Error() string {
	res := make([]string, len(c.Errors))
	for i, err := range c.Errors {
		res[i] = err.Error()
	}
	return fmt.Sprintf("%d errors: %s", len(c.Errors), strings.Join(res, ", "))
}
