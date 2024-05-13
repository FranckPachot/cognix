package ai

import (
	"cognix.ch/api/v2/core/proto"
	"cognix.ch/api/v2/core/utils"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var ChunkingModule = fx.Options(
	fx.Provide(func() (*ChunkingConfig, error) {
		cfg := ChunkingConfig{}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		if err := cfg.Validate(); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		newChunking,
	),
)

func newChunking(cfg *ChunkingConfig) Chunking {
	if cfg.Strategy == StrategyLLM {
		return NewLLMChunking()
	}
	return NewStaticChunking(cfg)
}

var EmbeddingModule = fx.Options(
	fx.Provide(func() (*EmbeddingConfig, error) {
		cfg := EmbeddingConfig{}
		if err := utils.ReadConfig(&cfg); err != nil {
			return nil, err
		}
		if err := cfg.Validate(); err != nil {
			return nil, err
		}
		return &cfg, nil
	},
		newEmbeddingGRPCClient),
)

func newEmbeddingGRPCClient(cfg *EmbeddingConfig) (proto.EmbeddServiceClient, error) {
	conn, err := grpc.Dial(cfg.EmbeddingURL)
	if err != nil {
		return nil, err
	}
	return proto.NewEmbeddServiceClient(conn), nil
}
