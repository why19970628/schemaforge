<template>
  <div class="app-shell" :class="{ dark: theme === 'dark' }">
    <aside class="sidebar">
      <div class="brand">
        <div class="brand-mark">SF</div>
        <div>
          <div class="brand-name">SchemaForge</div>
          <div class="brand-sub">{{ t.brandSub }}</div>
        </div>
      </div>

      <div class="nav-group">
        <span>{{ t.data }}</span>
      </div>
      <button
        v-for="item in dataModes"
        :key="item.mode"
        class="nav-item"
        :class="{ active: mode === item.mode }"
        @click="selectMode(item.mode)"
      >
        <span>{{ labels[item.mode] }}</span>
        <span class="nav-arrow">›</span>
      </button>

      <div class="nav-group">
        <span>SQL</span>
      </div>
      <button
        v-for="item in sqlModes"
        :key="item.mode"
        class="nav-item"
        :class="{ active: mode === item.mode }"
        @click="selectMode(item.mode)"
      >
        <span>{{ labels[item.mode] }}</span>
        <span class="nav-arrow">›</span>
      </button>
    </aside>

    <main class="workspace">
      <header class="topbar">
        <div class="title-block">
          <div class="mode-kicker">
            <span class="mode-dot"></span>
            {{ t.mode }}
          </div>
          <h1>{{ labels[mode] }}</h1>
          <p>{{ t.subtitle }}</p>
        </div>
        <div class="toolbar">
          <el-segmented
            v-if="isJsonDiffMode"
            v-model="jsonDiffFormat"
            :options="diffFormatOptions"
            class="diff-format-switch"
          />
          <el-button @click="loadSample">{{ t.sample }}</el-button>
          <el-button type="success" @click="convert(true)">{{ t.convert }}</el-button>
          <el-button @click="clearAll">{{ t.clear }}</el-button>
        </div>
        <div class="switches">
          <el-button size="small" text @click="toggleLanguage">{{ t.languageSwitch }}</el-button>
          <el-button size="small" text @click="toggleTheme">{{ theme === 'dark' ? t.light : t.dark }}</el-button>
        </div>
      </header>

      <section v-if="isGormMode" class="options-panel">
        <div class="options-head">
          <div>
            <div class="options-title">{{ t.gormOptions }}</div>
            <div class="options-subtitle">{{ t.gormOptionsSub }}</div>
          </div>
          <span class="options-pill">GORM</span>
        </div>
        <el-form :inline="true" label-position="left">
          <el-form-item :label="t.packageName">
            <el-input v-model="options.packageName" placeholder="models" />
          </el-form-item>
          <el-form-item :label="t.nullable">
            <el-select v-model="options.nullable" style="width: 140px">
              <el-option label="zero" value="zero" />
              <el-option label="pointer" value="pointer" />
              <el-option label="sql-null" value="sql-null" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-checkbox v-model="options.jsonTags">{{ t.jsonTags }}</el-checkbox>
          </el-form-item>
          <el-form-item>
            <el-checkbox v-model="options.comments">{{ t.comments }}</el-checkbox>
          </el-form-item>
        </el-form>
      </section>

      <section v-if="isJsonDiffMode" class="editors diff-editors">
        <div class="pane">
          <div class="pane-title">
            <div class="pane-title-main">
              <span>{{ t.leftJson }}</span>
              <span class="language-badge">{{ diffOutputLabel }}</span>
            </div>
            <div class="pane-actions">
              <el-button size="small" class="action-upload" @click="pickFile">{{ t.upload }}</el-button>
              <input ref="fileInput" type="file" hidden @change="readFile" />
            </div>
          </div>
          <CodeEditor
            v-model="input"
            :extensions="inputExtensions"
            @blur="convert(false)"
            @keydown="handleEditorKeydown"
          />
        </div>
        <div class="pane">
          <div class="pane-title">
            <div class="pane-title-main">
              <span>{{ t.rightJson }}</span>
              <span class="language-badge">JSON</span>
            </div>
            <div class="pane-actions">
              <el-button size="small" class="action-upload" @click="pickRightFile">{{ t.upload }}</el-button>
              <input ref="rightFileInput" type="file" hidden @change="readRightFile" />
            </div>
          </div>
          <CodeEditor
            v-model="rightInput"
            :extensions="inputExtensions"
            @blur="convert(false)"
            @keydown="handleEditorKeydown"
          />
        </div>
        <div class="pane">
          <div class="pane-title">
            <div class="pane-title-main">
              <span>{{ t.diffOutput }}</span>
              <span class="language-badge">JSON</span>
            </div>
            <div class="pane-actions">
              <el-button size="small" class="action-copy" @click="copyOutput">{{ t.copy }}</el-button>
            </div>
          </div>
          <CodeEditor
            v-model="output"
            :extensions="outputExtensions"
            :variant="jsonDiffFormat === 'unified' ? 'diff' : ''"
            readonly
          />
        </div>
      </section>

      <section v-else class="editors">
        <div class="pane">
          <div class="pane-title">
            <div class="pane-title-main">
              <span>{{ t.input }}</span>
              <span class="language-badge">{{ inputLanguageLabel }}</span>
            </div>
            <div class="pane-actions">
              <el-button v-if="isSqlMode" size="small" class="action-format" @click="formatSql">
                {{ t.formatSql }}
              </el-button>
              <el-button size="small" class="action-upload" @click="pickFile">{{ t.upload }}</el-button>
              <input ref="fileInput" type="file" hidden @change="readFile" />
            </div>
          </div>
          <CodeEditor
            v-model="input"
            :extensions="inputExtensions"
            @blur="convert(false)"
            @keydown="handleEditorKeydown"
          />
        </div>
        <div class="pane">
          <div class="pane-title">
            <div class="pane-title-main">
              <span>{{ t.output }}</span>
              <span class="language-badge">{{ outputLanguageLabel }}</span>
            </div>
            <div class="pane-actions">
              <el-button size="small" class="action-copy" @click="copyOutput">{{ t.copy }}</el-button>
            </div>
          </div>
          <CodeEditor
            v-model="output"
            :extensions="outputExtensions"
            readonly
          />
        </div>
      </section>

      <section class="bottom">
        <div class="status">
          <span class="status-dot"></span>
          {{ status }}
        </div>
        <div class="history">
          <span v-if="history.length" class="history-label">{{ t.recent }}</span>
          <el-tag
            v-for="item in history"
            :key="item.id"
            class="history-item"
            @click="restoreHistory(item)"
          >
            {{ labels[item.mode] }} · {{ item.time }}
          </el-tag>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { sql } from '@codemirror/lang-sql'
