// Copyright 2021-2023 Zenauth Ltd.
// SPDX-License-Identifier: Apache-2.0

package util

import (
	"crypto/md5" //nolint:gosec
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	pdpv1 "github.com/cerbos/cloud-api/genpb/cerbos/cloud/pdp/v1"
	"github.com/google/uuid"
)

var (
	AppName   = "cerbos"
	BuildDate = "unknown"
	Commit    = "unknown"
	Version   = "unknown"
)

func AppVersion() string {
	var sb strings.Builder
	_, _ = sb.WriteString(Version)
	_, _ = sb.WriteString(fmt.Sprintf("\nBuild timestamp: %s\n", BuildDate))
	_, _ = sb.WriteString(fmt.Sprintf("Build commit: %s\n", Commit))

	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Sum != "" {
			_, _ = sb.WriteString(fmt.Sprintf("Module version: %s, Module checksum: %s\n", info.Main.Version, info.Main.Sum))
		}

		_, _ = sb.WriteString(fmt.Sprintf("Go version: %s\n", info.GoVersion))
		for _, bs := range info.Settings {
			if strings.HasPrefix(bs.Key, "vcs") {
				_, _ = sb.WriteString(fmt.Sprintf("%s: %s\n", bs.Key, bs.Value))
			}
		}
	}

	return sb.String()
}

func AppShortVersion() string {
	var sb strings.Builder
	_, _ = sb.WriteString(Version)

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, bs := range info.Settings {
			if bs.Key == "vcs.revision" {
				_, _ = sb.WriteString("-")
				_, _ = sb.WriteString(bs.Value)
			}
		}
	}

	return sb.String()
}

func PDPIdentifier(instanceID string) *pdpv1.Identifier {
	if instanceID == "" {
		//nolint:gosec
		nodeID := md5.Sum(uuid.NodeID())
		instanceID = fmt.Sprintf("%X-%d", nodeID, os.Getpid())
	}

	return &pdpv1.Identifier{
		Instance: instanceID,
		Version:  AppShortVersion(),
	}
}
