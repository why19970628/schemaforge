<template>
  <div ref="host" class="code-editor"></div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { EditorState, Compartment } from '@codemirror/state'
import { Decoration, EditorView, keymap, placeholder } from '@codemirror/view'
import { defaultKeymap, history, historyKeymap } from '@codemirror/commands'
import { HighlightStyle, defaultHighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { tags as t } from '@lezer/highlight'

const props = defineProps({
  modelValue: { type: String, default: '' },
  extensions: { type: Array, default: () => [] },
  readonly: { type: Boolean, default: false },
  placeholder: { type: String, default: '' },
  variant: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue', 'blur', 'keydown'])

const host = ref(null)
let view
const languageCompartment = new Compartment()
const editableCompartment = new Compartment()
const variantCompartment = new Compartment()
const schemaforgeHighlight = HighlightStyle.define([
  { tag: [t.keyword, t.operatorKeyword, t.modifier], color: '#cf222e', fontWeight: '600' },
  { tag: [t.string, t.special(t.string)], color: '#0a3069' },
  { tag: [t.number, t.bool, t.null], color: '#0550ae' },
  { tag: [t.comment, t.lineComment, t.blockComment], color: '#6e7781', fontStyle: 'italic' },
  { tag: [t.definition(t.variableName), t.className, t.typeName], color: '#8250df' },
  { tag: [t.propertyName, t.attributeName], color: '#953800' },
  { tag: [t.function(t.variableName), t.function(t.propertyName)], color: '#8250df' },
])
const diffLineTheme = EditorView.baseTheme({
  '.cm-line.diff-added': { backgroundColor: '#dafbe1', color: '#116329' },
  '.cm-line.diff-removed': { backgroundColor: '#ffebe9', color: '#82071e' },
  '.cm-line.diff-meta': { backgroundColor: '#ddf4ff', color: '#0550ae', fontWeight: '700' },
  '&dark .cm-line.diff-added': { backgroundColor: '#0f2a17', color: '#7ee787' },
  '&dark .cm-line.diff-removed': { backgroundColor: '#3d1518', color: '#ffa198' },
  '&dark .cm-line.diff-meta': { backgroundColor: '#0d274d', color: '#79c0ff' },
})

onMounted(() => {
  view = new EditorView({
    parent: host.value,
    state: EditorState.create({
      doc: props.modelValue,
      extensions: baseExtensions(),
    }),
  })
})

onBeforeUnmount(() => {
  view?.destroy()
})

watch(
  () => props.modelValue,
  (value) => {
    if (!view || value === view.state.doc.toString()) return
    view.dispatch({
      changes: { from: 0, to: view.state.doc.length, insert: value },
    })
  },
)

watch(
  () => props.extensions,
  (extensions) => {
    if (!view) return
    view.dispatch({
      effects: languageCompartment.reconfigure(extensions),
    })
  },
  { deep: true },
)

watch(
  () => props.readonly,
  (readonly) => {
    if (!view) return
    view.dispatch({
      effects: editableCompartment.reconfigure(EditorView.editable.of(!readonly)),
    })
  },
)

watch(
  () => props.variant,
  (variant) => {
    if (!view) return
    view.dispatch({
      effects: variantCompartment.reconfigure(variantExtensions(variant)),
    })
  },
)

function baseExtensions() {
  return [
    history(),
    keymap.of([...defaultKeymap, ...historyKeymap]),
    placeholder(props.placeholder),
    syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
    syntaxHighlighting(schemaforgeHighlight),
    languageCompartment.of(props.extensions),
    editableCompartment.of(EditorView.editable.of(!props.readonly)),
    variantCompartment.of(variantExtensions(props.variant)),
    EditorView.lineWrapping,
    EditorView.updateListener.of((update) => {
      if (update.docChanged) emit('update:modelValue', update.state.doc.toString())
    }),
    EditorView.domEventHandlers({
      blur: () => {
        emit('blur')
      },
      keydown: (event) => {
        emit('keydown', event)
      },
    }),
  ]
}

function variantExtensions(variant) {
  if (variant !== 'diff') return []
  return [
    diffLineTheme,
    EditorView.decorations.compute(['doc'], (state) => {
      const decorations = []
      for (let lineNumber = 1; lineNumber <= state.doc.lines; lineNumber += 1) {
        const line = state.doc.line(lineNumber)
        const text = line.text
        if (text.startsWith('+++') || text.startsWith('---') || text.startsWith('@@')) {
          decorations.push(Decoration.line({ class: 'diff-meta' }).range(line.from))
        } else if (text.startsWith('+')) {
          decorations.push(Decoration.line({ class: 'diff-added' }).range(line.from))
        } else if (text.startsWith('-')) {
          decorations.push(Decoration.line({ class: 'diff-removed' }).range(line.from))
        }
      }
      return Decoration.set(decorations)
    }),
  ]
}
</script>
