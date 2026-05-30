package logger

import "log/slog"

type Package struct {
	Logger *slog.Logger
}

func MakePackage(
	cfg Config,
) *Package {
	return &Package{
		Logger: NewLogger(cfg),
	}
}

func (p *Package) Error(err error) {
	p.Logger.Error("Error",
		slog.Any("error", err),
	)
}
