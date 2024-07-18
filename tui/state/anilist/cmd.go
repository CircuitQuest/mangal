package anilist

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/luevano/libmangal/metadata/anilist"
	"github.com/luevano/mangal/log"
	"github.com/luevano/mangal/tui/base"
	"github.com/luevano/mangal/util/cache"
)

func (s *state) getHistoryCmd() tea.Msg {
	found, err := cache.GetAnilistSearchHistory(&s.history)
	if err != nil {
		return err
	}
	if found {
		s.search.SetSuggestions(s.history.Get())
	}
	return nil
}

func (s *state) updateHistoryCmd(query string) tea.Cmd {
	return func() tea.Msg {
		s.history.Add(query)
		s.history.Sort()
		s.search.SetSuggestions(s.history.Get())
		return cache.SetAnilistSearchHistory(s.history)
	}
}

func (s *state) setMetadataCmd(manga anilist.Manga) tea.Cmd {
	return func() tea.Msg {
		s.manga.SetMetadata(&manga)

		msg := fmt.Sprintf("Set Anilist %q", manga.String())
		log.Log(msg+" to manga %q", s.manga)
		return base.Notify(msg)()
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
				if manga.ID() == closest.ID() {
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
			s.updateKeybinds()
			return nil
		},
		base.Loaded,
	)
}
