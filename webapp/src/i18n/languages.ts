/**
 * 语言配置文件
 * 
 * 这个文件集中管理所有支持的语言配置，便于维护和扩展。
 * 
 * 如何新增语言：
 * 1. 在此文件的 languages 数组中添加新语言配置
 * 2. 在 src/i18n/locales/ 目录下创建对应的翻译文件（如：ja-JP.ts）
 * 3. 在 src/i18n/index.ts 中导入新的翻译文件并添加到 messages 中
 * 4. 完成！语言切换器会自动显示新语言选项
 */

export interface Language {
  /** 语言代码（BCP 47 格式），如：'zh-CN', 'en-US', 'ja-JP' */
  code: string
  /** 语言的本地化名称（用该语言自身的文字显示） */
  name: string
  /** 语言对应的旗帜 Emoji */
  flag: string
  /** 语言的英文名称（可选，用于调试或文档） */
  englishName?: string
}

/**
 * 支持的语言列表
 * 
 * 数组顺序决定了语言切换器中的显示顺序
 */
export const languages: Language[] = [
  {
    code: 'zh-CN',
    name: '简体中文',
    flag: '🇨🇳',
    englishName: 'Simplified Chinese',
  },
  {
    code: 'zh-TW',
    name: '繁體中文',
    flag: '🇨🇳',
    englishName: 'Traditional Chinese',
  },
  {
    code: 'en-US',
    name: 'English',
    flag: '🇺🇸',
    englishName: 'English',
  },
  {
    code: 'ja-JP',
    name: '日本語',
    flag: '🇯🇵',
    englishName: 'Japanese',
  },
  {
    code: 'ko-KR',
    name: '한국어',
    flag: '🇰🇷',
    englishName: 'Korean',
  },
  {
    code: 'ru-RU',
    name: 'Русский',
    flag: '🇷🇺',
    englishName: 'Russian',
  },
]

/**
 * 默认语言代码
 */
export const DEFAULT_LANGUAGE = 'zh-CN'

/**
 * 根据语言代码获取语言配置
 */
export function getLanguage(code: string): Language | undefined {
  return languages.find(lang => lang.code === code)
}

/**
 * 获取语言名称
 */
export function getLanguageName(code: string): string {
  const lang = getLanguage(code)
  return lang ? lang.name : code
}

/**
 * 检查语言代码是否受支持
 */
export function isLanguageSupported(code: string): boolean {
  return languages.some(lang => lang.code === code)
}
