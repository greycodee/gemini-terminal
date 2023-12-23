# Gemini Terminal
[[English]](./README.md) [[中文]](./README_zh.md) [[日本語]](./README_jp.md)

Gemini Terminal是一个基于Google AI的终端聊天应用。你可以通过以下步骤来使用它。

![Gemini Terminal](./628566.gif)

## 安装与使用

1. 克隆这个项目

```bash
git clone https://github.com/greycodee/gemini-terminal.git
```

2. 构建项目

```bash
cd gemini-terminal && go build .
```

3. 运行项目

```bash
./gemini-terminal
```

> 注意: 你需要在 `$HOME/.local/share/gemini/config.ini` 文件中设置你自己的Google AI密钥。

## 配置

默认的配置文件位于 `$HOME/.local/share/gemini/config.ini`，你可以在这个文件中设置你的Google AI密钥和Gemini模型名称。

```ini
[Gemini]
# 设置你自己的Google AI密钥
googleAIKey=
# 设置Gemini模型名称
model=gemini-pro
[SafetySetting]
# HarmBlockUnspecified HarmBlockThreshold = 0
# HarmBlockLowAndAbove means content with NEGLIGIBLE will be allowed.
# HarmBlockLowAndAbove HarmBlockThreshold = 1
# HarmBlockMediumAndAbove means content with NEGLIGIBLE and LOW will be allowed.
# HarmBlockMediumAndAbove HarmBlockThreshold = 2
# HarmBlockOnlyHigh means content with NEGLIGIBLE, LOW, and MEDIUM will be allowed.
# HarmBlockOnlyHigh HarmBlockThreshold = 3
# HarmBlockNone means all content will be allowed.
# HarmBlockNone HarmBlockThreshold = 4
level=4
```

## 聊天历史

默认的数据库文件位于 `$HOME/.local/share/gemini/gemini.db`，你可以在这个文件中查看你的聊天历史。
