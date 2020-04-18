package eventbus

import (
	"errors"
	"testing"

	"github.com/mishudark/eventhus"
)

type producerStubEntry struct {
	event  eventhus.Event
	bucket string
	subset string
}

type producerStub struct {
	entries []producerStubEntry
	err     error
}

func (s *producerStub) Publish(event eventhus.Event, bucket, subset string) error {
	s.entries = append(s.entries, producerStubEntry{
		event:  event,
		bucket: bucket,
		subset: subset,
	})

	return s.err
}

func Test_MultiPublisherError_String(t *testing.T) {
	sut := MultiPublisherError{}
	sut.Add(errors.New("testing"))
	sut.Add(errors.New("testing"))
	sut.Add(errors.New("testing"))

	e := `A few errors occurred:
	1) testing
	2) testing
	3) testing`

	if e != sut.Error() {
		t.Error(
			"incorrect output from MultiPublisherError.Error()\n-- expected\n",
			e,
			"\n-- actual\n",
			sut.Error(),
		)
	}
}

func Test_MultiPublisherError_Len(t *testing.T) {
	sut := MultiPublisherError{}
	sut.Add(errors.New("testing"))
	sut.Add(errors.New("testing"))
	sut.Add(errors.New("testing"))

	if sut.Len() != 3 {
		t.Error("incorrect length after adding 3 got", sut.Len())
	}
}

func Test_NewMultiPublisher(t *testing.T) {
	producerOne := &producerStub{}
	producerTwo := &producerStub{}
	sut := NewMultiPublisher(producerOne, producerTwo)

	err := sut.Publish(eventhus.Event{}, "banks", "accounts")

	if len(producerOne.entries) != 1 || len(producerTwo.entries) != 1 {
		t.Error("MultiPublisher did not give both producers 1 event")
	}

	if err != nil {
		t.Error("error provided when it shouldn't of been provided:", err)
	}
}

func Test_NewMultiPublisher_AggregatesErrors(t *testing.T) {
	producerOne := &producerStub{
		err: errors.New("expected error 1"),
	}
	producerTwo := &producerStub{
		err: errors.New("expected error 2"),
	}
	sut := NewMultiPublisher(producerOne, producerTwo)

	aerr := sut.Publish(eventhus.Event{}, "banks", "accounts")

	err, ok := aerr.(MultiPublisherError)

	if !ok {
		t.Error("error was expected to be type MultiPublisherError")
		t.FailNow()
	}

	if err.Len() != 2 {
		t.Error("error length was expected to be 2 got:", err.Len())
	}
}
