module github.com/tbellembois/gochimitheque-wasm

go 1.19

replace (
	github.com/tbellembois/gochimitheque v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque
	github.com/tbellembois/gochimitheque-utils v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-utils
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
)

require (
	github.com/go-ldap/ldap/v3 v3.4.4
	github.com/rocketlaunchr/react v1.0.9
	github.com/tbellembois/gochimitheque v0.0.0
	github.com/tbellembois/gochimitheque-utils v0.0.0
	golang.org/x/text v0.3.7
	honnef.co/go/js/dom/v2 v2.0.0-20200509013220-d4405f7ab4d8
)

require (
	github.com/Azure/go-ntlmssp v0.0.0-20220621081337-cb9428e4ac1e // indirect
	github.com/go-asn1-ber/asn1-ber v1.5.4 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
)
