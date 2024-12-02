// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package product

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/chushi-io/lf-install/internal/build"
	"github.com/hashicorp/go-version"
)

var (
	vaultVersionOutputRe = regexp.MustCompile(`OpenBao ` + simpleVersionRe)
)

var OpenBao = Product{
	Name: "vault",
	BinaryName: func() string {
		if runtime.GOOS == "windows" {
			return "vault.exe"
		}
		return "vault"
	},
	GetVersion: func(ctx context.Context, path string) (*version.Version, error) {
		cmd := exec.CommandContext(ctx, path, "version")

		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		stdout := strings.TrimSpace(string(out))

		submatches := vaultVersionOutputRe.FindStringSubmatch(stdout)
		if len(submatches) != 2 {
			return nil, fmt.Errorf("unexpected number of version matches %d for %s", len(submatches), stdout)
		}
		v, err := version.NewVersion(submatches[1])
		if err != nil {
			return nil, fmt.Errorf("unable to parse version %q: %w", submatches[1], err)
		}

		return v, err
	},
	BuildInstructions: &BuildInstructions{
		GitRepoURL:    "https://github.com/openbao/openbao.git",
		PreCloneCheck: &build.GoIsInstalled{},
		Build:         &build.GoBuild{},
	},
}
