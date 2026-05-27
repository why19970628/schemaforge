const labels = {
  en: {
    "json-go": "JSON to Go",
    "yaml-go": "YAML to Go",
    "xml-json": "XML to JSON",
    "sql-ent": "SQL to Ent",
    "sql-gorm": "SQL to GORM",
    "sql-es": "SQL to ES",
    "sql-mongo": "SQL to MongoDB"
  },
  zh: {
    "json-go": "JSON 转 Go",
    "yaml-go": "YAML 转 Go",
    "xml-json": "XML 转 JSON",
    "sql-ent": "SQL 转 Ent",
    "sql-gorm": "SQL 转 GORM",
    "sql-es": "SQL 转 ES",
    "sql-mongo": "SQL 转 MongoDB"
  }
};

const translations = {
  en: {
    brandSub: "schema cockpit",
    dataGroup: "Data",
    subtitle: "Paste a structure sample, brew production-ready schema code.",
    formatSQL: "Format SQL",
    sample: "Sample",
    upload: "Upload",
    convert: "Convert",
    copy: "Copy",
    clear: "Clear",
    dark: "Dark",
    light: "Light",
    language: "中文",
    inputLabel: "Input",
    outputLabel: "Output",
    ready: "Ready",
    cleared: "Cleared",
    copied: "Copied output",
    brewing: "Brewing...",
    converted: "Converted",
    sqlOnly: "SQL formatting is available in SQL modes",
    formatFailed: "Format failed: paste a complete CREATE TABLE statement",
    formatted: "SQL formatted",
    convertFailed: "Conversion failed",
    fileLoaded: "File loaded",
    fileFailed: "Failed to read file",
    shortcutHint: "Tip: Cmd/Ctrl + Enter converts the current input"
  },
  zh: {
    brandSub: "结构转换工作台",
    dataGroup: "数据",
    subtitle: "粘贴结构样例，生成可用的 Schema 或 Go 代码。",
    formatSQL: "格式化 SQL",
    sample: "示例",
    upload: "上传",
    convert: "转换",
    copy: "复制",
    clear: "清空",
    dark: "深色",
    light: "浅色",
    language: "English",
    inputLabel: "输入",
    outputLabel: "输出",
    ready: "就绪",
    cleared: "已清空",
    copied: "已复制输出",
    brewing: "转换中...",
    converted: "已转换",
    sqlOnly: "SQL 格式化仅支持 SQL 模式",
    formatFailed: "格式化失败：请粘贴完整的 CREATE TABLE 语句",
    formatted: "SQL 已格式化",
    convertFailed: "转换失败",
    fileLoaded: "文件已载入",
    fileFailed: "文件读取失败",
    shortcutHint: "提示：Cmd/Ctrl + Enter 可转换当前输入"
  }
};

const samples = {
  "json-go": '{"id": 1, "name": "Ada", "active": true}',
  "yaml-go": 'id: 1\nname: Ada\nactive: true',
  "xml-json": '<user id="1"><name>Ada</name></user>',
  "sql-ent": "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id', `email` varchar(255) NOT NULL COMMENT 'email', `age` int DEFAULT 0, PRIMARY KEY (`id`), UNIQUE KEY `uk_email` (`email`)) COMMENT='users';",
  "sql-gorm": "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id', `email` varchar(255) NOT NULL COMMENT 'email', `age` int DEFAULT 0, PRIMARY KEY (`id`), UNIQUE KEY `uk_email` (`email`)) COMMENT='users';",
  "sql-es": "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id', `email` varchar(255) NOT NULL COMMENT 'email', `age` int DEFAULT 0, PRIMARY KEY (`id`), UNIQUE KEY `uk_email` (`email`)) COMMENT='users';",
  "sql-mongo": "CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id', `email` varchar(255) NOT NULL COMMENT 'email', `age` int DEFAULT 0, PRIMARY KEY (`id`), UNIQUE KEY `uk_email` (`email`)) COMMENT='users';"
};

