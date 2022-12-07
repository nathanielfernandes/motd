package slp

import (
	"strings"
)

// https://wiki.vg/Chat
type Chat struct {
	Text  string `json:"text,omitempty"`
	Extra []Chat `json:"extra,omitempty"`

	// Shared between all components
	Bold          *bool  `json:"bold,omitempty"`
	Italic        *bool  `json:"italic,omitempty"`
	Underlined    *bool  `json:"underlined,omitempty"`
	Strikethrough *bool  `json:"strikethrough,omitempty"`
	Obfuscated    *bool  `json:"obfuscated,omitempty"`
	Font          string `json:"font,omitempty"`
	Color         string `json:"color,omitempty"`
}

func ColorToUnicode(color string) string {
	switch color {
	case "dark_red":
		return "\u00A74"
	case "red":
		return "\u00A7c"
	case "gold":
		return "\u00A76"
	case "yellow":
		return "\u00A7e"
	case "dark_green":
		return "\u00A72"
	case "green":
		return "\u00A7a"
	case "aqua":
		return "\u00A7b"
	case "dark_aqua":
		return "\u00A73"
	case "dark_blue":
		return "\u00A71"
	case "blue":
		return "\u00A79"
	case "light_purple":
		return "\u00A7d"
	case "dark_purple":
		return "\u00A75"
	case "white":
		return "\u00A7f"
	case "gray":
		return "\u00A77"
	case "dark_gray":
		return "\u00A78"
	case "black":
		return "\u00A70"
	case "reset":
		return "\u00A7r"
	default:
		return ""
	}
}

func (c Chat) buildString(parent Chat) string {
	str := strings.Builder{}

	if c.Text != "" {
		if c.Color != "" {
			str.WriteString(ColorToUnicode(c.Color))
		} else if parent.Color != "" {
			str.WriteString(ColorToUnicode(parent.Color))
		}

		if c.Bold != nil {
			if *c.Bold {
				str.WriteString("\u00A7l")
			}
		} else if parent.Bold != nil {
			if *parent.Bold {
				str.WriteString("\u00A7l")
			}
		}

		if c.Italic != nil {
			if *c.Italic {
				str.WriteString("\u00A7o")
			}
		} else if parent.Italic != nil {
			if *parent.Italic {
				str.WriteString("\u00A7o")
			}
		}

		if c.Underlined != nil {
			if *c.Underlined {
				str.WriteString("\u00A7n")
			}
		} else if parent.Underlined != nil {
			if *parent.Underlined {
				str.WriteString("\u00A7n")
			}
		}

		if c.Strikethrough != nil {
			if *c.Strikethrough {
				str.WriteString("\u00A7m")
			}
		} else if parent.Strikethrough != nil {
			if *parent.Strikethrough {
				str.WriteString("\u00A7m")
			}
		}

		if c.Obfuscated != nil {
			if *c.Obfuscated {
				str.WriteString("\u00A7k")
			}
		} else if parent.Obfuscated != nil {
			if *parent.Obfuscated {
				str.WriteString("\u00A7k")
			}
		}

		str.WriteString(c.Text)
	}

	for _, extra := range c.Extra {
		str.WriteString(extra.buildString(c))
	}

	return str.String()

}

func (c Chat) String() string {
	return c.buildString(Chat{})
}
