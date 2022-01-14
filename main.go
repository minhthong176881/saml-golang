package main

// import (
// 	"context"
// 	"crypto/rsa"
// 	"crypto/tls"
// 	"crypto/x509"
// 	"fmt"
// 	"net/http"
// 	"net/url"

// 	"github.com/crewjam/saml/samlsp"
// )

// func hello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello, %s!", samlsp.AttributeFromContext(r.Context(), "cn"))
// }

// func main() {
// 	keyPair, err := tls.LoadX509KeyPair("myservice.cert", "myservice.key") 
// 	if err != nil {
// 		panic(err)
// 	}
// 	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
// 	if err != nil {
// 		panic(err)
// 	}

// 	idpMetadataURL, err := url.Parse("https://dev-9788402.okta.com/app/dev-9788402_samlgolang_1/exk3lamm4nTIk1cFK5d7/sso/saml")
// 	if err != nil {
// 		panic(err)
// 	}
// 	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient, *idpMetadataURL)
// 	if err != nil {
// 		panic(err)
// 	}

// 	rootURL, err := url.Parse("http://localhost:8000")
// 	if err != nil {
// 		panic(err)
// 	}

// 	samlSP, _ := samlsp.New(samlsp.Options{
// 		URL:         *rootURL,
// 		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
// 		Certificate: keyPair.Leaf,
// 		IDPMetadata: idpMetadata,
// 	})
// 	app := http.HandlerFunc(hello)
// 	http.Handle("/hello", samlSP.RequireAccount(app))
// 	http.Handle("/saml/", samlSP)
// 	http.ListenAndServe(":8000", nil)
// }
