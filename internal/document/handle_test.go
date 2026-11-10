// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package document

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func TestHandleFromURI_escapedFilename(t *testing.T) {
	testCases := []struct {
		filename string
		escaped  string
	}{
		{filename: "dev#local.tfvars", escaped: "dev%23local.tfvars"},
		{filename: "main#local.tf", escaped: "main%23local.tf"},
		{filename: "dev local.tfvars", escaped: "dev%20local.tfvars"},
		{filename: "dev%23local.tfvars", escaped: "dev%2523local.tfvars"},
		{filename: "dev+local.tfvars", escaped: "dev%2Blocal.tfvars"},
		{filename: "café.tfvars", escaped: "caf%C3%A9.tfvars"},
	}
	for _, dirURI := range []string{"file:///random/project%23one", "file:///C%3A/project%23one"} {
		for _, tc := range testCases {
			t.Run(dirURI+"/"+tc.escaped, func(t *testing.T) {
				h := HandleFromURI(dirURI + "/" + tc.escaped)
				if h.Filename != tc.filename {
					t.Fatalf("expected filename %q, got %q", tc.filename, h.Filename)
				}
				expected := Handle{Dir: DirHandleFromURI(dirURI), Filename: tc.filename}
				if h != expected {
					t.Fatalf("expected handle %#v, got %#v", expected, h)
				}
				expectedPath := filepath.Join(expected.Dir.Path(), tc.filename)
				if h.FullPath() != expectedPath {
					t.Fatalf("expected path %q, got %q", expectedPath, h.FullPath())
				}
				if fromPath := HandleFromPath(expectedPath); fromPath != h {
					t.Fatalf("path and URI handles differ: %#v != %#v", fromPath, h)
				}
			})
		}
	}
}

func TestHandleFromURI(t *testing.T) {
	testCases := []struct {
		RawURI         string
		ExpectedHandle Handle
	}{
		{
			RawURI: "file:///random/path/to/config.tf",
			ExpectedHandle: Handle{
				Dir:      DirHandle{URI: "file:///random/path/to"},
				Filename: "config.tf",
			},
		},
		{
			RawURI: "file:///C:/random/path/to/config.tf",
			ExpectedHandle: Handle{
				Dir:      DirHandle{URI: "file:///C:/random/path/to"},
				Filename: "config.tf",
			},
		},
		{
			RawURI: "file:///C%3A/random/path/to/config.tf",
			ExpectedHandle: Handle{
				Dir:      DirHandle{URI: "file:///C:/random/path/to"},
				Filename: "config.tf",
			},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			h := HandleFromURI(tc.RawURI)
			if h != tc.ExpectedHandle {
				t.Fatalf("expected handle: %#v, given: %#v", tc.ExpectedHandle, h)
			}
		})
	}
}

func TestHandleFromPath(t *testing.T) {
	type testCase struct {
		RawURI         string
		ExpectedHandle Handle
	}

	testCases := []testCase{}
	if runtime.GOOS == "windows" {
		testCases = []testCase{
			{
				RawURI: `C:\random\path\to\config.tf`,
				ExpectedHandle: Handle{
					Dir:      DirHandle{URI: "file:///C:/random/path/to"},
					Filename: "config.tf",
				},
			},
		}
	} else {
		testCases = []testCase{
			{
				RawURI: "/random/path/to/config.tf",
				ExpectedHandle: Handle{
					Dir:      DirHandle{URI: "file:///random/path/to"},
					Filename: "config.tf",
				},
			},
		}
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			h := HandleFromPath(tc.RawURI)
			if h != tc.ExpectedHandle {
				t.Fatalf("expected handle: %#v, given: %#v", tc.ExpectedHandle, h)
			}
		})
	}
}

func TestDirHandleFromPath(t *testing.T) {
	type testCase struct {
		RawURI         string
		ExpectedHandle DirHandle
	}

	testCases := []testCase{}
	if runtime.GOOS == "windows" {
		testCases = []testCase{
			{
				RawURI:         `C:\random\path\to`,
				ExpectedHandle: DirHandle{URI: "file:///C:/random/path/to"},
			},
			{
				RawURI:         `C:\random\path\to\`,
				ExpectedHandle: DirHandle{URI: "file:///C:/random/path/to"},
			},
		}
	} else {
		testCases = []testCase{
			{
				RawURI:         "/random/path/to",
				ExpectedHandle: DirHandle{URI: "file:///random/path/to"},
			},
			{
				RawURI:         "/random/path/to/",
				ExpectedHandle: DirHandle{URI: "file:///random/path/to"},
			},
		}
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			h := DirHandleFromPath(tc.RawURI)
			if h != tc.ExpectedHandle {
				t.Fatalf("expected handle: %#v, given: %#v", tc.ExpectedHandle, h)
			}
		})
	}
}

func TestHandle_FullURI(t *testing.T) {
	type testCase struct {
		Handle      Handle
		ExpectedURI string
	}

	testCases := []testCase{
		{
			Handle: Handle{
				Dir:      DirHandle{URI: "file:///C:/random/path/to"},
				Filename: "config.tf",
			},
			ExpectedURI: "file:///C:/random/path/to/config.tf",
		},
		{
			Handle: Handle{
				Dir:      DirHandle{URI: "file:///random/path/to"},
				Filename: "config.tf",
			},
			ExpectedURI: "file:///random/path/to/config.tf",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if tc.ExpectedURI != tc.Handle.FullURI() {
				t.Fatalf("expected URI: %#v, given: %#v", tc.ExpectedURI, tc.Handle.FullURI())
			}
		})
	}
}

func TestHandle_FullPath(t *testing.T) {
	type testCase struct {
		Handle       Handle
		ExpectedPath string
	}

	testCases := []testCase{}
	if runtime.GOOS == "windows" {
		testCases = []testCase{
			{
				Handle: Handle{
					Dir:      DirHandle{URI: "file:///C:/random/path/to"},
					Filename: "config.tf",
				},
				ExpectedPath: `C:\random\path\to\config.tf`,
			},
		}
	} else {
		testCases = []testCase{
			{
				Handle: Handle{
					Dir:      DirHandle{URI: "file:///random/path/to"},
					Filename: "config.tf",
				},
				ExpectedPath: "/random/path/to/config.tf",
			},
		}
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if tc.ExpectedPath != tc.Handle.FullPath() {
				t.Fatalf("expected path: %#v, given: %#v", tc.ExpectedPath, tc.Handle.FullPath())
			}
		})
	}
}
