package storage

import (
	"fmt"
	"os"
	"sync"
)

// WAL represents a Write-Ahead Log for durability and crash recovery. Simply a hard disk storage.
type WAL struct {
	mu   sync.Mutex
	file *os.File
}

// NewWAL opens the log file, If it doesnt exists, it creates it.
func NewWAL(filePath string) (*WAL, error) {
	// os.O_APPEND: We want to append to the file instead of overwriting it.
	// os.O_SYNC asks the OS to write to disk immediatly.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %v", err)
	}
	return &WAL{file: file}, nil
}

// Append writes a new command to the hard drive BEFORE it goes into RAM.
func (w *WAL) Append(key, value string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Format the log entry with simple comma-separated strings.
	entry := fmt.Sprintf("PUT, %s, %s\n", key, value)

	// Writing the string to the file buffer.
	if _, err := w.file.WriteString(entry); err != nil {
		return fmt.Errorf("failed to write to WAL: %v", err)
	}

	// CRITICAL: fsync. This forces the operating system to flush the buffer
	// and actually spin the physical hard drive platter to save the data.
	return w.file.Sync()
}
