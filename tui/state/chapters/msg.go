package chapters

import (
	"github.com/luevano/libmangal"
	"github.com/zyedidia/generic/set"
)

type readChapterMsg struct {
	path    string
	item    *item
	options libmangal.ReadOptions
}

type downloadChapterMsg struct {
	item      *item
	options   libmangal.DownloadOptions
	readAfter bool
}

type downloadChaptersMsg struct {
	items   set.Set[*item]
	options libmangal.DownloadOptions
}
