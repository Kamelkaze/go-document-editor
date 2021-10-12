package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const DOCUMENT_FOLDER = "documents/"
const DEFAULT_PORT = 8080

//Use pointers to allow for nil values when updating a document, to note which fields should not be updated
type Payload struct {
	Title, Signee *string
	Content       Content
}

type Content struct {
	Header, Data *string
}

func main() {
	port := DEFAULT_PORT
	if len(os.Args) > 1 {
		arg := os.Args[1]
		intArg, err := strconv.Atoi(arg)
		if err == nil && intArg <= 65535 {
			port = intArg
		}
	}
	routeSetup()
	err := os.MkdirAll(DOCUMENT_FOLDER, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	strPort := ":" + strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(strPort, nil))
}

func routeSetup() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/new", newDocumentHandler)
	http.HandleFunc("/read", readDocumentHandler)
	http.HandleFunc("/update", updateDocumentHandler)
	http.HandleFunc("/delete", deleteDocumentHandler)
}

type FileWriter func(Payload, string) error
type DocumentRemover func(string) error

func writeToFile(payload Payload, fileName string) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return os.WriteFile(DOCUMENT_FOLDER+fileName, jsonPayload, 0777)
}

func newDocumentHandler(w http.ResponseWriter, r *http.Request) {
	fileNames, err := getExistingFileNames()
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	payload, err := extractPayload(r.Body)
	if err != nil {
		clientErrorResponse(w, err)
		return
	}
	createNewDocument(w, payload, writeToFile, fileNames)
}

func readDocumentHandler(w http.ResponseWriter, r *http.Request) {
	document, _, err := getDocumentByUrl(r.URL)
	if err != nil {
		clientErrorResponse(w, err)
		return
	}
	readDocument(w, document)
}

func updateDocumentHandler(w http.ResponseWriter, r *http.Request) {
	oldDocument, title, err := getDocumentByUrl(r.URL)
	if err != nil {
		clientErrorResponse(w, err)
		return
	}
	var oldContents Payload
	err = json.Unmarshal(oldDocument, &oldContents)
	if err != nil {
		clientErrorResponse(w, errors.New("existing document is incorrectly formatted, unable to read"))
		return
	}
	newContents, err := extractPayload(r.Body)
	if err != nil {
		clientErrorResponse(w, err)
		return
	}
	updatedContents := updateDocumentContents(oldContents, newContents)
	updateDocument(w, os.Remove, writeToFile, title, updatedContents)
}

func deleteDocumentHandler(w http.ResponseWriter, r *http.Request) {
	_, title, err := getDocumentByUrl(r.URL)
	if err != nil {
		clientErrorResponse(w, err)
		return
	}
	deleteDocument(w, title, os.Remove)
}

func getExistingFileNames() ([]string, error) {
	files, err := os.ReadDir(DOCUMENT_FOLDER)
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	if err != nil {
		return nil, err
	}
	return fileNames, nil
}

func extractPayload(body io.ReadCloser) (payload Payload, err error) {
	decoder := json.NewDecoder(body)
	err = decoder.Decode(&payload)
	if err != nil {
		return payload, err
	}
	return payload, nil
}

func createNewDocument(w http.ResponseWriter, payload Payload, writeContentToFile FileWriter, fileNames []string) {
	if payload.Title == nil {
		clientErrorResponse(w, errors.New("No title found"))
		return
	}

	for _, fileName := range fileNames {
		if *payload.Title == fileName {
			clientErrorResponse(w, errors.New("A document with that title already exists!"))
			return
		}
	}
	err := writeContentToFile(payload, *payload.Title)
	if err != nil {
		serverErrorResponse(w, err)
	}
}

func readDocument(w http.ResponseWriter, document []byte) {
	fmt.Fprint(w, string(document))
}

func updateDocument(w http.ResponseWriter, removeDocument DocumentRemover, writeContentToFile FileWriter, oldTitle string, updatedContent Payload) {
	err := removeDocument(DOCUMENT_FOLDER + oldTitle)
	if err != nil {
		serverErrorResponse(w, err)
		return
	}
	err = writeContentToFile(updatedContent, *updatedContent.Title)
	if err != nil {
		serverErrorResponse(w, err)
	}
}

func updateDocumentContents(oldContents, newContents Payload) (updatedContents Payload) {
	updatedContents.Title = oldContents.Title
	if newContents.Title != nil {
		updatedContents.Title = newContents.Title
	}

	updatedContents.Signee = oldContents.Signee
	if newContents.Signee != nil {
		updatedContents.Signee = newContents.Signee
	}

	updatedContents.Content.Header = oldContents.Content.Header
	if newContents.Content.Header != nil {
		updatedContents.Content.Header = newContents.Content.Header
	}

	updatedContents.Content.Data = oldContents.Content.Data
	if newContents.Content.Data != nil {
		updatedContents.Content.Data = newContents.Content.Data
	}

	return updatedContents
}

func deleteDocument(w http.ResponseWriter, title string, removeDocument DocumentRemover) {
	err := removeDocument(DOCUMENT_FOLDER + title)
	if err != nil {
		serverErrorResponse(w, err)
	}
}

func getDocumentByUrl(url *url.URL) ([]byte, string, error) {
	title, err := getTitleFromUrl(url)
	if err != nil {
		return nil, "", err
	}

	document, err := os.ReadFile(DOCUMENT_FOLDER + title)
	if err != nil {
		return nil, "", errors.New("no document with title \"" + title + "\" could be found")
	}

	return document, title, nil
}

func getTitleFromUrl(url *url.URL) (string, error) {
	title := url.Query().Get("title")
	if title == "" {
		return "", errors.New("no title was provided")
	}
	return title, nil
}

func clientErrorResponse(w http.ResponseWriter, errorMessage error) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, errorMessage.Error())
}

func serverErrorResponse(w http.ResponseWriter, errorMessage error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, errorMessage.Error())
}
