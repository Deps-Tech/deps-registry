package registry

var WellKnownAliases = map[string][]string{
	"luasocket": {
		"socket",
		"socket.http",
		"socket.ftp",
		"socket.smtp",
		"socket.url",
		"socket.headers",
		"socket.tp",
		"socket.core",
	},
	"cjson": {
		"cjson.safe",
	},
	"windows": {
		"windows.message",
	},
	"ssl": {
		"ssl.https",
	},
	"mimgui": {
		"mimgui.imgui",
		"mimgui.dx9",
		"mimgui.cdefs",
	},
	"socket": {
		"socket.http",
		"socket.ftp",
		"socket.smtp",
		"socket.url",
		"socket.headers",
		"socket.tp",
		"socket.core",
	},
	"mime": {
		"mime.core",
	},
	"ltn12": {
		"ltn12",
	},
	"xml": {
		"xml.core",
	},
	"lub": {
		"lub",
	},
}

func GetAliases(packageID string) []string {
	if aliases, ok := WellKnownAliases[packageID]; ok {
		return aliases
	}
	return nil
}

func ResolveAlias(alias string) string {
	for pkgID, aliases := range WellKnownAliases {
		for _, a := range aliases {
			if a == alias {
				return pkgID
			}
		}
	}
	return ""
}

