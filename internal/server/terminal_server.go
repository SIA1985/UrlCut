package server

import (
	"UrlCut/internal/logic"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

func openBrowser(fullUrl string) (err error) {
	switch os := runtime.GOOS; os {
	case "linux":
		err = exec.Command("x-www-browser", fullUrl).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", fullUrl).Start()
	case "darwin":
		err = exec.Command("open", fullUrl).Start()
	default:
		err = fmt.Errorf("Не поддерживаемая платформа!")
	}

	return
}

type Terminal struct {
	logic *logic.Logic
}

func (h *Terminal) Listen() {
	//todo: Context
	var err error

	var cutCmd *regexp.Regexp
	cutCmd, err = regexp.Compile(`cut [A-Za-z0-9:/.]*`)
	if err != nil {
		//todo: err
		return
	}

	var redirectCmd *regexp.Regexp
	redirectCmd, err = regexp.Compile(`redirect [A-Za-z0-9]*`)
	if err != nil {
		//todo: err
		return
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if cutCmd.MatchString(sc.Text()) {
			fmt.Println("cut")
			continue
		}

		if redirectCmd.MatchString(sc.Text()) {
			fmt.Println("redirect")
			continue
		}

		fmt.Println("Неизвестная комманда: " + sc.Text())

	}
}
