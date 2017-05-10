package process

import (
	"io/ioutil"
	"menu"
	"os"
	"path/filepath"
	"utils"
)

func Replace() {
	apps := make([]string, 0)
	list, err := ioutil.ReadDir(filepath.Join(sar_dir, "replace"))

	if err != nil {
		menu.Text("Ошибка чтения директории:", err)
		return
	}
	for _, l := range list {
		apps = append(apps, l.Name())
	}
	if len(apps) == 0 {
		menu.Text("Нет файлов для замены")
		return
	}

	itm := menu.List("Выберите версию", "", apps)
	if itm > 0 && itm <= len(apps) {
		appFile := filepath.Join(sar_dir, "replace", apps[itm-1])

		tmp_app := utils.GetPathCount("/system/app_old")
		tmp_priv_app := utils.GetPathCount("/system/priv-app_old")
		var err error
		defer func() {
			if err != nil {
				menu.Text("Восстановление после неудачного востанновления")
				utils.Busybox("cp", "-avrf", tmp_app, "/system/app")
				utils.Busybox("cp", "-avrf", tmp_priv_app, "/system/priv-app")
			}
			os.RemoveAll(tmp_app)
			os.RemoveAll(tmp_priv_app)
		}()

		if utils.Exists("/system/app") {
			menu.Text("/system/app ->", tmp_app)
			err = utils.Busybox("mv", "/system/app", tmp_app)
			if err != nil {
				menu.Text("***", "Error rename /system/app", "\""+err.Error()+"\"")
				return
			}
		}

		if utils.Exists("/system/priv-app") {
			menu.Text("/system/priv-app ->", tmp_priv_app)
			err = utils.Busybox("mv", "/system/priv-app", tmp_priv_app)
			if err != nil {
				menu.Text("***", "Error rename /system/priv-app", "\""+err.Error()+"\"")
				return
			}
		}

		err = utils.Busybox("tar", "-xzvf", appFile)
		if err != nil {
			menu.Text("***", "Ошибка восстановления:", "\""+err.Error()+"\"")
		} else {
			permAll("/system/app")
			permAll("/system/priv-app")
		}
	}
}
