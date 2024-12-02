// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package releases

import (
	"context"
	"testing"

	"github.com/chushi-io/lf-install/internal/testutil"
	"github.com/chushi-io/lf-install/product"
	"github.com/chushi-io/lf-install/src"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
)

func TestVersions_List(t *testing.T) {
	testutil.EndToEndTest(t)

	cons, err := version.NewConstraint(">= 1.0.0, < 1.0.10")
	if err != nil {
		t.Fatal(err)
	}

	versions := &Versions{
		Product:     product.OpenTofu,
		Constraints: cons,
	}

	ctx := context.Background()
	sources, err := versions.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expectedVersions := []string{
		"1.0.0",
		"1.0.1",
		"1.0.2",
		"1.0.3",
		"1.0.4",
		"1.0.5",
		"1.0.6",
		"1.0.7",
		"1.0.8",
		"1.0.9",
	}
	if diff := cmp.Diff(expectedVersions, sourcesToRawVersions(sources)); diff != "" {
		t.Fatalf("unexpected versions: %s", diff)
	}
}

func TestVersions_List_enterprise(t *testing.T) {
	testutil.EndToEndTest(t)

	cons, err := version.NewConstraint(">= 1.9.0, < 1.9.9")
	if err != nil {
		t.Fatal(err)
	}

	versions := &Versions{
		Product:     product.OpenBao,
		Constraints: cons,
		Install: InstallationOptions{
			LicenseDir: "/some/path",
		},
		Enterprise: &EnterpriseOptions{
			Meta: "hsm",
		},
	}

	ctx := context.Background()
	sources, err := versions.List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	expectedVersions := []string{
		"1.9.0+ent.hsm",
		"1.9.1+ent.hsm",
		"1.9.2+ent.hsm",
		"1.9.3+ent.hsm",
		"1.9.4+ent.hsm",
		"1.9.5+ent.hsm",
		"1.9.6+ent.hsm",
		"1.9.7+ent.hsm",
		"1.9.8+ent.hsm",
	}
	if diff := cmp.Diff(expectedVersions, sourcesToRawVersions(sources)); diff != "" {
		t.Fatalf("unexpected versions: %s", diff)
	}

	for _, source := range sources {
		if *source.(*ExactVersion).Enterprise != *versions.Enterprise {
			t.Fatalf("unexpected Enterprise data: %v", source.(*ExactVersion).Enterprise)
		}

		if source.(*ExactVersion).Enterprise == versions.Enterprise {
			t.Fatalf("the Enterprise data should be copied, not referenced")
		}
	}
}

func sourcesToRawVersions(srcs []src.Source) []string {
	rawVersions := make([]string, len(srcs))

	for idx, src := range srcs {
		source := src.(*ExactVersion)
		rawVersions[idx] = source.Version.String()
	}

	return rawVersions
}
