package decode

import "html"

// 数値文字参照をデコードする
func unescapeHtml(source string) (result string) {
	return html.UnescapeString(source)
}
