# redeam-rest: A sample RESTful deploy

# General Notes

This is an attempt to demonstrate a generalized approach to development of a simple CRUD REST api.

I have opted to:
- use https://github.com/golang-standards/project-layout for directory structure
- use protobuf and database versioning
- use protobuf REST annotations (swagger)

These decisions were made to emulate a more realistic production level approach. As such, the code is a bit more complicated than the most barebones of implementations, but it represents a service that is more akin to the sort of standard I hold myself to.

In retrospect, I might choose a simpler approach, because my time ended up being signficantly constrained due to "life". On the upside, this provides a pretty decent scaffolding (and is based upon prior work I had done in that regard) for future projects.

# Stack notes
This is built using golang, postgres and gRPC

# Implemented Endpoints

## Request: GET /v1/book/{id}
`curl -i -H 'Accept: application/json' http://localhost:8080/v1/book/1`
*Response:* ```HTTP/1.1 200 OK
Content-Type: application/json
Grpc-Metadata-Content-Type: application/grpc
Date: Wed, 06 Mar 2019 18:20:07 GMT
Content-Length: 269

{"api":"v1","book":{"id":"1","title":"30-Second Philosophies The 50 Most Thought-Provoking Philosophies, Each Explained in Half a Minute","author":"Barry (Editor) Loewer","publisher":"Metro Books","publish_date":"2002-10-02T15:00:00Z","rating":2,"status":"CHECKED_IN"}```

### Request: GET /v1/book/all
`curl -i -H 'Accept: application/json' http://localhost:8080/v1/book/all`
* Response:* ```HTTP/1.1 200 OK
Content-Type: application/json
Grpc-Metadata-Content-Type: application/grpc
Date: Wed, 06 Mar 2019 18:20:52 GMT
Content-Length: 763

{"api":"v1","books":[{"id":"1","title":"30-Second Philosophies The 50 Most Thought-Provoking Philosophies, Each Explained in Half a Minute","author":"Barry (Editor) Loewer","publisher":"Metro Books","publish_date":"2002-10-02T15:00:00Z","rating":2,"status":"CHECKED_IN"},{"id":"3","title":"30-Second Philosophies The 50 Most Thought-Provoking Philosophies, Each Explained in Half a Minute","author":"Barry (Editor) Loewer","publisher":"Metro Books","publish_date":"2002-10-02T15:00:00Z","rating":2,"status":"CHECKED_IN"},{"id":"2","title":"30-Second Philosophies The 50 Most Thought-Provoking Philosophies, Each Explained in Half a Minute","author":"Barry Loewer","publisher":"Metro Books","publish_date":"2002-10-02T15:00:00Z","rating":2,"status":"CHECKED_IN"}]}```

### Request: POST /v1/book
`curl -i -H 'Accept: application/json' http://localhost:8080/v1/book --data '{"api": "v1","book": {"title": "30-Second Philosophies The 50 Most Thought-Provoking Philosophies, Each Explained in Half a Minute","author": "Barry Loewer","publisher": "Metro Books","publishDate": "2002-10-02T15:00:00Z","rating": 2.0,"status": 1}}'`
*Body:* ```{
	"api": "v1",
	"book": {
		"title": "30-Second Philosophies The 50 Most Thought-Provoking Philosophies, Each Explained in Half a Minute",
		"author": "Barry Loewer",
		"publisher": "Metro Books",
		"publishDate": "2002-10-02T15:00:00Z",
		"rating": 2.0,
		"status": 1
	}
}```

*Response:* ```HTTP/1.1 200 OK
Content-Type: application/json
Grpc-Metadata-Content-Type: application/grpc
Date: Wed, 06 Mar 2019 18:27:00 GMT
Content-Length: 21

{"api":"v1","id":"6"}```

### Request: PUT /v1/book/{id}
`curl -i -H 'Accept: application/json' http://localhost:8080/v1/book --data '{"api": "v1","book": {"title": "30-Second Philosophies The 50 Most Thought-Provoking Philosophies, Each Explained in Half a Minute","author": "Barry Loewer","publisher": "Metro Books","publishDate": "2002-10-02T15:00:00Z","rating": 2.0,"status": 1}}'`
*Body:* ```{
	"api": "v1",
	"book": {
		"title": "30-Second Philosophies The 50 Most Thought-Provoking Philosophies, Each Explained in Half a Minute",
		"author": "Barry Loewer",
		"publisher": "Metro Books",
		"publishDate": "2002-10-02T15:00:00Z",
		"rating": 2.0,
		"status": 1
	}
}```

*Response:* ```HTTP/1.1 200 OK
Content-Type: application/json
Grpc-Metadata-Content-Type: application/grpc
Date: Wed, 06 Mar 2019 18:27:00 GMT
Content-Length: 21

{"api":"v1","updated":"1"}```

### Request: DELETE /v1/book/{id}
`curl -i -H 'Accept: application/json' http://localhost:8080/v1/book/1 --request DELETE `
*Response:* ```HTTP/1.1 200 OK
Content-Type: application/json
Grpc-Metadata-Content-Type: application/grpc
Date: Wed, 06 Mar 2019 18:27:00 GMT
Content-Length: 21

{"api":"v1","deleted":"1"}```

