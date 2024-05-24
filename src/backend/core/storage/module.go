package storage

import (
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
)

//var MinioModule = fx.Options(
//	fx.Provide(func() (*MinioConfig, error) {
//		cfg := MinioConfig{}
//		err := utils.ReadConfig(&cfg)
//		return &cfg, err
//	},
//		newMinioClient,
//	),
//)

var MinioModule = fx.Options(
	fx.Provide(
		func() (*MinioConfig, error) {
			cfg := MinioConfig{}
			err := utils.ReadConfig(&cfg)
			return &cfg, err
		},
		newMinioClient,
	),
)

func newMinioClient(cfg *MinioConfig) (MinIOClient, error) {
	if cfg.Mocked {
		return NewMinIOMockClient()
	}
	return NewMinIOClient(cfg)
}