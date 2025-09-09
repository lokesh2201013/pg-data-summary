package models

import "time"


type ConnectionDetails struct {
	Host     string `json:"host"`
	Port     *int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}
type Summary struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name"`
    SyncedAt  time.Time `json:"synced_at"`
    SourceInfo ConnectionDetails `json:"source_info" gorm:"embedded;embeddedPrefix:source_"`
    Schemas   []Schema  `json:"schemas" gorm:"foreignKey:SummaryID;references:ID"`
}

type Schema struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    SummaryID string    `json:"summary_id" gorm:"index"`
    Name      string    `json:"name"`
    SyncedAt  time.Time `json:"synced_at"`

    Tables []Table `json:"tables" gorm:"foreignKey:SchemaID;references:ID"`
}

type Table struct {
    ID       string `json:"id" gorm:"primaryKey"`
    SchemaID string `json:"schema_id" gorm:"index"`
    Name     string  `json:"name"`
    RowCount int64   `json:"row_count"`
    SizeMB   float64 `json:"size_mb"`
}
