package confirm

func New(message string, onResponse OnResponseFunc) *State {
	return &State{
		message:    message,
		keyMap:     newKeyMap(),
		onResponse: onResponse,
	}
}
