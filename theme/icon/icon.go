package icon

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/color"
	"github.com/luevano/mangal/theme/style"
)

type icon struct {
	color   lipgloss.TerminalColor
	symbols symbols
}

type symbols map[Type]string

func (i icon) String() string {
	return style.Bold.Base.Foreground(i.color).Render(i.symbols[currentType])
}

var (
	Confirm = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: "?",
			TypeNerd:  "\uEB32",
		},
	}

	Progress = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: "@",
			TypeNerd:  "\U000F0997",
		},
	}

	Mark = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: "*",
			TypeNerd:  "\U000F0F22",
		},
	}

	Download = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: "#",
			TypeNerd:  "\uF019",
		},
	}

	Check = icon{
		color: color.Success,
		symbols: symbols{
			TypeASCII: "~",
			TypeNerd:  "\uF05D",
		},
	}

	Cross = icon{
		color: color.Error,
		symbols: symbols{
			TypeASCII: "x",
			TypeNerd:  "\uF05C",
		},
	}

	Search = icon{
		color: color.Accent,
		symbols: symbols{
			TypeASCII: ">",
			TypeNerd:  "\uF002",
		},
	}

	Recent = icon{
		color: color.Secondary,
		symbols: symbols{
			TypeASCII: "~",
			TypeNerd:  "\uF017",
		},
	}

	Read = icon{
		color: color.Secondary,
		symbols: symbols{
			TypeASCII: "r",
			TypeNerd:  "\U000F0447",
		},
	}

	Available = icon{
		color: color.Secondary,
		symbols: symbols{
			TypeASCII: "a",
			TypeNerd:  "\uEB28",
		},
	}
)
