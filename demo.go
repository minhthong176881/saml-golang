package main

import (
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"io/ioutil"

	"encoding/base64"
	"encoding/xml"

	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

func mainDemo() {
	// res, err := http.Get("http://idp.oktadev.com/metadata")
	// if err != nil {
	// 	panic(err)
	// }

	xmlFile, err := os.Open("metadata.xml")
	if err != nil {
		panic(err)
	}

	rawMetadata, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		panic(err)
	}

	metadata := &types.EntityDescriptor{}
	err = xml.Unmarshal(rawMetadata, metadata)
	if err != nil {
		panic(err)
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				panic(fmt.Errorf("metadata certificate(%d) must not be empty", idx))
			}
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				panic(err)
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				panic(err)
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}

	// We sign the AuthnRequest with a random key because Okta doesn't seem
	// to verify these.
	randomKeyStore := dsig.RandomKeyStoreForTest()

	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      metadata.IDPSSODescriptor.SingleSignOnServices[0].Location,
		IdentityProviderIssuer:      metadata.EntityID,
		ServiceProviderIssuer:       "urn:example:idp",
		AssertionConsumerServiceURL: "http://localhost:3001/msr/saml/auth",
		SignAuthnRequests:           true,
		AudienceURI:                 "urn:example:idp",
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  randomKeyStore,
	}

	http.HandleFunc("/msr/saml/auth", func(rw http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		assertionInfo, err := sp.RetrieveAssertionInfo(req.FormValue("SAMLResponse"))
		if err != nil {
			fmt.Println("here1")
			fmt.Println(err)
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		fmt.Printf("AssertionInfo: %+v\n", assertionInfo.WarningInfo)

		if assertionInfo.WarningInfo.InvalidTime {
			fmt.Println("Here2")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.NotInAudience {
			fmt.Println("Here3")
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		fmt.Fprintf(rw, "NameID: %s\n", assertionInfo.NameID)
		fmt.Printf("NameID: %s\n", assertionInfo.NameID)

		fmt.Fprintf(rw, "Assertions:\n")
		fmt.Println("Assertions:")

		for key, val := range assertionInfo.Values {
			fmt.Fprintf(rw, "  %s: %+v\n", key, val)
			fmt.Printf("  %s: %+v\n", key, val)
		}

		fmt.Fprintf(rw, "\n")

		fmt.Fprintf(rw, "Warnings:\n")
		fmt.Fprintf(rw, "%+v\n", assertionInfo.WarningInfo)
	})

	println("Visit this URL To Authenticate:")
	authURL, err := sp.BuildAuthURL("")
	if err != nil {
		panic(err)
	}

	println(authURL)

	println("Supply:")
	fmt.Printf("  SP ACS URL      : %s\n", sp.AssertionConsumerServiceURL)

	err = http.ListenAndServe(":3001", nil)
	if err != nil {
		panic(err)
	}
}