package stonaldoc

import (
	"context"
	"io"
)

// Asset represents asset information with external IDs.
type Asset struct {
	ExternalIDs map[string]string `json:"externalIds"`
}

// Documentation represents documentation information.
type Documentation struct {
	UID string `json:"uid"`
}

// Folder represents folder information.
type Folder struct {
	Template string `json:"template"`
}

// FileManifest represents the flexible structure that can handle different JSON formats.
type FileManifest struct {
	// Fields for the first JSON format
	Asset         *Asset         `json:"asset,omitempty"`
	Documentation *Documentation `json:"documentation,omitempty"`
	Folder        *Folder        `json:"folder,omitempty"`

	// Field for the second JSON format
	Disconnected *bool `json:"disconnected,omitempty"`
}

// OrgContext holds the organization context (env and organization ID).
type OrgContext struct {
	Env     string `json:"env"`
	Stack   string `json:"stack"`
	OrgCode string `json:"orgCode"`
}

// UploadResponse is the response from the API.
type UploadResponse struct {
	// DocumentID is the ID of the document that was uploaded.
	DocumentID string `json:"documentId"`
	// DuplicateDocumentIDs is the IDs of the documents that are duplicates of the uploaded document.
	DuplicateDocumentIDs []string `json:"duplicateDocumentIds"`
}

// Authenticator interface for authentication providers.
type Authenticator interface {
	GetKeycloakAuthToken() (*JWT, error)
}

// JWT represents a JWT token response from Keycloak.
type JWT struct {
	AccessToken string `json:"access_token"`
}

// Client is the interface for the document storage API client.
type Client interface {
	// Upload uploads a file to the document storage API
	Upload(
		ctx context.Context,
		org *OrgContext,
		manifest *FileManifest,
		fileStream io.ReadCloser,
		fileName string,
	) (*UploadResponse, error)
}