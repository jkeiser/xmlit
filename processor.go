package xmlchan

import "encoding/xml"
import "errors"
import "io"
import "reflect"

var FINISHED = errors.New("FINISHED")

type TokenProcessor func(decoder *Decoder) error

func Process(reader io.Reader, process TokenProcessor) (<-chan Result, chan<- struct{}) {
	results := make(chan Result)
	done := make(chan struct{})
	ProcessTo(reader, results, done, process)
	return results, done
}

func ProcessTo(reader io.Reader, results chan<- Result, done <-chan struct{}, process TokenProcessor) {
	go func() {
		defer close(results)

		decoder := Decoder{
			decoder: xml.NewDecoder(reader),
			results: results,
			done:    done,
		}

		for !decoder.Done() {
			err := process(&decoder)
			if err == FINISHED || err == io.EOF {
				return
			} else if err != nil {
				decoder.SendError(err)
				return
			}
		}
	}()
}

func ToList(reader io.Reader, slice interface{}, process TokenProcessor) error {
	results, done := Process(reader, process)
	defer close(done)

	// Get []<type> (dereference *[]<type>)
	sliceVal := reflect.ValueOf(slice).Elem()

	for result := range results {
		value, err := result.Get()
		if err != nil {
			return err
		}
		sliceVal.Set(reflect.Append(sliceVal, reflect.ValueOf(value)))
	}
	return nil
}

type Result struct {
	Result interface{}
	Error  error
}

func (result Result) Get() (interface{}, error) {
	return result.Result, result.Error
}

func NewResult(result interface{}) Result {
	return Result{Result: result}
}

func NewError(err error) Result {
	return Result{Error: err}
}
