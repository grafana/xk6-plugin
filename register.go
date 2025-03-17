// Package deno is a xk6-plugin extension module.
package deno

import (
	"io/fs"
	"log/slog"
	"path/filepath"
	"regexp"

	"github.com/grafana/xk6-plugin/module"
	"go.k6.io/k6/js/modules"
)

func registerPlugin(plugin string) {
	modules.Register(module.ImportPath+"/./"+plugin, module.New("./"+plugin))
}

var reScript = regexp.MustCompile(`\.plugin\.(ts|js|mjs|cjs)$`)

func findPlugins() []string {
	plugins := make([]string, 0)

	err := filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if reScript.MatchString(info.Name()) {
			plugins = append(plugins, filepath.ToSlash(path))
		}

		return nil
	})
	if err != nil {
		slog.Error(err.Error())
	}

	return plugins
}

func registerPlugins() {
	for _, plugin := range findPlugins() {
		registerPlugin(plugin)
	}
}

func init() {
	registerPlugins()
}
