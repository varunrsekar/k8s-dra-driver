/*
Copyright The Kubernetes Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flags

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/urfave/cli/v3"
	"k8s.io/klog/v2"
	"k8s.io/utils/dump"

	"sigs.k8s.io/dra-driver-nvidia-gpu/pkg/featuregates"
)

// cliValueWrapper wraps pflag.Value to implement the interfaces of cli.Value (flag.Value + flag.Getter).
type cliValueWrapper struct {
	flagVal pflag.Value
}

// Set implements the Set method from the flag.Value interface.
func (w *cliValueWrapper) Set(s string) error { return w.flagVal.Set(s) }

// String implements the String method from the flag.Value interface.
func (w *cliValueWrapper) String() string { return w.flagVal.String() }

// Get implements the Get method of the flag.Getter interface. It checks if the underlying flag value implements the
// Getter interface. If yes, it explicitly casts the type to the Getter interface and invokes the Get method. Otherwise, it
// defaults to the String method.
func (w *cliValueWrapper) Get() interface{} {
	if g, ok := w.flagVal.(interface{ Get() interface{} }); ok {
		return g.Get()
	}
	return w.flagVal.String()
}

func pflagToCLI(flag *pflag.Flag, category string) cli.Flag {
	var cliValue cli.Value = &cliValueWrapper{flagVal: flag.Value}
	return &cli.GenericFlag{
		Name:        flag.Name,
		Category:    category,
		Usage:       flag.Usage,
		Value:       cliValue,
		Destination: &cliValue,
		Sources:     cli.EnvVars(strings.ToUpper(strings.ReplaceAll(flag.Name, "-", "_"))),
	}
}

func LogStartupConfig(parsedFlags any, loggingConfig *LoggingConfig) {
	// Always log component startup config (level 0).
	klog.Infof("\nFeature gates: %#v\nFlags: %s",
		// Flat boolean map -- no pretty-printing needed.
		featuregates.ToMap(),
		// Based on go-spew's Sdump(), with indentation. Type information is
		// always displayed (cannot be disabled). Rely on `parsedFlags` to also
		// contain the klog log verbosity (we want to log this here in any case,
		// so that the log output explicitly mentions that).
		dump.Pretty(parsedFlags),
	)

	// This is a complex object, comprised of largely static default klog
	// component configuration. Various parts can be overridden via environment
	// variables or CLI flags: it makes sense to log the interpolated config,
	// but only on a high verbosity level.
	klog.V(6).Infof("Logging config: %s", dump.Pretty(loggingConfig))
}
