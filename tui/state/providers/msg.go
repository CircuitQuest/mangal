package providers

import "github.com/luevano/libmangal"

type loadProviderMsg struct {
	item *item
}

type searchMangasMsg struct {
	client *libmangal.Client
}
