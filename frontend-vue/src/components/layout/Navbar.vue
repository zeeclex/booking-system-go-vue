<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import BaseModal from '../ui/BaseModal.vue';

const router = useRouter();

// State User
const userName = ref('Guest');
const userMajor = ref('Informatics Engineering ITI'); // Translated

// State Modal
const showLogoutModal = ref(false);

const loadUserData = () => {
    const savedName = localStorage.getItem('user_name');
    if (savedName) {
        userName.value = savedName;
    }
};

const handleLogout = () => {
    localStorage.clear();
    router.push('/login');
};

onMounted(loadUserData);
</script>

<template>
    <nav class="sticky top-0 z-40 w-full bg-white/80 backdrop-blur-md border-b border-slate-200 transition-all duration-300">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 md:px-8 h-16 md:h-20 flex justify-between items-center transition-all duration-300">
            
            <div class="flex items-center gap-3 md:gap-4">
                <img src="@/assets/image/logo.png" alt="ITIDoorz Logo" class="w-9 h-9 md:w-11 md:h-11 object-contain drop-shadow-sm transition-all" />
                
                <div>
                    <span class="block text-lg md:text-xl font-black text-slate-800 uppercase tracking-tighter leading-none transition-all">ITIDoorz</span>
                    <span class="block text-[9px] md:text-[10px] font-bold text-indigo-500 uppercase tracking-[0.2em] leading-none mt-0.5 md:mt-1 transition-all">User Panel</span>
                </div>
            </div>

            <div class="flex items-center gap-3 md:gap-5">
                
                <div class="hidden md:block text-right">
                    <p class="text-sm font-black text-slate-700 uppercase tracking-tight truncate max-w-37.5">{{ userName }}</p>
                    <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ userMajor }}</p>
                </div>

                <div class="hidden md:block h-8 w-px bg-slate-200"></div>

                <div class="flex items-center gap-2 md:gap-3">
                    <div class="w-9 h-9 md:w-10 md:h-10 rounded-xl bg-slate-100 border border-slate-200 flex items-center justify-center text-slate-500 transition-all">
                        <i class="pi pi-user text-sm md:text-base"></i>
                    </div>

                    <button 
                        @click="showLogoutModal = true"
                        class="w-9 h-9 md:w-10 md:h-10 rounded-xl bg-red-50 text-red-500 hover:bg-red-100 hover:text-red-600 transition-colors flex items-center justify-center border border-red-100 active:scale-95"
                        title="Logout"
                    >
                        <i class="pi pi-sign-out text-sm md:text-base font-bold"></i>
                    </button>
                </div>
            </div>

        </div>

        <BaseModal 
            :show="showLogoutModal" 
            title="Confirm Logout" 
            icon="pi-exclamation-triangle" 
            @close="showLogoutModal = false"
        >
            <div class="text-center py-4 space-y-2">
                <p class="text-slate-600 font-medium text-sm md:text-base">Are you sure you want to log out?</p>
                <p class="text-xs text-slate-400">Your session will end.</p>
            </div>
            <template #footer>
                <div class="flex flex-col-reverse md:flex-row gap-2 md:gap-3 w-full">
                    <Button label="Cancel" severity="secondary" @click="showLogoutModal = false" class="w-full md:flex-1 font-bold rounded-xl" text />
                    <Button label="Yes, Logout" severity="danger" @click="handleLogout" class="w-full md:flex-1 font-bold shadow-lg shadow-red-100 rounded-xl" />
                </div>
            </template>
        </BaseModal>

    </nav>
</template>