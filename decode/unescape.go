package decode

import "html"

func unescapeHtml(source string) (result string) {
	return html.UnescapeString(source)
}