import { json } from '@codemirror/lang-json'
import { xml } from '@codemirror/lang-xml'
import { html } from '@codemirror/lang-html'
import { go } from '@codemirror/lang-go'
import CodeEditor from './components/CodeEditor.vue'

const translations = {
  en: {
    brandSub: 'schema cockpit',
    data: 'Data',
    subtitle: 'Paste structured input, tune options, and convert on blur.',
    formatSql: 'Format SQL',
    sample: 'Sample',
    upload: 'Upload',
    convert: 'Convert',
    copy: 'Copy',
    clear: 'Clear',
    input: 'Input',
    output: 'Output',
    leftJson: 'Left JSON',
    rightJson: 'Right JSON',
    diffOutput: 'Diff',
    structured: 'Structured',
    unified: 'Unified Diff',
    ready: 'Ready',
    converted: 'Converted',
    converting: 'Converting...',
    copied: 'Copied output',
    cleared: 'Cleared',
    fileLoaded: 'File loaded',
    fileFailed: 'Failed to read file',
    mode: 'Current mode',
    recent: 'Recent',
    formatFailed: 'Format failed: paste a complete CREATE TABLE statement',
    formatted: 'SQL formatted',
    languageSwitch: '中文',
    dark: 'Dark',
    light: 'Light',
    gormOptions: 'GORM options',
    gormOptionsSub: 'Applied only when generating GORM structs.',
    packageName: 'Package',
    nullable: 'Nullable',
    jsonTags: 'JSON tags',
    comments: 'Comments',
  },
  zh: {
    brandSub: '结构转换工作台',
    data: '数据',
    subtitle: '粘贴结构输入，调整配置，失焦自动转换。',
    formatSql: '格式化 SQL',
    sample: '示例',
    upload: '上传',
    convert: '转换',
    copy: '复制',
    clear: '清空',
    input: '输入',
    output: '输出',
    leftJson: '左侧 JSON',
    rightJson: '右侧 JSON',
    diffOutput: '差异结果',
    structured: '结构化',
    unified: 'Git Diff',
    ready: '就绪',
    converted: '已转换',
    converting: '转换中...',
    copied: '已复制输出',
    cleared: '已清空',
    fileLoaded: '文件已载入',
    fileFailed: '文件读取失败',
    mode: '当前模式',
    recent: '最近',
    formatFailed: '格式化失败：请粘贴完整 CREATE TABLE 语句',
    formatted: 'SQL 已格式化',
    languageSwitch: 'English',
    dark: '深色',
    light: '浅色',
    gormOptions: 'GORM 配置',
    gormOptionsSub: '仅在生成 GORM 结构体时生效。',
    packageName: '包名',
    nullable: '空值策略',
    jsonTags: 'JSON 标签',
    comments: '注释',
  },
}

