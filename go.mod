module github.com/tbellembois/gochimitheque-wasm

go 1.18

replace (
	github.com/tbellembois/gochimitheque v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque
	github.com/tbellembois/gochimitheque-utils v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-utils
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
)

require (
	github.com/go-ldap/ldap/v3 v3.4.2
	github.com/rocketlaunchr/react v1.0.9
	github.com/tbellembois/gochimitheque v0.0.0
	github.com/tbellembois/gochimitheque-utils v0.0.0
	golang.org/x/text v0.3.7
	honnef.co/go/js/dom/v2 v2.0.0-20200509013220-d4405f7ab4d8
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20211209120228-48547f28849e // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.3 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20220224003255-dbe011f71a99 // indirect
)
