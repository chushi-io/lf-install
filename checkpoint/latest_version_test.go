// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package checkpoint

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
	_ src.Installable    = &LatestVersion{}
	_ src.Removable      = &LatestVersion{}
	_ src.LoggerSettable = &LatestVersion{}
)

func TestLatestVersion(t *testing.T) {
	testutil.EndToEndTest(t)

	lv := &LatestVersion{
		Product: product.OpenTofu,
	}
	lv.SetLogger(testutil.TestLogger())

	ctx := context.Background()

	execPath, err := lv.Install(ctx)
	if err != nil {
		t.Fatal(err)
	}

	licensePath := filepath.Join(filepath.Dir(execPath), "LICENSE.txt")
	t.Cleanup(func() {
		lv.Remove(ctx)
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

func TestLatestVersionValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		lv          LatestVersion
		expectedErr error
	}{
		"Product-incorrect-binary-name": {
			lv: LatestVersion{
				Product: product.Product{
					BinaryName: func() string { return "invalid!" },
					Name:       product.OpenTofu.Name,
				},
			},
			expectedErr: fmt.Errorf("invalid binary name: \"invalid!\""),
		},
		"Product-incorrect-name": {
			lv: LatestVersion{
				Product: product.Product{
					BinaryName: product.OpenTofu.BinaryName,
					Name:       "invalid!",
				},
			},
			expectedErr: fmt.Errorf("invalid product name: \"invalid!\""),
		},
		"Product-valid": {
			lv: LatestVersion{
				Product: product.OpenTofu,
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := testCase.lv.Validate()

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
