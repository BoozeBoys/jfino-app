package testutils

type StatusReaderMock struct {
	status [][]byte
	err    error
}

func (sr *StatusReaderMock) Status() ([][]byte, error) {
	return sr.status, sr.err
}

func (sr *StatusReaderMock) SetStatus(status [][]byte) {
	sr.status = status
}

func (sr *StatusReaderMock) SetError(err error) {
	sr.err = err
}
