package model

import "time"

type OwnerPremium struct {
	Id             string     `gorm:"column:id;primaryKey" json:"id"`
	Idowner        *string    `gorm:"column:idowner" json:"idowner"`
	NamaLengkap    *string    `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	Telp           *string    `gorm:"column:telp" json:"telp"`
	Ktp            *string    `gorm:"column:ktp" json:"ktp"`
	Npwp           *string    `gorm:"column:npwp" json:"npwp"`
	Alamat         *string    `gorm:"column:alamat" json:"alamat"`
	FotoKtp        *string    `gorm:"column:foto_ktp" json:"foto_ktp"`
	FotoSelfieKtp  *string    `gorm:"column:foto_selfie_ktp" json:"foto_selfie_ktp"`
	FotoNpwp       *string    `gorm:"column:foto_npwp" json:"foto_npwp"`
	FotoTtd        *string    `gorm:"column:foto_ttd" json:"foto_ttd"`
	JenisKelamin   *string    `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	TempatLahir    *string    `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	TanggalLahir   *time.Time `gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	NamaIbuKandung *string    `gorm:"column:nama_ibu_kandung" json:"nama_ibu_kandung"`
	Approved       *string    `gorm:"column:approved" json:"approved"`
	ApprovedBy     *string    `gorm:"column:approved_by" json:"approved_by"`
	Reason         *string    `gorm:"column:reason" json:"reason"`
	ScheduledAt    *time.Time `gorm:"column:scheduled_at" json:"scheduled_at"`
	CreatedAt      *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"column:updated_at" json:"updated_at"`
	Hapus          *string    `gorm:"column:hapus" json:"hapus"`
}

func (OwnerPremium) TableName() string {
	return "owner_premium"
}
