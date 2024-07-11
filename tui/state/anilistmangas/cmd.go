package anilistmangas

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lmanilist "github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/tui/base"
)

func (s *state) searchCmd(ctx context.Context, query string) tea.Cmd {
	return tea.Sequence(
		base.Loading(fmt.Sprintf("Searching %q on Anilist", query)),
		func() tea.Msg {
			var mangas []lmanilist.Manga

			// to keep the closest on top
			closest, found, err := s.anilist.FindClosestManga(ctx, query)
			if err != nil {
				return err
			}
			if found {
				mangas = append(mangas, closest)
			}

			// the rest of the results
			mangaSearchResults, err := s.anilist.SearchMangas(ctx, query)
			if err != nil {
				return nil
			}
			for _, manga := range mangaSearchResults {
				// except the closest
				if manga.ID == closest.ID {
					continue
				}
				mangas = append(mangas, manga)
			}

			items := make([]list.Item, len(mangas))
			for i, m := range mangas {
				items[i] = &item{manga: m}
			}
			s.list.SetItems(items)

			s.searched = true
			return nil
		},
		base.Loaded,
	)
}
