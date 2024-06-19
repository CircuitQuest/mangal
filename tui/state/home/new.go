package home

func New() *state {
	return &state{
		keyMap: newKeyMap(),
	}
}
