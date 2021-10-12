curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
"title": '"$1"',
"content": {
"header": '"$2"',
"data": '"$3"'
},
"signee": '"$4"'
}' \
  http://localhost:8080/new

