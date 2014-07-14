package xmlit

import "io"
import "encoding/xml"
import "github.com/jkeiser/iter"

type elementProcessor func(decoder *xml.Decoder, start *xml.StartElement) (interface{}, error)

func iterateElements(reader io.ReadCloser, processor elementProcessor) iter.Iterator {
	var decoder = xml.NewDecoder(reader)
	return iter.Iterator{
		Next: func() (interface{}, error) {
			for {
				token, err := decoder.Token()
				if err != nil {
					if err == io.EOF {
						return nil, iter.FINISHED
					}
					return nil, err
				}
				start, ok := token.(xml.StartElement)
				if ok {
					var result interface{}
					result, err = processor(decoder, &start)
					if err != nil {
						return nil, err
					}
					if result != nil {
						return result, nil
					}
				}
			}
		},
		Close: func() { reader.Close() },
	}
}
