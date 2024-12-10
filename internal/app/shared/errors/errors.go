package shared

import "fmt"

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}
