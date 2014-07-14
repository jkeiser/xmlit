package xmlit

import "io"
import "encoding/xml"
import "reflect"
import "github.com/jkeiser/iter"

type CreatorMap map[xml.Name]func() interface{}

// iterator := xmlit.DecodeElements(reader, Customer{}, Order{}, Account{})
func DecodeElements(reader io.ReadCloser, exemplars ...interface{}) iter.Iterator {
	return iterateElements(reader, decodeElementsFunc(exemplars))
}

// iterator := xmlit.DecodeElementsByXmlName(reader, CreatorMap{
//   "Customer": func() { Customer{} },
//   "Order": func() { Order{} },
//   "DelinquentOrder": func() { DelinquentOrder{} },
//   "Account": func() { Account{} },
// })
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
