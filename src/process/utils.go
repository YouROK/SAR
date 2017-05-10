package process

import (
	"io/ioutil"
	"menu"
	"os"
	"path/filepath"
)

func findDirs(dirs []string) []string {
	retdirs := make([]string, 0)
	for _, d := range dirs {
		list, _ := ioutil.ReadDir(d)
		for _, dd := range list {
			if dd.IsDir() {
				retdirs = append(retdirs, filepath.Join(d, dd.Name()))
			}
		}
	}
	return retdirs
}

func permAll(path string) {
	info, err := os.Stat(path)
	if err != nil {
		menu.Text("Ошибка применения прав:", err)
		return
	}
	if filepath.Ext(path) == ".apk" && info.Mode().IsRegular() {
		os.Chown(path, 0, 0)
		os.Chmod(path, 0644)
	}
	if info.Mode().IsDir() {
		os.Chown(path, 0, 0)
		os.Chmod(path, 0755)
		list, _ := ioutil.ReadDir(path)
		for _, d := range list {
			permAll(filepath.Join(path, d.Name()))
		}
	}
}
