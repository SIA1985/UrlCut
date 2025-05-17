package server

import (
	"UrlCut/internal/logic"
	"UrlCut/internal/webbrowser"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

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

		select {
		case <-srvCtx.Done():
			return
		default:

		}

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

			webbrowser.Open(fullUrl)

			continue
		}

		log.Println("Неизвестная комманда: '" + cmdTxt + "'")

	}
}
