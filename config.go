package main

import (
	"os"
	"path"

	"github.com/coderconvoy/lazyf"
)

func confLocs() []string {
	home := os.Getenv("HOME")
	return []string{
		"test_data/.sitemanconf",
		".sitemanconf",
		path.Join(home, ".sitemanconf"),
		path.Join(home, ".config/siteman/init"),
	}
}

func getConfig(cloc string, usedef bool) (lazyf.LZ, error) {
	if usedef {
		return lazyf.LZ{}, nil
	}
	if cloc == "" {
		return lazyf.GetConfigN(0, confLocs()...)
	}
	return lazyf.GetConfigN(0, cloc)
}
