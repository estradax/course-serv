curl -v -X POST \
  -H 'Connection: close' \
  -H 'Content-Type: application/json' \
  -d '{"email": "johndoe5@mail.com", "password": "foobar1234"}' \
  http://localhost:8080/api/v1/login
