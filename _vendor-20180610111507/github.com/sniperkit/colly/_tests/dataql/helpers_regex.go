package main

import (
	"regexp"
)

// reIsNumeric check if the string contains only numbers. Empty string is valid.
func reIsNumeric(val string) bool {
	if len(val) == 0 {
		return true
	}
	regexNumeric := regexp.MustCompile("^[0-9]+$")
	return regexNumeric.MatchString(val)
}

// reIsFloat check if the string is a float.
func reIsFloat(val string) bool {
	regexFloat := regexp.MustCompile("^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$")
	return regexFloat.MatchString(val)
}

// reIsAlpha check if the string contains only letters. Empty string is valid.
func reIsAlpha(val string) bool {
	regexAlpha := regexp.MustCompile("^[a-zA-Z]+$")
	return regexAlpha.MatchString(val)
}

// IsAlphanumeric check if the string contains only letters and numbers. Empty string is valid.
func reIsAlphanumeric(val string) bool {
	if len(val) == 0 {
		return true
	}
	regexAlphanumeric := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return regexAlphanumeric.MatchString(val)
}

// reIsAlphanumeric check if the string contains only letters, numbers, space, hypens, and underscore.
// Only allow letters and numbers at the start and the end.
// Empty string is valid.
func reIsAlphanumSpaceHyphenUnderscore(val string) bool {
	if len(val) == 0 {
		return true
	}
	regex := regexp.MustCompile("^[a-zA-Z0-9]+[a-zA-Z0-9-_ ]*[a-zA-Z0-9]$")
	return regex.MatchString(val)
}

func IsEmail(vdata string) bool {
	email := regexp.MustCompile(`^[\w\.\_]{2,}@([\w\-]+\.){1,}\.[a-z]$`)
	return email.MatchString(vdata)
}

func IsPhone(vdata string) bool {
	phone := regexp.MustCompile(`^(\+86)?1[3-9][0-9]{9}$`)
	return phone.MatchString(vdata)
}

func IsNumeric(vdata string) bool {
	numeric := regexp.MustCompile(`^[0-9\.]{1,}$`)
	return numeric.MatchString(vdata)
}

func IsAlphaNumeric(vdata string) bool {
	alpha := regexp.MustCompile(`^[a-zA-Z0-9]{1,}$`)
	return alpha.MatchString(vdata)
}

func IsIp(vdata string) bool {
	ip := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`)
	return ip.MatchString(vdata)
}

func IsDateTime(vdata string) bool {
	dttmrex := regexp.MustCompile(`^[0-9]{4}\-[0-9]{2}\-[0-9]{2}[\s]{1,4}[0-9]{1,2}\:[0-9]{1,2}\:[0-9]{1,2}$`)
	return dttmrex.MatchString(vdata)
}

func IsDate(vdata string) bool {
	date := regexp.MustCompile(`^[0-9]{4}\-[0-9]{2}\-[0-9]{2}$`)
	return date.MatchString(vdata)
}

func IsTime(vdata string) bool {
	tmrex := regexp.MustCompile(`^[0-9]{1,2}\:[0-9]{1,2}\:[0-9]{1,2}$`)
	return tmrex.MatchString(vdata)
}

func IsIdCard(vdata string) bool {
	idcard := regexp.MustCompile(`(^\d{15}$)|(^\d{17}[\d|x|X]$)`)
	return idcard.MatchString(vdata)
}

func IsUserName(vdata string) bool {
	user := regexp.MustCompile(`^[\w\@\-\.]{3,}$`)
	return user.MatchString(vdata)
}

func IsPasswd(vdata string) bool { //要求六位以上且还有英文和字母
	isnull := regexp.MustCompile(`^[^\s]{6,}$`)
	isalpha := regexp.MustCompile(`[a-zA-Z]`)
	isnumeric := regexp.MustCompile(`[0-9]`)
	return isnull.MatchString(vdata) && isnumeric.MatchString(vdata) && isalpha.MatchString(vdata)
}

func IsChinese(vdata string) bool {
	chinese := regexp.MustCompile(`^[\u4e00-\u9fa5]{1,}$`)
	return chinese.MatchString(vdata)
}
