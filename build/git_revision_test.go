// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package build

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/chushi-io/lf-install/internal/testutil"
	"github.com/chushi-io/lf-install/product"
	"github.com/chushi-io/lf-install/src"
	"github.com/hashicorp/go-version"
)

var (
	_ src.Buildable      = &GitRevision{}
	_ src.Removable      = &GitRevision{}
	_ src.LoggerSettable = &GitRevision{}
)

func TestGitRevision_tofu(t *testing.T) {
	testutil.EndToEndTest(t)

	tempDir := t.TempDir()

	gr := &GitRevision{
		Product:    product.OpenTofu,
		LicenseDir: tempDir,
	}
	gr.SetLogger(testutil.TestLogger())

	ctx := context.Background()

	execPath, err := gr.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}

	licensePath := filepath.Join(tempDir, dstLicenseFileName)
	t.Cleanup(func() {
		gr.Remove(ctx)
		// check if license was deleted
		if _, err := os.Stat(licensePath); !os.IsNotExist(err) {
			t.Fatalf("license file not deleted at %q: %s", licensePath, err)
		}
	})

	v, err := product.OpenTofu.GetVersion(ctx, execPath)
	if err != nil {
		t.Fatal(err)
	}

	latestConstraint, err := version.NewConstraint(">= 1.0")
	if err != nil {
		t.Fatal(err)
	}
	if !latestConstraint.Check(v.Core()) {
		t.Fatalf("versions don't match (expected: %s, installed: %s)",
			latestConstraint, v)
	}

	// check if license was copied
	if _, err := os.Stat(licensePath); err != nil {
		t.Fatalf("expected license file not found at %q: %s", licensePath, err)
	}
}

func TestGitRevision_consul(t *testing.T) {
	testutil.EndToEndTest(t)

	gr := &GitRevision{Product: product.OpenTofu}
	gr.SetLogger(testutil.TestLogger())

	ctx := context.Background()

	execPath, err := gr.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { gr.Remove(ctx) })

	v, err := product.OpenTofu.GetVersion(ctx, execPath)
	if err != nil {
		t.Fatal(err)
	}

	latestConstraint, err := version.NewConstraint(">= 1.0")
	if err != nil {
		t.Fatal(err)
	}
	if !latestConstraint.Check(v.Core()) {
		t.Fatalf("versions don't match (expected: %s, installed: %s)",
			latestConstraint, v)
	}
}

func TestGitRevision_vault(t *testing.T) {
	testutil.EndToEndTest(t)

	gr := &GitRevision{Product: product.OpenBao}
	gr.SetLogger(testutil.TestLogger())

	ctx := context.Background()

	execPath, err := gr.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { gr.Remove(ctx) })

	v, err := product.OpenBao.GetVersion(ctx, execPath)
	if err != nil {
		t.Fatal(err)
	}

	latestConstraint, err := version.NewConstraint(">= 1.0")
	if err != nil {
		t.Fatal(err)
	}
	if !latestConstraint.Check(v.Core()) {
		t.Fatalf("versions don't match (expected: %s, installed: %s)",
			latestConstraint, v)
	}
}

func TestGitRevisionValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		gr          GitRevision
		expectedErr error
	}{
		"Product-incorrect-binary-name": {
			gr: GitRevision{
				Product: product.Product{
					BinaryName: func() string { return "invalid!" },
					Name:       product.OpenTofu.Name,
				},
			},
			expectedErr: fmt.Errorf("invalid binary name: \"invalid!\""),
		},
		"Product-incorrect-name": {
			gr: GitRevision{
				Product: product.Product{
					BinaryName: product.OpenTofu.BinaryName,
					Name:       "invalid!",
				},
			},
			expectedErr: fmt.Errorf("invalid product name: \"invalid!\""),
		},
		"Product-missing-build-instructions": {
			gr: GitRevision{
				Product: product.Product{
					BinaryName: product.OpenTofu.BinaryName,
					Name:       product.OpenTofu.Name,
				},
			},
			expectedErr: fmt.Errorf("no build instructions"),
		},
		"Product-missing-build-instructions-build": {
			gr: GitRevision{
				Product: product.Product{
					BinaryName: product.OpenTofu.BinaryName,
					BuildInstructions: &product.BuildInstructions{
						GitRepoURL: product.OpenTofu.BuildInstructions.GitRepoURL,
					},
					Name: product.OpenTofu.Name,
				},
			},
			expectedErr: fmt.Errorf("missing build instructions"),
		},
		"Product-missing-build-instructions-gitrepourl": {
			gr: GitRevision{
				Product: product.Product{
					BinaryName: product.OpenTofu.BinaryName,
					BuildInstructions: &product.BuildInstructions{
						Build: product.OpenTofu.BuildInstructions.Build,
					},
					Name: product.OpenTofu.Name,
				},
			},
			expectedErr: fmt.Errorf("missing repository URL"),
		},
		"Product-valid": {
			gr: GitRevision{
				Product: product.OpenTofu,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.gr.Validate()

			if err == nil && testCase.expectedErr != nil {
				t.Fatalf("expected error: %s, got no error", testCase.expectedErr)
			}

			if err != nil && testCase.expectedErr == nil {
				t.Fatalf("expected no error, got error: %s", err)
			}

			if err != nil && testCase.expectedErr != nil && err.Error() != testCase.expectedErr.Error() {
				t.Fatalf("expected error: %s, got error: %s", testCase.expectedErr, err)
			}
		})
	}
}
