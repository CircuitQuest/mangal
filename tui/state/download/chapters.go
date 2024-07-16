package download

import (
	"github.com/luevano/libmangal/mangadata"
	"github.com/luevano/libmangal/metadata"
	"github.com/luevano/mangal/theme/style"
	stringutil "github.com/luevano/mangal/util/string"
)

var renderChapterFullInfo bool

type chapter struct {
	chapter mangadata.Chapter
	state   chapterState
	down    *metadata.DownloadedChapter
}

func (c chapter) render(styles styles) string {
	str := stringutil.FormatFloa32(c.chapter.Info().Number) + " - " + c.chapter.Info().Title
	if c.down == nil {
		return styles.toDownload.Render(str)
	}

	s := styles.toDownload
	status := "s: " + string(c.down.ChapterStatus)
	switch c.down.ChapterStatus {
	case metadata.DownloadStatusNew:
		s = style.Normal.Base
		status = styles.succeed.Render(status)
	case metadata.DownloadStatusMissingMetadata, metadata.DownloadStatusOverwritten, metadata.DownloadStatusSkip:
		status = styles.warning.Render(status)
	case metadata.DownloadStatusFailed:
		status = styles.failed.Render(status)
	default: // exists or other
		status = styles.toDownload.Render(status)
	}
	str = s.Render(str+" - ") + status + s.Render(" - f: "+c.down.Filename)
	return str
}

type chapters []*chapter

func (c chapters) getEach() (d, s, f chapters) {
	for _, ch := range c {
		switch ch.state {
		case cSToDownload:
			d = append(d, ch)
		case cSSucceed:
			s = append(s, ch)
		case cSFailed:
			f = append(f, ch)
		default:
			panic("unexpected chapter download state")
		}
	}
	return d, s, f
}

func (c chapters) toDownload() (d chapters) {
	for _, ch := range c {
		if ch.state == cSToDownload {
			d = append(d, ch)
		}
	}
	return d
}

func (c chapters) succeed() (s chapters) {
	for _, ch := range c {
		if ch.state == cSSucceed {
			s = append(s, ch)
		}
	}
	return s
}

func (c chapters) failed() (f chapters) {
	for _, ch := range c {
		if ch.state == cSFailed {
			f = append(f, ch)
		}
	}
	return f
}
