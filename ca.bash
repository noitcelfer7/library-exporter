openssl req -x509 -newkey rsa:4096 \
  -keyout ca-key.pem -out ca-cert.pem \
  -days 365 -nodes -subj "/CN=My CA"
