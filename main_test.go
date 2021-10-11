package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

type MockDocumentRemover struct {
	ShouldFail bool
}

func (remover MockDocumentRemover) removeDocument(title string) error {
	if remover.ShouldFail {
		return errors.New("Fail")
	}
	return nil
}

type MockFileWriter struct {
	ShouldFail bool
}

func (writer MockFileWriter) writeToFile(payload Payload, title string) error {
	if writer.ShouldFail {
		return errors.New("Fail")
	}
	return nil
}

type NewDocTest struct {
	Payload Payload
}

var newDocTests = []NewDocTest{
	{Payload{nil, nil, Content{nil, nil}}},
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
}

func TestNewDocument(t *testing.T) {
	for _, test := range newDocTests {
		var bfr bytes.Buffer
		err := json.NewEncoder(&bfr).Encode(test.Payload)
		if err != nil {
			t.Fatal(err)
		}
		//req := httptest.NewRequest("GET", "http://localhost:8080/new", nil)
		w := httptest.NewRecorder()
		mockWriter := MockFileWriter{true}
		fileNames := []string{"asd"}
		createNewDocument(w, test.Payload, mockWriter.writeToFile, fileNames)
		if w.Result().StatusCode != 400 {
			t.Errorf("Should have gotten a status code 400, instead got %s", strconv.Itoa(w.Result().StatusCode))
		}
	}

}
func TestReadDocument(t *testing.T) {
	w := httptest.NewRecorder()
	readDocument(w, []byte("hello"))
}

func TestUpdateDocument(t *testing.T) {
	w := httptest.NewRecorder()
	remover := MockDocumentRemover{false}
	writer := MockFileWriter{false}
	str := "test"
	payload := Payload{&str, &str, Content{&str, &str}}
	updateDocument(w, remover.removeDocument, writer.writeToFile, "old", payload)
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
}
