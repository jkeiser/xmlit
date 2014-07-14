package xmlit

import "testing"
import "strings"
import "github.com/jkeiser/iter"
import "encoding/xml"
import "io"

func TestDecodeElements(t *testing.T) {
	iterator := DecodeElements(StringReader(XMLARIFFIC), Abominable{}, SuperTramp{})
	expectItem(t, iterator, SuperTramp{DogDays: 1, FunTimes: "ridgemont"})
	expectItem(t, iterator, SuperTramp{DogDays: 2})
	expectItem(t, iterator, Abominable{Serious: "snowman"})
	expectItem(t, iterator, SuperTramp{DogDays: 3})
	expectDoneIterating(t, iterator)
}

func TestDecodeElementsByXmlName(t *testing.T) {
	iterator := DecodeElementsByXmlName(StringReader(XMLARIFFIC), CreatorMap{
		xml.Name{Local: "LessAbominable"}: func() interface{} { return new(Abominable) },
		xml.Name{Local: "SuperTramp"}:     func() interface{} { return new(SuperTramp) },
	})
	expectItem(t, iterator, SuperTramp{DogDays: 1, FunTimes: "ridgemont"})
	expectItem(t, iterator, SuperTramp{DogDays: 2})
	expectItem(t, iterator, Abominable{Serious: "iceman"})
	expectItem(t, iterator, SuperTramp{DogDays: 3})
	expectDoneIterating(t, iterator)
}

func TestDecodeElementsWithErrorXml(t *testing.T) {
	str := `<?xml version="1.0" encoding="UTF-8" ?>
<blah><SuperTramp><DogDays>3</DogDays></SuperTramp><biddle`
	iterator := DecodeElements(StringReader(str), SuperTramp{})
	expectItem(t, iterator, SuperTramp{DogDays: 3})
	expectError(t, iterator)
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

func expectItem(t *testing.T, iterator iter.Iterator, expected interface{}) {
	item, err := iterator.Next()
	if item != expected || err != nil {
		t.Errorf("Got %+v / %+v, expected item %+v", item, err, expected)
	}
}

func expectError(t *testing.T, iterator iter.Iterator) {
	item, err := iterator.Next()
	if err == nil || item != nil {
		t.Errorf("Got %+v / %+v, expected error", item, err)
	}
}

func expectDoneIterating(t *testing.T, iterator iter.Iterator) {
	item, err := iterator.Next()
	if err != iter.FINISHED || item != nil {
		t.Errorf("Got %+v / %+v, expected iteration to be finished", item, err)
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
