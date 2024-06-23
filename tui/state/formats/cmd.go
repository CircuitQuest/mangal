package formats

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal"
	"github.com/luevano/mangal/config"
	"github.com/luevano/mangal/log"
)

func setFormatForCmd(what forWhat, format libmangal.Format) tea.Cmd {
	log.Log("setFormatFor %s %s", format, what)
	return func() tea.Msg {
		log.Log("actually setting format %s for %s", format, what)
		var err error
		switch what {
		case forRead:
			err = config.Read.Format.Set(format)
		case forDownload:
			err = config.Download.Format.Set(format)
		case forBoth:
			err = config.Read.Format.Set(format)
			if err != nil {
				return err
			}
			err = config.Download.Format.Set(format)
		}
		if err != nil {
			return err
		}
		log.Log("formatsUpdatedMsg")
		return formatsUpdatedMsg{}
	}
}

func writeConfigCmd() tea.Msg {
	log.Log("writeConfig")
	return config.Write()
}
