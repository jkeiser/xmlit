package xmlchan

import "io"
import "encoding/xml"
import "reflect"

type elementProcessor func(decoder *xml.Decoder, start *xml.StartElement) (interface{}, error)

func iterateElements(decoder *xml.Decoder, processor elementProcessor) Iterator {
	return elementIterator{decoder: decoder, processor: processor}
}

type elementIterator struct {
	decoder   *xml.Decoder
	processor elementProcessor
}

func (iter elementIterator) Next() (interface{}, error) {
	for {
		token, err := iter.decoder.Token()
		if err == io.EOF {
			return nil, FINISHED
		}
		if err != nil {
			return nil, err
		}
		start, ok := token.(xml.StartElement)
		if ok {
			var result interface{}
			result, err = iter.processor(iter.decoder, &start)
			if err != nil {
				return nil, err
			}
			if result != nil {
				return result, nil
			}
		}
	}
}

type CreatorMap map[xml.Name]func() interface{}

// iterator := xmlchan.DecodeElements(reader, Customer{}, Order{}, Account{})
func DecodeElements(reader io.Reader, exemplars ...interface{}) Iterator {
	var decoder = xml.NewDecoder(reader)
	return iterateElements(decoder, decodeElementsFunc(exemplars))
}

// iterator := xmlchan.DecodeElementsByXmlName(reader, CreatorMap{
//   "Customer": func() { Customer{} },
//   "Order": func() { Order{} },
//   "DelinquentOrder": func() { DelinquentOrder{} },
//   "Account": func() { Account{} },
// })
func DecodeElementsByXmlName(reader io.Reader, creators CreatorMap) Iterator {
	var decoder = xml.NewDecoder(reader)
	return iterateElements(decoder, decodeElementsByXmlNameFunc(creators))
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
