package testdata

import (
	"path/filepath"
	"runtime"
)

// basepath is the root directory of this package.
var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

// Path returns the absolute path corresponding to the given path that
// is taken as relative to the testdata directory.
func Path(rel string) string {
	return filepath.Join(basepath, rel)
}

// X509Path returns the absolute path corresponding to the given path that
// is taken as relative to the testdata/x509 directory.
func X509Path(rel string) string {
	return Path(filepath.Join("x509", rel))
}
