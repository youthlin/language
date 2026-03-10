package language

import (
	"slices"
	"sort"
	"strings"
)

var (
	// TagList 系统加载的语言标签列表
	TagList      = []Tag{}
	LanguageList = []string{}
	ScriptList   = []string{}
	RegionList   = []string{}
	VariantList  = []string{}
)

func init() {
	var (
		mLang    = map[string]bool{}
		mScript  = map[string]bool{}
		mRegion  = map[string]bool{}
		mVariant = map[string]bool{}
	)
	for _, item := range cldrItems {
		TagList = append(TagList, item.tag)
		for lang := range item.languages {
			if !mLang[lang] && !strings.Contains(lang, "-") {
				mLang[lang] = true
				LanguageList = append(LanguageList, lang)
			}
		}
		for script := range item.scripts {
			if !mScript[script] && !strings.Contains(script, "-") {
				mScript[script] = true
				ScriptList = append(ScriptList, script)
			}
		}
		for region := range item.regions {
			if !mRegion[region] && !strings.Contains(region, "-") {
				mRegion[region] = true
				RegionList = append(RegionList, region)
			}
		}
		for variant := range item.variants {
			if !mVariant[variant] && !strings.Contains(variant, "-") {
				mVariant[variant] = true
				VariantList = append(VariantList, variant)
			}
		}
		sort.Strings(LanguageList)
		sort.Strings(ScriptList)
		sort.Strings(RegionList)
		sort.Strings(VariantList)
	}
}

// Tag 语言标签
//
// 如 zh-cmn-Hans-CN-pinyin-u-co-phonebk-x-abc.UTF-8 对应:
// Language: zh
// ExtLang: cmn
// Script: Hans
// Region: CN
// Variant: pinyin
// Extension: u-co-phonebk
// PrivateUse: x-abc
// Encoding: UTF-8
//
// @see https://www.rfc-editor.org/info/bcp47
type Tag struct {
	Language   string
	ExtLang    string
	Script     string
	Region     string
	Variant    []string
	Extension  string
	PrivateUse string
	Encoding   string
}

// Of 从字符串创建语言标签实例, 如 "zh-CN", "zh_Hans_CN.UTF-8"
// 解析符合 BCP47 规范的语言标签为 Tag 结构体
func Of(tag string) (t Tag) {
	tag = strings.TrimSpace(tag)
	if s, ok := special[tag]; ok {
		tag = s
	}
	// 去掉 .UTF-8 后缀, 如 zh_CN.UTF-8 只保留 zh_CN
	parts := strings.Split(tag, ".")
	if len(parts) > 1 {
		t.Encoding = parts[1]
	}
	tag = parts[0]

	// 替换下划线为短横线, 如 zh_CN 替换为 zh-CN
	tag = strings.ReplaceAll(tag, "_", "-")
	// 按短横线分割, 如 zh-CN 分割为 zh, CN
	parts = strings.Split(tag, "-")
	parts = slices.DeleteFunc(parts, isSpace)

	// language–extlang–script–region–variant–extension–privateuse
	// 来源: https://youthlin.com/?p=1843
	// https://www.rfc-editor.org/rfc/rfc5646.html#section-2.1
	// language: 2 或 3 字母 不可省略
	// extlang: 3 字母 可选 最多 1 个
	// script: 4 字母 可选 最多 1 个
	// region: 2 字母或 3 数字 可选 最多 1 个
	// variant: 字母开头则5-8长度 数字开头则4字符 可选 可多个
	// extension: 单字母开头
	// privateuse: x 开头

	for len(parts) > 0 {
		sub := parts[0]
		switch len(sub) {
		case 1:
			// 单字符: extension 或 privateuse
			end := len(parts)
			for i, sub := range parts {
				if sub == "x" || sub == "X" {
					t.PrivateUse = strings.Join(parts[i:], "-")
					end = i
				}
			}
			t.Extension = strings.Join(parts[:end], "-")
			return
		case 2:
			// 2 字母: language 或 region
			switch {
			case t.Language == "":
				t.Language = strings.ToLower(sub)
			case t.Region == "":
				t.Region = strings.ToUpper(sub)
			}
		case 3:
			if t.Region == "" {
				// 3 字母: language 或 extlang
				// 3 数字: region
				switch {
				case isAlpha(sub[0]):
					// lang, extlang 应当在 region 之前
					// Region 已经有值的话 lang, extlang 就不能再赋值了
					switch {
					case t.Language == "":
						t.Language = strings.ToLower(sub)
					case t.ExtLang == "":
						t.ExtLang = strings.ToLower(sub)
					}
				case isDigit(sub[0]):
					t.Region = sub
				}
			}
		case 4:
			// 4 字母: script
			// 4 字符且数字开头: variant
			switch {
			case isAlpha(sub[0]) && t.Script == "" && t.Region == "":
				sub = strings.ToLower(sub)
				sub = strings.ToUpper(sub[:1]) + sub[1:]
				t.Script = sub
			case isDigit(sub[0]):
				// sub = strings.ToUpper(sub)
				t.Variant = append(t.Variant, sub)
			}
		default:
			// sub = strings.ToUpper(sub)
			sub = strings.TrimSpace(sub)
			t.Variant = append(t.Variant, sub)
		}
		parts = parts[1:]
	}
	return
}

