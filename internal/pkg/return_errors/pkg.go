package rerrors

func InterfaceIsNil() error {
	return New("Interfaces didnt set")
}
