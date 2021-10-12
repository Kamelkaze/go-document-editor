To start the server run "go run main.go", the server will listen on port 8080 by default.
You can specify port by providing a port number e.g. "go run main 8000".

The curl folder contains useful scripts for using the program. Check info.txt for information about which arguments needs to be passed to each script.

Alternatively, make request in the form of http://localhost:8080/$command, where $command could be: 'new', 'update', 'read', or 'delete'.
The commands 'update', 'read', and 'delete' require a url parameter named title e.g.: new?title=document1. The 'new' command requires a payload in the form specified. When using the 'update' command the payload should include those fields which are to be updated (in the form specified), with thier new values. Exluded fields will not be changed. So to create a new document,:http://localhost:8080/new + payload. To delete a document: http://localhost:8080/delete?title=doc25

To run the tests run "go test".

Regarding the testing I used the handlers (eg. newDocumentsHandler) to extract as many dependencies as I could (io, json, url) so that we can focus on testing without being dependant on externalities. In this case the unit tests become somwhat trivial since there isn't much logic in the application. There are certainly good arguments for testing these handlers as well, or simply mocking all of the dependencies, but a bit overkill for this task I think.

There are some limitations regarding title name. The title of the document is set as the filename, so you can't use characters that your OS does not allow in filenames. You also cannot have two documents with the same title. 

Another limiations is that the documents are expected to be created and edited only by the program, so there are no guarantees about what could happen if you manually edit them.