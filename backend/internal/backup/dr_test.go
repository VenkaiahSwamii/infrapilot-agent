package backup

import (
	"testing"
)

func TestDRTestEngineInit(t *testing.T) {
	rest := NewRestoreEngine("./backups")
	if rest == nil {
		t.Fatalf("expected non-nil RestoreEngine")
	}

	pg := NewPostgresBackup()
	if pg == nil {
		t.Fatalf("expected non-nil PostgresBackup")
	}
}
