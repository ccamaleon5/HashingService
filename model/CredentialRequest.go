package model

import (
	strfmt "github.com/go-openapi/strfmt"
)

type CredentialRequest struct{
	Credentials []CredentialSubject
}

type CredentialSubject struct {

	// The claims that will be generated with the credential
	Content interface{} `json:"content,omitempty"`

	// The evidence obtained from the validation of the claims, may be photos, physical documents, links, etc
	Evidence interface{} `json:"evidence,omitempty"`

	// credential expiration date
	// Format: date-time
	ExpirationDate strfmt.DateTime `json:"expirationDate,omitempty"`

	// credential issuance date
	// Format: date-time
	IssuanceDate strfmt.DateTime `json:"issuanceDate,omitempty"`

	// Credential Type
	Type string `json:"type,omitempty"`
}