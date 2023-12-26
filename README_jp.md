# Gemini Terminal

[[English]](./README.md) [[中文]](./README_zh.md) [[日本語]](./README_jp.md)

Gemini Terminal は Google AI に基づくチャットアプリケーションです。以下の手順で使用することができます。

![](./tui.png)

## インストールと使用方法

1. このリポジトリをクローンします

```bash
git clone https://github.com/greycodee/gemini-terminal.git
```

2. プロジェクトをビルドします

```bash
cd gemini-terminal && go build .
```

3. プロジェクトを実行します

```bash
./gemini-terminal
```

> 注意: `$HOME/.local/share/gemini/config.ini` ファイルで自分の Google AI キーを設定する必要があります。

## 設定

デフォルトの設定ファイルは `$HOME/.local/share/gemini/config.ini` にあります。このファイルで Google AI キーと Gemini モデル名を設定することができます。

```ini
[Gemini]
# 自分のGoogle AIキーを設定します
googleAIKey=
# Geminiモデル名を設定します
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

## チャット履歴

デフォルトのデータベースファイルは `$HOME/.local/share/gemini/gemini.db` にあります。このファイルでチャット履歴を見ることができます。
