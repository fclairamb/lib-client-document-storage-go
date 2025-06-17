package clientdoc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

// APIClient is the interface for the API client.
type APIClient interface {
	// Upload uploads a file to the API
	Upload(
		ctx context.Context,
		org *OrgContext,
		manifest *FileManifest,
		fileStream io.ReadCloser,
		fileName string,
	) (*UploadResponse, error)
}

// UploadResponse is the response from the API.
type UploadResponse struct {
	// DocumentID is the ID of the document that was uploaded.
	DocumentID string `json:"documentId"`
	// DuplicateDocumentIDs is the IDs of the documents that are duplicates of the uploaded document.
	DuplicateDocumentIDs []string `json:"duplicateDocumentIds"`
}

// OrgContext holds the organization context (env and organization ID).
type OrgContext struct {
	Env     string `json:"env"`
	Stack   string `json:"stack"`
	OrgCode string `json:"orgCode"`
}

var ErrFailedToSplitPath = errors.New("failed to split path")

// ParseOrgContext parses an organization context from a path.
func ParseOrgContext(s string) (OrgContext, string, error) {
	const expectedParts = 4
	spl := strings.SplitN(s, "/", expectedParts)
	if len(spl) < expectedParts {
		return OrgContext{}, "", fmt.Errorf("%w: %s", ErrFailedToSplitPath, s)
	}
	return OrgContext{
		Env:     spl[0],
		Stack:   spl[1],
		OrgCode: spl[2],
	}, spl[3], nil
}
