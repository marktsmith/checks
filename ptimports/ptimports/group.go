// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ptimports

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/palantir/pkg/pkgpath"
)

func repoForFile(filename string) (string, error) {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	relative := abs
	if goPathSrcRel, err := pkgpath.NewAbsPkgPath(abs).GoPathSrcRel(); err == nil {
		relative = goPathSrcRel
	}
	segments := strings.Split(relative, "/")
	if len(segments) < 3 {
		return "", fmt.Errorf("expected repo to be located under at least 3 subdirectories but received relative filepath: %v", relative)
	}
	// append trailing / to prevent matches on repos with superstring names
	return filepath.Join(segments[:3]...) + "/", nil
}

type importGrouper interface {
	importGroup(importPath string) int
}

func newVendoredGrouper(repoPath string) importGrouper {
	return vendoredGrouper{repoPath}
}

// vendoredGrouper groups packages by standard library, vendored, an in-repo
// packages.
type vendoredGrouper struct {
	repoPath string
}

func (g vendoredGrouper) importGroup(importPath string) int {
	switch {
	case inStandardLibrary(importPath):
		return 0
	case !g.inThisRepo(importPath):
		return 1
	default:
		return 2
	}
}

func (g vendoredGrouper) inThisRepo(importPath string) bool {
	if !strings.HasSuffix(importPath, "/") {
		importPath += "/"
	}
	return strings.HasPrefix(importPath, g.repoPath)
}

func inStandardLibrary(importPath string) bool {
	return !strings.Contains(importPath, ".")
}
