package process

import (
	"io/ioutil"
	"menu"
	"os"
	"path/filepath"
	"strings"
	"utils"
)

func InstallApps() {
	menu.Clear()
	apps := make([]string, 0, 3)
	list, err := ioutil.ReadDir(filepath.Join(sar_dir, "apps"))
	if err != nil {
		menu.Text("Ошибка чтения директории:", err)
		return
	}
	for _, l := range list {
		apps = append(apps, l.Name())
	}
	if len(apps) == 0 {
		menu.Text("Нет приложений для установки")
		return
	}
	itm := menu.List("Выберите приложение", "", apps)
	if itm > 0 && itm <= len(apps) {
		pos := menu.List("", "", []string{"Установить", "Удалить"})
		if pos == 1 {
			installApps(apps[itm-1])
		} else if pos == 2 {
			deleteApps(apps[itm-1])
		}
	}
}

func installApps(name string) {
	appDir := filepath.Join(sar_dir, "apps", name)
	if utils.IsEmptyDir(appDir) {
		menu.Text("Папка", appDir, "пустая")
		return
	}
	dirs := findDirs([]string{filepath.Join(appDir, "system/app"), filepath.Join(appDir, "system/priv-app")})
	for _, d := range dirs {
		ds := filepath.Dir(strings.TrimPrefix(d, appDir))
		err := utils.Busybox("cp", "-avrf", d, ds)
		if err != nil {
			menu.Text("Ошибка при копировании:", d, err)
		} else {
			permAll(ds)
		}
	}
}

func deleteApps(name string) {
	appDir := filepath.Join(sar_dir, "apps", name)
	if utils.IsEmptyDir(appDir) {
		menu.Text("Папка", appDir, "пустая")
		return
	}
	dirs := findDirs([]string{filepath.Join(appDir, "system/app"), filepath.Join(appDir, "system/priv-app")})
	for _, d := range dirs {
		ds := strings.TrimPrefix(d, appDir)
		menu.Text("Удаление:", ds)
		err := os.RemoveAll(ds)
		if err != nil {
			menu.Text("Ошибка при удалении:", ds, err)
		}
	}
}
