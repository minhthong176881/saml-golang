# saml-golang

Options:
- main.go:
  - Run:
      `openssl req -x509 -newkey rsa:2048 -keyout myservice.key -out myservice.cert -days 365 -nodes -subj "/CN=myservice.example.com" to create key and cert files`.
  - Run `go run main.go`
  - Get service provider metadata: 

      `mdpath=saml-test-$USER-$HOST.xml`      
      `curl localhost:8000/saml/metadata > $mdpath`
  - Upload mdpath file to https://samltest.id/upload.php
  - Browse to localhost:8000/hello
  
- demo.go:
  - Setup okta application
  - Run `go run demo.go`
  - Browse to localhost:3001
