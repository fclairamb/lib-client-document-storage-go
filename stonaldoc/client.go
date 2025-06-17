package stonaldoc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"github.com/stonal-tech/lib-httpclient-go/stonalhttpclient"
)

var (
	// ErrAPIRequestFailed is returned when an API request fails.
	ErrAPIRequestFailed = errors.New("API request failed")
)

type client struct {
	logger     *slog.Logger
	auth       Authenticator
	httpClient stonalhttpclient.Doer
}

// New creates a new document storage API client.
func New(
	logger *slog.Logger,
	httpClient stonalhttpclient.Doer,
	auth Authenticator,
) Client {
	return &client{
		logger:     logger,
		auth:       auth,
		httpClient: httpClient,
	}
}

func getAPIURL(org *OrgContext) string {
	var domainName string
	switch org.Env {
	case "prod":
		domainName = "stonal.io"
	case "staging":
		domainName = "stonal-staging.io"
	case "test":
		domainName = "stonal-test.io"
	case "dev":
		domainName = "stonal-dev.io"
	default:
		domainName = "?" // Or handle error appropriately
	}

	host := "api." + domainName
	return "https://" + host + "/document-storage/v1/organizations/" + org.OrgCode
}

// Upload uploads a file and its manifest to the Stonal API.
func (c *client) Upload(
	ctx context.Context,
	org *OrgContext,
	manifest *FileManifest,
	file io.ReadCloser,
	fileName string,
) (*UploadResponse, error) {
	defer func() {
		if errClose := file.Close(); errClose != nil {
			c.logger.Error("failed to close file", slog.Any("error", errClose))
		}
	}()

	body, writer, err := c.createMultipartForm(manifest, file, fileName)
	if err != nil {
		return nil, err
	}

	resp, err := c.sendRequest(ctx, org, body, writer)
	if err != nil {
		return nil, err
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			c.logger.Error("failed to close response body", slog.Any("error", errClose))
		}
	}()

	return c.handleResponse(resp)
}

func (c *client) createMultipartForm(manifest *FileManifest, file io.ReadCloser,
	fileName string) (*bytes.Buffer, *multipart.Writer, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add the manifest as JSON
	manifestField, errManifestField := writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": []string{`form-data; name="manifest"; filename="manifest.json"`},
		"Content-Type":        []string{"application/json"},
	})
	if errManifestField != nil {
		return nil, nil, fmt.Errorf("failed to create manifest field: %w", errManifestField)
	}

	manifestJSON, errMarshal := json.Marshal(manifest)
	if errMarshal != nil {
		return nil, nil, fmt.Errorf("failed to convert manifest to JSON: %w", errMarshal)
	}

	if _, errWrite := manifestField.Write(manifestJSON); errWrite != nil {
		return nil, nil, fmt.Errorf("failed to write manifest: %w", errWrite)
	}

	// Add the file
	fileField, errFileField := writer.CreateFormFile("file", fileName)
	if errFileField != nil {
		return nil, nil, fmt.Errorf("failed to create file field: %w", errFileField)
	}

	if _, errCopy := io.Copy(fileField, file); errCopy != nil {
		return nil, nil, fmt.Errorf("failed to copy file data: %w", errCopy)
	}

	if errWriterClose := writer.Close(); errWriterClose != nil {
		return nil, nil, fmt.Errorf("failed to close multipart writer: %w", errWriterClose)
	}

	return &body, writer, nil
}

func (c *client) sendRequest(
	ctx context.Context,
	org *OrgContext,
	body *bytes.Buffer,
	writer *multipart.Writer,
) (*http.Response, error) {
	apiURL := getAPIURL(org) + "/files/upload"

	req, errReq := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, body)
	if errReq != nil {
		return nil, fmt.Errorf("failed to create request: %w", errReq)
	}

	authToken, errAuth := c.auth.GetKeycloakAuthToken()
	if errAuth != nil {
		return nil, fmt.Errorf("failed to get oidc auth token: %w", errAuth)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken.AccessToken)

	resp, errDo := c.httpClient.Do(req)
	if errDo != nil {
		return nil, fmt.Errorf("failed to send request: %w", errDo)
	}

	return resp, nil
}

func (c *client) handleResponse(resp *http.Response) (*UploadResponse, error) {
	respBody, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil, fmt.Errorf("%w: status %d, failed to read response body: %w",
			ErrAPIRequestFailed, resp.StatusCode, errRead)
	}

	c.logger.Info("response", "status", resp.StatusCode, "body", string(respBody))

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrAPIRequestFailed, resp.StatusCode, string(respBody))
	}

	var uploadResponse UploadResponse
	if errUnmarshal := json.Unmarshal(respBody, &uploadResponse); errUnmarshal != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", errUnmarshal)
	}

	c.logger.Info("upload response", "id", uploadResponse.DocumentID)

	return &uploadResponse, nil
}