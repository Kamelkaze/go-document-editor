package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func isError(shouldPass bool, w *httptest.ResponseRecorder) bool {
	return !shouldPass && w.Result().StatusCode < 300 || shouldPass && w.Result().StatusCode >= 300
}

type MockDocumentRemover struct {
	shouldPass bool
}

func (remover MockDocumentRemover) removeDocument(title string) error {
	if !remover.shouldPass {
		return errors.New("Fail")
	}
	return nil
}

type MockFileWriter struct {
	shouldPass bool
}

func (writer MockFileWriter) writeToFile(payload Payload, title string) error {
	if !writer.shouldPass {
		return errors.New("Fail")
	}
	return nil
}

type NewDocTest struct {
	Payload    Payload
	shouldPass bool
}

var title string = "title"
var signee string = "signee"
var header string = "header"
var text string = "text"

var existingTitle string = "bad"

var newDocTests = []NewDocTest{
	{Payload{&title, &signee, Content{&header, &text}}, true},
	{Payload{&title, &signee, Content{nil, nil}}, true},
	{Payload{&title, nil, Content{nil, nil}}, true},
	{Payload{&title, nil, Content{&header, &text}}, true},
	{Payload{nil, nil, Content{nil, nil}}, false},
	{Payload{nil, &signee, Content{&header, &text}}, false},
	{Payload{&existingTitle, &signee, Content{&header, &text}}, false},
}

func TestNewDocument(t *testing.T) {
	for i := 0; i < 2; i++ { //Checks that tests fail if writer throws error
		writerShouldPass := i == 0
		for _, test := range newDocTests {
			var bfr bytes.Buffer
			err := json.NewEncoder(&bfr).Encode(test.Payload)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()
			mockWriter := MockFileWriter{writerShouldPass}
			fileNames := []string{"bad", "good"}
			createNewDocument(w, test.Payload, mockWriter.writeToFile, fileNames)
			if isError(test.shouldPass && writerShouldPass, w) {
				t.Errorf("Tests result was %t, got %t", (test.shouldPass && writerShouldPass), !(test.shouldPass && writerShouldPass))
			}
		}
	}
}
func TestReadDocument(t *testing.T) {
	w := httptest.NewRecorder()
	content := []byte("hello")
	readDocument(w, content)
	if !bytes.Equal(content, w.Body.Bytes()) {
		t.Errorf("Content not equal")
	}
}

type UpdateDocumentTest struct {
	readerShouldPass bool
	writerShouldPass bool
	shouldPass       bool
}

var updateDocumentTest = []UpdateDocumentTest{
	{true, true, true},
	{true, false, false},
	{false, true, false},
	{false, false, false},
}

func TestUpdateDocument(t *testing.T) {
	for _, test := range updateDocumentTest {
		w := httptest.NewRecorder()
		remover := MockDocumentRemover{test.readerShouldPass}
		writer := MockFileWriter{test.writerShouldPass}
		str := "test"
		payload := Payload{&str, &str, Content{&str, &str}}
		updateDocument(w, remover.removeDocument, writer.writeToFile, "old", payload)
		if isError(test.shouldPass, w) {
			t.Errorf("Failed, should not have failed")
		}
	}
}

type UpdateNewDocContentTest struct {
	new    Payload
	old    Payload
	result Payload
}

var n string = "new"
var o string = "old"

var updateNewDocContentTest = []UpdateNewDocContentTest{
	{
		Payload{&o, &o, Content{&o, &o}},
		Payload{&n, &n, Content{nil, nil}},
		Payload{&n, &n, Content{&o, &o}},
	},
	{
		Payload{&o, &o, Content{&o, &o}},
		Payload{nil, nil, Content{&n, &n}},
		Payload{&o, &o, Content{&n, &n}},
	},
	{
		Payload{&o, &o, Content{&o, &o}},
		Payload{nil, nil, Content{nil, nil}},
		Payload{&o, &o, Content{&o, &o}},
	},
	{
		Payload{&o, &o, Content{&o, &o}},
		Payload{&n, &n, Content{&n, &n}},
		Payload{&n, &n, Content{&n, &n}},
	},
}

func TestUpdateDocumentContents(t *testing.T) {
	for _, test := range updateNewDocContentTest {
		res := updateDocumentContents(test.new, test.old)
		if !reflect.DeepEqual(res, test.result) {
			t.Errorf("Updated payload not equal to expected payload.")
		}
	}
}

func TestDeleteDocument(t *testing.T) {
	w := httptest.NewRecorder()
	remover := MockDocumentRemover{false}
	deleteDocument(w, "doc1", remover.removeDocument)
	if isError(false, w) {
		t.Errorf("Did fail, should not have failed")
	}
	w = httptest.NewRecorder()
	remover = MockDocumentRemover{true}
	deleteDocument(w, "doc1", remover.removeDocument)
	if isError(true, w) {
		t.Errorf("Did not fail, should have failed")
	}
}

func TestGetTitleFromUrl(t *testing.T) {
	val := "fish"
	url := url.URL{
		Scheme:   "https",
		Host:     "fakewebsite.com",
		Path:     "",
		RawQuery: "title=" + val,
	}
	res, err := getTitleFromUrl(&url)
	if err != nil {
		t.Errorf("Recived error " + err.Error())
	}
	if res != val {
		t.Errorf("Value of parameter should have been %s, recived %s instead", val, res)
	}
	url.RawQuery = "titl=" + val
	if err != nil {
		t.Errorf("Should have failed, did not")
	}
}
