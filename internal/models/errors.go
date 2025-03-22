package models

import "errors"

var ErrRecordNotFound error = errors.New("Kayıt Bulunamadi")

var ErrTokenIsExpired error = errors.New("Token süresi dolmuş")

var ErrInvalidToken error = errors.New("Geçersiz token")

var ErrMissingAuthorization error = errors.New("Authorization başlığı eksik")

var ErrKeyIsEmpty error = errors.New("Key Boş")
