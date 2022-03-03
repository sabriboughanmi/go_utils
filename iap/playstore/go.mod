module github.com/sabriboughanmi/go_utils/iap/playstore

go 1.16

require (
	github.com/sabriboughanmi/go_utils/iap/playstore/androidpublisher v0.0.0-20220303110457-40385dcba362
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b
	google.golang.org/api v0.70.0
)

replace github.com/sabriboughanmi/go_utils/iap/playstore/androidpublisher => ./androidpublisher
