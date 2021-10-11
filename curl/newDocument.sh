curl --header "Content-Type: application/json" \
  --request POST \
  --data '{
"title": "doc3",
"content": {
"header": "Employee contract",
"data": "I agree to the terms"
},
"signee": "John McClane"
}' \
  http://localhost:8080/new

