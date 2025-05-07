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
	var err error

	cutCmd := regexp.MustCompile(`cut\s\s*[A-Za-z0-9:\/.]*`)

	redirectCmd := regexp.MustCompile(`redirect\s\s*[A-Za-z0-9]*`)

	url := regexp.MustCompile(`(.*\s\s*)([A-Za-z0-9:\/.]*)`)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		cmdTxt := sc.Text()
		if cutCmd.MatchString(cmdTxt) {
			fullUrl := url.FindStringSubmatch(cmdTxt)[2]

			var cutUrl string
			cutUrl, err = h.logic.CutUrl(fullUrl)
			if err != nil {
				log.Println(err)
				continue
			}

			fmt.Println(cutUrl)

			continue
		}

		if redirectCmd.MatchString(cmdTxt) {
			cutUrl := url.FindStringSubmatch(cmdTxt)[2]

			var fullUrl string
			fullUrl, err = h.logic.GetFullUrl(cutUrl)
			if err != nil {
				log.Println(err)
				continue
			}

			err = openBrowser(fullUrl)
			if err != nil {
				log.Println(err)
				continue
			}

			continue
		}

		log.Println("Неизвестная комманда: " + cmdTxt)

	}
}
