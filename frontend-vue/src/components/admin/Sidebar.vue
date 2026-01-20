<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import BaseModal from '../ui/BaseModal.vue'; // Path adjusted to match your structure

const route = useRoute();
const router = useRouter();
const showLogoutModal = ref(false);
const userName = ref('Administrator');
const userRole = ref('admin');

// Props for sidebar open state
defineProps({
    isOpen: {
        type: Boolean,
        default: false
    }
});

// Emit for closing sidebar
defineEmits(['close-sidebar']);

// Check Active Menu
const isActive = (path) => route.path === path;

// Load user info from storage
onMounted(() => {
    const storedName = localStorage.getItem('user_name');
    const storedRole = localStorage.getItem('role');
    if (storedName) userName.value = storedName;
    if (storedRole) userRole.value = storedRole;
});

// Trigger Logout Modal
const confirmLogout = () => {
    showLogoutModal.value = true;
};

// Execute Logout
const handleLogout = () => {
    localStorage.clear();
    router.push('/login');
};
</script>

<template>
  <div 
    v-if="isOpen" 
    @click="$emit('close-sidebar')"
    class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-40 lg:hidden transition-opacity"
  ></div>
  <aside 
    :class="[
      'w-64 bg-white text-slate-600 flex flex-col h-screen border-r border-slate-200 font-sans shadow-[4px_0_24px_rgba(0,0,0,0.02)] transition-transform duration-300 z-50',
      'fixed top-0 left-0',
      'lg:sticky',
      isOpen ? 'translate-x-0' : '-translate-x-full',
      'lg:translate-x-0'       
    ]"
  >
    <div class="h-20 flex items-center px-6 gap-4 border-b border-slate-100 shrink-0 bg-white">
      <img src="@/assets/image/logo.png" alt="ITIDoorz Logo" class="w-10 h-10 object-contain shrink-0 drop-shadow-sm" />
      <div class="flex flex-col justify-center">
        <span class="text-xl font-black text-slate-800 tracking-tighter uppercase leading-none">ITIDoorz</span>
        <span class="text-[10px] font-bold text-indigo-500 tracking-[0.2em] uppercase block mt-1 leading-none">
          Admin Panel
        </span>
      </div>
    </div>
    <nav class="flex-1 mt-6 px-4 space-y-1 overflow-y-auto custom-scrollbar">
      <div class="px-3 mb-2 mt-2">
        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Main Menu</span>
      </div>
      <router-link to="/admin/dashboard"
        @click="$emit('close-sidebar')"
        :class="['flex items-center gap-3 px-3 py-3 rounded-xl transition-all font-bold text-sm',
        isActive('/admin/dashboard') ? 'bg-indigo-50 text-indigo-700' : 'hover:bg-slate-50 hover:text-indigo-600']">
        <i :class="['pi', isActive('/admin/dashboard') ? 'pi-chart-bar text-indigo-600' : 'pi-chart-bar text-slate-400']"></i>
        <span>Dashboard</span>
      </router-link>
      <div class="px-3 mb-2 mt-6">
        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">System Management</span>
      </div>
      <router-link to="/admin/users" 
        @click="$emit('close-sidebar')"
        :class="['flex items-center gap-3 px-3 py-3 rounded-xl transition-all font-bold text-sm', 
        isActive('/admin/users') ? 'bg-indigo-50 text-indigo-700' : 'hover:bg-slate-50 hover:text-indigo-600']">
        <i :class="['pi', isActive('/admin/users') ? 'pi-users text-indigo-600' : 'pi-users text-slate-400']"></i>
        <span>User Management</span>
      </router-link>
      <router-link to="/admin/rooms" 
        @click="$emit('close-sidebar')"
        :class="['flex items-center gap-3 px-3 py-3 rounded-xl transition-all font-bold text-sm', 
        isActive('/admin/rooms') ? 'bg-indigo-50 text-indigo-700' : 'hover:bg-slate-50 hover:text-indigo-600']">
        <i :class="['pi', isActive('/admin/rooms') ? 'pi-building text-indigo-600' : 'pi-building text-slate-400']"></i>
        <span>Room Inventory</span>
      </router-link>
      <div class="px-3 mb-2 mt-6">
        <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Operations</span>
      </div>
      <router-link to="/admin/bookings" 
        @click="$emit('close-sidebar')"
        :class="['flex items-center gap-3 px-3 py-3 rounded-xl transition-all font-bold text-sm', 
        isActive('/admin/bookings') ? 'bg-indigo-50 text-indigo-700' : 'hover:bg-slate-50 hover:text-indigo-600']">
        <i :class="['pi', isActive('/admin/bookings') ? 'pi-calendar-plus text-indigo-600' : 'pi-calendar-plus text-slate-400']"></i>
        <span>Booking Requests</span>
      </router-link>
      <router-link to="/admin/reports" 
        @click="$emit('close-sidebar')"
        :class="['flex items-center gap-3 px-3 py-3 rounded-xl transition-all font-bold text-sm', 
        isActive('/admin/reports') ? 'bg-indigo-50 text-indigo-700' : 'hover:bg-slate-50 hover:text-indigo-600']">
        <i :class="['pi', isActive('/admin/reports') ? 'pi-file-excel text-indigo-600' : 'pi-file-excel text-slate-400']"></i>
        <span>Report Generator</span>
      </router-link>
    </nav>
    <div class="p-4 border-t border-slate-100 bg-slate-50/50 shrink-0">
      <div class="flex items-center gap-3 mb-4 px-2">
        <div class="w-9 h-9 rounded-lg bg-white flex items-center justify-center border border-slate-200 shadow-sm shrink-0">
          <i class="pi pi-user text-slate-500 text-sm"></i>
        </div>
        <div class="overflow-hidden">
            <p class="text-xs font-bold text-slate-700 truncate capitalize">{{ userName }}</p>
            <p class="text-[9px] text-slate-400 truncate font-mono uppercase">{{ userRole }}</p>
        </div>
      </div>
      <button 
        @click="confirmLogout"
        class="flex items-center justify-center gap-2 p-3 rounded-xl transition-all font-bold text-xs w-full text-red-500 bg-white hover:bg-red-50 hover:text-red-600 border border-slate-200 hover:border-red-100 shadow-sm"
      >
        <i class="pi pi-sign-out"></i>
        <span>LOGOUT</span>
      </button>
    </div>
    <BaseModal 
        :show="showLogoutModal" 
        title="Confirm Logout" 
        icon="pi-exclamation-triangle" 
        @close="showLogoutModal = false"
    >
        <div class="text-center py-4 space-y-2">
            <p class="text-slate-600 font-medium text-sm">Are you sure you want to log out?</p>
            <p class="text-xs text-slate-400">Your session will be ended.</p>
        </div>
        <template #footer>
            <div class="flex gap-2 w-full">
                <Button label="Cancel" severity="secondary" @click="showLogoutModal = false" class="flex-1 font-bold rounded-xl" text />
                <Button label="Yes, Logout" severity="danger" @click="handleLogout" class="flex-1 font-bold shadow-lg shadow-red-100 rounded-xl" />
            </div>
        </template>
    </BaseModal>
  </aside>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 20px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background-color: #94a3b8;
}
</style>