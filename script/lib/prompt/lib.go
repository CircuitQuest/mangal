package prompt

import (
	"github.com/luevano/mangal/script/lib/prompt/fzf"
	luadoc "github.com/mangalorg/gopher-luadoc"
)

const libName = "prompt"

func Lib() *luadoc.Lib {
	return &luadoc.Lib{
		Name:        libName,
		Description: "Various prompts for interracting with the user",
		Libs: []*luadoc.Lib{
			fzf.Lib(),
		},
	}
}
