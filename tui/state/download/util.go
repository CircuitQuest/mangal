package download

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/luevano/mangal/theme/icon"
	"github.com/luevano/mangal/theme/style"
	stringutil "github.com/luevano/mangal/util/string"
)

// updateKeyMap enables open/retry keybinds if necessary
func (s *state) updateKeyMap() {
	s.keyMap.open.SetEnabled(len(s.succeed) > 0 && !s.downloading)
	s.keyMap.retry.SetEnabled(len(s.failed) > 0 && !s.downloading)
}

func (s *state) viewDownloading() string {
	spinnerView := s.spinner.View()
	return fmt.Sprintf(`%s Downloading %s - %d/%d

%s

%s %s`,
		icon.Progress,
		s.toDownload[s.currentIdx].String(),
		s.currentIdx+1, len(s.toDownload),
		s.progress.ViewAs(float64(s.currentIdx)/float64(len(s.toDownload))),
		spinnerView,
		style.Normal.Secondary.Render(stringutil.Trim(s.message, s.size.Width-lipgloss.Width(spinnerView)-1)),
	)
}

func (s *state) viewDownloaded() string {
	var (
		succeed = len(s.succeed)
		failed  = len(s.failed)
	)

	if failed == 0 {
		return style.Normal.Success.
			Render(fmt.Sprintf(
				"%s downloaded successfully!",
				stringutil.Quantify(succeed, "chapter", "chapters"),
			))
	}

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(
		"%s downloaded successfully, %d failed.",
		stringutil.Quantify(succeed, "chapter", "chapters"),
		failed,
	))

	sb.WriteString("\n\nFailed:\n")

	if failed <= 3 {
		for _, chapter := range s.failed {
			sb.WriteString(fmt.Sprintf("\n%s", chapter))
		}
	} else {
		indices := make([]float32, failed)
		for i, c := range s.failed {
			indices[i] = c.Info().Number
		}

		sb.WriteString(stringutil.FormatRanges(indices))
	}

	return sb.String()
}
