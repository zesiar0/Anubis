package handler

import (
	"Anubis/pkg/common"
	"Anubis/pkg/logger"
	"Anubis/pkg/model"
	"Anubis/pkg/proxy"
	"Anubis/pkg/utils"
	"fmt"
	"strconv"
	"strings"
)

func (h *InteractiveHandler) displayAssets(searchHeader string) {
	idLabel := "ID"
	hostLabel := "Hostname"
	ipLabel := "IP"
	protocolsLabel := "Protocols"
	platformLabel := "Platform"
	commentLabel := "Comment"

	labels := []string{idLabel, hostLabel, ipLabel, protocolsLabel, platformLabel, commentLabel}
	fields := []string{"ID", "Hostname", "IP", "Protocols", "Platform", "Comment"}
	fieldsSize := map[string][3]int{
		"ID":       {0, 0, 5},
		"Hostname": {0, 40, 0},
		"IP":       {0, 8, 40},
		"Protocol": {0, 8, 0},
		"Platform": {0, 8, 0},
		"Comment":  {0, 0, 0},
	}

	generateRowFunc := func(i int, item *model.Asset) map[string]string {
		row := make(map[string]string)
		row["ID"] = strconv.Itoa(i + 1)
		row["Hostname"] = item.Name
		row["IP"] = item.Address
		row["Protocol"] = item.Protocol
		row["Platform"] = item.Platform
		row["Comment"] = joinMultiLineString(item.Comment)

		return row
	}

	assetDisplay := "the asset"
	data := make([]map[string]string, len(h.assets))
	for i := range h.assets {
		data[i] = generateRowFunc(i, &h.assets[i])
	}

	h.displayResult(searchHeader, assetDisplay, labels, fields, fieldsSize, generateRowFunc)

}

func joinMultiLineString(lines string) string {
	lines = strings.ReplaceAll(lines, "\r", "\n")
	lines = strings.ReplaceAll(lines, "\n\n", "\n")
	lineArray := strings.Split(strings.TrimSpace(lines), "\n")
	lineSlice := make([]string, 0, len(lineArray))
	for _, item := range lineArray {
		cleanLine := strings.TrimSpace(item)
		if cleanLine == "" {
			continue
		}
		lineSlice = append(lineSlice, strings.ReplaceAll(cleanLine, " ", ","))
	}
	return strings.Join(lineSlice, "|")
}

type createRowFunc func(i int, item *model.Asset) map[string]string

func (h *InteractiveHandler) displayResult(searchHeader, assetDisplay string,
	labels, fields []string, fieldSize map[string][3]int,
	generateRowFunc createRowFunc) {
	term := h.term
	data := make([]map[string]string, len(h.assets))
	for i := range h.assets {
		data[i] = generateRowFunc(i, &h.assets[i])
	}

	w, _ := term.GetSize()
	caption := utils.WrapperString("", utils.Green)
	table := common.WrapperTable{
		Fields:      fields,
		Labels:      labels,
		FieldsSize:  fieldSize,
		Data:        data,
		TotalSize:   w,
		Caption:     caption,
		TruncPolicy: common.TruncMiddle,
	}
	table.Initial()
	loginTip := "Enter ID number directly login %s, multiple search use // + field, such as: //16"
	loginTip = fmt.Sprintf(loginTip, assetDisplay)
	pageActionTip := "Page up: b Page down: n"
	actionTip := fmt.Sprintf("%s %s", loginTip, pageActionTip)
	_, _ = term.Write([]byte(utils.CharClear))
	_, _ = term.Write([]byte(table.Display()))
	utils.IgnoreErrWriteString(term, utils.WrapperString(actionTip, utils.Green))
	utils.IgnoreErrWriteString(term, utils.CharNewLine)
	utils.IgnoreErrWriteString(term, utils.WrapperString(searchHeader, utils.Green))
	utils.IgnoreErrWriteString(term, utils.CharNewLine)
}

func (h *InteractiveHandler) proxyAsset(asset int) {
	account := &model.Account{
		Name:       "a3bz-ssh",
		Username:   "root",
		Secret:     "zengjiahua..123",
		SecretType: "pass",
	}

	proxyOpts := make([]proxy.ConnectionOption, 0, 3)
	proxyOpts = append(proxyOpts, proxy.ConnectProtocol("ssh"))
	proxyOpts = append(proxyOpts, proxy.ConnectUser("root"))
	proxyOpts = append(proxyOpts, proxy.ConnectAccount(account))
	srv, err := proxy.NewServer(h.sess, proxyOpts...)
	if err != nil {
		logger.Errorf("create proxy server err: %s", err)
		return
	}
	srv.Proxy()
}
