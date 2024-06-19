package chapsdownloaded

func New(options Options) *state {
	state := &state{
		options: options,
	}
	state.keyMap = newKeyMap(state)
	return state
}
