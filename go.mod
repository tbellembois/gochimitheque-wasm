module github.com/tbellembois/gochimitheque-wasm

go 1.16

replace (
	github.com/tbellembois/gochimitheque v2.0.6+incompatible => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
)

require (
	github.com/rocketlaunchr/react v1.0.9
	github.com/tbellembois/gochimitheque v2.0.6+incompatible
	github.com/tbellembois/gochimitheque-utils v0.0.0-20210130182605-75000d2f72a5
	golang.org/x/text v0.3.5
	honnef.co/go/js/dom/v2 v2.0.0-20200509013220-d4405f7ab4d8
)
