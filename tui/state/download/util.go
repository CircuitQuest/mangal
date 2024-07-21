package download

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	"github.com/luevano/mangal/util/chapter"
	stringutil "github.com/luevano/mangal/util/string"
)

// updateKeybinds enables open/retry keybinds if necessary
func (s *state) updateKeybinds() {
	_, succ, fail := s.chapters.GetEach()
	hasSucc := len(succ) > 0
	hasFail := len(fail) > 0
	downloading := s.downloading != dSDownloaded
	s.keyMap.open.SetEnabled(hasSucc && !downloading)
	s.keyMap.retry.SetEnabled(hasFail && !downloading)
}

func (s *state) viewDownloading() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		s.viewSummary(),
		" ",
		s.progress.ViewAs(float64(s.currentIdx)/float64(len(s.toDownload))),
		" ",
		s.spinner.View()+" "+style.Normal.Secondary.Render(s.message),
	)
}

func (s *state) viewSummary() string {
	ch := s.chapters[s.currentIdx].Chapter
	if s.retrying {
		s.message = fmt.Sprintf("429 Too Many Requests (retry #%d)", s.retryCount)
		return fmt.Sprintf(
			`%s Retrying chapter %s "%s" in %s%s%d/%d`,
			icon.Timer.Colored(),
			s.styles.accent.Render(stringutil.FormatFloa32(ch.Info().Number)),
			s.styles.accent.Render(ch.String()),
			s.timer.View(),
			s.sep, s.currentIdx+1, len(s.toDownload),
		)
	}
	return fmt.Sprintf(
		`%s Downloading chapter %s "%s"%s%d/%d`,
		icon.Progress.Colored(),
		s.styles.accent.Render(stringutil.FormatFloa32(ch.Info().Number)),
		s.styles.accent.Render(ch.String()),
		s.sep, s.currentIdx+1, len(s.toDownload),
	)
}

func (s *state) viewDownloaded() string {
	down, succ, fail := s.chapters.GetEach()
	var lines []string
	if len(down) != 0 {
		lines = append(lines, s.viewChapterList("To download:", s.styles.toDownload, down))
	}
	if len(succ) != 0 {
		lines = append(lines, s.viewChapterList("Succeeded:", s.styles.succeed, succ))
	}
	if len(fail) != 0 {
		lines = append(lines, s.viewChapterList("Failed:", s.styles.failed, fail))
	}
	expl := "\nThe lists follow the template:\n" +
		"`<number> - <title> - s: <download status> - f: <filename>`"
	lines = append(lines, s.styles.toDownload.Render(expl))
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func (s *state) viewChapterList(header string, headerStyle lipgloss.Style, chaps chapter.Chapters) string {
	var sb strings.Builder
	sb.Grow(200)

	sb.WriteString(headerStyle.Render(header))
	sb.WriteByte('\n')

	item := s.styles.itemEnum.Foreground(headerStyle.GetForeground()).Render(icon.Item.Raw())
	subItm := s.styles.subItemEnum.Foreground(headerStyle.GetForeground()).Render(icon.SubItem.Raw())
	var lastDir string
	for i, ch := range chaps {
		if ch.Down != nil && ch.Down.Directory != lastDir {
			lastDir = ch.Down.Directory
			sb.WriteString(item)
			sb.WriteString(s.styles.toDownload.Render(lastDir))
			sb.WriteByte('\n')
		}
		sb.WriteString(subItm)
		sb.WriteString(ch.Render(style.Normal.Base, s.styles.toDownload, s.styles.succeed, s.styles.failed, s.styles.warning))
		if i < len(chaps) {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
