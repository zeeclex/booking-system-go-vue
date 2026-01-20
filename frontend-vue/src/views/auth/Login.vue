<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { AuthService } from '@/service/AuthService';
import BaseModal from '../../components/ui/BaseModal.vue';

// --- STATE ---
const router = useRouter();
const email = ref('');
const password = ref('');
const loading = ref(false);
const showErrorModal = ref(false);
const errorMessage = ref('');

// --- HANDLE LOGIN ---
const handleLogin = async () => {
  if (!email.value || !password.value) {
    errorMessage.value = "Email and Password are required!";
    showErrorModal.value = true;
    return;
  }
  loading.value = true;
  try {
    // 1. Request to Go Backend
    const res = await AuthService.login({
      email: email.value,
      password: password.value
    });

    // 2. Save session to LocalStorage
    const { token, role, id, name } = res.data;
    localStorage.setItem('token', token);
    localStorage.setItem('role', role); 
    localStorage.setItem('user_id', id);
    localStorage.setItem('user_name', name);

    // 3. Redirect based on role
    if (role === 'admin') {
      await router.push('/admin/dashboard'); 
    } else {
      await router.push('/user/dashboard'); 
    }
  } catch (err) {
    console.error("Login Error:", err);
    if (err.response && err.response.data && err.response.data.message) {
        errorMessage.value = err.response.data.message;
    } else {
        errorMessage.value = "Failed to connect to server. Please check your internet connection.";
    }
    showErrorModal.value = true;
    localStorage.clear();
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-slate-50 font-sans p-4 md:p-6 relative overflow-hidden">
    <div class="absolute top-[-10%] left-[-10%] w-64 h-64 md:w-96 md:h-96 bg-indigo-200 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-64 h-64 md:w-96 md:h-96 bg-blue-200 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000"></div>
    <div class="w-full max-w-sm bg-white rounded-2xl border border-slate-200 shadow-xl shadow-slate-200/50 overflow-hidden relative z-10 transition-all">
      <div class="text-center pt-8 pb-6 px-6 md:pt-12 md:pb-8 md:px-8 bg-white">
        <div class="inline-flex items-center justify-center mb-4 md:mb-6 transform hover:scale-105 transition-transform duration-300">
            <img src="@/assets/image/logo.png" alt="ITIDoorz Logo" class="w-20 h-20 md:w-24 md:h-24 object-contain drop-shadow-md" />
        </div>
        <h2 class="text-2xl md:text-3xl font-black text-slate-800 uppercase tracking-tighter">ITIDoorz</h2>
        <div class="flex items-center justify-center gap-2 mt-2 opacity-60">
          <div class="h-px w-6 md:w-8 bg-slate-400"></div>
          <p class="text-[9px] md:text-[10px] font-black uppercase tracking-[0.2em] text-slate-500">Access System</p>
          <div class="h-px w-6 md:w-8 bg-slate-400"></div>
        </div>
      </div>
      <div class="px-6 pb-8 md:px-10 md:pb-12 space-y-4 md:space-y-5">
        <div class="space-y-1.5">
          <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Identity</label>
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <i class="pi pi-envelope text-slate-400 group-focus-within:text-indigo-500 transition-colors"></i>
            </div>
            <input 
              v-model="email" 
              type="email" 
              placeholder="name@iti.ac.id"
              class="w-full py-3 md:py-3.5 pl-11 pr-4 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent font-bold text-slate-700 text-sm transition-all placeholder:text-slate-300 hover:bg-slate-100" 
            />
          </div>
        </div>
        <div class="space-y-1.5">
          <label class="block text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Credential</label>
          <div class="relative group">
            <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
              <i class="pi pi-lock text-slate-400 group-focus-within:text-indigo-500 transition-colors"></i>
            </div>
            <input 
              v-model="password" 
              type="password" 
              placeholder="••••••••"
              @keyup.enter="handleLogin"
              class="w-full py-3 md:py-3.5 pl-11 pr-4 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:bg-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent font-bold text-slate-700 text-sm transition-all placeholder:text-slate-300 hover:bg-slate-100" 
            />
          </div>
        </div>
        <Button 
          @click="handleLogin" 
          :loading="loading"
          class="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-3 md:py-3.5 rounded-xl shadow-lg shadow-indigo-200 hover:shadow-indigo-300 transition-all transform active:scale-95 flex justify-center gap-2 mt-4"
        >
          <span class="font-bold tracking-wide text-sm">{{ loading ? 'AUTHENTICATING...' : 'SIGN IN' }}</span>
          <i v-if="!loading" class="pi pi-arrow-right font-bold text-xs"></i>
        </Button>
      </div>
    </div>
    <div class="mt-8 text-center relative z-10">
      <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">
        Informatics Engineering <br> Institut Teknologi Indonesia
      </p>
      <p class="text-[10px] text-slate-300 mt-2 font-mono">© 2026 ITIDoorz v1.0</p>
    </div>
    <BaseModal :show="showErrorModal" title="Access Denied" icon="pi-lock" @close="showErrorModal = false">
      <div class="text-center space-y-4 py-2">
        <div class="w-16 h-16 bg-red-50 text-red-500 rounded-full flex items-center justify-center mx-auto border-4 border-red-50">
          <i class="pi pi-times text-2xl font-bold"></i>
        </div>
        <div>
          <h3 class="font-bold text-slate-800 text-lg">Login Failed</h3>
          <p class="text-slate-500 text-sm mt-1 px-4 leading-relaxed">{{ errorMessage }}</p>
        </div>
      </div>
      <template #footer>
        <Button label="Try Again" severity="danger" @click="showErrorModal = false" class="w-full font-bold shadow-lg shadow-red-100 py-3 rounded-xl" />
      </template>
    </BaseModal>
  </div>
</template>

<style scoped>
@keyframes blob {
  0% { transform: translate(0px, 0px) scale(1); }
  33% { transform: translate(30px, -50px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
  100% { transform: translate(0px, 0px) scale(1); }
}
.animate-blob {
  animation: blob 7s infinite;
}
.animation-delay-2000 {
  animation-delay: 2s;
}
</style>