package validators

import (
	"math/big"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
)

// Require a non-empty string
func RequireString(field string, val string, errors *map[string]string, errorMsg string) {
	if len(val) <= 0 {
		(*errors)[field] = errorMsg
	}
}

// Require a non-empty string
func MaxLength(field string, val string, errors *map[string]string, errorMsg string, max int) {
	if len(val) > max {
		(*errors)[field] = errorMsg
	}
}

func MinLength(field string, val string, errors *map[string]string, errorMsg string, min int) {
	if len(val) < min {
		(*errors)[field] = errorMsg
	}
}

// Require a sane email, which means an @ sign with a . after it.
func RequireSaneEmail(field string, val string, errors *map[string]string, errorMsg string) {
	if !strings.Contains(val, "@") || strings.LastIndex(val, ".") < strings.LastIndex(val, "@") {
		(*errors)[field] = errorMsg
	}
}

// Require string to match a regexp
func RequireMatchRegexp(field string, val string, errors *map[string]string, errorMsg string, re *regexp.Regexp) {
	if !re.MatchString(val) {
		(*errors)[field] = errorMsg
	}
}

// Require string to be a number
func RequireNumber(field string, val string, errors *map[string]string, errorMsg string) {
	var r big.Rat
	if _, ok := r.SetString(val); !ok {
		(*errors)[field] = errorMsg
	}
}

// Require nofile or png, gif, jpeg type format
func RequireImageExtension(field string, header *multipart.FileHeader, errors *map[string]string, errorMsg string) {
	if header == nil {
		return
	}
	ext := filepath.Ext(filepath.Ext(header.Filename))
	valid := map[string]bool{".png": true, ".gif": true, ".jpg": true, ".jpeg": true}
	if _, ok := valid[ext]; !ok {
		(*errors)[field] = errorMsg
	}
}
