package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// UnixTime is a custom time type that can unmarshal from numeric Unix timestamps.
// Kraken API returns timestamps as floating point numbers (e.g., 1616663618.2637).
type UnixTime struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler for UnixTime.
// It handles both numeric timestamps (int or float) and string timestamps.
func (t *UnixTime) UnmarshalJSON(data []byte) error {
	// Handle null
	if string(data) == "null" || string(data) == "0" {
		t.Time = time.Time{}
		return nil
	}

	// Try to unmarshal as a number (most common for Kraken)
	var timestamp float64
	if err := json.Unmarshal(data, &timestamp); err == nil {
		sec := int64(timestamp)
		nsec := int64((timestamp - float64(sec)) * 1e9)
		t.Time = time.Unix(sec, nsec)
		return nil
	}

	// Try to unmarshal as a string
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		// Try parsing as numeric string first
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			sec := int64(f)
			nsec := int64((f - float64(sec)) * 1e9)
			t.Time = time.Unix(sec, nsec)
			return nil
		}
		// Try RFC3339
		if parsed, err := time.Parse(time.RFC3339, str); err == nil {
			t.Time = parsed
			return nil
		}
		// Try RFC3339Nano
		if parsed, err := time.Parse(time.RFC3339Nano, str); err == nil {
			t.Time = parsed
			return nil
		}
	}

	return fmt.Errorf("UnixTime: cannot unmarshal %s", string(data))
}

// MarshalJSON implements json.Marshaler for UnixTime.
func (t UnixTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("0"), nil
	}
	// Return as floating point timestamp
	timestamp := float64(t.Unix()) + float64(t.Nanosecond())/1e9
	return json.Marshal(timestamp)
}

// IsZero returns true if the time is zero.
func (t UnixTime) IsZero() bool {
	return t.Time.IsZero()
}
