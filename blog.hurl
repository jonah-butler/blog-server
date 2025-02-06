# TEST: get blog collection from provided offset
# TAG: positive
GET http://127.0.0.1:8080/blog/?offset=0

HTTP 200
[Asserts]
status == 200

# TEST: get blog collection from provided offset
# TAG: positive
GET http://127.0.0.1:8080/blog/random

HTTP 200
[Asserts]
status == 200

# TEST: get blog with invalid slug and json properties are null
# TAG: negative
GET http://127.0.0.1:8080/blog/sorting-in

HTTP 404
[Asserts]
status == 404

# TEST: get blog with valid slug and check additional properties
# TAG: positive
GET http://127.0.0.1:8080/blog/sorting-in-go

HTTP 200
[Asserts]
status == 200
jsonpath "$.post1" != null
jsonpath "$.previousPost" != null
