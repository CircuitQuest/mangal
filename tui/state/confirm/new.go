package confirm

func New(message string, onResponse onResponseFunc) *state {
	return &state{
		message:    message,
		keyMap:     newKeyMap(),
		onResponse: onResponse,
	}
}
