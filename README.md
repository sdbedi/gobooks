
# Gobooks - About

This API tracks books, with the following attributes:
Title (string) (required)\
Author (string) (required)
Publisher (string) \
Publish Date (string)\
Rating (1-3) (required)\
Status (must CheckedIn or Checkout; defaults to CheckedIn)\

Currently, you must supply an author, a title, and a Rating of 1-3 on creating the book. Books are created with a default status of "CheckedIn", unless explcitly supplied with a status of "CheckedOut". Attempts to supply any other status will trigger an error. 

# Getting Started
You'll need to have Docker, Postgres and Go install on your system. Otherwise, here are the steps:

1. Clone the repo:
```git clone https://github.com/sdbedi/gobooks.git```
2. Cd into the root directory, ie gobooks
3. Execute: 
```docker-compose up```

The docker-compose up command will download all dependencies, compile the source code, intialize the Postgres database, and run the compiled code.

# To Test

First start the api:
```docker-compose up```

Then, in a separate terminal window, navigate to the root directory:
```go test -v ```

# Endpoints

Here are the endpoints - note that the book's ID is created at the same time as the book; it will be different from the examples here.

**Get a book**
```http request
GET http://localhost:8080/api/v1/books?id=123456789
```

**List books**
```http request
GET http://localhost:8080/api/v1/books/list
```

**List books w/ limit**
```http request
GET http://localhost:8080/api/v1/books/list?limit=1
```

**Query books by title**
```http request
GET http://localhost:8080/api/v1/books/list?title=e
```

**Create a book**
```http request
POST http://localhost:8080/api/v1/books
Content-Type: application/json

{
    "author": "Zadie Smith",
    "title": "White Teeth",
    "publisher" : "Penguin",
    "publishdate": "2002"
    "rating": 2,
    "status": "CheckedOut"
}
```

**Update book's general details**
```http request
PUT http://localhost:8080/api/v1/books/update
Content-Type: application/json

{
    "id": "123456789",
    "author": "Autograph Man"
}
```

**Delete the book**
```http request
DELETE http://localhost:8080/api/v1/books?id=123456>
Content-Type: application/json
```

# Known Issues/TODOS
1. Testing Requires GCC (GNU Compiler Collection). If you encounter of this type:
```runtime/cgo cgo: exec gcc: exec: "gcc": executable file not found```
    You'll need to install GCC with the following:
    ```apt-get install build-essential```

2. Currently, the build process creates and saves a docker image. If you make any updates to the source code, you'll need to delete the gobooks_app Docker image before running docker compose again.

3. Current, there's no guarding of input for the ratings field, ie it's possible to enter a rating of 4 via update. Implementing the safety check will require some refactoring. 

4. The date field will currently accept any string input. I imagine that there would be considerable variety in the form of this data, depending on the source/type of book. 

5. There's currently no , ie, if you submit the same info over and over, you'll end up creating different book entities with different ids but similar underlying info. 


