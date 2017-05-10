package process

import (
	"io/ioutil"
	"menu"
	"strings"
)

const buildPropPath = "/system/build.prop"

var (
	menuLocale = []string{
		"ru-RU",
		"en-US",
		"zh-CN",
		"Своё значение",
	}
)

func ChangeLocal() {
	menu.Clear()
	buf, err := ioutil.ReadFile(buildPropPath)
	if err != nil {
		menu.Text("Error read build.prop:", err)
		menu.Read()
		return
	}

	text := string(buf)
	locale := "не найдено"
	pos := strings.Index(text, "ro.product.locale=")

	if pos >= 0 {
		pos += 18
		end := strings.Index(text[pos:], "\n")
		locale = text[pos : pos+end]
	}
	menu.Text("Текущая локаль:", locale)
	sel := menu.List("Выберите:", "0 для выхода", menuLocale)
	switch sel {
	case 1:
		{ //ru-RU
			locale = "ru-RU"
		}
	case 2:
		{ //en-US
			locale = "en-US"
		}
	case 3:
		{ //zh-CN
			locale = "zh-CN"
		}
	case 4:
		{ //edit
			locale = menu.Edit("Введите локаль:")
		}
	default:
		return
	}

	if menu.Question("Записать \"" + locale + "\" в build.prop?") {
		lines := strings.Split(string(text), "\n")
		if pos != -1 {
			for i, line := range lines {
				if strings.Contains(line, "ro.product.locale=") {
					lines[i] = "ro.product.locale=" + locale
					break
				}
			}
		} else {
			lines = append(lines, "ro.product.locale="+locale)
		}
		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(buildPropPath, []byte(output), 0644)
		if err != nil {
			menu.Text("Ошибка при записи:", err, "\nEnter чтобы продолжить")
			menu.Read()
			return
		}
		menu.Text("Локаль изменена")
	}
}
