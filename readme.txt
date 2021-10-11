To start the server run "go run main", the server will listen on port 8080 by default.
You can specify port by providing a port number e.g. "go run main 8000".

Make request in the form of http://localhost:8080/$command, where $command could be: 'new', 'update', 'read', or 'delete'.
The commands 'update', 'read', and 'delete' require a url parameter named title e.g.: new?title=document1. The 'new' command requires a payload in the form specified. When using the 'update' command the payload should include those fields which are to be updated (in the form specified), with thier new values. Exluded fields will not be changed. So to create a new document,:http://localhost:8080/new + payload. To delete a document: http://localhost:8080/delete?title=doc25

To run tests run "go test".