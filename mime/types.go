package mime

import (
	"net/http"
)

const (
	Plain = "text/plain"
	HTML  = "text/html"
	XML2  = "text/xml"
	XML   = "application/xml"
	XHTML = "application/xhtml+xml"

	FormData = "multipart/form-data"
	GIF      = "image/gif"
	JPEG     = "image/jpeg"
	PNG      = "image/png"
	WEBP     = "image/webp"

	MPEG = "video/mpeg"

	FormURLEncoded = "application/x-www-form-urlencoded"
	OctetStream    = "application/octet-stream"
	JSON           = "application/json"
	PDF            = "application/pdf"
	MSWord         = "application/msword"
	GZIP           = "application/x-gzip"
)

const (
	ContentType        = "Content-Type"
	ContentDisposition = "Content-Disposition"
	CharsetUTF8        = "charset=utf-8"

	charsetSuffix = "; " + CharsetUTF8

	PlainUTF8 = Plain + charsetSuffix

	// Hope this style is better than HTMLUTF8, etc.
	HtmlUTF8 = HTML + charsetSuffix
	JsonUTF8 = JSON + charsetSuffix
	XmlUTF8  = XML + charsetSuffix
)

func GetContentType(h http.Header) string {
	t := h.Get(ContentType)
	for i, ch := range t {
		if ch == ' ' || ch == ';' {
			t = t[:i]
			break
		}
	}
	return t
}
