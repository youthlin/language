package language_test

import (
	"reflect"
	"testing"

	"github.com/youthlin/language"
)

var zhHansCN = language.Of("zh-Hans-CN")

func TestTag(t *testing.T) {
	tag := language.Of("zh-cmn-Hans-CN")
	t.Logf("%v", tag)
	// zh-cmn-Hans-CN

	t.Logf("%v", tag.GetTagName(tag))
	// 简体中文（中国）

	name, code := tag.GetTagLangNameCode(tag)
	t.Logf("%v, %v", name, code)
	// 简体中文, zh-Hans

	t.Logf("%v", tag.GetLangName("zh"))        // 中文
	t.Logf("%v", tag.GetLangName("hak"))       // 客家语
	t.Logf("%v", tag.GetScriptName("Hans"))    // 简体
	t.Logf("%v", tag.GetRegionName("CN"))      // 中国
	t.Logf("%v", tag.GetVariantName("PINYIN")) // 拼音罗马字

	t.Logf("%v", tag.FormatKV(
		tag.String(),
		tag.FormatLocale(
			tag.GetLangName("zh"),
			tag.GetScriptName("Hans"),
			tag.GetRegionName("CN"),
		),
	)) // zh-cmn-Hans-CN：中文（简体，中国）

	t.Logf("%v", tag.GetLocalePattern())        // {0}（{1}）
	t.Logf("%v", tag.GetLocaleSeparator())      // {0}，{1}
	t.Logf("%v", tag.GetLocaleKeyTypePattern()) // {0}：{1}
}

func TestTagList(t *testing.T) {
	for _, tag := range language.TagList() {
		t.Logf("[%v]", tag)
	}
}

func TestOf(t *testing.T) {
	tests := []struct {
		name string
		want language.Tag
	}{
		{name: "", want: language.Tag{}},
		{name: ".", want: language.Tag{}},
		{name: "_", want: language.Tag{}},
		{name: "-", want: language.Tag{}},
		{name: "--", want: language.Tag{}},
		{name: "a", want: language.Tag{Extension: "a"}},
		{name: "zh", want: language.Tag{Language: "zh"}},
		{name: "zh-cmn", want: language.Tag{Language: "zh", ExtLang: "cmn"}},
		// invalid: extlang 重复 只取第一个
		{name: "zh-cmn-hak", want: language.Tag{Language: "zh", ExtLang: "cmn"}},
		{name: "zh--cmn", want: language.Tag{Language: "zh", ExtLang: "cmn"}},
		{name: "zh-yue-HK", want: language.Tag{Language: "zh", ExtLang: "yue", Region: "HK"}},
		{name: "cmn", want: language.Tag{Language: "cmn"}},
		{name: "be_TARASK", want: language.Tag{Language: "be", Variant: []string{"TARASK"}}},
		{name: "zh-Hans", want: language.Tag{Language: "zh", Script: "Hans"}},
		{name: "zh-Hans-CN", want: language.Tag{Language: "zh", Script: "Hans", Region: "CN"}},
		// Script 首字母自动大小 Region 自动大写
		{name: "ZH-hans-cn", want: language.Tag{Language: "zh", Script: "Hans", Region: "CN"}},
		{name: "zh-cmn-Hans-CN", want: language.Tag{Language: "zh", ExtLang: "cmn", Script: "Hans", Region: "CN"}},
		{name: "zh-cMn-Hans-CN-pinyin", want: language.Tag{Language: "zh", ExtLang: "cmn", Script: "Hans", Region: "CN", Variant: []string{"pinyin"}}},
		{name: "x-myext", want: language.Tag{PrivateUse: "x-myext"}},
		{name: "es-005", want: language.Tag{Language: "es", Region: "005"}},
		{name: "de-1901", want: language.Tag{Language: "de", Variant: []string{"1901"}}},
		// invalid: region 重复 只取第一个
		{name: "de-419-DE", want: language.Tag{Language: "de", Region: "419"}},
		{name: "gd-oxendict-1990", want: language.Tag{Language: "gd", Variant: []string{"oxendict", "1990"}}},
		{name: "zh_CN.UTF-8", want: language.Tag{Language: "zh", Region: "CN", Encoding: "UTF-8"}},
		{name: "en-US-1994-u-ca-gregorian", want: language.Tag{Language: "en", Region: "US", Variant: []string{"1994"}, Extension: "u-ca-gregorian"}},
		{name: "en-US-1994-u-ca-gregorian-x-abc", want: language.Tag{Language: "en", Region: "US", Variant: []string{"1994"}, Extension: "u-ca-gregorian", PrivateUse: "x-abc"}},
		{name: "en-US-1994-x-abc", want: language.Tag{Language: "en", Region: "US", Variant: []string{"1994"}, PrivateUse: "x-abc"}},
		// 一些 grandfathered 标记
		{name: "en-GB-oed", want: language.Tag{Language: "en", Region: "GB", Variant: []string{"oxendict"}}},
		{name: "i-ami", want: language.Tag{Language: "ami"}},
		{name: "i-bnn", want: language.Tag{Language: "bnn"}},
		{name: "i-default", want: language.Tag{}},
		{name: "i-enochian", want: language.Tag{Extension: "i-enochian"}},
		{name: "i-hak", want: language.Tag{Language: "hak"}},
		{name: "i-klingon", want: language.Tag{Language: "tlh"}},
		{name: "i-lux", want: language.Tag{Language: "lb"}},
		{name: "i-navajo", want: language.Tag{Language: "nv"}},
		{name: "i-pwn", want: language.Tag{Language: "pwn"}},
		{name: "i-tao", want: language.Tag{Language: "tao"}},
		{name: "i-tay", want: language.Tag{Language: "tay"}},
		{name: "i-tsu", want: language.Tag{Language: "tsu"}},
		{name: "sgn-BE-FR", want: language.Tag{Language: "sfb"}},
		{name: "sgn-BE-NL", want: language.Tag{Language: "vgt"}},
		{name: "sgn-CH-DE", want: language.Tag{Language: "sgg"}},
		{name: "art-lojban", want: language.Tag{Language: "jbo"}},
		{name: "cel-gaulish", want: language.Tag{Language: "cel", Variant: []string{"gaulish"}}},
		{name: "no-bok", want: language.Tag{Language: "nb"}},
		{name: "no-nyn", want: language.Tag{Language: "nn"}},
		{name: "zh-guoyu", want: language.Tag{Language: "zh", ExtLang: "cmn", Variant: []string{"guoyu"}}},
		{name: "zh-hakka", want: language.Tag{Language: "hak"}},
		{name: "zh-min", want: language.Tag{Language: "zh", Variant: []string{"min"}}},
		{name: "zh-min-nan", want: language.Tag{Language: "nan"}},
		{name: "zh-xiang", want: language.Tag{Language: "hsn"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := language.Of(tt.name)
			t.Logf("%s", tag.String())
			if !reflect.DeepEqual(tag, tt.want) {
				t.Errorf("Makelanguage.Tag() = %#v, want %#v", tag, tt.want)
			}
		})
	}
}

