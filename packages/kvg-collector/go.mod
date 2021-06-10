module github.com/kiel-live/kiel-live/packages/kvg-collector

go 1.16

replace (
	github.com/kiel-live/kiel-live/packages/client => ../client
	github.com/kiel-live/kiel-live/packages/pub-sub-proto => ../pub-sub-proto
)

require (
	github.com/go-co-op/gocron v1.6.2
	github.com/kiel-live/kiel-live/packages/client v2.0.0+incompatible
	github.com/kiel-live/kiel-live/packages/pub-sub-proto v2.0.0+incompatible
	github.com/thoas/go-funk v0.8.0
)
