package logger

import "log/slog"

type LogRecord struct {
	logger *slog.Logger

	isDiscarded bool

	args []any
}

func NewLogRecord(logger *slog.Logger) *LogRecord {
	return &LogRecord{
		logger:      logger,
		isDiscarded: false,
	}
}

func (r *LogRecord) AddAttrs(attrs ...slog.Attr) {
	if r.isDiscarded {
		return
	}

	for _, attr := range attrs {
		r.args = append(r.args, attr)
	}
}

func (r *LogRecord) Info(msg string, args ...any) {
	if r.isDiscarded {
		return
	}

	res := make([]any, 0, len(r.args)+len(args))
	res = append(res, r.args...)
	res = append(res, args...)

	r.logger.Info(msg, res...)
}

func (r *LogRecord) Discard() {
	r.isDiscarded = true
	r.logger = NewDiscard()
	r.args = nil
}
