// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package document

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-ls/internal/uri"
)

// Handle represents a document location
//
// This may be received via LSP from the client (as URI)
// or constructed from a file path on OS FS.
type Handle struct {
	Dir      DirHandle
	Filename string
}

// HandleFromURI creates a Handle from a given URI.
//
// docURI is expected to be a document URI (rather than dir).
// It is however outside the scope of the function to verify
// this is actually the case or whether the file exists.
func HandleFromURI(docUri string) Handle {
	filename := path.Base(docUri)
	dirUri := strings.TrimSuffix(docUri, "/"+filename)

	return Handle{
		Dir:      DirHandleFromURI(dirUri),
		Filename: filepath.Base(uri.MustPathFromURI(docUri)),
	}
}

// HandleFromPath creates a Handle from a given path.
//
// docPath is expected to be a document path (rather than dir).
// It is however outside the scope of the function to verify
// this is actually the case or whether the file exists.
func HandleFromPath(docPath string) Handle {
	filename := filepath.Base(docPath)
	dirPath := strings.TrimSuffix(docPath, fmt.Sprintf("%c%s", os.PathSeparator, filename))

	return Handle{
		Dir:      DirHandleFromPath(dirPath),
		Filename: filename,
	}
}

func (h Handle) FullPath() string {
	return filepath.Join(h.Dir.Path(), h.Filename)
}

func (h Handle) FullURI() string {
	return h.Dir.URI + "/" + h.Filename
}
