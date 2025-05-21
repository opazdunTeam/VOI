package logger

import (
	"backend/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger представляет собой обертку вокруг zap.Logger
type Logger struct {
	*zap.Logger
}

// NewLogger создает новый экземпляр логгера
func NewLogger(cfg *config.Config) (*Logger, error) {
	var zapConfig zap.Config

	if cfg.IsDevelopment() {
		// Конфигурация для разработки (более подробные логи)
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// Конфигурация для продакшена (структурированные JSON логи)
		zapConfig = zap.NewProductionConfig()
	}

	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		Logger: zapLogger.With(
			zap.String("app", cfg.App.Name),
			zap.String("mode", cfg.App.Mode),
		),
	}, nil
}

// Named создает именованный логгер для конкретного компонента
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		Logger: l.Logger.Named(name),
	}
}

// With добавляет постоянные поля к логгеру
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}
