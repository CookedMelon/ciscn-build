module jkscan

go 1.18

require (
	github.com/atotto/clipboard v0.1.4
	github.com/huin/asn1ber v0.0.0-20120622192748-af09f62e6358
	github.com/icodeface/tls v0.0.0-20190904083142-17aec93c60e5
	github.com/lcvvvv/appfinger v0.1.1

	//gonmap
	github.com/lcvvvv/gonmap v1.3.4
	github.com/lcvvvv/pool v0.0.0-00010101000000-000000000000
	github.com/lcvvvv/simplehttp v0.1.1 // indirect

	//grdp
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40
	github.com/miekg/dns v1.1.50 // indirect
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5

	//chinese
	golang.org/x/text v0.3.7 // indirect
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/lcvvvv/stdio v0.1.2 // indirect
	github.com/twmb/murmur3 v1.1.6 // indirect
	golang.org/x/mod v0.4.2 // indirect
	golang.org/x/net v0.0.0-20210916014120-12bc252f5db8 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.6-0.20210726203631-07bc1bf47fb2 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)

replace github.com/lcvvvv/pool => ./lib/pool

//replace github.com/lcvvvv/gonmap => ../go-github/gonmap
