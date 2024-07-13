package errorstate

func New(err error) *state {
	return &state{
		error:  err,
		keyMap: newKeyMap(),
	}
}
