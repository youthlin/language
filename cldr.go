package language

import (
	"strings"
)

// cldrItem 不同语言的 cldr 详情
type cldrItem struct {
	tag      Tag         // 当前语言标签
	parent   *cldrItem   // 父标签
	children []*cldrItem // 子标签
	// 当前语境下 语言名称的格式 {0}=lang {1}=region 如：中文使用 {0}（{1}） 全角括号
	localePattern string
	// 当前语境下 两个字段并列的格式, 如：中文使用 {0}，{1} 全角逗号并列
	localeSeparator string
	// 当前语境下 key-value 并列的格式，如：中文使用{0}：{1} 全角冒号分隔
	localeKeyTypePattern string
	// 当前语境下 各语言的名称
	languages map[string]string
	// 当前语境下 书写系统的名称
	scripts map[string]string
	// 当前语境下 各地区的名称
	regions map[string]string
	// 当前语境下 各种变体的名称
	variants map[string]string
}

// String 打印格式, 直接输出 Tag
func (item cldrItem) String() string {
	return item.tag.String()
}

// Parent 获取父项目, 如: zh_Hans_CN -> zh_Hans -> zh -> root
func (item *cldrItem) Parent() *cldrItem {
	tag := item.tag
	if len(tag.Variant) > 0 {
		// 有 variant, 去掉变体, 如 be_TARASK -> be
		// 方法接收者是结构体不是指针 这里修改不会影响外层调用方
		tag.Variant = nil
		parentCode := tag.String()
		if parent, ok := cldrItemMap[parentCode]; ok {
			return parent
		}
	}
	if tag.Region != "" {
		// 有 region, 去掉后是 parent
		// lang-Script-Region -> lang-Script
		// lang-Region -> lang
		tag.Region = ""
		parentCode := tag.String()
		if parent, ok := cldrItemMap[parentCode]; ok {
			return parent
		}
	}
	if tag.Script != "" {
		// 没有 Region 但有 script, 去掉后是 parent
		// lang-Script -> lang
		tag.Script = ""
		parentCode := tag.String()
		if parent, ok := cldrItemMap[parentCode]; ok {
			return parent
		}
	}
	if tag.Language != "root" {
		// 没有 Region/Script 的 lang 其父标签是 root
		tag.Language = "root"
		parentCode := tag.String()
		if parent, ok := cldrItemMap[parentCode]; ok {
			return parent
		}
	}
	return _root
}

// getTagName 获取语言标签的名称
func (loc *cldrItem) getTagName(show Tag) string {
	if loc == nil {
		return ""
	}
	var (
		nameRes    = loc.getLangName(show)
		lang       = nameRes.name
		suffixPart = []string{} // 后缀部分, 如: （拉丁文）、（简体，中国）
	)
	// 获取 zh-Hans-CN 的语言名称时, 会直接使用 zh_Hans 得到 简体中文
	// 此时无需重复获取 Hans 脚本名称
	if show.Script != "" && !nameRes.containsScript {
		if script := getScriptName(loc, show.Script); script != "" {
			suffixPart = append(suffixPart, script)
		}
	}
	// 同理: ar_001 -> 现代标准阿拉伯语
	if show.Region != "" && !nameRes.containsRegion {
		if region := getRegionName(loc, show.Region); region != "" {
			suffixPart = append(suffixPart, region)
		}
	}
	// be_TARASK
	for _, variant := range show.Variant {
		if variant != nameRes.containsVariant {
			if name := getVariantName(loc, variant); name != "" {
				suffixPart = append(suffixPart, name)
			}
		}
	}
	return loc.formatLocale(lang, suffixPart...)
}

func (loc *cldrItem) formatLocale(lang string, suffixPart ...string) string {
	var suffix string
	if size := len(suffixPart); size > 0 {
		suffix = suffixPart[0]
		pattern := getSeparator(loc) // 多个字段并列的格式
		for i := 1; i < size; i++ {
			suffix = strings.Replace(pattern, "{0}", suffix, 1)
			suffix = strings.Replace(suffix, "{1}", suffixPart[i], 1)
		}
	}
	if suffix != "" {
		pattern := getPattern(loc) // 语言名称的格式
		lang = strings.Replace(pattern, "{0}", lang, 1)
		lang = strings.Replace(lang, "{1}", suffix, 1)
	}
	return lang
}

// getLangNameResult 语言标签转为字符串名称的结果
type getLangNameResult struct {
	name            string // 语言名称
	containsLang    bool   // 名称是否包含了 lang 部分, 比如 zh-CN 使用了 zh
	containsExtLang bool   // 名称是否包含了 extLang 部分, 比如 zh-yue 使用的是 yue 而不是 zh
	containsScript  bool   // 名称是否包含了 script 部分, 比如 zh-Hans 对应 简体中文
	containsRegion  bool   // 名称是否包含了 region 部分, 比如 ar_001 对应 现代标准阿拉伯语
	containsVariant string // 名称是否包含了变体字段, 比如 be_TARASK=白俄罗斯语（传统正写法）
}

