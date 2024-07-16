package download

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	stringutil "github.com/luevano/mangal/util/string"
)

// updateKeybinds enables open/retry keybinds if necessary
func (s *state) updateKeybinds() {
	_, succ, fail := s.chapters.getEach()
	hasSucc := len(succ) > 0
	hasFail := len(fail) > 0
	downloading := s.downloading != dSDownloaded
	s.keyMap.open.SetEnabled(hasSucc && !downloading)
	s.keyMap.retry.SetEnabled(hasFail && !downloading)
}

func (s *state) viewDownloading() string {
	ch := s.chapters[s.currentIdx].chapter
	summary := fmt.Sprintf(
		`%s Downloading chapter %s "%s" %s %d/%d`,
		icon.Progress.Colored(),
		s.styles.accent.Render(stringutil.FormatFloa32(ch.Info().Number)),
		s.styles.accent.Render(ch.String()),
		s.sep, s.currentIdx+1, len(s.toDownload),
	)
	return lipgloss.JoinVertical(
		lipgloss.Left,
		summary,
		" ",
		s.progress.ViewAs(float64(s.currentIdx)/float64(len(s.toDownload))),
		" ",
		s.spinner.View()+" "+style.Normal.Secondary.Render(s.message),
	)
}

func (s *state) viewDownloaded() string {
	down, succ, fail := s.chapters.getEach()
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

func (s *state) viewChapterList(header string, headerStyle lipgloss.Style, chaps chapters) string {
	var sb strings.Builder
	sb.Grow(200)

	sb.WriteString(headerStyle.Render(header))
	sb.WriteByte('\n')

	item := s.styles.itemEnum.Foreground(headerStyle.GetForeground()).Render(icon.Item.Raw())
	subItm := s.styles.subItemEnum.Foreground(headerStyle.GetForeground()).Render(icon.SubItem.Raw())
	var lastDir string
	for i, ch := range chaps {
		if ch.down != nil && ch.down.Directory != lastDir {
			lastDir = ch.down.Directory
			sb.WriteString(item)
			sb.WriteString(s.styles.toDownload.Render(lastDir))
			sb.WriteByte('\n')
		}
		sb.WriteString(subItm)
		sb.WriteString(ch.render(s.styles))
		if i < len(chaps) {
			sb.WriteByte('\n')
		}
	}

	return sb.String()
}