const mode = ref('json-go')
const language = ref(localStorage.getItem('schemaforge-language') || 'en')
const theme = ref(localStorage.getItem('schemaforge-theme') || 'light')
const input = ref('')
const rightInput = ref('')
const output = ref('')
const jsonDiffFormat = ref('unified')
const status = ref('')
const fileInput = ref(null)
const rightFileInput = ref(null)
const lastConverted = ref({ mode: '', input: '', right: '', format: '', options: '' })
const history = ref(JSON.parse(localStorage.getItem('schemaforge-history') || '[]'))
const options = ref({
  packageName: 'models',
  nullable: 'zero',
  jsonTags: true,
  comments: true,
})

const samples = {
  'json-go': '{"id": 1, "name": "Ada", "active": true}',
  'json-diff': '{\n  "user": {\n    "id": 1,\n    "name": "Ada",\n    "age": 18,\n    "tags": ["go"]\n  }\n}',
  'yaml-go': 'id: 1\nname: Ada\nactive: true',
  'xml-json': '<user id="1"><name>Ada</name></user>',
  'sql-ent': "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id', `email` varchar(255) NOT NULL COMMENT 'email', `age` int DEFAULT 0, PRIMARY KEY (`id`), UNIQUE KEY `uk_email` (`email`)) COMMENT='users';",
  'sql-gorm': "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id', `email` varchar(255) NOT NULL COMMENT 'email', `created_at` datetime NULL, PRIMARY KEY (`id`), UNIQUE KEY `uk_email` (`email`)) COMMENT='users';",
  'sql-es': "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id', `email` varchar(255) NOT NULL COMMENT 'email', PRIMARY KEY (`id`)) COMMENT='users';",
  'sql-mongo': "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id', `email` varchar(255) NOT NULL COMMENT 'email', PRIMARY KEY (`id`)) COMMENT='users';",
}

const rightSamples = {
  'json-diff': '{\n  "user": {\n    "id": "1",\n    "name": "Grace",\n    "email": "ada@example.com",\n    "tags": ["go", "sql"]\n  }\n}',
}

const modeLabels = {
  en: {
    'json-go': 'JSON to Go',
    'json-diff': 'JSON Diff',
    'yaml-go': 'YAML to Go',
    'xml-json': 'XML to JSON',
    'sql-ent': 'SQL to Ent',
    'sql-gorm': 'SQL to GORM',
    'sql-es': 'SQL to ES',
    'sql-mongo': 'SQL to MongoDB',
  },
  zh: {
    'json-go': 'JSON 转 Go',
    'json-diff': 'JSON 对比',
    'yaml-go': 'YAML 转 Go',
    'xml-json': 'XML 转 JSON',
    'sql-ent': 'SQL 转 Ent',
    'sql-gorm': 'SQL 转 GORM',
    'sql-es': 'SQL 转 ES',
    'sql-mongo': 'SQL 转 MongoDB',
  },
}

