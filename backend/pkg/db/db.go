package db

import (
	"database/sql"
	"fmt"
	"time"

	"backend/pkg/config"
	"backend/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Database представляет обертку для соединения с БД
type Database struct {
	DB     *gorm.DB
	Config *config.Config
	Logger *logger.Logger
}

// NewDatabase создает новое подключение к базе данных
func NewDatabase(cfg *config.Config, log *logger.Logger) (*Database, error) {
	log = log.Named("database")
	log.Info("Initializing database connection")

	if cfg.IsDevelopment() {
		// В режиме разработки автоматически создаем базу данных, если её нет
		if err := ensureDatabaseExists(cfg, log); err != nil {
			return nil, fmt.Errorf("failed to ensure database exists: %w", err)
		}
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Используем имена таблиц в единственном числе
		},
		// Другие настройки GORM
		PrepareStmt:            true, // Кешировать подготовленные запросы
		SkipDefaultTransaction: true, // Отключаем транзакции по умолчанию для повышения производительности
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info("Database connection established")

	return &Database{
		DB:     db,
		Config: cfg,
		Logger: log,
	}, nil
}

// Close закрывает соединение с базой данных
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	return sqlDB.Close()
}

// AutoMigrate запускает автоматическую миграцию для моделей
func (d *Database) AutoMigrate(models ...interface{}) error {
	d.Logger.Info("Running database migrations")

	if err := d.DB.AutoMigrate(models...); err != nil {
		return fmt.Errorf("database migration failed: %w", err)
	}

	d.Logger.Info("Database migrations completed")
	return nil
}

// ensureDatabaseExists проверяет существование БД и создает её при отсутствии
func ensureDatabaseExists(cfg *config.Config, log *logger.Logger) error {
	// Только для dev-режима
	if !cfg.IsDevelopment() {
		return nil
	}

	log.Info("Checking if database exists")

	// Подключаемся к postgres для создания нашей БД
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=%s dbname=postgres",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	// Проверяем существование БД
	var exists bool
	query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", cfg.Database.Name)
	err = db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	// Создаем БД если не существует
	if !exists {
		log.Info("Creating database", zap.String("name", cfg.Database.Name))
		_, err = db.Exec("CREATE DATABASE " + cfg.Database.Name)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Info("Database created successfully")
	} else {
		log.Debug("Database already exists")
	}

	return nil
}