// special 将一些 grandfathered 标记特殊识别
var special = map[string]string{
	"en-GB-oed": "en-GB-oxendict", // 英语（英国，《牛津英语词典》拼法）
	"i-ami":     "ami",
	"i-bnn":     "bnn",
	"i-default": "",
	// "i-enochian": "",// 以诺语 人造虚拟语言
	"i-hak":     "hak", // 客家语
	"i-klingon": "tlh", // 克林贡语
	"i-lux":     "lb",  // 卢森堡语
	// "i-mingo": "", // 明戈语?
	"i-navajo":   "nv", // 纳瓦霍语
	"i-pwn":      "pwn",
	"i-tao":      "tao",
	"i-tay":      "tay",
	"i-tsu":      "tsu",
	"sgn-BE-FR":  "sfb", // 比利时-法语的手语
	"sgn-BE-NL":  "vgt", // 手语
	"sgn-CH-DE":  "sgg", // 手语
	"art-lojban": "jbo", // 逻辑语
	// "cel-gaulish": "", // 高卢语(xcg, xga, xtg)
	"no-bok":   "nb",           // 书面挪威语
	"no-nyn":   "nn",           // 挪威尼诺斯克语
	"zh-guoyu": "zh-cmn-guoyu", // 国语 -> 普通话
	"zh-hakka": "hak",          // 客家语
	// 视为 zh+闽语变体 添加空格到5字符以便能够被识别为变体
	"zh-min":     "zh- min ", // 闽语(cdo闽东, cpx蒲县, czo闽中, mnp闽北, nan闽南)
	"zh-min-nan": "nan",      // 闽南语
	"zh-xiang":   "hsn",      // 湘语
}

// String 返回语言标签的字符串表示, 符合 BCP47 规范
func (t Tag) String() string {
	var sb strings.Builder
	if t.Language != "" {
		sb.WriteString(t.Language)
	}
	if t.ExtLang != "" {
		if sb.Len() > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(t.ExtLang)
	}
	if t.Script != "" {
		if sb.Len() > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(t.Script)
	}
	if t.Region != "" {
		if sb.Len() > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(t.Region)
	}
	if len(t.Variant) > 0 {
		if sb.Len() > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(strings.Join(t.Variant, "-"))
	}
	if t.Extension != "" {
		if sb.Len() > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(t.Extension)
	}
	if t.PrivateUse != "" {
		if sb.Len() > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(t.PrivateUse)
	}
	if t.Encoding != "" {
		sb.WriteString("." + t.Encoding)
	}
	return sb.String()
}

func isSpace(s string) bool {
	return strings.TrimSpace(s) == ""
}

func isAlpha(s byte) bool {
	return (s >= 'a' && s <= 'z') || (s >= 'A' && s <= 'Z')
}

func isDigit(s byte) bool {
	return s >= '0' && s <= '9'
}

