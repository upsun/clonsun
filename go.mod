module github.com/upsun/clonsun

go 1.22.6

require (
	github.com/spf13/pflag v1.0.5
	github.com/upsun/convsun v0.3.2
	github.com/upsun/lib-sun v0.3.4
)

require (
	github.com/getsentry/sentry-go v0.28.1 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/upsun/lib-sun => ../lib-sun
//replace github.com/upsun/convsun => ../convsun