const dataModes = [{ mode: 'json-go' }, { mode: 'json-diff' }, { mode: 'yaml-go' }, { mode: 'xml-json' }]
const sqlModes = [{ mode: 'sql-ent' }, { mode: 'sql-gorm' }, { mode: 'sql-es' }, { mode: 'sql-mongo' }]
const t = computed(() => translations[language.value])
const labels = computed(() => modeLabels[language.value])
const isSqlMode = computed(() => mode.value.startsWith('sql-'))
const isGormMode = computed(() => mode.value === 'sql-gorm')
const isJsonDiffMode = computed(() => mode.value === 'json-diff')
const activeOptions = computed(() => (isGormMode.value ? options.value : {}))
const optionKey = computed(() => JSON.stringify(activeOptions.value))
const inputExtensions = computed(() => extensionsFor(inputLanguageForMode(mode.value)))
const outputExtensions = computed(() => extensionsFor(outputLanguageForMode(mode.value)))
const inputLanguageLabel = computed(() => inputLanguageForMode(mode.value).toUpperCase())
const outputLanguageLabel = computed(() => outputLanguageForMode(mode.value).toUpperCase())
const diffOutputLabel = computed(() => (jsonDiffFormat.value === 'unified' ? 'DIFF' : 'JSON'))
const diffFormatOptions = computed(() => [
  { label: t.value.unified, value: 'unified' },
  { label: t.value.structured, value: 'structured' },
])

function extensionsFor(kind) {
  if (kind === 'sql') return [sql()]
  if (kind === 'json') return [json()]
  if (kind === 'xml') return [xml()]
  if (kind === 'html') return [html()]
  if (kind === 'go') return [go()]
  return []
}

function inputLanguageForMode(value) {
  if (value.startsWith('sql-')) return 'sql'
  if (value === 'json-go' || value === 'json-diff') return 'json'
  if (value === 'yaml-go') return 'yaml'
  if (value === 'xml-json') return 'xml'
  return 'text'
}

function outputLanguageForMode(value) {
  if (value === 'json-diff') return jsonDiffFormat.value === 'unified' ? 'text' : 'json'
  if (value === 'sql-es' || value === 'sql-mongo' || value === 'xml-json') return 'json'
  return 'go'
}

function selectMode(next) {
  mode.value = next
  input.value = samples[next] || ''
  rightInput.value = rightSamples[next] || ''
  output.value = ''
  resetConverted()
  status.value = t.value.ready
}

function loadSample() {
  input.value = samples[mode.value] || ''
  rightInput.value = rightSamples[mode.value] || ''
  output.value = ''
  resetConverted()
}

function clearAll() {
  input.value = ''
  rightInput.value = ''
  output.value = ''
  resetConverted()
  status.value = t.value.cleared
}

function pickFile() {
  fileInput.value?.click()
}

function pickRightFile() {
  rightFileInput.value?.click()
}

async function readFile(event) {
  const file = event.target.files?.[0]
  if (!file) return
  try {
    input.value = await file.text()
    output.value = ''
    resetConverted()
    status.value = t.value.fileLoaded
  } catch {
    status.value = t.value.fileFailed
  } finally {
    event.target.value = ''
  }
}

async function readRightFile(event) {
  const file = event.target.files?.[0]
  if (!file) return
  try {
    rightInput.value = await file.text()
    output.value = ''
    resetConverted()
    status.value = t.value.fileLoaded
  } catch {
    status.value = t.value.fileFailed
  } finally {
    event.target.value = ''
  }
}

async function convert(force) {
  const value = input.value.trim()
  if (!value) return
  const rightValue = rightInput.value.trim()
  if (isJsonDiffMode.value && !rightValue) return
  if (
    !force &&
    lastConverted.value.mode === mode.value &&
    lastConverted.value.input === value &&
    lastConverted.value.right === rightValue &&
    lastConverted.value.format === jsonDiffFormat.value &&
    lastConverted.value.options === optionKey.value
  ) return
  status.value = t.value.converting
  output.value = ''
  const response = await fetch('/api/convert', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      mode: mode.value,
      input: value,
      right: rightValue,
      format: isJsonDiffMode.value ? jsonDiffFormat.value : '',
      options: activeOptions.value,
    }),
  })
  const payload = await response.json()
  if (!response.ok) {
    status.value = payload.error || 'Conversion failed'
    return
  }
  output.value = payload.output
  lastConverted.value = { mode: mode.value, input: value, right: rightValue, format: jsonDiffFormat.value, options: optionKey.value }
  status.value = t.value.converted
  pushHistory()
}