// GetTagName 获取 show 标签在 t 语境中整个语言标签的名称, 数据来源于 CLDR.
//
//   - 一般情况下 `lang-Script-Region-Variant` 在中文语境下会返回 `语言名称（书写格式，区域，变体）` 的格式
//   - 特殊地, 如 `zh-Hans` 会整体当作语言识别为“简体中文”(在 cldr 数据中定义), 而不是 "中文（简体）"
//
// 比如
//
//	zh.GetTagName(zh) -> 中文
//	zh.GetTagName(zh-Hans) -> 简体中文
//	zh.GetTagName(zh-Hans-CN) -> 简体中文（中国）
//	zh.GetTagName(yue-Hant-HK) -> 粤语（繁体，中国香港特别行政区）
func (t Tag) GetTagName(show Tag) string {
	loc := t.getCldrItem()
	return loc.getTagName(show)
}

// GetTagLangNameCode 获取 show 标签在 t 语境中的语言名称和对应的代码, 数据来源于 CLDR.
//
//   - 一般情况下, `lang-Script-Region-Variant` 在中文语境下会返回 lang 部分对应的 `语言名称`
//   - 特殊地, 比如 `zh-Hans` 会整体当作语言识别为“简体中文”(在 cldr 数据中定义), 而不仅仅返回 "中文"
//
// 即, 返回的名称可能是识别了 Script, Region 的.
//
// 扩展语言如有单独定义, 则优先使用, 如 zh-yue 直接显示为 "粤语"(因为 CLDR 的数据里直接就有 yue 的详情)
//
// 比如
//
//	zh.GetTagLangNameCode(zh) -> 中文, zh
//	zh.GetTagLangNameCode(zh-Latn) -> 中文, zh
//	zh.GetTagLangNameCode(zh-Hans-CN) -> 简体中文, zh-Hans
//	zh.GetTagLangNameCode(yue) -> 粤语, yue
//	zh.GetTagLangNameCode(zh-yue) -> 粤语, yue
//	zh.GetTagLangNameCode(ar) -> 阿拉伯语, ar
//	zh.GetTagLangNameCode(ar-001) -> 现代标准阿拉伯语, ar-001
//	zh.GetTagLangNameCode(en) -> 英语, en
//	zh.GetTagLangNameCode(en_US) -> 美国英语, en_US
func (t Tag) GetTagLangNameCode(show Tag) (name, code string) {
	loc := t.getCldrItem()
	res := loc.getLangName(show)
	sub := []string{}
	if res.containsLang {
		sub = append(sub, show.Language)
	}
	if res.containsExtLang {
		sub = append(sub, show.ExtLang)
	}
	if res.containsScript {
		sub = append(sub, show.Script)
	}
	if res.containsRegion {
		sub = append(sub, show.Region)
	}
	if res.containsVariant != "" {
		sub = append(sub, res.containsVariant)
	}
	return res.name, strings.Join(sub, "-")
}

// GetLangName 获取 lang 语言对应的名称
// 比如 zh.GetLangName("zh") -> 中文
func (t Tag) GetLangName(lang string) string {
	return getLangName(t.getCldrItem(), lang)
}

// GetScriptName 获取 script 书写格式在 t 语境中的名称
// 比如 zh.GetScriptName("Hans") -> 简体
func (t Tag) GetScriptName(script string) string {
	return getScriptName(t.getCldrItem(), script)
}

// GetRegionName 获取 region 区域在 t 语境中的名称
// 比如 zh.GetRegionName("CN") -> 中国
func (t Tag) GetRegionName(region string) string {
	return getRegionName(t.getCldrItem(), region)
}

// GetVariantName 获取 variant 在 t 语境中的变体名称
// 比如 zh.GetVariantName("PINYIN") -> 拼音罗马字
func (t Tag) GetVariantName(variant string) string {
	return getVariantName(t.getCldrItem(), variant)
}

