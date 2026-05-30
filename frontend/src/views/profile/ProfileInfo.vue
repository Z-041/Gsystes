<template>
  <div class="profile-page">
    <div class="page-header">
      <h2>个人中心</h2>
    </div>

    <el-row :gutter="20">
      <el-col :span="14">
        <el-card shadow="never" class="profile-card">
          <template #header>
            <span class="card-title">个人信息</span>
          </template>
          <el-form ref="formRef" :model="form" :rules="rules" label-width="80px" style="max-width: 400px">
            <el-form-item label="昵称" prop="nickname">
              <el-input v-model="form.nickname" />
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="form.email" />
            </el-form-item>
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="form.phone" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="handleSave">保存修改</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :span="10">
        <el-card shadow="never" class="profile-card">
          <template #header>
            <span class="card-title">头像</span>
          </template>
          <div class="avatar-section">
            <el-avatar :size="100" style="background: linear-gradient(135deg, var(--color-primary), var(--color-primary-light)); font-size: 40px">
              {{ avatarLetter }}
            </el-avatar>
            <el-upload
              :show-file-list="false"
              :before-upload="beforeAvatarUpload"
              :http-request="handleUpload"
              style="margin-top: 12px"
            >
              <el-button size="small">上传头像</el-button>
            </el-upload>
          </div>
        </el-card>

        <el-card shadow="never" class="profile-card" style="margin-top: 16px">
          <template #header>
            <span class="card-title">修改密码</span>
          </template>
          <el-button @click="passwordRef?.open()">修改密码</el-button>
        </el-card>
      </el-col>
    </el-row>

    <PasswordDialog ref="passwordRef" />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { getUserProfileApi, updateProfileApi, uploadAvatarApi } from '@/api/user'
import PasswordDialog from '@/components/common/PasswordDialog.vue'

defineOptions({ name: 'ProfileInfo' })

const authStore = useAuthStore()

const formRef = ref<FormInstance>()
const passwordRef = ref<InstanceType<typeof PasswordDialog>>()
const saving = ref(false)

const form = reactive({
  nickname: '',
  email: '',
  phone: '',
})

const rules: FormRules = {
  email: [{ type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }],
}

const avatarLetter = computed(() => {
  const u = authStore.userInfo
  return (u?.nickname?.charAt(0) || u?.username?.charAt(0) || 'U')
})

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    await updateProfileApi(form)
    ElMessage.success('保存成功')
    authStore.fetchUserInfo()
  } finally {
    saving.value = false
  }
}

function beforeAvatarUpload(file: File) {
  const isImage = file.type.startsWith('image/')
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB')
    return false
  }
  return true
}

async function handleUpload(options: any) {
  try {
    const res = await uploadAvatarApi(options.file)
    ElMessage.success('头像上传成功')
    authStore.fetchUserInfo()
  } catch {
    // handled by interceptor
  }
}

onMounted(async () => {
  try {
    const res = await getUserProfileApi()
    const p = res.data
    form.nickname = p.nickname || ''
    form.email = p.email || ''
    form.phone = p.phone || ''
  } catch { /* */ }
})
</script>

<style scoped lang="scss">
.page-header {
  margin-bottom: 16px;
  h2 { font-size: var(--font-size-xl); font-weight: 700; color: var(--color-text); }
}

.profile-card {
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);
  margin-bottom: 16px;
}

.card-title {
  font-weight: 600;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 0;
}
</style>
