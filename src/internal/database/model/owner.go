package model

import "time"

type Owner struct {
	Idowner                    string    `gorm:"column:idowner;primaryKey" json:"idowner"`
	Nama                       *string   `gorm:"column:nama" json:"nama"`
	Email                      *string   `gorm:"column:email" json:"email"`
	Password                   *string   `gorm:"column:password" json:"password"`
	Kota                       *string   `gorm:"column:kota" json:"kota"`
	Telp                       *string   `gorm:"column:telp" json:"telp"`
	Hapus                      *string   `gorm:"column:hapus" json:"hapus"`
	Salt                       *string   `gorm:"column:salt" json:"salt"`
	RememberToken              *string   `gorm:"column:remember_token" json:"remember_token"`
	VerifiedEmail              int64     `gorm:"column:verified_email" json:"verified_email"`
	VerifiedTelp               int64     `gorm:"column:verified_telp" json:"verified_telp"`
	KodeVerifikasiTelp         *string   `gorm:"column:kode_verifikasi_telp" json:"kode_verifikasi_telp"`
	UserOwnerIduserTakHapusRek *string   `gorm:"column:user_owner_iduser_tak_hapus_rek" json:"user_owner_iduser_tak_hapus_rek"`
	LevelAfiliasi              *string   `gorm:"column:level_afiliasi" json:"level_afiliasi"`
	KodeAfiliasi               string    `gorm:"column:kode_afiliasi" json:"kode_afiliasi"`
	KodeAfiliator              string    `gorm:"column:kode_afiliator" json:"kode_afiliator"`
	PersentaseAfiliator        string    `gorm:"column:persentase_afiliator" json:"persentase_afiliator"`
	SaldoAfiliator             *string   `gorm:"column:saldo_afiliator" json:"saldo_afiliator"`
	TutorNumber                *string   `gorm:"column:tutor_number" json:"tutor_number"`
	DataGenerated              *string   `gorm:"column:data_generated" json:"data_generated"`
	CreatedAt                  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt                  time.Time `gorm:"column:updated_at" json:"updated_at"`
	ExpiryTrialDate            string    `gorm:"column:expiry_trial_date" json:"expiry_trial_date"`
	Active                     int64     `gorm:"column:active" json:"active"`
	LastLogin                  *string   `gorm:"column:last_login" json:"last_login"`
	EmailReminders             *string   `gorm:"column:email_reminders" json:"email_reminders"`
	Versi                      string    `gorm:"column:versi" json:"versi"`
	DateReset                  string    `gorm:"column:date_reset" json:"date_reset"`
	SmsPremium                 *string   `gorm:"column:sms_premium" json:"sms_premium"`
	UploadToken                *string   `gorm:"column:upload_token" json:"upload_token"`
	LastKoinNotif              *int64    `gorm:"column:last_koin_notif" json:"last_koin_notif"`
	Idfotoprofil               *string   `gorm:"column:idfotoprofil" json:"idfotoprofil"`
	NamaUsaha                  *string   `gorm:"column:nama_usaha" json:"nama_usaha"`
	LogoUrl                    *string   `gorm:"column:logo_url" json:"logo_url"`
	IsPremium                  *string   `gorm:"column:is_premium" json:"is_premium"`
	Walkthrough                *string   `gorm:"column:walkthrough" json:"walkthrough"`
	Subscribe                  *string   `gorm:"column:subscribe" json:"subscribe"`
}

func (Owner) TableName() string {
	return "owner"
}