// FormatLocale 按 t 语境下的格式拼接字段
// 如果 more 字段有多个, 先用字段并列格式拼接 (中文下就是逗号分隔)
// 然后再用语言展示格式拼接 lang 和拼接好的 more (中文下就是将 more 放在括号中)
//
// 比如 zh.FormatLocale("中文","简体","中国") -> 中文（简体，中国）
// en.FormatLocale("Chinese", "Traditional", "Hong Kong SAR China") -> Chinese (Traditional, Hong Kong SAR China)
func (t Tag) FormatLocale(lang string, more ...string) string {
	loc := t.getCldrItem()
	return loc.formatLocale(lang, more...)
}

// FormatKV 按 t 语境下的键值对格式拼接 key, value 字段
func (t Tag) FormatKV(key, value string) string {
	pattern := t.GetLocaleKeyTypePattern()
	pattern = strings.Replace(pattern, "{0}", key, 1)
	pattern = strings.Replace(pattern, "{1}", value, 1)
	return pattern
}

// GetLocalePattern 获取 t 语境中的语言标签模版
// 如 zh 语言的模板为 {0}（{1}） 使用全角括号, 其中 {0} 为语言名称, {1} 补充说明
func (t Tag) GetLocalePattern() string {
	return getPattern(t.getCldrItem())
}

// GetLocaleSeparator 获取 t 语境中的字段分隔模板
// 如 zh 语言的字段分隔模板为 {0}，{1} 使用全角逗号
func (t Tag) GetLocaleSeparator() string {
	return getSeparator(t.getCldrItem())
}

// GetLocaleKeyTypePattern 获取 t 语境中的键值对模版
// 如 zh 语言的键值对模版为 {0}：{1} 使用全角冒号分隔键值对
func (t Tag) GetLocaleKeyTypePattern() string {
	return getKeyTypePattern(t.getCldrItem())
}

// getCldrItem 获取对应的 CldrItem 详情
func (t Tag) getCldrItem() (loc *cldrItem) {
	// fast path: 尝试直接从 localMap 查找
	if loc = searchCldrItem(t); loc != nil {
		return
	}
	// zh-yue-Hant-CN 提升为 yue-Hant-CN
	if t.ExtLang != "" {
		if loc = searchCldrItemWithLang(t, t.ExtLang); loc != nil {
			return
		}
	}
	// 去掉 ExtLang 再找
	if loc = searchCldrItemWithLang(t, t.Language); loc != nil {
		return
	}
	// 返回兜底结果
	return _root // defined in table.gen.go
}

// searchCldrItem 直接用语言标签搜索对应的 CldrItem
func searchCldrItem(t Tag) (loc *cldrItem) {
	code := t.String()
	return cldrItemMap[code]
}

// searchCldrItemWithLang 把传入的语言标签改成指定的语言 再搜索对应的 CldrItem
func searchCldrItemWithLang(t Tag, lang string) (loc *cldrItem) {
	// lang-extLang-Script-Region-Variant-Extension-PrivateUse-Encoding
	// 去掉 extlang, Extension-PrivateUse-Encoding
	t.Language = lang
	t.ExtLang = ""
	t.Extension = ""
	t.PrivateUse = ""
	t.Encoding = ""

	// lang-Script-Region-Variants
	if loc = searchCldrItem(t); loc != nil {
		return
	}

	// 单一变体
	for _, variant := range t.Variant {
		if loc = searchLocaleWithVariant(t, variant); loc != nil {
			return
		}
	}

	// 去掉变体
	return searchLocaleWithVariant(t, "")
}

// searchLocaleWithVariant 指定变体或者去掉变体 再搜索对应的 CldrItem
func searchLocaleWithVariant(t Tag, variant string) (loc *cldrItem) {
	if variant == "" {
		t.Variant = nil
	} else {
		t.Variant = []string{variant}
	}

	if loc = searchCldrItem(t); loc != nil {
		return
	}

	// 去掉 Region 再试
	{
		tt := t
		tt.Region = ""
		if loc = searchCldrItem(tt); loc != nil {
			return
		}
	}
	// 去掉 Script 再试
	{
		tt := t
		tt.Script = ""
		if loc = searchCldrItem(tt); loc != nil {
			return
		}
	}
	// 去掉 Script/Region 再试
	{
		tt := t
		tt.Script = ""
		tt.Region = ""
		if loc = searchCldrItem(tt); loc != nil {
			return
		}
	}
	return
}
