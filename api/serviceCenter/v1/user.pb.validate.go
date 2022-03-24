// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/serviceCenter/v1/user.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on RegisterRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *RegisterRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// RegisterRequestMultiError, or nil if none found.
func (m *RegisterRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetUser() == nil {
		err := RegisterRequestValidationError{
			field:  "User",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetUser()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, RegisterRequestValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, RegisterRequestValidationError{
					field:  "User",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetUser()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return RegisterRequestValidationError{
				field:  "User",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetDeviceConfigRegisterInfos() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RegisterRequestValidationError{
						field:  fmt.Sprintf("DeviceConfigRegisterInfos[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RegisterRequestValidationError{
						field:  fmt.Sprintf("DeviceConfigRegisterInfos[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RegisterRequestValidationError{
					field:  fmt.Sprintf("DeviceConfigRegisterInfos[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(m.GetDeviceStateRegisterInfos()) < 1 {
		err := RegisterRequestValidationError{
			field:  "DeviceStateRegisterInfos",
			reason: "value must contain at least 1 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetDeviceStateRegisterInfos() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, RegisterRequestValidationError{
						field:  fmt.Sprintf("DeviceStateRegisterInfos[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, RegisterRequestValidationError{
						field:  fmt.Sprintf("DeviceStateRegisterInfos[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return RegisterRequestValidationError{
					field:  fmt.Sprintf("DeviceStateRegisterInfos[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return RegisterRequestMultiError(errors)
	}

	return nil
}

// RegisterRequestMultiError is an error wrapping multiple validation errors
// returned by RegisterRequest.ValidateAll() if the designated constraints
// aren't met.
type RegisterRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterRequestMultiError) AllErrors() []error { return m }

// RegisterRequestValidationError is the validation error returned by
// RegisterRequest.Validate if the designated constraints aren't met.
type RegisterRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterRequestValidationError) ErrorName() string { return "RegisterRequestValidationError" }

// Error satisfies the builtin error interface
func (e RegisterRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterRequestValidationError{}

// Validate checks the field values on RegisterReply with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *RegisterReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on RegisterReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in RegisterReplyMultiError, or
// nil if none found.
func (m *RegisterReply) ValidateAll() error {
	return m.validate(true)
}

func (m *RegisterReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	// no validation rules for Token

	if len(errors) > 0 {
		return RegisterReplyMultiError(errors)
	}

	return nil
}

// RegisterReplyMultiError is an error wrapping multiple validation errors
// returned by RegisterReply.ValidateAll() if the designated constraints
// aren't met.
type RegisterReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m RegisterReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m RegisterReplyMultiError) AllErrors() []error { return m }

// RegisterReplyValidationError is the validation error returned by
// RegisterReply.Validate if the designated constraints aren't met.
type RegisterReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RegisterReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RegisterReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RegisterReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RegisterReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RegisterReplyValidationError) ErrorName() string { return "RegisterReplyValidationError" }

// Error satisfies the builtin error interface
func (e RegisterReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRegisterReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RegisterReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RegisterReplyValidationError{}

// Validate checks the field values on LoginReply with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *LoginReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on LoginReply with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in LoginReplyMultiError, or
// nil if none found.
func (m *LoginReply) ValidateAll() error {
	return m.validate(true)
}

func (m *LoginReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	// no validation rules for Token

	if len(errors) > 0 {
		return LoginReplyMultiError(errors)
	}

	return nil
}

// LoginReplyMultiError is an error wrapping multiple validation errors
// returned by LoginReply.ValidateAll() if the designated constraints aren't met.
type LoginReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m LoginReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m LoginReplyMultiError) AllErrors() []error { return m }

// LoginReplyValidationError is the validation error returned by
// LoginReply.Validate if the designated constraints aren't met.
type LoginReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e LoginReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e LoginReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e LoginReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e LoginReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e LoginReplyValidationError) ErrorName() string { return "LoginReplyValidationError" }

// Error satisfies the builtin error interface
func (e LoginReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sLoginReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = LoginReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = LoginReplyValidationError{}

// Validate checks the field values on UnregisterReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UnregisterReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UnregisterReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UnregisterReplyMultiError, or nil if none found.
func (m *UnregisterReply) ValidateAll() error {
	return m.validate(true)
}

func (m *UnregisterReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Success

	if len(errors) > 0 {
		return UnregisterReplyMultiError(errors)
	}

	return nil
}

// UnregisterReplyMultiError is an error wrapping multiple validation errors
// returned by UnregisterReply.ValidateAll() if the designated constraints
// aren't met.
type UnregisterReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UnregisterReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UnregisterReplyMultiError) AllErrors() []error { return m }

// UnregisterReplyValidationError is the validation error returned by
// UnregisterReply.Validate if the designated constraints aren't met.
type UnregisterReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UnregisterReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UnregisterReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UnregisterReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UnregisterReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UnregisterReplyValidationError) ErrorName() string { return "UnregisterReplyValidationError" }

// Error satisfies the builtin error interface
func (e UnregisterReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUnregisterReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UnregisterReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UnregisterReplyValidationError{}

// Validate checks the field values on DownloadClientCodeRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DownloadClientCodeRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DownloadClientCodeRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DownloadClientCodeRequestMultiError, or nil if none found.
func (m *DownloadClientCodeRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DownloadClientCodeRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Username

	if len(errors) > 0 {
		return DownloadClientCodeRequestMultiError(errors)
	}

	return nil
}

// DownloadClientCodeRequestMultiError is an error wrapping multiple validation
// errors returned by DownloadClientCodeRequest.ValidateAll() if the
// designated constraints aren't met.
type DownloadClientCodeRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DownloadClientCodeRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DownloadClientCodeRequestMultiError) AllErrors() []error { return m }

// DownloadClientCodeRequestValidationError is the validation error returned by
// DownloadClientCodeRequest.Validate if the designated constraints aren't met.
type DownloadClientCodeRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DownloadClientCodeRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DownloadClientCodeRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DownloadClientCodeRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DownloadClientCodeRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DownloadClientCodeRequestValidationError) ErrorName() string {
	return "DownloadClientCodeRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DownloadClientCodeRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDownloadClientCodeRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DownloadClientCodeRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DownloadClientCodeRequestValidationError{}

// Validate checks the field values on File with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *File) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on File with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in FileMultiError, or nil if none found.
func (m *File) ValidateAll() error {
	return m.validate(true)
}

func (m *File) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Content

	// no validation rules for Name

	if len(errors) > 0 {
		return FileMultiError(errors)
	}

	return nil
}

// FileMultiError is an error wrapping multiple validation errors returned by
// File.ValidateAll() if the designated constraints aren't met.
type FileMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FileMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FileMultiError) AllErrors() []error { return m }

// FileValidationError is the validation error returned by File.Validate if the
// designated constraints aren't met.
type FileValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FileValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FileValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FileValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FileValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FileValidationError) ErrorName() string { return "FileValidationError" }

// Error satisfies the builtin error interface
func (e FileValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFile.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FileValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FileValidationError{}