func TestGetTagName(t *testing.T) {
	for _, tag := range language.TagList() {
		t.Logf("[%v]\t对应的中文名是 [%v]\t它自身名称是 [%v]",
			tag, zhHansCN.GetTagName(tag), tag.GetTagName(tag))
	}
}

func TestGetTagLangNameCode(t *testing.T) {
	for _, s := range language.LanguageList() {
		tag := language.Of(s)
		name, code := zhHansCN.GetTagLangNameCode(tag)
		t.Logf("[%v]: %v, %v", s, name, code)
	}
}

func TestGetLangName(t *testing.T) {
	for _, s := range language.LanguageList() {
		t.Logf("[%v]: %v", s, zhHansCN.GetLangName(s))
	}
}

func TestGetScriptName(t *testing.T) {
	for _, s := range language.ScriptList() {
		t.Logf("[%v]: %v", s, zhHansCN.GetScriptName(s))
	}
}

func TestGetRegionName(t *testing.T) {
	for _, s := range language.RegionList() {
		t.Logf("[%v]: %v", s, zhHansCN.GetRegionName(s))
	}
}

func TestGetVariantName(t *testing.T) {
	for _, s := range language.VariantList() {
		t.Logf("[%v]: %v", s, zhHansCN.GetVariantName(s))
	}
}

func TestFormatLocale(t *testing.T) {
	name := zhHansCN.FormatLocale(
		zhHansCN.GetLangName("zh"),
		zhHansCN.GetLangName("yue"),
		zhHansCN.GetScriptName("Hant"),
		zhHansCN.GetRegionName("HK"),
	)
	t.Logf("%v", name)
	en := language.Of("en")
	t.Logf("%v", en.FormatLocale(
		en.GetLangName("zh"),
		en.GetLangName("yue"),
		en.GetScriptName("Hant"),
		en.GetRegionName("HK"),
	))
}

func TestFormatKV(t *testing.T) {
	t.Logf("%v", zhHansCN.FormatKV(
		zhHansCN.String(),
		zhHansCN.GetTagName(zhHansCN),
	))
}
