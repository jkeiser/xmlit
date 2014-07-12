package xmlchan

import "encoding/xml"

type Decoder struct {
	decoder *xml.Decoder
	results chan<- Result
	done    <-chan struct{}
}

func (decoder Decoder) Next() (token xml.Token, err error) {
	return decoder.decoder.Token()
}

func (decoder Decoder) Skip() error {
	return decoder.decoder.Skip()
}

func (decoder Decoder) DecodeElement(v interface{}, start *xml.StartElement) error {
	return decoder.decoder.DecodeElement(v, start)
}

func (decoder Decoder) Done() bool {
	select {
	case _, _ = <-decoder.done:
		// When the done channel closes, this case will trigger immediately
		return true
	default:
		return false
	}
}

func (decoder Decoder) Send(result Result) {
	select {
	case decoder.results <- result:
	case _, _ = <-decoder.done:
		// If we are told we're done then we silently drop the send.  Done() will
		// eventually catch us.
	}
}

func (decoder Decoder) SendResult(result interface{}) {
	decoder.Send(NewResult(result))
}

func (decoder Decoder) SendError(err error) {
	decoder.Send(NewError(err))
}
