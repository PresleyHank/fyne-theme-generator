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
	if !t.needToGenerateFont() {
		return themeFile, "", nil
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
	return dstFile, nil
}

func generateTheme(t *Setting) (string, error) {
	source, err := buildThemeSource(t)
	if err != nil {
		return "", err
	}
	return generateSource(source, dstThemeFile)
}

func buildThemeSource(t *Setting) ([]byte, error) {
	if t.ExportForV2() {
		return buildThemeSourceForV2(t)
	}
	return buildThemeSourceForV1(t)
}

func buildThemeSourceForV2(t *Setting) ([]byte, error) {
	buf := newBufferWrapper()

	buf.writeln("// Code generated by fyne-theme-generator")
	buf.writeln("")
	buf.writeln("package %s", t.packageName)
	buf.writeln("")
	buf.writeln("")
	buf.writeln("import (")
	buf.writeln("\"image/color\"")
	buf.writeln("")
	buf.writeln("\"fyne.io/fyne/v2\"")
	buf.writeln("\"fyne.io/fyne/v2/theme\"")
	buf.writeln(")")
	buf.writeln("")
	buf.writeln("type %s struct{}", t.themeStructName)
	buf.writeln("")
	buf.writeln("func (%s) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {", t.themeStructName)
	buf.writeln("	switch c {")
	buf.writeln("	case theme.ColorNameBackground:")
	buf.writeln("		return %#v", t.backgroundColor)
	buf.writeln("	case theme.ColorNameButton:")
	buf.writeln("		return %#v", t.buttonColor)
	buf.writeln("	case theme.ColorNameDisabledButton:")
	buf.writeln("		return %#v", t.disabledButtonColor)
	buf.writeln("	case theme.ColorNameDisabled:")
	buf.writeln("		return %#v", t.disabledColor)
	buf.writeln("	case theme.ColorNameError:")
	buf.writeln("		return %#v", t.errorColor)
	buf.writeln("	case theme.ColorNameFocus:")
	buf.writeln("		return %#v", t.focusColor)
	buf.writeln("	case theme.ColorNameForeground:")
	buf.writeln("		return %#v", t.foregroundColor)
	buf.writeln("	case theme.ColorNameHover:")
	buf.writeln("		return %#v", t.hoverColor)
	buf.writeln("	case theme.ColorNameInputBackground:")
	buf.writeln("		return %#v", t.inputBackgroundColor)
	buf.writeln("	case theme.ColorNamePlaceHolder:")
	buf.writeln("		return %#v", t.placeHolderColor)
	buf.writeln("	case theme.ColorNamePressed:")
	buf.writeln("		return %#v", t.pressedColor)
	buf.writeln("	case theme.ColorNamePrimary:")
	buf.writeln("		return %#v", t.primaryColor)
	buf.writeln("	case theme.ColorNameScrollBar:")
	buf.writeln("		return %#v", t.scrollBarColor)
	buf.writeln("	case theme.ColorNameShadow:")
	buf.writeln("		return %#v", t.shadowColor)
	buf.writeln("	default:")
	buf.writeln("		return theme.DefaultTheme().Color(c, v)")
	buf.writeln("	}")
	buf.writeln("}")
	buf.writeln("")
	buf.writeln("func (%s) Font(s fyne.TextStyle) fyne.Resource {", t.themeStructName)
	buf.writeln("	if s.Monospace {")
	if t.isSetTextMonospaceFont() {
		buf.writeln("return %s", sanitiseName(t.textMonospaceFont.Name()))
	} else {
		buf.writeln("return theme.DefaultTheme().Font(s)")
	}
	buf.writeln("	}")
	buf.writeln("	if s.Bold {")
	buf.writeln("		if s.Italic {")
	if t.isSetTextBoldItalicFont() {
		buf.writeln("return %s", sanitiseName(t.textBoldItalicFont.Name()))
	} else {
		buf.writeln("return theme.DefaultTheme().Font(s)")
	}
	buf.writeln("		}")
	if t.isSetTextBoldFont() {
		buf.writeln("return %s", sanitiseName(t.textBoldFont.Name()))
	} else {
		buf.writeln("return theme.DefaultTheme().Font(s)")
	}
	buf.writeln("	}")
	buf.writeln("	if s.Italic {")
	if t.isSetTextItalicFont() {
		buf.writeln("return %s", sanitiseName(t.textItalicFont.Name()))
	} else {
		buf.writeln("return theme.DefaultTheme().Font(s)")
	}
	buf.writeln("	}")
	if t.isSetTextFont() {
		buf.writeln("return %s", sanitiseName(t.textFont.Name()))
	} else {
		buf.writeln("return theme.DefaultTheme().Font(s)")
	}
	buf.writeln("}")
	buf.writeln("")
	buf.writeln("func (%s) Icon(n fyne.ThemeIconName) fyne.Resource {", t.themeStructName)
	buf.writeln("	return theme.DefaultTheme().Icon(n)")
	buf.writeln("}")
	buf.writeln("")
	buf.writeln("func (%s) Size(s fyne.ThemeSizeName) float32 {", t.themeStructName)
	buf.writeln("	switch s {")
	buf.writeln("	case theme.SizeNameCaptionText:")
	buf.writeln("		return %#v", t.captionTextSize)
	buf.writeln("	case theme.SizeNameInlineIcon:")
	buf.writeln("		return %#v", t.inlineIconSize)
	buf.writeln("	case theme.SizeNamePadding:")
	buf.writeln("		return %#v", t.paddingSize)
	buf.writeln("	case theme.SizeNameScrollBar:")
	buf.writeln("		return %#v", t.scrollBarSize)
	buf.writeln("	case theme.SizeNameScrollBarSmall:")
	buf.writeln("		return %#v", t.scrollBarSmallSize)
	buf.writeln("	case theme.SizeNameSeparatorThickness:")
	buf.writeln("		return %#v", t.separatorThicknessSize)
	buf.writeln("	case theme.SizeNameText:")
	buf.writeln("		return %#v", t.textSize)
	buf.writeln("	case theme.SizeNameInputBorder:")
	buf.writeln("		return %#v", t.inputBorderSize)
	buf.writeln("	default:")
	buf.writeln("		return theme.DefaultTheme().Size(s)")
	buf.writeln("	}")
	buf.writeln("}")

	return format.Source(buf.Bytes())
}

