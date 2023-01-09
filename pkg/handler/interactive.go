package handler

import (
	"Anubis/pkg/logger"
	"Anubis/pkg/utils"
	"github.com/gliderlabs/ssh"
	"strconv"
	"strings"
)

type InteractiveHandler struct {
	sess *WrapperSession
	term *utils.Terminal
}

func NewInteractiveHandler(sess ssh.Session) *InteractiveHandler {
	wrapperSess := NewWrapperSession(sess)
	term := utils.NewTerminal(wrapperSess, "Opt> ")
	handler := &InteractiveHandler{
		sess: wrapperSess,
		term: term,
	}

	return handler
}

type selectType int

const (
	TypeAsset = iota + 1
)

func (h *InteractiveHandler) Dispatch() {
	defer logger.Infof("Request %s: User stop interactive", h.sess.ID())
    h.displayHelp()
	for {
		line, err := h.term.ReadLine()
		if err != nil {
			logger.Fatalf("User close connect %s", err)
			break
		}

		line = strings.TrimSpace(line)
		switch len(line) {
		case 1:
			switch line {
			case "p":
				h.SetSelectType(TypeAsset)
				h.displayAssets("")
				continue
			case "h":
				h.displayHelp()
				continue
			case "q":
				logger.Infof("User enter %s to exit", line)
				return
			}
		}
		h.Proxy(line)
	}
}

func (h *InteractiveHandler) Proxy(line string) {
	if indexNum, err := strconv.Atoi(line); err == nil {
		if indexNum > 0 {
			h.proxyAsset(indexNum)
			return
		}
	}
}

func (h *InteractiveHandler) displayHelp() {
	h.term.SetPrompt("Opt> ")
	h.displayBanner()
}

func (h *InteractiveHandler) SetSelectType(s selectType) {
	switch s {
	case TypeAsset:
		h.term.SetPrompt("[Host]> ")
	}
}
