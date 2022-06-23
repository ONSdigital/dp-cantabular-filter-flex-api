module github.com/ONSdigital/dp-cantabular-filter-flex-api

go 1.18

// The following replacement is for module versions that have known vulnerabilities (from nancy sleuth)
replace github.com/spf13/cobra => github.com/spf13/cobra v1.4.0

// to avoid 'sonatype-2021-4899' non-CVE Vulnerability
exclude github.com/gorilla/sessions v1.2.1

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.143.0
	github.com/ONSdigital/dp-authorisation v0.2.0
	github.com/ONSdigital/dp-component-test v0.7.0
	github.com/ONSdigital/dp-healthcheck v1.3.0
	github.com/ONSdigital/dp-kafka/v3 v3.3.2
	github.com/ONSdigital/dp-mongodb/v3 v3.0.2
	github.com/ONSdigital/dp-net v1.4.1
	github.com/ONSdigital/dp-net/v2 v2.4.0
	github.com/ONSdigital/log.go/v2 v2.2.0
	github.com/cucumber/godog v0.12.5
	github.com/go-chi/chi/v5 v5.0.7
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.8
	github.com/google/uuid v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/maxcnunes/httpfake v1.2.4
	github.com/pkg/errors v0.9.1
	github.com/rdumont/assistdog v0.0.0-20201106100018-168b06230d14
	github.com/smartystreets/goconvey v1.7.2
	go.mongodb.org/mongo-driver v1.9.1
)

require (
	github.com/ONSdigital/dp-api-clients-go v1.43.0 // indirect
	github.com/ONSdigital/dp-mongodb-in-memory v1.3.1 // indirect
	github.com/ONSdigital/dp-rchttp v1.0.0 // indirect
	github.com/ONSdigital/go-ns v0.0.0-20210916104633-ac1c1c52327e // indirect
	github.com/Shopify/sarama v1.34.1 // indirect
	github.com/aws/aws-sdk-go v1.44.40 // indirect
	github.com/chromedp/cdproto v0.0.0-20220620034043-93c6adfcda6a // indirect
	github.com/chromedp/chromedp v0.8.2 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/go-avro/avro v0.0.0-20171219232920-444163702c11 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.3 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.0.0 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.2 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/justinas/alice v1.2.0 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/shurcooL/graphql v0.0.0-20220606043923-3cf50f8a0a29 // indirect
	github.com/smartystreets/assertions v1.13.0 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/square/mongo-lock v0.0.0-20220601164918-701ecf357cd7 // indirect
	github.com/stretchr/testify v1.7.4 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.0.0-20220622184535-263ec571b305 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220622161953-175b2fd9d664 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
