<template>
  <div ref="host" class="code-editor"></div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { EditorState, Compartment } from '@codemirror/state'
import { EditorView, keymap, placeholder } from '@codemirror/view'
import { defaultKeymap, history, historyKeymap } from '@codemirror/commands'
import { HighlightStyle, defaultHighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { tags as t } from '@lezer/highlight'

const props = defineProps({
  modelValue: { type: String, default: '' },
  extensions: { type: Array, default: () => [] },
  readonly: { type: Boolean, default: false },
  placeholder: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue', 'blur', 'keydown'])

const host = ref(null)
let view
const languageCompartment = new Compartment()
const editableCompartment = new Compartment()
const schemaforgeHighlight = HighlightStyle.define([
  { tag: [t.keyword, t.operatorKeyword, t.modifier], color: '#cf222e', fontWeight: '600' },
  { tag: [t.string, t.special(t.string)], color: '#0a3069' },
  { tag: [t.number, t.bool, t.null], color: '#0550ae' },
  { tag: [t.comment, t.lineComment, t.blockComment], color: '#6e7781', fontStyle: 'italic' },
  { tag: [t.definition(t.variableName), t.className, t.typeName], color: '#8250df' },
  { tag: [t.propertyName, t.attributeName], color: '#953800' },
  { tag: [t.function(t.variableName), t.function(t.propertyName)], color: '#8250df' },
])

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

function baseExtensions() {
  return [
    history(),
    keymap.of([...defaultKeymap, ...historyKeymap]),
    placeholder(props.placeholder),
    syntaxHighlighting(defaultHighlightStyle, { fallback: true }),
    syntaxHighlighting(schemaforgeHighlight),
    languageCompartment.of(props.extensions),
    editableCompartment.of(EditorView.editable.of(!props.readonly)),
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
</script>
