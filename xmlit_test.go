package xmlit

import "testing"
import "strings"
import "encoding/xml"
import "io"

func TestDecodeElements(t *testing.T) {
	decoder := DecodeElements(StringReader(XMLARIFFIC), Abominable{}, SuperTramp{})
	expectItem(t, decoder, SuperTramp{DogDays: 1, FunTimes: "ridgemont"})
	expectItem(t, decoder, SuperTramp{DogDays: 2})
	expectItem(t, decoder, Abominable{Serious: "snowman"})
	expectItem(t, decoder, SuperTramp{DogDays: 3})
	expectDoneIterating(t, decoder)
}

func TestDecodeElementsByXmlName(t *testing.T) {
	decoder := ElementDecoder{
		Decoder: xml.NewDecoder(StringReader(XMLARIFFIC)),
		Creators: CreatorMap{
			xml.Name{Local: "LessAbominable"}: func() interface{} { return new(Abominable) },
			xml.Name{Local: "SuperTramp"}:     func() interface{} { return new(SuperTramp) },
		},
	}
	expectItem(t, &decoder, SuperTramp{DogDays: 1, FunTimes: "ridgemont"})
	expectItem(t, &decoder, SuperTramp{DogDays: 2})
	expectItem(t, &decoder, Abominable{Serious: "iceman"})
	expectItem(t, &decoder, SuperTramp{DogDays: 3})
	expectDoneIterating(t, &decoder)
}

func TestDecodeElementsWithErrorXml(t *testing.T) {
	str := `<?xml version="1.0" encoding="UTF-8" ?>
<blah><SuperTramp><DogDays>3</DogDays></SuperTramp><biddle`
	decoder := DecodeElements(StringReader(str), SuperTramp{})
	expectItem(t, decoder, SuperTramp{DogDays: 3})
	expectError(t, decoder)
}

const XMLARIFFIC = `
<?xml version="1.0" encoding="UTF-8" ?>
<list>
  <anotherlist>
    <SuperTramp>
      <DogDays>1</DogDays>
      <FunTimes>ridgemont</FunTimes>
    </SuperTramp>
    <yetanotherlist>
      <SuperTramp>
        <DogDays>2</DogDays>
      </SuperTramp>
    </yetanotherlist>
    <Yeti>brown fur</Yeti>
    <Abominable><Serious>snowman</Serious></Abominable>
    <LessAbominable><Serious>iceman</Serious></LessAbominable>
    <Yeti>blue</Yeti>
  </anotherlist>
  <Yeti>black</Yeti>
  <SuperTramp>
    <DogDays>3</DogDays>
  </SuperTramp>
  <Yeti>turquoise</Yeti>
</list>
`

type SuperTramp struct {
	DogDays  int
	FunTimes string
}

type Abominable struct {
	Serious string
}

func expectItem(t *testing.T, decoder *ElementDecoder, expected interface{}) {
	item := decoder.Next()
	err := decoder.Error()
	if item != expected || err != nil {
		t.Errorf("Got %+v / %+v, expected item %+v", item, err, expected)
	}
}

func expectError(t *testing.T, decoder *ElementDecoder) {
	item := decoder.Next()
	err := decoder.Error()
	if err == nil || item != nil {
		t.Errorf("Got %+v / %+v, expected error", item, err)
	}
}

func expectDoneIterating(t *testing.T, decoder *ElementDecoder) {
	hasNext := decoder.HasNext()
	if hasNext {
		item := decoder.Next()
		err := decoder.Error()
		t.Errorf("Got hasNext=%v, %+v / %+v, expected iteration to be finished", hasNext, item, err)
	}
}

type EmptyCloser struct {
	Reader io.Reader
}

func (e EmptyCloser) Read(p []byte) (n int, err error) {
	return e.Reader.Read(p)
}

func (e EmptyCloser) Close() error { return nil }

func StringReader(str string) io.ReadCloser {
	return EmptyCloser{strings.NewReader(str)}
}
