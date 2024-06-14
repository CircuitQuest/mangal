package chapsdownloaded

func New(options Options) *State {
	state := &State{
		options: options,
	}
	state.keyMap = newKeyMap(state)
	return state
}
