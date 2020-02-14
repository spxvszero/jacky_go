package jk_utils

import (
	"os"
	"path/filepath"
	"strings"
)

//byte
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func AppendPath(pre string, subffix string) string {
	var res string
	pre = strings.TrimSpace(pre)
	subffix = strings.TrimSpace(subffix)
	if pre[len(pre)-1] == '/' {
		if subffix[0] == '/' {
			res = pre + subffix[1:]
		}else {
			res = pre + subffix
		}
	}else {
		if subffix[0] == '/' {
			res = pre + subffix
		}else {
			res = pre + "/" + subffix
		}
	}
	return res
}