function pushHistory() {
  const item = {
    id: `${Date.now()}-${Math.random().toString(16).slice(2)}`,
    mode: mode.value,
    input: input.value,
    right: rightInput.value,
    output: output.value,
    options: { ...options.value },
    time: new Date().toLocaleTimeString(),
  }
  history.value = [item, ...history.value].slice(0, 10)
  localStorage.setItem('schemaforge-history', JSON.stringify(history.value))
}

function restoreHistory(item) {
  mode.value = item.mode
  input.value = item.input
  rightInput.value = item.right || ''
  output.value = item.output
  options.value = { ...options.value, ...item.options }
  lastConverted.value = {
    mode: item.mode,
    input: item.input.trim(),
    right: rightInput.value.trim(),
    format: jsonDiffFormat.value,
    options: JSON.stringify(activeOptions.value),
  }
}

function copyOutput() {
  navigator.clipboard.writeText(output.value)
  ElMessage.success(t.value.copied)
}

function formatSql() {
  const formatted = formatCreateTableSQL(input.value)
  if (!formatted) {
    status.value = t.value.formatFailed
    return
  }
  input.value = formatted
  resetConverted()
  status.value = t.value.formatted
}

function handleEditorKeydown(event) {
  if (event.key === 'Enter' && (event.metaKey || event.ctrlKey)) {
    event.preventDefault()
    convert(true)
  }
}

function toggleLanguage() {
  language.value = language.value === 'zh' ? 'en' : 'zh'
  localStorage.setItem('schemaforge-language', language.value)
}

function toggleTheme() {
  theme.value = theme.value === 'dark' ? 'light' : 'dark'
  localStorage.setItem('schemaforge-theme', theme.value)
}

function resetConverted() {
  lastConverted.value = { mode: '', input: '', right: '', format: '', options: '' }
}

function formatCreateTableSQL(sqlText) {
  const trimmed = sqlText.trim()
  const open = trimmed.indexOf('(')
  const close = trimmed.lastIndexOf(')')
  if (!/^create\s+table/i.test(trimmed) || open < 0 || close < open || !trimmed.endsWith(';')) return ''
  const header = normalizeSQLKeywords(trimmed.slice(0, open).trim())
  const body = trimmed.slice(open + 1, close)
  const suffix = normalizeSQLKeywords(trimmed.slice(close + 1).trim().replace(/;$/, ''))
  const definitions = splitSQLDefinitions(body)
  if (definitions.length === 0) return ''
  const lines = definitions.map((part) => `  ${normalizeSQLKeywords(part.trim().replace(/,$/, ''))},`)
  return `${header} (\n${lines.join('\n')}\n)${suffix ? ` ${suffix}` : ''};`
}

function splitSQLDefinitions(body) {
  const parts = []
  let current = ''
  let depth = 0
  let quote = ''
  for (let index = 0; index < body.length; index += 1) {
    const ch = body[index]
    if (quote) {
      current += ch
      if (ch === quote && body[index - 1] !== '\\') quote = ''
      continue
    }
    if (ch === "'" || ch === '"' || ch === '`') {
      quote = ch
      current += ch
      continue
    }
    if (ch === '(') depth += 1
    if (ch === ')') depth -= 1
    if (ch === ',' && depth === 0) {
      if (current.trim()) parts.push(current)
      current = ''
      continue
    }
    current += ch
  }
  if (current.trim()) parts.push(current)
  return parts
}

function normalizeSQLKeywords(value) {
  return value
    .replace(/\bcreate\s+table\b/gi, 'CREATE TABLE')
    .replace(/\bprimary\s+key\b/gi, 'PRIMARY KEY')
    .replace(/\bunique\s+key\b/gi, 'UNIQUE KEY')
    .replace(/\bnot\s+null\b/gi, 'NOT NULL')
    .replace(/\bdefault\b/gi, 'DEFAULT')
    .replace(/\bauto_increment\b/gi, 'AUTO_INCREMENT')
    .replace(/\bcomment\b/gi, 'COMMENT')
    .replace(/\bengine\b/gi, 'ENGINE')
    .replace(/\bcharset\b/gi, 'CHARSET')
    .replace(/\bcollate\b/gi, 'COLLATE')
}