// getLangName 获取语言名称
func (loc *cldrItem) getLangName(show Tag) (result getLangNameResult) {
	if show.ExtLang != "" {
		lang := show.Language
		// lang-extLang 的, 先尝试使用 extLang: zh-yue 直接使用 yue
		show.Language = show.ExtLang
		if result = getTagLangName(loc, show); result.name != "" {
			result.containsExtLang = true
			return
		}
		show.Language = lang // 还原
	}

	// zh-cmn 直接使用 cmn 无结果, 使用 zh
	if result = getTagLangName(loc, show); result.name != "" {
		result.containsLang = true
		return
	}

	// 没有对应名称 直接使用传入的代码
	result.name = show.Language
	result.containsLang = true
	if show.ExtLang != "" {
		result.containsExtLang = true
		result.name += "-" + show.ExtLang
	}
	return
}

// getTagLangName 获取语言名称. 如:
// - zh.getTagLangName(en) -> 英语
// - zh.getTagLangName(en_US) -> 美国英语, useRegion=true
// - zh.getTagLangName(zh_Hans) -> 简体中文, useScript=true
// - zh.getTagLangName(be_TARASK) -> 白俄罗斯语（传统正写法）, useVariant=TARASK
func getTagLangName(loc *cldrItem, show Tag) (result getLangNameResult) {
	if loc == nil {
		return
	}
	for loc.tag.Language != "root" {
		code := strings.Join([]string{show.Language, show.Script, show.Region}, "_")
		if name, ok := loc.languages[code]; ok {
			result.name = name
			result.containsScript = true
			result.containsRegion = true
			result.containsVariant = ""
			return
		}
		code = strings.Join([]string{show.Language, show.Region}, "_")
		if name, ok := loc.languages[code]; ok {
			result.name = name
			result.containsScript = false
			result.containsRegion = true
			result.containsVariant = ""
			return
		}
		code = strings.Join([]string{show.Language, show.Script}, "_")
		if name, ok := loc.languages[code]; ok {
			result.name = name
			result.containsScript = true
			result.containsRegion = false
			result.containsVariant = ""
			return
		}
		for _, variant := range show.Variant {
			// be_TARASK
			code = strings.Join([]string{show.Language, variant}, "_")
			if name, ok := loc.languages[code]; ok {
				result.name = name
				result.containsScript = false
				result.containsRegion = false
				result.containsVariant = variant
				return
			}
		}
		if name, ok := loc.languages[show.Language]; ok {
			result.name = name
			result.containsScript = false
			result.containsRegion = false
			result.containsVariant = ""
			return
		}
		loc = loc.parent
	}
	return
}

// getLangName 获取语言名称
func getLangName(loc *cldrItem, lang string) string {
	if lang == "" {
		return ""
	}
	for loc.tag.Language != "root" {
		if name, ok := loc.languages[lang]; ok {
			return name
		}
		loc = loc.parent
	}
	return lang
}

// getScriptName 获取书写系统的名称
func getScriptName(loc *cldrItem, script string) string {
	if len(script) != 4 {
		return ""
	}
	script = strings.ToLower(script)
	script = strings.ToUpper(script[:1]) + script[1:]
	for loc.tag.Language != "root" {
		if name, ok := loc.scripts[script]; ok {
			return name
		}
		loc = loc.parent
	}
	return script
}

// getRegionName 获取区域的名称
func getRegionName(loc *cldrItem, region string) string {
	if region == "" {
		return ""
	}
	region = strings.ToUpper(region)
	for loc.tag.Language != "root" {
		if name, ok := loc.regions[region]; ok {
			return name
		}
		loc = loc.parent
	}
	return region
}

// getVariantName 获取变体的名称
func getVariantName(loc *cldrItem, variant string) string {
	if variant == "" {
		return ""
	}
	variant = strings.ToUpper(variant)
	for loc.tag.Language != "root" {
		if name, ok := loc.variants[variant]; ok {
			return name
		}
		loc = loc.parent
	}
	return variant
}

// getPattern 获取语言的展示格式
func getPattern(loc *cldrItem) string {
	for loc != nil && loc.tag.Language != "root" {
		if loc.localePattern != "" {
			return loc.localePattern
		}
		loc = loc.parent
	}
	return "{0} ({1})"
}

// getSeparator 获取字段并列的展示格式
func getSeparator(loc *cldrItem) string {
	for loc != nil && loc.tag.Language != "root" {
		if loc.localeSeparator != "" {
			return loc.localeSeparator
		}
		loc = loc.parent
	}
	return "{0}, {1}"
}

// getKeyTypePattern 获取键值对的展示格式
func getKeyTypePattern(loc *cldrItem) string {
	for loc != nil && loc.tag.Language != "root" {
		if loc.localeKeyTypePattern != "" {
			return loc.localeKeyTypePattern
		}
		loc = loc.parent
	}
	return "{0}: {1}"
}
