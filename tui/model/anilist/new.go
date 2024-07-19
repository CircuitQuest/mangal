package anilist

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/tui/model/help"
)

func New(anilist *anilist.Anilist, standalone bool) *Model {
	_styles := defaultStyles()
	input := textinput.New()
	input.Width = 60
	input.PromptStyle = _styles.prompt
	input.TextStyle = _styles.text
	idInput := input
	secretInput := input
	codeInput := input

	idInput.Prompt = "ID:     "
	secretInput.Prompt = "Secret: "
	codeInput.Prompt = "Code:   "

	m := &Model{
		idInput:              idInput,
		secretInput:          secretInput,
		codeInput:            codeInput,
		help:                 help.New(),
		anilist:              anilist,
		standalone:           standalone,
		notificationDuration: 2 * time.Second,
		title:                _styles.title.Render("Anilist"),
		selectCursor:         _styles.prompt.Render(icon.Item.Raw()),
		state:                Uninitialized,
		current:              ID,
		styles:               _styles,
		keyMap:               newKeyMap(),
	}
	m.updateCurrent()
	return m
}
