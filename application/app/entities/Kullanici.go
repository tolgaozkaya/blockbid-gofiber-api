package entities

import "time"

type User struct {
	KullaniciID  string    `json:"kullaniciID" validate:"required,uuid4"`
	KullaniciAdi string    `json:"kullaniciAdi" validate:"required,min=3,max=32"`
	Sifre        string    `json:"sifre" validate:"required,min=6"`
	Isim         string    `json:"isim" validate:"required"`
	Soyisim      string    `json:"soyisim" validate:"required"`
	Eposta       string    `json:"eposta" validate:"required,email"`
	Telefon      string    `json:"telefon" validate:"required"`
	TCKimlikNo   string    `json:"tcKimlikNo" validate:"required,len=11"`
	DogumTarihi  time.Time `json:"dogumTarihi" validate:"required"`
}

type CreateUserRequest struct {
	KullaniciAdi string    `json:"kullaniciAdi" validate:"required,min=3,max=32"`
	Sifre        string    `json:"sifre" validate:"required,min=6"`
	Isim         string    `json:"isim" validate:"required"`
	Soyisim      string    `json:"soyisim" validate:"required"`
	Eposta       string    `json:"eposta" validate:"required,email"`
	Telefon      string    `json:"telefon" validate:"required"`
	TCKimlikNo   string    `json:"tcKimlikNo" validate:"required,len=11"`
	DogumTarihi  time.Time `json:"dogumTarihi" validate:"required"` // YYYY-MM-DD formatında
}

type UpdateUserRequest struct {
	KullaniciID  string    `json:"kullaniciID" validate:"required,uuid4"` // Required to identify the user
	KullaniciAdi string    `json:"kullaniciAdi,omitempty" validate:"omitempty,min=3,max=32"`
	Sifre        string    `json:"sifre,omitempty" validate:"omitempty,min=6"`
	Isim         string    `json:"isim,omitempty" validate:"omitempty"`
	Soyisim      string    `json:"soyisim,omitempty" validate:"omitempty"`
	Eposta       string    `json:"eposta,omitempty" validate:"omitempty,email"`
	Telefon      string    `json:"telefon,omitempty" validate:"omitempty"`
	DogumTarihi  time.Time `json:"dogumTarihi" validate:"required"` // YYYY-MM-DD formatında
}

type QueryUserRequest struct {
	KullaniciID  string    `json:"kullaniciID,omitempty" validate:"omitempty,uuid4"`
	KullaniciAdi string    `json:"kullaniciAdi,omitempty" validate:"omitempty,min=3,max=32"`
	Isim         string    `json:"isim,omitempty" validate:"omitempty"`
	Soyisim      string    `json:"soyisim,omitempty" validate:"omitempty"`
	Eposta       string    `json:"eposta,omitempty" validate:"omitempty,email"`
	Telefon      string    `json:"telefon,omitempty" validate:"omitempty"`
	TCKimlikNo   string    `json:"tcKimlikNo,omitempty" validate:"omitempty,len=11"`
	DogumTarihi  time.Time `json:"dogumTarihi" validate:"required"` // YYYY-MM-DD formatında
}

type DeleteUserRequest struct {
	KullaniciID string `json:"kullaniciID" validate:"required,uuid4"`
}
