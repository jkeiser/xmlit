package xmlchan

import "testing"
import "strings"

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

func TestDecodeElements(t *testing.T) {
	iter := DecodeElements(strings.NewReader(XMLARIFFIC), Abominable{}, SuperTramp{})
	expectItem(t, iter, SuperTramp{DogDays: 1, FunTimes: "ridgemont"})
	expectItem(t, iter, SuperTramp{DogDays: 2})
	expectItem(t, iter, Abominable{Serious: "snowman"})
	expectItem(t, iter, SuperTramp{DogDays: 3})
	expectDoneIterating(t, iter)
}

func expectItem(t *testing.T, iter Iterator, expected interface{}) {
	item, err := iter.Next()
	if item != expected {
		t.Errorf("Got %+v / %+v, expected item %+v", item, err, expected)
	}
}

func expectDoneIterating(t *testing.T, iter Iterator) {
	item, err := iter.Next()
	if err != FINISHED {
		t.Errorf("Got %+v / %+v, expected iteration to be finished", item, err)
	}
}
