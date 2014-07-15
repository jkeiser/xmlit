package xmlit

import "io"
import "encoding/xml"
import "reflect"

// The CreatorMap lets you specify a map of element names to value factories: the
// value factories return pointers to new objects which will be passed to
// xml.Decoder.DecodeElement() when <TypeName> is seen in the XML.  Generally you
// will use new() to get the pointer.
//
// The input stream will be closed after iteration.
//
//   iterator := xmlit.DecodeElementsByXmlName(reader, CreatorMap{
//     "Customer":        func() interface{} { new(Customer) },
//     "Order":           func() interface{} { new(Order) },
//     "DelinquentOrder": func() interface{} { new(DelinquentOrder) },
//     "Account":         func() interface{} { new(Account) },
//   })
type CreatorMap map[xml.Name]func() interface{}

type ElementDecoder struct {
	Decoder  *xml.Decoder
	Creators CreatorMap
	next     interface{}
	err      error
}

func (decoder *ElementDecoder) HasNext() bool {
	if decoder.err != nil {
		return false
	}
	if decoder.next != nil {
		return true
	}
	decoder.next, decoder.err = decodeNextElement(decoder.Decoder, decoder.Creators)
	return decoder.err == nil
}

func (decoder *ElementDecoder) Next() interface{} {
	decoder.HasNext()
	result := decoder.next
	decoder.next = nil
	return result
}

func (decoder *ElementDecoder) Error() error {
	if decoder.err == io.EOF {
		return nil
	}
	return decoder.err
}

// Returns an iterator that uses xml.Decoder to decode a stream of the given types.
// Will ignore any containing tags, skipping straight to the actual elements in
// question.
//
// The input is a list of "exemplars."  DecodeElements will analyze their types
// and use the type names to decide what XML element tags to read.  A new value
// of the type will be created for each decoded element.
//
// The input stream will be closed after iteration.
//
//   iterator := xmlit.DecodeElements(reader, Customer{}, Order{}, Account{})
//   iterator.Each(func (item interface{}) { fmt.Println(item) })
func DecodeElements(reader io.ReadCloser, exemplars ...interface{}) *ElementDecoder {
	creators := CreatorMap{}
	for _, exemplar := range exemplars {
		exemplarType := reflect.TypeOf(exemplar)
		exemplarName := xml.Name{Local: exemplarType.Name()}
		creators[exemplarName] = func() interface{} { return reflect.New(exemplarType).Interface() }
	}
	result := ElementDecoder{
		Decoder:  xml.NewDecoder(reader),
		Creators: creators,
	}
	return &result
}

func decodeNextElement(decoder *xml.Decoder, creators CreatorMap) (interface{}, error) {
	for {
		token, err := decoder.Token()
		if err != nil {
			return nil, err
		}

		start, ok := token.(xml.StartElement)
		if !ok {
			continue
		}

		creator, ok := creators[start.Name]
		if !ok {
			continue
		}

		result := creator()
		err = decoder.DecodeElement(result, &start)
		if err != nil {
			return nil, err
		}

		// This craziness is because DecodeElement wants a pointer-to-thing-as-interface,
		// and reflection is the only way to address or dereference such a beast.
		return reflect.ValueOf(result).Elem().Interface(), nil
	}
}
