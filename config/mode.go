package config

//go:generate enumer -type=Mode -trimprefix=Mode -json -yaml -text -transform=lower
type Mode uint8

const (
	ModeNone Mode = iota + 1
	ModeTUI
	ModeWeb
	ModeScript
	ModeInline
)
