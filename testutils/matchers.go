package testutils

import (
	"database/sql/driver"
	"time"
)

// AnyTime is a custom matcher for time.Time arguments in SQL queries
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface for time.Time
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
