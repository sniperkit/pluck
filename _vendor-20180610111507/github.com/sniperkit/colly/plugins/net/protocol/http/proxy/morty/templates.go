package morty

import (
	"encoding/base64"
)

func init() {
	FaviconBase64 := "iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII"
	FAVICON_BYTES, _ = base64.StdEncoding.DecodeString(FaviconBase64)
}

var FAVICON_BYTES []byte

var (
	HTML_FORM_EXTENSION string = `<input type="hidden" name="mortyurl" value="%s" /><input type="hidden" name="mortyhash" value="%s" />`

	HTML_BODY_EXTENSION string = `
<input type="checkbox" id="mortytoggle" autocomplete="off" />
<div id="mortyheader">
  <p>This is a <a href="https://github.com/asciimoo/morty">proxified and sanitized</a> view of the page,<br />visit <a href="%s" rel="noreferrer">original site</a>.</p><p><label for="mortytoggle">hide</label></p>
</div>
<style>
#mortyheader { position: fixed; margin: 0; box-sizing: border-box; -webkit-box-sizing: border-box; top: 15%%; left: 0; max-width: 140px; overflow: hidden; z-index: 2147483647 !important; font-size: 12px; line-height: normal; border-width: 4px 4px 4px 0; border-style: solid; border-color: #1abc9c; background: #FFF; padding: 12px 12px 8px 8px; color: #444; }
#mortyheader * { box-sizing: content-box; margin: 0; border: none; padding: 0; overflow: hidden; z-index: 2147483647 !important; line-height: 1em; font-size: 12px !important; font-family: sans !important; font-weight: normal; text-align: left; text-decoration: none; }
#mortyheader p { padding: 0 0 0.7em 0; display: block; }
#mortyheader a { color: #3498db; font-weight: bold; display: inline; }
#mortyheader label { text-align: right; cursor: pointer; display: block; color: #444; }
input[type=checkbox]#mortytoggle { display: none; }
input[type=checkbox]#mortytoggle:checked ~ div { display: none; visibility: hidden; }
</style>
`

	HTML_HEAD_CONTENT_TYPE string = `<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="referrer" content="no-referrer">
`

	MORTY_HTML_PAGE_START string = `<!doctype html>
<html>
<head>
<title>MortyProxy</title>
<meta name="viewport" content="width=device-width, initial-scale=1 , maximum-scale=1.0, user-scalable=1" />
<style>
html { height: 100%; }
body { min-height : 100%; display: flex; flex-direction:column; font-family: 'Garamond', 'Georgia', serif; text-align: center; color: #444; background: #FAFAFA; margin: 0; padding: 0; font-size: 1.1em; }
input { border: 1px solid #888; padding: 0.3em; color: #444; background: #FFF; font-size: 1.1em; }
input[placeholder] { width:80%; }
a { text-decoration: none; #2980b9; }
h1, h2 { font-weight: 200; margin-bottom: 2rem; }
h1 { font-size: 3em; }
.container { flex:1; min-height: 100%; margin-bottom: 1em; }
.footer { margin: 1em; }
.footer p { font-size: 0.8em; }
</style>
</head>
<body>
	<div class="container">
		<h1>MortyProxy</h1>
`

	MORTY_HTML_PAGE_END string = `
	</div>
	<div class="footer">
		<p>Morty rewrites web pages to exclude malicious HTML tags and CSS/HTML attributes. It also replaces external resource references to prevent third-party information leaks.<br />
		<a href="https://github.com/asciimoo/morty">view on github</a>
		</p>
	</div>
</body>
</html>`
)
