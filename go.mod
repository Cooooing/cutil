module github.com/Cooooing/cutil

go 1.24

replace (
	github.com/Cooooing/cutil/common => ./common
	github.com/Cooooing/cutil/excel => ./excel
	github.com/Cooooing/cutil/query => ./query
	github.com/Cooooing/cutil/stream => ./stream
	github.com/Cooooing/cutil/timewheel => ./timewheel
)

require (
	github.com/Cooooing/cutil/common v0.0.0-20250826093745-1ed80c0aa9d7
	github.com/Cooooing/cutil/excel v0.0.0-20250825024611-76e9bd7621a7
	github.com/Cooooing/cutil/query v0.0.0-20250825024611-76e9bd7621a7
	github.com/Cooooing/cutil/stream v0.0.0-20250825024611-76e9bd7621a7
	github.com/Cooooing/cutil/timewheel v0.0.0-20250825024611-76e9bd7621a7
)

require (
	github.com/orcaman/concurrent-map/v2 v2.0.1 // indirect
	github.com/panjf2000/ants/v2 v2.11.3 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/tiendc/go-deepcopy v1.6.1 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/excelize/v2 v2.9.1 // indirect
	github.com/xuri/nfp v0.0.1 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.28.0 // indirect
)
