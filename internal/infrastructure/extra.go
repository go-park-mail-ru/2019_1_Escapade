package infrastructure

// WithExtraI defines groups of objects
//  where you can set a callback action
type WithExtraI interface {
	Extra() error
}

// WithExtra action as callback, set by user
type WithExtra struct {
	CallExtra func() error
}

// Extra execute extra action
func (extra *WithExtra) Extra() error {
	if extra.CallExtra == nil {
		return nil
	}
	return extra.CallExtra()
}
