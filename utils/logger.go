package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logDir string, debug bool) (*zap.Logger, error) {
	// Buat directory logs jika belum ada
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Generate filename dengan format: app-2026-01-04.log
	filename := fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))
	logPath := filepath.Join(logDir, filename)

	// Encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	if debug {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Set format log
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	// File sink dengan rotasi log harian
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,    // MB - maksimal ukuran file sebelum rotasi
		MaxBackups: 30,    // Simpan 30 file backup
		MaxAge:     30,    // days - hapus file lebih dari 30 hari
		Compress:   true,  // Compress old log files
		LocalTime:  true,  // Gunakan waktu lokal
	})

	// Console sink (stdout)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// Set log level
	var logLevel zapcore.Level
	if debug {
		logLevel = zap.DebugLevel
	} else {
		logLevel = zap.InfoLevel
	}

	// Gabungkan file dan console output
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileWriter, logLevel),    // Log ke file (JSON format)
		zapcore.NewCore(consoleEncoder, consoleWriter, logLevel), // Log ke console (readable format)
	)

	// Buat logger dengan options
	logger := zap.New(core, 
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	
	logger.Info("Logger initialized", 
		zap.String("log_file", logPath),
		zap.Bool("debug_mode", debug),
	)

	return logger, nil
}
