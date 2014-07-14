package xmlit

import "io"
import "encoding/xml"
import "reflect"
import "github.com/jkeiser/iter"

// Used to specify the list of XML names and value factories used to parse the
// stream of XML.
type CreatorMap map[xml.Name]func() interface{}

// Returns an iterator that uses xml.Decoder to decode a strema of the given types.
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
func DecodeElements(reader io.ReadCloser, exemplars ...interface{}) iter.Iterator {
	return iterateElements(reader, decodeElementsFunc(exemplars))
}

// Returns an iterator that uses xml.Decoder to decode a strema of the given types.
// Will ignore any containing tags, skipping straight to the actual elements in
// question.
//
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
//   iterator.Each(func (item interface{}) { fmt.Println(item) })
func DecodeElementsByXmlName(reader io.ReadCloser, creators CreatorMap) iter.Iterator {
	return iterateElements(reader, decodeElementsByXmlNameFunc(creators))
}

func decodeElementsFunc(exemplars []interface{}) elementProcessor {
	creators := CreatorMap{}
	for _, exemplar := range exemplars {
		exemplarType := reflect.TypeOf(exemplar)
		exemplarName := xml.Name{Local: exemplarType.Name()}
		creators[exemplarName] = func() interface{} { return reflect.New(exemplarType).Interface() }
	}
	return decodeElementsByXmlNameFunc(creators)
}

func decodeElementsByXmlNameFunc(creators CreatorMap) elementProcessor {
	return func(decoder *xml.Decoder, start *xml.StartElement) (interface{}, error) {
		creator, ok := creators[start.Name]
		if ok {
			result := creator()
			err := decoder.DecodeElement(result, start)
			if err != nil {
				return nil, err
			}
			// This craziness is because DecodeElement wants a pointer-to-thing-as-interface,
			// and reflection is the only way to address or dereference such a beast.
			return reflect.ValueOf(result).Elem().Interface(), nil
		}
		return nil, nil
	}
}
