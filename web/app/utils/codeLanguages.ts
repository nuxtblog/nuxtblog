/**
 * Default code languages shown in the editor code block dropdown.
 * Users can customize this list via Writing Settings.
 */
export const DEFAULT_CODE_LANGUAGES: Array<{ label: string; value: string }> = [
  { label: 'Plain', value: '' },
  { label: 'JavaScript', value: 'javascript' },
  { label: 'TypeScript', value: 'typescript' },
  { label: 'Python', value: 'python' },
  { label: 'Go', value: 'go' },
  { label: 'Rust', value: 'rust' },
  { label: 'HTML', value: 'html' },
  { label: 'CSS', value: 'css' },
  { label: 'JSON', value: 'json' },
  { label: 'YAML', value: 'yaml' },
  { label: 'Bash', value: 'bash' },
  { label: 'SQL', value: 'sql' },
  { label: 'Markdown', value: 'markdown' },
  { label: 'Java', value: 'java' },
  { label: 'C++', value: 'cpp' },
  { label: 'C#', value: 'csharp' },
  { label: 'PHP', value: 'php' },
  { label: 'Ruby', value: 'ruby' },
  { label: 'Swift', value: 'swift' },
  { label: 'Kotlin', value: 'kotlin' },
  { label: 'Lua', value: 'lua' },
  { label: 'Dockerfile', value: 'dockerfile' },
  { label: 'XML', value: 'xml' },
  { label: 'TOML', value: 'toml' },
]

/**
 * Capitalize first letter of a language name for display.
 */
export function formatLanguageLabel(value: string): string {
  if (!value) return 'Plain'
  // Known display names for languages whose hljs id differs from the common name
  const labelMap: Record<string, string> = {
    cpp: 'C++',
    csharp: 'C#',
    objectivec: 'Objective-C',
    plaintext: 'Plain Text',
    vbnet: 'VB.NET',
    fsharp: 'F#',
  }
  if (labelMap[value]) return labelMap[value]
  return value.charAt(0).toUpperCase() + value.slice(1)
}
