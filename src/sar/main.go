package main

import (
	"config"
	"fmt"
	"menu"
	"process"
	utils2 "process/utils"
	"utils"
)

var (
	mainmenu = []string{
		"Сделать резевную копию системных apk",
		"Восстановить системные apk из резервной копии",
		"Изменить локаль в build.prop",
		"Установить/удалить дополнительные программы",
		"Заменить системные apk",
		"Перезапуск меню",
	}
)

const (
	backup = iota + 1
	restore
	buildprop
	install_app
	replace
	reboot
)

func main() {
	fmt.Println("Version 1.0.1")

	db, err := utils2.OpenDB("/home/yourok/tmp/db.dbolt")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.AddDir("/home/yourok/tmp/video/config/")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(db.GetNames("/home/yourok/tmp/video/config/"))

	db.PrintAll(nil, 0)
	return

	utils.CheckSu()
	config.Set("busybox", utils.CheckBusyBox())
	utils.Remount()
	exit := false

	for !exit {
		menu.Clear()
		sel := menu.List("Выберите:", "", mainmenu)
		switch sel {
		case buildprop:
			process.ChangeLocal()
		case backup:
			process.Backup()
		case restore:
			process.Restore()
		case install_app:
			process.InstallApps()
		case replace:
			process.Replace()
		case reboot:
			process.Reboot()
		default:
			exit = true
		}
	}
}
