package testutils

type HeloWaiterMock struct {
	Err error
}

func (h *HeloWaiterMock) WaitHelo() error {
	return h.Err
}
