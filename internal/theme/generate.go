package theme

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"regexp"
	"strings"
)

const (
	dstThemeFile = "./theme_gen.go"
	dstFontFile  = "./font_gen.go"
)

func Generate(t *Setting) (string, string, error) {
	themeFile, err := generateTheme(t)
	if err != nil {
		return "", "", err
	}
	fontFile, err := generateFont(t)
	if err != nil {
		return "", "", err
	}
	return themeFile, fontFile, nil
}

func generateSource(source []byte, dstFile string) (string, error) {
	dst, err := os.Create(dstFile)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	dst.Write(source)
	return dstFontFile, nil
}

func generateTheme(t *Setting) (string, error) {
	source, err := buildThemeSource(t)
	if err != nil {
		return "", err
	}
	return generateSource(source, dstThemeFile)
}

func buildThemeSource(t *Setting) ([]byte, error) {
	buf := newBufferWrapper()

	buf.writeln("// Code generated by fyne-theme-generator")
	buf.writeln("")
	buf.writeln("package %s", t.packageName)
	buf.writeln("")
	buf.writeln("")
	buf.writeln("import (")
	buf.writeln("\"image/color\"")
	buf.writeln("")
	buf.writeln("\"fyne.io/fyne\"")
	buf.writeln("\"fyne.io/fyne/theme\"")
	buf.writeln(")")
	buf.writeln("")
	buf.writeln("type %s struct{}", t.themeStructName)
	buf.writeln("")
	buf.writeln("func (%s) BackgroundColor() color.Color      { return %#v }", t.themeStructName, t.backgroundColor)
	buf.writeln("func (%s) ButtonColor() color.Color          { return %#v }", t.themeStructName, t.buttonColor)
	buf.writeln("func (%s) DisabledButtonColor() color.Color  { return %#v }", t.themeStructName, t.disabledButtonColor)
	buf.writeln("func (%s) TextColor() color.Color            { return %#v }", t.themeStructName, t.textColor)
	buf.writeln("func (%s) DisabledTextColor() color.Color    { return %#v }", t.themeStructName, t.disabledTextColor)
	buf.writeln("func (%s) IconColor() color.Color            { return %#v }", t.themeStructName, t.iconColor)
	buf.writeln("func (%s) DisabledIconColor() color.Color    { return %#v }", t.themeStructName, t.disabledIconColor)
	buf.writeln("func (%s) HyperlinkColor() color.Color       { return %#v }", t.themeStructName, t.hyperlinkColor)
	buf.writeln("func (%s) PlaceHolderColor() color.Color     { return %#v }", t.themeStructName, t.placeHolderColor)
	buf.writeln("func (%s) PrimaryColor() color.Color         { return %#v }", t.themeStructName, t.primaryColor)
	buf.writeln("func (%s) HoverColor() color.Color           { return %#v }", t.themeStructName, t.hoverColor)
	buf.writeln("func (%s) FocusColor() color.Color           { return %#v }", t.themeStructName, t.focusColor)
	buf.writeln("func (%s) ScrollBarColor() color.Color       { return %#v }", t.themeStructName, t.scrollBarColor)
	buf.writeln("func (%s) ShadowColor() color.Color          { return %#v }", t.themeStructName, t.shadowColor)
	buf.writeln("func (%s) TextSize() int                     { return %#v }", t.themeStructName, t.textSize)
	// Note: Currently, font cannot be changed from this application.
	buf.writeln("func (%s) TextFont() fyne.Resource           { return theme.LightTheme().TextFont() }", t.themeStructName)
	buf.writeln("func (%s) TextBoldFont() fyne.Resource       { return theme.LightTheme().TextBoldFont() }", t.themeStructName)
	buf.writeln("func (%s) TextItalicFont() fyne.Resource     { return theme.LightTheme().TextItalicFont() }", t.themeStructName)
	buf.writeln("func (%s) TextBoldItalicFont() fyne.Resource { return theme.LightTheme().TextBoldItalicFont() }", t.themeStructName)
	buf.writeln("func (%s) TextMonospaceFont() fyne.Resource  { return theme.LightTheme().TextMonospaceFont() }", t.themeStructName)
	buf.writeln("func (%s) Padding() int                      { return %#v }", t.themeStructName, t.padding)
	buf.writeln("func (%s) IconInlineSize() int               { return %#v }", t.themeStructName, t.iconInlineSize)
	buf.writeln("func (%s) ScrollBarSize() int                { return %#v }", t.themeStructName, t.scrollBarSize)
	buf.writeln("func (%s) ScrollBarSmallSize() int           { return %#v }", t.themeStructName, t.scrollBarSmallSize)

	return format.Source(buf.Bytes())
}

func generateFont(t *Setting) (string, error) {
	source, err := buildFontSource(t)
	if err != nil {
		return "", err
	}
	return generateSource(source, dstFontFile)
}

func buildFontSource(t *Setting) ([]byte, error) {
	buf := newBufferWrapper()

	buf.writeln("// Code generated by fyne-theme-generator; DO NOT EDIT.")
	buf.writeln("")
	buf.writeln("package %s", t.packageName)
	buf.writeln("")
	buf.writeln("import \"fyne.io/fyne\"")
	buf.writeln("")
	buf.writeln("var %s = %#v\n", sanitiseName(t.textFont.Name()), t.textFont)
	buf.writeln("")
	buf.writeln("var %s = %#v\n", sanitiseName(t.textBoldFont.Name()), t.textBoldFont)
	buf.writeln("")
	buf.writeln("var %s = %#v\n", sanitiseName(t.textItalicFont.Name()), t.textItalicFont)
	buf.writeln("")
	buf.writeln("var %s = %#v\n", sanitiseName(t.textBoldItalicFont.Name()), t.textBoldItalicFont)
	buf.writeln("")
	buf.writeln("var %s = %#v\n", sanitiseName(t.textMonospaceFont.Name()), t.textMonospaceFont)
	buf.writeln("")

	return buf.Bytes(), nil
}

func sanitiseName(file string) string {
	titled := strings.Title(file)

	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	name := reg.ReplaceAllString(titled, "")

	return "font" + name
}

type bufferWrapper struct {
	*bytes.Buffer
}

func newBufferWrapper() *bufferWrapper {
	return &bufferWrapper{&bytes.Buffer{}}
}

func (b *bufferWrapper) writeln(s string, a ...interface{}) {
	b.WriteString(fmt.Sprintf(s+"\n", a...))
}
