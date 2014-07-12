package xmlchan

import "encoding/xml"
import "reflect"

type ElementProcessor func(decoder *Decoder, start *xml.StartElement) (interface{}, error)

func ProcessElements(process ElementProcessor) TokenProcessor {
	return func(decoder *Decoder) error {
		token, err := decoder.Next()
		if err != nil {
			return err
		}
		start, ok := token.(*xml.StartElement)
		if ok {
			var result interface{}
			result, err = process(decoder, start)
			if err != nil {
				return err
			}
			if result != nil {
				decoder.SendResult(result)
			}
		}
		return nil
	}
}

func decodeElement(decoder *Decoder, start *xml.StartElement, t reflect.Type) (interface{}, error) {
	result := reflect.Zero(t).Interface()
	err := decoder.DecodeElement(result, start)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// results, errors = xmlchan.ProcessTypes(Actor{}, Series{}).Process(reader)
// for result := range results {
//
// }
func ProcessTypes(sentinels ...interface{}) TokenProcessor {
	return ProcessElements(func(decoder *Decoder, start *xml.StartElement) (interface{}, error) {
		for _, s := range sentinels {
			t := reflect.TypeOf(s)
			if start.Name.Local == t.Name() {
				return decodeElement(decoder, start, t)
			}
		}
		return nil, nil
	})
}

// results, errors = xmlchan.ProcessTypesByName(map[string]reflect.Type{"Foo": Actor, "Bar": Series}).Process(reader)
func ProcessTypesByName(types map[string]reflect.Type) TokenProcessor {
	return ProcessElements(func(decoder *Decoder, start *xml.StartElement) (interface{}, error) {
		for name, t := range types {
			if start.Name.Local == name {
				return decodeElement(decoder, start, t)
			}
		}
		return nil, nil
	})
}

func ProcessTypesByXmlName(types map[xml.Name]reflect.Type) TokenProcessor {
	return ProcessElements(func(decoder *Decoder, start *xml.StartElement) (interface{}, error) {
		for name, t := range types {
			if start.Name == name {
				return decodeElement(decoder, start, t)
			}
		}
		return nil, nil
	})
}
