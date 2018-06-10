package morty

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"mime"
	"net/url"
	"path/filepath"
)

func mergeURIs(u1, u2 *url.URL) *url.URL {
	if u2 == nil {
		return u1
	}
	return u1.ResolveReference(u2)
}

// force content-disposition to attachment
func contentDispositionForceAttachment(contentDispositionBytes []byte, url *url.URL) []byte {
	var contentDispositionParams map[string]string
	if contentDispositionBytes != nil {
		var err error
		_, contentDispositionParams, err = mime.ParseMediaType(string(contentDispositionBytes))
		if err != nil {
			contentDispositionParams = make(map[string]string)
		}
	} else {
		contentDispositionParams = make(map[string]string)
	}
	_, fileNameDefined := contentDispositionParams["filename"]
	if !fileNameDefined {
		// TODO : sanitize filename
		contentDispositionParams["fileName"] = filepath.Base(url.Path)
	}
	return []byte(mime.FormatMediaType("attachment", contentDispositionParams))
}

func inArray(b []byte, a [][]byte) bool {
	for _, b2 := range a {
		if bytes.Equal(b, b2) {
			return true
		}
	}
	return false
}

func hash(msg string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func verifyRequestURI(uri, hashMsg, key []byte) bool {
	h := make([]byte, hex.DecodedLen(len(hashMsg)))
	_, err := hex.Decode(h, hashMsg)
	if err != nil {
		log.Println("hmac error:", err)
		return false
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(uri)
	return hmac.Equal(h, mac.Sum(nil))
}
