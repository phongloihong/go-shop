package valueobject

import (
	"fmt"
	"slices"
	"time"
)

// enums
type TimeFormat string

const (
	FormatDefault  TimeFormat = "2006-01-02 15:04:05"
	FormatISO      TimeFormat = "2006-01-02T15:04:05Z07:00"
	FormatDateOnly TimeFormat = "2006-01-02"
	FormatTimeOnly TimeFormat = "15:04:05"
)

func (tf TimeFormat) IsValid() bool {
	validFormats := []TimeFormat{
		FormatDefault, FormatISO, FormatDateOnly, FormatTimeOnly,
	}

	return slices.Contains(validFormats, tf)
}

// Functional Options Pattern for DateTime parsing options
type ParseOptions struct {
	Timezone string
	Format   TimeFormat
}

type ParseOption func(*ParseOptions)

func WithTimezone(timezone string) ParseOption {
	return func(opts *ParseOptions) {
		opts.Timezone = timezone
	}
}

func WithFormat(format string) ParseOption {
	return func(opts *ParseOptions) {
		opts.Format = TimeFormat(format)
	}
}

// Value Object for DateTime
type DateTime int64

func NewTime(t int64) DateTime {
	return DateTime(t)
}

func (dt DateTime) String(options ...ParseOption) (string, error) {
	opts := &ParseOptions{
		Timezone: "UTC",
		Format:   time.RFC3339,
	}

	// Apply options
	for _, option := range options {
		option(opts)
	}

	if !opts.Format.IsValid() {
		return "", fmt.Errorf("invalid time format: %s", opts.Format)
	}

	// parse timezone
	loc, err := time.LoadLocation(opts.Timezone)
	if err != nil {
		return "", fmt.Errorf("invalid timezone: %v", err)
	}

	// convert timesamp to time.Time
	t := time.Unix(int64(dt), 0).In(loc)

	return t.Format(string(opts.Format)), nil
}

func (dt DateTime) Time() time.Time {
	return time.Unix(int64(dt), 0)
}

func (dt DateTime) Unix() int64 {
	return int64(dt)
}
