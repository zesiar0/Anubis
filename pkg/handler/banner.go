package handler

import (
	"Anubis/pkg/logger"
	"Anubis/pkg/utils"
	"fmt"
	"io"
	"text/template"
)

type MenuItem struct {
	id       int
	instruct string
	helpText string
}

type Menu []MenuItem

type ColorMeta struct {
	GreenBoldColor string
	ColorEnd       string
}

func (h *InteractiveHandler) displayBanner() {
	defaultTitle := utils.WrapperTitle("Welcome to anubis open source system")
	menu := Menu{
		{id: 1, instruct: "p", helpText: "display the host you have permission"},
		{id: 2, instruct: "h", helpText: "print help"},
		{id: 3, instruct: "q", helpText: "exit"},
	}

	title := defaultTitle
	prefix := utils.CharClear + utils.CharTab + utils.CharTab
	suffix := utils.CharNewLine + utils.CharNewLine
	welcomeMsg := prefix + "  " + title + suffix
	_, err := io.WriteString(h.sess, welcomeMsg)
	if err != nil {
		logger.Errorf("Send to client error, %s", err)
		return
	}

	cm := ColorMeta{GreenBoldColor: "\033[1;32m", ColorEnd: "\033[0m"}
	for _, v := range menu {
		line := fmt.Sprintf("\t%d) Enter {{.GreenBoldColor}}%s{{.ColorEnd}} to %s.%s",
			v.id, v.instruct, v.helpText, "\r\n")
		tmpl := template.Must(template.New("item").Parse(line))
		if err := tmpl.Execute(h.sess, cm); err != nil {
			logger.Error(err)
		}
	}
}
