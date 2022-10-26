//go:build wireinject
// +build wireinject

package ingector

// For further explanations about the wire and viper packages please refer to this repository Wiki page:
// https://openappsec.io/smartsync-shared-files/-/wikis/Go-Template#wire-and-dependency-injection
import (
	"context"

	"github.com/google/wire"
	"openappsec.io/smartsync-shared-files/internal/app"
	"openappsec.io/smartsync-shared-files/internal/app/drivers/http/rest"
	"openappsec.io/smartsync-shared-files/internal/app/sharedfiles"
	"openappsec.io/smartsync-shared-files/internal/pkg/filesdb/filesystem"
	"openappsec.io/configuration"
	"openappsec.io/configuration/viper"
	"openappsec.io/health"
)

// InitializeApp is an Adapter injector
func InitializeApp(ctx context.Context) (*app.App, error) {
	wire.Build(
		viper.NewViper,
		wire.Bind(new(configuration.Repository), new(*viper.Adapter)),

		configuration.NewConfigurationService,
		wire.Bind(new(rest.Configuration), new(*configuration.Service)),
		wire.Bind(new(app.Configuration), new(*configuration.Service)),
		wire.Bind(new(filesystem.Configuration), new(*configuration.Service)),

		sharedfiles.NewSharedFilesService,
		wire.Bind(new(rest.SharedFilesService), new(*sharedfiles.Service)),

		filesystem.NewAdapter,
		wire.Bind(new(sharedfiles.FileSystem), new(*filesystem.Adapter)),

		health.NewService,
		wire.Bind(new(rest.HealthService), new(*health.Service)),
		wire.Bind(new(app.HealthService), new(*health.Service)),

		rest.NewHTTPAdapter,
		wire.Bind(new(app.RestAdapter), new(*rest.Adapter)),

		app.NewApp,
	)

	return &app.App{}, nil
}