let mode = "json-go";
let language = localStorage.getItem("schemaforge-language") || "en";
let lastConverted = { mode: "", input: "" };
const input = document.querySelector("#input");
const output = document.querySelector("#output");
const status = document.querySelector("#status");
const formatButton = document.querySelector("#format-sql");
const themeButton = document.querySelector("#theme");
const languageButton = document.querySelector("#language");
const fileInput = document.querySelector("#file-input");

function t(key) {
  return translations[language][key] || translations.en[key] || key;
}

function isSQLMode(value) {
  return value === "sql-ent" || value === "sql-gorm" || value === "sql-es" || value === "sql-mongo";
}

function applyLanguage() {
  document.documentElement.lang = language === "zh" ? "zh-CN" : "en";
  document.querySelectorAll("[data-i18n]").forEach((node) => {
    node.textContent = t(node.dataset.i18n);
  });
  document.querySelectorAll("nav button[data-mode]").forEach((button) => {
    button.textContent = labels[language][button.dataset.mode];
  });
  document.querySelector("#title").textContent = labels[language][mode];
  document.querySelector("#subtitle").textContent = t("subtitle");
  formatButton.textContent = t("formatSQL");
  document.querySelector("#sample").textContent = t("sample");
  document.querySelector("#upload").textContent = t("upload");
  document.querySelector("#convert").textContent = t("convert");
  document.querySelector("#copy").textContent = t("copy");
  document.querySelector("#clear").textContent = t("clear");
  languageButton.textContent = t("language");
  updateThemeButton();
}

function updateThemeButton() {
  themeButton.textContent = document.body.dataset.theme === "dark" ? t("light") : t("dark");
}

function setMode(next) {
  mode = next;
  document.querySelector("#title").textContent = labels[language][mode];
  document.querySelectorAll("nav button").forEach((button) => {
    button.classList.toggle("active", button.dataset.mode === mode);
  });
  const showFormatter = isSQLMode(mode);
  formatButton.classList.toggle("hidden", !showFormatter);
  formatButton.disabled = !showFormatter;
  input.value = samples[mode] || "";
  output.value = "";
  lastConverted = { mode: "", input: "" };
  status.textContent = t("ready");
  resizeTextareas();
}

document.querySelectorAll("nav button").forEach((button) => {
  button.addEventListener("click", () => setMode(button.dataset.mode));
});

document.querySelector("#sample").addEventListener("click", () => {
  input.value = samples[mode] || "";
  lastConverted = { mode: "", input: "" };
  resizeTextareas();
});

document.querySelector("#upload").addEventListener("click", () => {
  fileInput.click();
});

fileInput.addEventListener("change", async () => {
  const file = fileInput.files && fileInput.files[0];
  if (!file) return;
  try {
    input.value = await file.text();
    output.value = "";
    lastConverted = { mode: "", input: "" };
    status.textContent = t("fileLoaded");
    resizeTextareas();
  } catch {
    status.textContent = t("fileFailed");
  } finally {
    fileInput.value = "";
  }
});

document.querySelector("#clear").addEventListener("click", () => {
  input.value = "";
  output.value = "";
  lastConverted = { mode: "", input: "" };
  status.textContent = t("cleared");
  resizeTextareas();
});

formatButton.addEventListener("click", () => {
  if (!isSQLMode(mode)) {
    status.textContent = t("sqlOnly");
    return;
  }
  const formatted = formatCreateTableSQL(input.value);
  if (!formatted) {
    status.textContent = t("formatFailed");
    return;
  }
  input.value = formatted;
  lastConverted = { mode: "", input: "" };
  status.textContent = t("formatted");
  resizeTextareas();
});

document.querySelector("#copy").addEventListener("click", async () => {
  await navigator.clipboard.writeText(output.value);
  status.textContent = t("copied");
});

async function runConvert({ force = false } = {}) {
  const value = input.value.trim();
  if (!value) return;
  if (!force && lastConverted.mode === mode && lastConverted.input === value) return;
  status.textContent = t("brewing");
  output.value = "";
  const response = await fetch("/api/convert", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ mode, input: value })
  });
  const payload = await response.json();
  if (!response.ok) {
    status.textContent = payload.error || t("convertFailed");
    return;
  }
  output.value = payload.output;
  lastConverted = { mode, input: value };
  status.textContent = t("converted");
  resizeTextareas();
}

