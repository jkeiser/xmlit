package xmlchan

import "testing"
import "strings"
import "encoding/xml"
import "reflect"

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

func retrieveDogDaysStrings(d *Decoder) error {
	token, err := d.Next()
	if err != nil {
		return err
	}
	start, ok := token.(xml.StartElement)
	if ok && start.Name.Local == "DogDays" {
		var str xml.Token
		str, err = d.Next()
		d.SendResult(string(str.(xml.CharData)))
	}
	return nil
}

func TestProcess(t *testing.T) {
	results, done := Process(strings.NewReader(XMLARIFFIC), retrieveDogDaysStrings)
	defer close(done)
	expectResult(t, results, "1")
	expectResult(t, results, "2")
	expectResult(t, results, "3")
	expectStreamClosed(t, results)
}

func TestProcessTo(t *testing.T) {
	results := make(chan Result, 10)
	done := make(chan struct{})
	ProcessTo(strings.NewReader(XMLARIFFIC), results, done, retrieveDogDaysStrings)
	defer close(done)
	expectResult(t, results, "1")
	expectResult(t, results, "2")
	expectResult(t, results, "3")
	expectStreamClosed(t, results)
}

func TestProcessWhenCallerEndsEarly(t *testing.T) {
	results, done := Process(strings.NewReader(XMLARIFFIC), retrieveDogDaysStrings)
	expectResult(t, results, "1")
	close(done)
	// We can't guarantee the other side will receive the close before we grab a value, but we *know* there won't be 2 more--at most 1.
	// Throw away one value.
	<-results
	expectStreamClosed(t, results)
}

func TestProcessFinished(t *testing.T) {
	results, done := Process(strings.NewReader(XMLARIFFIC), func(d *Decoder) error {
		d.SendResult("1")
		return FINISHED
	})
	defer close(done)
	expectResult(t, results, "1")
	expectStreamClosed(t, results)
}

func TestProcessWithErrorXml(t *testing.T) {
	str := `<?xml version="1.0" encoding="UTF-8" ?>
<blah><DogDays>3</DogDays><biddle`
	results, done := Process(strings.NewReader(str), retrieveDogDaysStrings)
	defer close(done)
	expectResult(t, results, "3")
	expectError(t, results)
	expectStreamClosed(t, results)
}

func TestToList(t *testing.T) {
	var list []string
	ToList(strings.NewReader(XMLARIFFIC), &list, retrieveDogDaysStrings)
	if len(list) != 3 {
		t.Errorf("Expected list to be 3 long, was %d", len(list))
	}
	if list[0] != "1" {
		t.Errorf("Expected list[0] to be 1, was %v", list[0])
	}
	if list[1] != "2" {
		t.Errorf("Expected list[1] to be 2, was %v", list[0])
	}
	if list[2] != "3" {
		t.Errorf("Expected list[2] to be 3, was %v", list[0])
	}
}

func TestProcessTypes(t *testing.T) {
	results, done := Process(strings.NewReader(XMLARIFFIC),
		ProcessTypes(Abominable{}, SuperTramp{}))
	defer close(done)
	expectResult(t, results, SuperTramp{DogDays: 1, FunTimes: "ridgemont"})
	expectResult(t, results, SuperTramp{DogDays: 2})
	expectResult(t, results, Abominable{Serious: "snowman"})
	expectResult(t, results, SuperTramp{DogDays: 3})
	expectStreamClosed(t, results)
}

func TestProcessTypesByName(t *testing.T) {
	// results, done := Process(strings.NewReader(XMLARIFFIC),
	// 	ProcessTypes(map[string]reflect.Type{
	// 		"LessAbominable": reflect.TypeOf(Abominable{}),
	// 		"SuperTramp":     reflect.TypeOf(SuperTramp{}),
	// 	}))
	// defer close(done)
	// expectResult(t, results, SuperTramp{DogDays: 1, FunTimes: "ridgemont"})
	// expectResult(t, results, SuperTramp{DogDays: 2})
	// expectResult(t, results, Abominable{Serious: "iceman"})
	// expectResult(t, results, SuperTramp{DogDays: 3})
	// expectStreamClosed(t, results)
}

func TestProcessTypesByXmlName(t *testing.T) {

}

func expectResult(t *testing.T, results <-chan Result, expected interface{}) {
	x := <-results
	if x != NewResult(expected) {
		t.Errorf("Got %+v, expected %+v", x, expected)
	}
}

func expectStreamClosed(t *testing.T, results <-chan Result) {
	x, ok := <-results
	if ok {
		t.Errorf("Results stream was not closed! Received %+v instead.", x)
	}
}

func expectError(t *testing.T, results <-chan Result) {
	x := <-results
	if x.Error == nil {
		t.Errorf("Got %+v, expected error", x)
	}
}
