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
	github.com/lcvvvv/simplehttp v0.1.1
	github.com/lcvvvv/stdio v0.1.2

	//grdp
	github.com/lunixbochs/struc v0.0.0-20200707160740-784aaebc1d40
	github.com/miekg/dns v1.1.50 // indirect
	golang.org/x/crypto v0.14.0

	//chinese
	golang.org/x/text v0.13.0 // indirect
)

require (
	github.com/PuerkitoBio/goquery v1.8.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/twmb/murmur3 v1.1.6 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
)

replace github.com/lcvvvv/pool => ./lib/pool

//replace github.com/lcvvvv/gonmap => ../go-github/gonmap
