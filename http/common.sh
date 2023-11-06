make_url() {
  echo "http://localhost:8080$1"
}

get() {
  local ret=$(make_url $1)

  curl -v -X GET \
    -H 'Connection: close' \
    $ret
}

get_with_token() {
  local ret=$(make_url $1)

  curl -v -X GET \
    -H 'Connection: close' \
    -H 'Authorization: Bearer '$2 \
    $ret
}
