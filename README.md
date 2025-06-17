# lib-client-document-storage-go

Go client library for Stonal document storage API.

## Overview

This library provides a Go client for the Stonal document storage API, allowing you to upload files with manifests to the document storage service.

## Installation

```bash
go get github.com/stonal-tech/lib-client-document-storage-go/stonaldoc
```

## Usage

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    
    "github.com/stonal-tech/lib-client-document-storage-go/stonaldoc"
    "github.com/stonal-tech/lib-httpclient-go/stonalhttpclient"
)

func main() {
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
    httpClient := stonalhttpclient.New(http.DefaultClient)
    
    // Implement your authenticator
    auth := &MyAuthenticator{}
    
    client := stonaldoc.New(logger, httpClient, auth)
    
    org := &stonaldoc.OrgContext{
        Env:     "dev",
        Stack:   "stack1",
        OrgCode: "org123",
    }
    
    manifest := &stonaldoc.FileManifest{
        Asset: &stonaldoc.Asset{
            ExternalIDs: map[string]string{
                "id1": "value1",
            },
        },
    }
    
    file, err := os.Open("document.pdf")
    if err != nil {
        panic(err)
    }
    
    response, err := client.Upload(context.Background(), org, manifest, file, "document.pdf")
    if err != nil {
        panic(err)
    }
    
    logger.Info("Upload successful", "documentId", response.DocumentID)
}

type MyAuthenticator struct{}

func (a *MyAuthenticator) GetKeycloakAuthToken() (*stonaldoc.JWT, error) {
    // Implement your authentication logic here
    return &stonaldoc.JWT{
        AccessToken: "your-token",
    }, nil
}
```

## Types

### OrgContext

Represents the organization context containing environment, stack, and organization code.

### FileManifest

Flexible structure that can handle different JSON formats for file manifests, including asset information, documentation metadata, and folder templates.

### UploadResponse

Response from the upload API containing the document ID and any duplicate document IDs.

## License

This project is proprietary to Stonal Technologies.