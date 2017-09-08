package main

import (
	"os"
	"path"

	"github.com/coderconvoy/lazyf"
)

func confLocs() {
	home := os.Getenv("HOME")
	return []string{
		"test_data/.sitemanconf",
		".sitemanconf",
		path.Join(home, ".sitemanconf"),
		path.Join(home, ".config/siteman/init"),
	}
}

func getConfig(cloc string, usedef bool) lazyf.LZ {
	if usedef {
		return lazyf.LZ{}
	}
	if cloc == "" {
		return lazyf.GetConfigN(0, confLocs())
	}
	return lazyf.GetConfigN(0, cloc)
}
