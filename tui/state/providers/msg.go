package providers

import "github.com/luevano/libmangal"

type loadProviderMsg struct {
	item *Item
}

type searchMangasMsg struct {
	client *libmangal.Client
}
