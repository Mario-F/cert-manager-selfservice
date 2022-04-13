// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.9.1 DO NOT EDIT.
package api

import (
	"time"
)

// Defines values for StatusMessagesSeverity.
const (
	StatusMessagesSeverityError StatusMessagesSeverity = "error"

	StatusMessagesSeverityInfo StatusMessagesSeverity = "info"

	StatusMessagesSeverityWarning StatusMessagesSeverity = "warning"
)

// Status defines model for Status.
type Status struct {
	// This is a overview of all certificates status
	Certificates *struct {
		Failed  int `json:"failed"`
		Pending int `json:"pending"`
		Ready   int `json:"ready"`
		Total   int `json:"total"`
		Unknown int `json:"unknown"`
	} `json:"certificates,omitempty"`

	// This are event messages produced by the cert-manager-selfservice
	Messages *[]struct {
		Message  string                 `json:"message"`
		Severity StatusMessagesSeverity `json:"severity"`
		Time     time.Time              `json:"time"`
	} `json:"messages,omitempty"`
}

// StatusMessagesSeverity defines model for Status.Messages.Severity.
type StatusMessagesSeverity string