watch(options, resetConverted, { deep: true })
watch(jsonDiffFormat, resetConverted)
selectMode(mode.value)
</script>

<style scoped>
.app-shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 272px 1fr;
  background:
    radial-gradient(circle at 70% -10%, rgba(9, 105, 218, 0.08), transparent 28%),
    linear-gradient(180deg, #f6f8fa 0%, #eef2f6 100%);
  color: #24292f;
}

.app-shell.dark {
  background:
    radial-gradient(circle at 70% -10%, rgba(88, 166, 255, 0.14), transparent 28%),
    linear-gradient(180deg, #0d1117 0%, #090c10 100%);
  color: #e6edf3;
}

.sidebar {
  padding: 24px 16px;
  background: rgba(255, 255, 255, 0.9);
  border-right: 1px solid #d0d7de;
  box-shadow: 1px 0 0 rgba(27, 31, 36, 0.04);
}

.dark .sidebar,
.dark .pane,
.dark .options-panel,
.dark .bottom,
.dark .switches {
  background: rgba(22, 27, 34, 0.94);
  border-color: #30363d;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 28px;
}

.brand-mark {
  display: grid;
  place-items: center;
  width: 42px;
  height: 42px;
  border-radius: 12px;
  color: #fff;
  font-weight: 800;
  background: linear-gradient(135deg, #0969da, #2da44e);
  box-shadow: 0 10px 24px rgba(9, 105, 218, 0.22);
}

.brand-name {
  font-size: 16px;
  font-weight: 800;
  letter-spacing: 0;
}

.brand-sub,
.nav-group,
.topbar p,
.pane-title,
.options-title,
.options-subtitle,
.status {
  color: #57606a;
}

.dark .brand-sub,
.dark .nav-group,
.dark .topbar p,
.dark .pane-title,
.dark .options-title,
.dark .options-subtitle,
.dark .status {
  color: #8b949e;
}

.nav-group {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin: 24px 8px 9px;
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
}

.nav-item {
  width: 100%;
  margin: 5px 0;
  padding: 10px 12px;
  border: 1px solid #d0d7de;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.72);
  color: inherit;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  cursor: pointer;
  transition: background 0.15s ease, border-color 0.15s ease, box-shadow 0.15s ease, transform 0.15s ease;
}

.nav-arrow {
  opacity: 0;
  color: #0969da;
  transform: translateX(-4px);
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.nav-item.active,
.nav-item:hover {
  border-color: #0969da;
  background: #ddf4ff;
  box-shadow: 0 6px 18px rgba(9, 105, 218, 0.1);
  transform: translateX(2px);
}

.nav-item.active .nav-arrow,
.nav-item:hover .nav-arrow {
  opacity: 1;
  transform: translateX(0);
}

.dark .nav-item {
  border-color: #30363d;
  background: rgba(13, 17, 23, 0.72);
}

.dark .nav-item.active,
.dark .nav-item:hover {
  border-color: #58a6ff;
  background: #13233a;
}

.workspace {
  min-width: 0;
  padding: 24px;
  display: grid;
  grid-template-rows: auto auto 1fr auto;
  gap: 16px;
}

.topbar {
  display: grid;
  grid-template-columns: minmax(260px, 1fr) auto minmax(160px, 1fr);
  align-items: center;
  gap: 18px;
  padding: 4px 2px 2px;
}

.mode-kicker {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  margin-bottom: 4px;
  color: #57606a;
  font-size: 12px;
  font-weight: 700;
}

.mode-dot,
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #2da44e;
  box-shadow: 0 0 0 3px rgba(45, 164, 78, 0.13);
}

.dark .mode-kicker {
  color: #8b949e;
}

.topbar h1 {
  margin: 0;
  font-size: 28px;
  line-height: 1.15;
}

.topbar p {
  margin: 4px 0 0;
}

.toolbar {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  justify-content: center;
  padding: 8px;
  border: 1px solid rgba(208, 215, 222, 0.8);
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.74);
  box-shadow: 0 10px 30px rgba(27, 31, 36, 0.06);
}

.diff-format-switch {
  --el-segmented-item-selected-bg-color: #0969da;
  --el-segmented-item-selected-color: #ffffff;
  --el-segmented-bg-color: #f6f8fa;
}

.dark .diff-format-switch {
  --el-segmented-item-selected-bg-color: #1f6feb;
  --el-segmented-bg-color: #0d1117;
}

.dark .toolbar {
  border-color: #30363d;
  background: rgba(22, 27, 34, 0.78);
  box-shadow: none;
}

.switches {
  justify-self: end;
  display: flex;
  gap: 4px;
  padding: 2px;
  border: 1px solid #d0d7de;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.82);
}

.options-panel,
.bottom,
.pane {
  border: 1px solid #d0d7de;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.94);
  box-shadow: 0 12px 32px rgba(27, 31, 36, 0.07);
}

.options-panel {
  padding: 12px 14px 0;
}

.options-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 10px;
}

