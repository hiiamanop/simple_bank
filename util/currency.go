package util

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// Currency is the type for handling currency enums
type Currency string

// List of supported currencies
const (
	USD Currency = "USD"
	EUR Currency = "EUR"
)

// Value implements the driver.Valuer interface
func (c Currency) Value() (driver.Value, error) {
	return string(c), nil
}

// Scan implements the sql.Scanner interface
func (c *Currency) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("currency cannot be null")
	}

	// Handle both string and []byte
	switch v := value.(type) {
	case []byte:
		*c = Currency(v)
		return nil
	case string:
		*c = Currency(v)
		return nil
	default:
		return fmt.Errorf("cannot scan type %T into Currency", value)
	}
}

// MarshalJSON implements the json.Marshaler interface
func (c Currency) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (c *Currency) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*c = Currency(s)
	return nil
}

// Validate checks if the currency is supported
func (c Currency) Validate() error {
	switch c {
	case USD, EUR:
		return nil
	default:
		return fmt.Errorf("invalid currency: %s", c)
	}
}
