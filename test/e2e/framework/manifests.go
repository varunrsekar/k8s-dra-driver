// Copyright The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package framework

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

// Specs are the embedded YAML / YAML-template sources for DRA workloads.
// text/template lets CEL selectors reference runtime values detected from
// the cluster's ResourceSlice.
//
//go:embed specs
var Specs embed.FS

// Render looks up specs/<name>, specs/<name>.yaml, or specs/<name>.yaml.tmpl
// and returns the rendered bytes. `vars` is the template data.
func Render(name string, vars any) ([]byte, error) {
	for _, candidate := range []string{name, name + ".yaml", name + ".yaml.tmpl"} {
		data, err := Specs.ReadFile("specs/" + candidate)
		if err != nil {
			continue
		}
		tmpl, err := template.New(name).Option("missingkey=error").Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", candidate, err)
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, vars); err != nil {
			return nil, fmt.Errorf("execute %s: %w", candidate, err)
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("spec not found: %s (tried %s, %s.yaml, %s.yaml.tmpl)", name, name, name, name)
}

// ApplyYAML shells out to `kubectl apply -f -`. Multi-document YAML and the
// server-side-apply subtleties are left to kubectl.
func ApplyYAML(ctx context.Context, yaml []byte) error {
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Stdin = bytes.NewReader(yaml)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