.options-title {
  margin-bottom: 2px;
  font-size: 12px;
  font-weight: 700;
}

.options-subtitle {
  font-size: 12px;
}

.options-pill,
.language-badge {
  display: inline-flex;
  align-items: center;
  height: 22px;
  padding: 0 8px;
  border: 1px solid #d0d7de;
  border-radius: 999px;
  color: #0969da;
  background: #ddf4ff;
  font-size: 11px;
  font-weight: 800;
}

.dark .options-pill,
.dark .language-badge {
  border-color: #30363d;
  color: #58a6ff;
  background: #13233a;
}

.editors {
  min-height: 0;
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(0, 1fr);
  gap: 18px;
}

.diff-editors {
  grid-template-columns: minmax(0, 0.9fr) minmax(0, 0.9fr) minmax(0, 1.1fr);
}

.pane {
  min-width: 0;
  overflow: hidden;
  box-shadow: 0 16px 38px rgba(27, 31, 36, 0.08);
}

.pane-title {
  padding: 11px 12px;
  border-bottom: 1px solid #d0d7de;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-weight: 700;
  background: linear-gradient(180deg, #ffffff 0%, #f6f8fa 100%);
}

.pane-title-main,
.pane-actions {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.pane-actions {
  flex-wrap: wrap;
  justify-content: flex-end;
}

.pane-actions :deep(.el-button) {
  border-radius: 999px;
  font-weight: 700;
}

.pane-actions :deep(.action-format) {
  border-color: #bf8700;
  color: #7d4e00;
  background: #fff8c5;
}

.pane-actions :deep(.action-format:hover) {
  border-color: #9a6700;
  color: #633c01;
  background: #fae17d;
}

.pane-actions :deep(.action-upload) {
  border-color: #0969da;
  color: #0969da;
  background: #ddf4ff;
}

.pane-actions :deep(.action-upload:hover) {
  border-color: #0550ae;
  color: #0550ae;
  background: #b6e3ff;
}

.pane-actions :deep(.action-copy) {
  border-color: #2da44e;
  color: #1a7f37;
  background: #dafbe1;
}

.pane-actions :deep(.action-copy:hover) {
  border-color: #1a7f37;
  color: #116329;
  background: #aceebb;
}

.dark .pane-actions :deep(.action-format) {
  border-color: #9e6a03;
  color: #f0b72f;
  background: #2d2105;
}

.dark .pane-actions :deep(.action-upload) {
  border-color: #1f6feb;
  color: #79c0ff;
  background: #0d274d;
}

.dark .pane-actions :deep(.action-copy) {
  border-color: #238636;
  color: #7ee787;
  background: #0f2a17;
}

.dark .pane-title {
  border-color: #30363d;
  background: linear-gradient(180deg, #161b22 0%, #0d1117 100%);
}

.bottom {
  min-height: 42px;
  padding: 9px 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.status {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
}

.history {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: flex-end;
}

.history-label {
  color: #57606a;
  font-size: 12px;
  font-weight: 700;
}

.dark .history-label {
  color: #8b949e;
}

.history-item {
  cursor: pointer;
}

@media (max-width: 1100px) {
  .app-shell {
    grid-template-columns: 1fr;
  }

  .topbar,
  .editors,
  .diff-editors {
    grid-template-columns: 1fr;
  }

  .switches,
  .toolbar {
    justify-self: start;
    justify-content: flex-start;
  }

  .toolbar {
    width: 100%;
  }
}
</style>
