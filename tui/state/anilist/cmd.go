package anilist

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
)

func (s *state) setMetadataCmd(manga anilist.Manga) tea.Cmd {
	return func() tea.Msg {
		s.manga.SetMetadata(manga.Metadata())

		msg := fmt.Sprintf("Set Anilist %q (%d)", manga, manga.ID)
		log.Log(msg+" to manga %q", s.manga)
		return base.NotifyWithDuration(msg, 3*time.Second)()
	}
}

func (s *state) searchCmd(ctx context.Context, query string) tea.Cmd {
	return tea.Sequence(
		base.Loading(fmt.Sprintf("Searching %q on Anilist", query)),
		func() tea.Msg {
			var mangas []anilist.Manga

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
