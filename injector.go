//go:build wireinject
// +build wireinject

package main

import (
	"permasalahanService/app"

	"permasalahanService/controller"
	"permasalahanService/repository"
	"permasalahanService/service"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

var permasalahanSet = wire.NewSet(
	repository.NewPermasalahanRepositoryImpl,
	wire.Bind(new(repository.PermasalahanRepository), new(*repository.PermasalahanRepositoryImpl)),
	service.NewPermasalahanServiceImpl,
	wire.Bind(new(service.PermasalahanService), new(*service.PermasalahanServiceImpl)),
	controller.NewPermasalahanControllerImpl,
	wire.Bind(new(controller.PermasalahanController), new(*controller.PermasalahanControllerImpl)),
)

var permasalahanTerpilihSet = wire.NewSet(
	repository.NewPermasalahanTerpilihRepositoryImpl,
	wire.Bind(new(repository.PermasalahanTerpilihRepository), new(*repository.PermasalahanTerpilihRepositoryImpl)),
	service.NewPermasalahanTerpilihServiceImpl,
	wire.Bind(new(service.PermasalahanTerpilihService), new(*service.PermasalahanTerpilihServiceImpl)),
	controller.NewPermasalahanTerpilihControllerImpl,
	wire.Bind(new(controller.PermasalahanTerpilihController), new(*controller.PermasalahanTerpilihControllerImpl)),
)

var isuStrategisSet = wire.NewSet(
	repository.NewIsuStrategisRepositoryImpl,
	wire.Bind(new(repository.IsuStrategisRepository), new(*repository.IsuStrategisRepositoryImpl)),
	service.NewIsuStrategisServiceImpl,
	wire.Bind(new(service.IsuStrategisService), new(*service.IsuStrategisServiceImpl)),
	controller.NewIsuStrategisControllerImpl,
	wire.Bind(new(controller.IsuStrategisController), new(*controller.IsuStrategisControllerImpl)),
)

func InitializedServer() *echo.Echo {
	wire.Build(
		app.GetConnection,
		wire.Value([]validator.Option{}),
		validator.New,
		permasalahanSet,
		permasalahanTerpilihSet,
		isuStrategisSet,
		app.NewRouter,
	)
	return nil
}
