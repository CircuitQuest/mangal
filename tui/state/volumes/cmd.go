package volumes

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/tui/state/chapters"
)

func (s *state) searchVolumeChapters(ctx context.Context, item *item) tea.Cmd {
	return tea.Sequence(
		base.Loading(fmt.Sprintf("Searching chapters for volume %s", item.volume)),
		func() tea.Msg {
			chapterList, err := s.client.VolumeChapters(ctx, item.volume)
			if err != nil {
				return err
			}

			return chapters.New(s.client, s.manga, item.volume, chapterList)
		},
		base.Loaded,
	)
}
