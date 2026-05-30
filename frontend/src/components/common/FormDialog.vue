<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    :width="config.width || '520px'"
    :close-on-click-modal="false"
    destroy-on-close
    @closed="onDialogClosed"
  >
    <div v-loading="detailLoading" class="dialog-body">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        :label-width="config.labelWidth || '80px'"
      >
        <template v-for="field in config.fields" :key="field.prop">
          <el-form-item
            v-if="isFieldVisible(field)"
            :label="field.label"
            :prop="field.prop"
          >
            <el-input
              v-if="field.type === 'input'"
              v-model="formData[field.prop]"
              :disabled="isFieldDisabled(field)"
              :placeholder="field.placeholder"
              v-bind="field.props || {}"
            />
            <el-input
              v-else-if="field.type === 'password'"
              v-model="formData[field.prop]"
              type="password"
              show-password
              :disabled="isFieldDisabled(field)"
              :placeholder="field.placeholder"
              v-bind="field.props || {}"
            />
            <el-input
              v-else-if="field.type === 'textarea'"
              v-model="formData[field.prop]"
              type="textarea"
              :rows="3"
              :placeholder="field.placeholder"
              v-bind="field.props || {}"
            />
            <el-select
              v-else-if="field.type === 'select'"
              v-model="formData[field.prop]"
              :placeholder="field.placeholder || '请选择'"
              style="width: 100%"
              v-bind="field.props || {}"
            >
              <el-option
                v-for="opt in field.options"
                :key="opt.value"
                :label="opt.label"
                :value="opt.value"
              />
            </el-select>
            <el-input-number
              v-else-if="field.type === 'number'"
              v-model="formData[field.prop]"
              :min="0"
              :max="999"
              v-bind="field.props || {}"
            />
            <el-radio-group
              v-else-if="field.type === 'radio'"
              v-model="formData[field.prop]"
              v-bind="field.props || {}"
            >
              <el-radio
                v-for="opt in field.options"
                :key="opt.value"
                :value="opt.value"
              >
                {{ opt.label }}
              </el-radio>
            </el-radio-group>
            <el-switch
              v-else-if="field.type === 'switch'"
              v-model="formData[field.prop]"
              v-bind="field.props || {}"
            />
            <slot
              v-else-if="field.type === 'custom'"
              :name="`field-${field.prop}`"
              :form="formData"
              :mode="mode"
            />
          </el-form-item>
        </template>
      </el-form>
    </div>
    <template #footer>
      <el-button @click="handleCancel">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        {{ mode === 'create' ? '创建' : '保存' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" generic="T extends Record<string, any>">
import { ref, reactive, computed, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormItemRule } from 'element-plus'

interface FormField {
  prop: string
  label: string
  type: 'input' | 'password' | 'textarea' | 'select' | 'number' | 'radio' | 'switch' | 'custom'
  rules?: FormItemRule[]
  props?: Record<string, any>
  options?: { label: string; value: any }[]
  disabledOnEdit?: boolean
  visible?: (form: Record<string, any>) => boolean
  placeholder?: string
  span?: number
}

interface FormDialogConfig<T = Record<string, any>> {
  title: string | ((mode: 'create' | 'edit') => string)
  width?: string
  labelWidth?: string
  fields: FormField[]
  fetchDetail?: (id: number) => Promise<T>
  onSubmit: (payload: any, id?: number) => Promise<void>
  onSuccess?: () => void
  onClosed?: () => void
}

const props = defineProps<{
  config: FormDialogConfig<T>
}>()

defineOptions({ name: 'FormDialog' })

const visible = ref(false)
const mode = ref<'create' | 'edit'>('create')
const editingId = ref<number | undefined>()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const detailLoading = ref(false)

const dialogTitle = computed(() => {
  const t = props.config.title
  return typeof t === 'function' ? t(mode.value) : t
})

const formData = reactive<Record<string, any>>({})
const formRules = reactive<Record<string, FormItemRule[]>>({})

props.config.fields.forEach((f) => {
  formData[f.prop] = undefined
  if (f.rules) {
    formRules[f.prop] = f.rules
  }
})

function isFieldVisible(field: FormField): boolean {
  if (field.visible) {
    return field.visible(formData)
  }
  return true
}

function isFieldDisabled(field: FormField): boolean {
  return !!(field.disabledOnEdit && mode.value === 'edit')
}

async function open(m: 'create' | 'edit', id?: number) {
  mode.value = m
  editingId.value = id
  visible.value = true

  props.config.fields.forEach((f) => {
    formData[f.prop] = undefined
  })

  if (m === 'edit' && id && props.config.fetchDetail) {
    detailLoading.value = true
    try {
      const detail = await props.config.fetchDetail(id)
      await nextTick()
      Object.keys(detail).forEach((key) => {
        if (key in formData) {
          formData[key] = detail[key]
        }
      })
    } catch {
      ElMessage.error('加载数据失败')
      visible.value = false
    } finally {
      detailLoading.value = false
    }
  }

  await nextTick()
  formRef.value?.clearValidate()
}

async function handleSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    const payload = { ...formData }
    Object.keys(payload).forEach((k) => {
      if (payload[k] === undefined || payload[k] === '') {
        delete payload[k]
      }
    })
    await props.config.onSubmit(payload, editingId.value)
    ElMessage.success(mode.value === 'create' ? '创建成功' : '保存成功')
    visible.value = false
    props.config.onSuccess?.()
  } catch {
    // handled by axios interceptor
  } finally {
    submitting.value = false
  }
}

function handleCancel() {
  visible.value = false
}

function onDialogClosed() {
  props.config.onClosed?.()
}

function close() {
  visible.value = false
}

defineExpose({ open, close })
</script>

<style scoped lang="scss">
.dialog-body {
  min-height: 60px;
  max-height: 60vh;
  overflow-y: auto;
  padding: 4px 0;
}
</style>
