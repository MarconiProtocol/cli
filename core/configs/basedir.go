package configs

import (
  "path/filepath"
  "strings"
)

var baseDir = ""

func GetBaseDir() string {
  return baseDir
}

func GetFullPath(childPath string) string {
  return filepath.Join(baseDir, childPath)
}

func SetBaseDir(newBaseDir string) {
  baseDir = strings.TrimRight(newBaseDir, "/")
}
