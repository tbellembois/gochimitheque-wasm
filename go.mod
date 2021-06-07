module github.com/tbellembois/gochimitheque-wasm

go 1.16

replace (
	github.com/tbellembois/gochimitheque v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque
	github.com/tbellembois/gochimitheque-utils v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-utils
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
)

require (
	github.com/rocketlaunchr/react v1.0.9
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tbellembois/gochimitheque v0.0.0
	github.com/tbellembois/gochimitheque-utils v0.0.0
	golang.org/x/text v0.3.5
	honnef.co/go/js/dom/v2 v2.0.0-20200509013220-d4405f7ab4d8
)
