curl --header "Content-Type: application/json" \
  --request POST \
 --data '{
"title": '"$2"',
"content": {
"header": '"$3"',
"data": '"$4"'
},
"signee": '"$5"'
}' \
  http://localhost:8080/update?title=$1
