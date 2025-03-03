module github.com/tbellembois/gochimitheque-wasm

go 1.22.0

toolchain go1.23.3

replace github.com/tbellembois/gochimitheque v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque

replace gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1

require (
	github.com/go-ldap/ldap/v3 v3.4.6
	github.com/rocketlaunchr/react v1.0.9
	github.com/tbellembois/gochimitheque v0.0.0
	golang.org/x/text v0.14.0
	honnef.co/go/js/dom/v2 v2.0.0-20231112215516-51f43a291193
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20221128193559-754e69321358 // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.5 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
)
