package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ini/ini"
	"github.com/google/generative-ai-go/genai"
)

const (
	CONFIG_PATH = "/.local/share/gemini/"
	CONFIG_NAME = "config.ini"
)

type GeminiChatConfig struct {
	ModelName     string
	GoogleAIKey   string
	SafetySetting []*genai.SafetySetting
}

func GetConfig() GeminiChatConfig {
	cfg := loadConfig()
	if cfg == nil {
		log.Fatal("load config failed")
	}

	level := cfg.Section("SafetySetting").Key("level").MustInt(0)
	googleAIKey := cfg.Section("Gemini").Key("googleAIKey").Value()
	if googleAIKey == "" {
		log.Fatalf("\x1b[31mPlease set your own google ai key in %s \n\x1b[0m",
			HOME_PATH+CONFIG_PATH+CONFIG_NAME)
	}

	var safetyLevel genai.HarmBlockThreshold = genai.HarmBlockThreshold(level)

	return GeminiChatConfig{
		ModelName:   cfg.Section("Gemini").Key("model").Value(),
		GoogleAIKey: googleAIKey,
		SafetySetting: []*genai.SafetySetting{
			{
				Category:  genai.HarmCategorySexuallyExplicit,
				Threshold: safetyLevel,
			},
			{
				Category:  genai.HarmCategoryHarassment,
				Threshold: safetyLevel,
			},
			{
				Category:  genai.HarmCategoryHateSpeech,
				Threshold: safetyLevel,
			},
			{
				Category:  genai.HarmCategoryDangerousContent,
				Threshold: safetyLevel,
			},
		},
	}
}

func loadConfig() *ini.File {
	fullPath := HOME_PATH + CONFIG_PATH + CONFIG_NAME
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		os.MkdirAll(HOME_PATH+CONFIG_PATH, os.ModePerm)
		initConfig := []byte(`
[Gemini]
# set your own google ai key
googleAIKey=
# set gemini model name
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

`)
		err := os.WriteFile(fullPath, initConfig, 0644)
		if err != nil {
			fmt.Println("初始化配置文件失败：", err)
			log.Fatal(err)
		}
		log.Fatalf("\x1b[31mPlease set your own google ai key in %s \n\x1b[0m", fullPath)
	} else {
		cfg, err := ini.Load(fullPath)
		if err != nil {
			// 处理错误
			log.Fatal(err)
		}
		return cfg
	}
	return nil
}
