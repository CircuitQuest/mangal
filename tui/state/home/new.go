package home

func New() *State {
	return &State{
		keyMap: newKeyMap(),
	}
}
