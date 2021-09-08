module github.com/tbellembois/gochimitheque-wasm

go 1.17

replace (
	github.com/tbellembois/gochimitheque v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque
	github.com/tbellembois/gochimitheque-utils v0.0.0 => /home/thbellem/workspace/workspace_go/src/github.com/tbellembois/gochimitheque-utils
	gopkg.in/russross/blackfriday.v2 => github.com/russross/blackfriday/v2 v2.0.1
)

require (
	github.com/rocketlaunchr/react v1.0.9
	github.com/tbellembois/gochimitheque v0.0.0
	github.com/tbellembois/gochimitheque-utils v0.0.0
	golang.org/x/text v0.3.6
	honnef.co/go/js/dom/v2 v2.0.0-20200509013220-d4405f7ab4d8
)

require (
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
)
