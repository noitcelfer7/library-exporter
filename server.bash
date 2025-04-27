openssl req -newkey rsa:4096 \
  -keyout server-key.pem -out server-req.pem \
  -config ssl.conf -nodes

openssl x509 -req -in server-req.pem \
  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial \
  -out server-cert.pem -days 365 \
  -extfile ssl.conf -extensions v3_req
