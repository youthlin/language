# language

- BCP 47 语言标签解析
- 语言标签本地化名称

## 安装

```bash
go get github.com/youthlin/language
```

## 语言标签解析

```go
import (
	"github.com/youthlin/language"
)

tag := language.Of("zh-Hans-CN")
// language.Tag{Language: "zh", Script: "Hans", Region: "CN"}

```

Tag 结构:

```go
// Tag 语言标签
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
```

## 获取语言标签名称

```go
zh := language.Of("zh-Hans-CN")
zh.GetTagName(zh) // 简体中文（中国）
zh.GetLangName("en") // 英语
zh.GetScriptName("Latn") // 拉丁文
zh.GetRegionName("US") // 美国
zh.GetVariantName("PINYIN") // 拼音罗马字
```
