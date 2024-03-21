package model

import "time"

type Revenue struct {
	Date                 time.Time `gorm:"autoCreateTime" json:"-" `
	TotalPendapatanKotor float64   `json:"topor"`
	ModalAwal            float64   `json:"modal"`
	PendapatanBersih     float64   `json:"perih"`
	MarginKeuntungan     float64   `json:"maun"`
}