func buildThemeSourceForV1(t *Setting) ([]byte, error) {
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
	buf.writeln("func (%s) TextColor() color.Color            { return %#v }", t.themeStructName, t.foregroundColor)
	buf.writeln("func (%s) DisabledTextColor() color.Color    { return %#v }", t.themeStructName, t.disabledColor)
	buf.writeln("func (%s) IconColor() color.Color            { return %#v }", t.themeStructName, t.foregroundColor)
	buf.writeln("func (%s) DisabledIconColor() color.Color    { return %#v }", t.themeStructName, t.disabledColor)
	buf.writeln("func (%s) HyperlinkColor() color.Color       { return %#v }", t.themeStructName, t.primaryColor)
	buf.writeln("func (%s) PlaceHolderColor() color.Color     { return %#v }", t.themeStructName, t.placeHolderColor)
	buf.writeln("func (%s) PrimaryColor() color.Color         { return %#v }", t.themeStructName, t.primaryColor)
	buf.writeln("func (%s) HoverColor() color.Color           { return %#v }", t.themeStructName, t.hoverColor)
	buf.writeln("func (%s) FocusColor() color.Color           { return %#v }", t.themeStructName, t.focusColor)
	buf.writeln("func (%s) ScrollBarColor() color.Color       { return %#v }", t.themeStructName, t.scrollBarColor)
	buf.writeln("func (%s) ShadowColor() color.Color          { return %#v }", t.themeStructName, t.shadowColor)
	buf.writeln("func (%s) TextSize() int                     { return %#v }", t.themeStructName, t.textSize)
	if t.isSetTextFont() {
		buf.writeln("func (%s) TextFont() fyne.Resource           { return %s }", t.themeStructName, sanitiseName(t.textFont.Name()))
	} else {
		buf.writeln("func (%s) TextFont() fyne.Resource           { return theme.LightTheme().TextFont() }", t.themeStructName)
	}
	if t.isSetTextBoldFont() {
		buf.writeln("func (%s) TextBoldFont() fyne.Resource       { return %s }", t.themeStructName, sanitiseName(t.textBoldFont.Name()))
	} else {
		buf.writeln("func (%s) TextBoldFont() fyne.Resource       { return theme.LightTheme().TextBoldFont() }", t.themeStructName)
	}
	if t.isSetTextItalicFont() {
		buf.writeln("func (%s) TextItalicFont() fyne.Resource     { return %s }", t.themeStructName, sanitiseName(t.textItalicFont.Name()))
	} else {
		buf.writeln("func (%s) TextItalicFont() fyne.Resource     { return theme.LightTheme().TextItalicFont() }", t.themeStructName)
	}
	if t.isSetTextBoldItalicFont() {
		buf.writeln("func (%s) TextBoldItalicFont() fyne.Resource { return %s }", t.themeStructName, sanitiseName(t.textBoldItalicFont.Name()))
	} else {
		buf.writeln("func (%s) TextBoldItalicFont() fyne.Resource { return theme.LightTheme().TextBoldItalicFont() }", t.themeStructName)
	}
	if t.isSetTextMonospaceFont() {
		buf.writeln("func (%s) TextMonospaceFont() fyne.Resource  { return %s }", t.themeStructName, sanitiseName(t.textMonospaceFont.Name()))
	} else {
		buf.writeln("func (%s) TextMonospaceFont() fyne.Resource  { return theme.LightTheme().TextMonospaceFont() }", t.themeStructName)
	}
	buf.writeln("func (%s) Padding() int                      { return %#v }", t.themeStructName, t.paddingSize)
	buf.writeln("func (%s) IconInlineSize() int               { return %#v }", t.themeStructName, t.inlineIconSize)
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
	if t.isSetTextFont() {
		buf.writeln("var %s = %#v\n", sanitiseName(t.textFont.Name()), t.textFont)
		buf.writeln("")
	}
	if t.isSetTextBoldFont() {
		buf.writeln("var %s = %#v\n", sanitiseName(t.textBoldFont.Name()), t.textBoldFont)
		buf.writeln("")
	}
	if t.isSetTextItalicFont() {
		buf.writeln("var %s = %#v\n", sanitiseName(t.textItalicFont.Name()), t.textItalicFont)
		buf.writeln("")
	}
	if t.isSetTextBoldItalicFont() {
		buf.writeln("var %s = %#v\n", sanitiseName(t.textBoldItalicFont.Name()), t.textBoldItalicFont)
		buf.writeln("")
	}
	if t.isSetTextMonospaceFont() {
		buf.writeln("var %s = %#v\n", sanitiseName(t.textMonospaceFont.Name()), t.textMonospaceFont)
		buf.writeln("")
	}

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
