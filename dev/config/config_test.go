package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorLoadConfig(t *testing.T) {
	const testConfigPath = "../ngaco"

	err := LoadConfig(testConfigPath)

	assert.Error(t, err)

	assert.Empty(t, Configuration)

}

func TestLoadConfig(t *testing.T) {
	const testConfigPath = "../"

	err := LoadConfig(testConfigPath)

	// Memastikan tidak ada error saat memuat konfigurasi
	assert.NoError(t, err, "Tidak seharusnya terjadi error saat memuat konfigurasi")

	// Memeriksa nilai yang diambil dari konfigurasi dengan nilai yang diharapkan
	expectedHost := "localhost"
	actualHost := Configuration.MongoDB.Host
	assert.Equal(t, expectedHost, actualHost, "Host harus sesuai")

}
