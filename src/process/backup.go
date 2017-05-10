package process

import (
	"io/ioutil"
	"menu"
	"os"
	"path/filepath"
	"time"
	"utils"
)

func Backup() {
	menu.Clear()
	bdir := getBackupDir()
	if bdir == "" {
		return
	}
	menu.Text("Сделать резервную копию всех системных apk, это может занять продолжительное время?\n" + bdir)
	if menu.Question("") {
		if bdir != "" {
			menu.Text("Backup app:")
			err := utils.Busybox("tar", "-zcvf", bdir+"_backup.tar.gz", "/system/app", "/system/priv-app")
			if err != nil {
				menu.Text("***", "Error backup apps", "\""+err.Error()+"\"")
			}
		}
	}
}

func Restore() {
	menu.Clear()
	dirs, err := ioutil.ReadDir(filepath.Join(sar_dir, "backup"))
	if err != nil {
		menu.Text("Ошибка чтения backup директории:", err)
		return
	}
	list := make([]string, 0)
	for _, d := range dirs {
		list = append(list, d.Name())
	}
	itm := menu.List("Восстановить из резервной копии:", "", list)
	if itm > 0 && itm <= len(list) && menu.Question("Осторожно, можно убить систему!!! Продолжить?") {
		bfile := filepath.Join(sar_dir, "backup", list[itm-1])
		menu.Text("Восстановление из:", bfile)
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

		err = utils.Busybox("tar", "-xzvf", bfile)
		if err != nil {
			menu.Text("***", "Error restore apps", "\""+err.Error()+"\"")
		}
	}
}

func getBackupDir() string {
	if !utils.Exists(sar_dir) {
		menu.Text("Создание SAR директории")
		err := os.MkdirAll(sar_dir, 0771)
		if err != nil {
			menu.Text("Ошибка при создании директории ", sar_dir, err)
			return ""
		}
	}

	backup_dir := filepath.Join(sar_dir, "backup")

	if !utils.Exists(backup_dir) {
		menu.Text("Create backup dir")
		err := os.MkdirAll(backup_dir, 0777)
		if err != nil {
			menu.Text("Ошибка при создании директории ", backup_dir, err)
			return ""
		}
	}

	datedir := time.Now().Local().Format("2006-01-02")
	return utils.GetPathCount(filepath.Join(backup_dir, datedir))
}