document.querySelector("#convert").addEventListener("click", async () => {
  await runConvert({ force: true });
});

input.addEventListener("blur", async () => {
  await runConvert();
});

input.addEventListener("keydown", async (event) => {
  if (event.key === "Tab") {
    event.preventDefault();
    insertAtCursor(input, "  ");
    resizeTextareas();
    return;
  }
  if (event.key === "Enter" && (event.metaKey || event.ctrlKey)) {
    event.preventDefault();
    await runConvert({ force: true });
  }
});

input.addEventListener("input", resizeTextareas);

themeButton.addEventListener("click", () => {
  const next = document.body.dataset.theme === "dark" ? "light" : "dark";
  document.body.dataset.theme = next;
  localStorage.setItem("schemaforge-theme", next);
  updateThemeButton();
});

languageButton.addEventListener("click", () => {
  language = language === "zh" ? "en" : "zh";
  localStorage.setItem("schemaforge-language", language);
  applyLanguage();
  status.textContent = t("ready");
});

function insertAtCursor(textarea, text) {
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  textarea.value = textarea.value.slice(0, start) + text + textarea.value.slice(end);
  textarea.selectionStart = textarea.selectionEnd = start + text.length;
}

function resizeTextareas() {
  [input, output].forEach((textarea) => {
    textarea.style.height = "auto";
    textarea.style.height = Math.max(560, textarea.scrollHeight) + "px";
  });
}

function formatCreateTableSQL(sql) {
  const trimmed = sql.trim();
  const open = trimmed.indexOf("(");
  const close = trimmed.lastIndexOf(")");
  if (!/^create\s+table/i.test(trimmed) || open < 0 || close < open || !trimmed.endsWith(";")) {
    return "";
  }
  const header = normalizeSQLKeywords(trimmed.slice(0, open).trim());
  const body = trimmed.slice(open + 1, close);
  const suffix = normalizeSQLKeywords(trimmed.slice(close + 1).trim().replace(/;$/, ""));
  const definitions = splitSQLDefinitions(body);
  if (definitions.length === 0) return "";
  const lines = definitions.map((part) => "  " + normalizeSQLKeywords(part.trim().replace(/,$/, "")) + ",");
  return `${header} (\n${lines.join("\n")}\n)${suffix ? " " + suffix : ""};`;
}

function splitSQLDefinitions(body) {
  const parts = [];
  let current = "";
  let depth = 0;
  let quote = "";
  for (let index = 0; index < body.length; index += 1) {
    const ch = body[index];
    if (quote) {
      current += ch;
      if (ch === quote && body[index - 1] !== "\\") quote = "";
      continue;
    }
    if (ch === "'" || ch === '"' || ch === "`") {
      quote = ch;
      current += ch;
      continue;
    }
    if (ch === "(") depth += 1;
    if (ch === ")") depth -= 1;
    if (ch === "," && depth === 0) {
      if (current.trim()) parts.push(current);
      current = "";
      continue;
    }
    current += ch;
  }
  if (current.trim()) parts.push(current);
  return parts;
}

function normalizeSQLKeywords(value) {
  return value
    .replace(/\bcreate\s+table\b/gi, "CREATE TABLE")
    .replace(/\bprimary\s+key\b/gi, "PRIMARY KEY")
    .replace(/\bunique\s+key\b/gi, "UNIQUE KEY")
    .replace(/\bunique\s+index\b/gi, "UNIQUE INDEX")
    .replace(/\bnot\s+null\b/gi, "NOT NULL")
    .replace(/\bdefault\b/gi, "DEFAULT")
    .replace(/\bauto_increment\b/gi, "AUTO_INCREMENT")
    .replace(/\bcomment\b/gi, "COMMENT")
    .replace(/\bengine\b/gi, "ENGINE")
    .replace(/\bcharset\b/gi, "CHARSET")
    .replace(/\bcollate\b/gi, "COLLATE");
}

const savedTheme = localStorage.getItem("schemaforge-theme") || "light";
document.body.dataset.theme = savedTheme;
applyLanguage();
setMode(mode);
status.textContent = t("shortcutHint");
