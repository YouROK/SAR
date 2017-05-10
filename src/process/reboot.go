package process

import (
	"menu"
	"os"
	"utils"
)

func Reboot() {
	menu.Clear()
	itm := menu.List("", "", []string{"Перезапуск Dalvik", "Перезапуск", "Перезапуск в recovery"})
	if itm == 2 {
		utils.Run("pkill", "zygote")
		os.Exit(0)
	} else if itm == 2 {
		utils.Run("reboot")
		os.Exit(0)
	} else if itm == 3 {
		utils.Run("reboot", "recovery")
		os.Exit(0)
	}

}
