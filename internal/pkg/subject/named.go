package subject

// Named wraps a Subject with its name.
type Named struct {
	// Name is the name of the subject.
	Name string

	// Subject embeds the subject itself.
	Subject
}
