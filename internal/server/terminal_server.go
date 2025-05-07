package server

import (
	"UrlCut/internal/logic"
	"bufio"
	"fmt"
	"log"
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
	cutCmd := regexp.MustCompile(`cut\s\s*[A-Za-z0-9:\/.]*`)

	redirectCmd := regexp.MustCompile(`redirect\s\s*[A-Za-z0-9]*`)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		if cutCmd.MatchString(sc.Text()) {

			continue
		}

		if redirectCmd.MatchString(sc.Text()) {

			continue
		}

		log.Println("Неизвестная комманда: " + sc.Text())

	}
}
