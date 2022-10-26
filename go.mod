module openappsec.io/smartsync-shared-files

go 1.18

require (
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/google/wire v0.5.0
	openappsec.io/configuration v0.5.5
	openappsec.io/ctxutils v0.5.0
	openappsec.io/errors v0.7.0
	openappsec.io/health v0.2.1
	openappsec.io/httputils v0.11.1
	openappsec.io/log v0.9.0
	openappsec.io/tracer v0.5.0
)

require (
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/afero v1.8.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.11.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.9.0 // indirect
	golang.org/x/sys v0.0.0-20220513210249-45d2b4557a2a // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	openappsec.io/configuration => ./dependencies/openappsec.io/configuration
	openappsec.io/ctxutils => ./dependencies/openappsec.io/ctxutils
	openappsec.io/errors => ./dependencies/openappsec.io/errors
	openappsec.io/health => ./dependencies/openappsec.io/health
	openappsec.io/httputils => ./dependencies/openappsec.io/httputils
	openappsec.io/kafka => ./dependencies/openappsec.io/kafka
	openappsec.io/log => ./dependencies/openappsec.io/log
	openappsec.io/redis => ./dependencies/openappsec.io/redis
	openappsec.io/tracer => ./dependencies/openappsec.io/tracer
)
