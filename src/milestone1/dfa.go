package milestone1

type TransitionKey struct {
	State string
	Input string
}

type DFA struct {
	StartState string
	FinalState []string
	Transition map[TransitionKey]string
}
