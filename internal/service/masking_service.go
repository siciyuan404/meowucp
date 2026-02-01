package service

import (
	"strings"
	"unicode"
)

func MaskEmail(email string) string {
	if email == "" {
		return ""
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***"
	}
	username := parts[0]
	domain := parts[1]
	if len(username) <= 2 {
		username = username
	} else {
		username = string(username[0]) + strings.Repeat("*", 1) + string(username[len(username)-1])
	}
	return username + "@" + domain
}

func MaskPhone(phone string) string {
	if phone == "" {
		return ""
	}
	runes := []rune(phone)
	length := len(runes)
	if length <= 4 {
		return strings.Repeat("*", length)
	}
	visibleSuffix := 4
	if length < 8 {
		visibleSuffix = 2
	}
	return string(runes[:2]) + "****" + string(runes[length-visibleSuffix:])
}

func MaskCreditCard(card string) string {
	if card == "" {
		return ""
	}
	digits := make([]rune, 0, len(card))
	for _, r := range card {
		if unicode.IsDigit(r) {
			digits = append(digits, r)
		}
	}
	if len(digits) <= 4 {
		return string(digits)
	}
	return strings.Repeat("*", len(digits)-4) + string(digits[len(digits)-4:])
}
