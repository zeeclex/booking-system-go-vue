<script setup>
import { ref, onMounted, computed } from 'vue';
import { UserService } from '@/service/UserService';
import AdminSidebar from '../../components/admin/Sidebar.vue';
import BaseModal from '../../components/ui/BaseModal.vue';
import StatCard from '../../components/ui/StatCard.vue';

// --- STATE MANAGEMENT ---
const users = ref([]);
const loading = ref(true);
const successMessage = ref('');

// Modal State
const showFormModal = ref(false);
const showDeleteModal = ref(false);
const showSuccessModal = ref(false);

// Form & Selection
const isEditing = ref(false);
const selectedUser = ref(null);
const form = ref({
    id: null,
    name: '',
    email: '',
    role: 'user',
    password: ''
});

// Mobile Sidebar State
const isSidebarOpen = ref(false);

// --- COMPUTED STATS ---
const userStats = computed(() => {
    return {
        total: users.value.length,
        admins: users.value.filter(u => u.role === 'admin').length,
        general: users.value.filter(u => u.role === 'user').length
    };
});

// --- HELPER: Role Styling ---
const getRoleSeverity = (role) => {
    switch (role) {
        case 'admin': return 'primary';
        case 'user': return 'info';
        default: return 'contrast';
    }
};

// --- DATA LOGIC ---
const loadUsers = async () => {
    loading.value = true;
    try {
        const res = await UserService.getUsers();
        users.value = res.data;
    } catch (err) {
        console.error("Failed to fetch users:", err);
    } finally {
        loading.value = false;
    }
};

// --- ACTION HANDLERS ---
const openAddModal = () => {
    isEditing.value = false;
    form.value = { name: '', email: '', role: 'user', password: '' };
    showFormModal.value = true;
};

const openEditModal = (user) => {
    isEditing.value = true;
    form.value = { 
        id: user.id,
        name: user.name,
        email: user.email,
        role: user.role,
        password: ''
    };
    showFormModal.value = true;
};

const openDeleteConfirm = (user) => {
    const currentUserId = localStorage.getItem('user_id');
    if (user.id == currentUserId) {
        alert("Action Denied: You cannot delete your own account.");
        return;
    }
    selectedUser.value = user;
    showDeleteModal.value = true;
};

const triggerSuccess = (msg) => {
    successMessage.value = msg;
    showSuccessModal.value = true;
    setTimeout(() => { showSuccessModal.value = false; }, 2000); 
};

// --- API CALLS ---
const handleSave = async () => {
    try {
        if (!form.value.name || !form.value.email) {
            alert("Name and Email are required!");
            return;
        }
        if (!isEditing.value && !form.value.password) {
            alert("Password is required for new users!");
            return;
        }
        if (isEditing.value) {
            await UserService.updateUser(form.value.id, form.value);
            triggerSuccess('User profile updated successfully!');
        } else {
            await UserService.createUser(form.value);
            triggerSuccess('New user account created!');
        }
        showFormModal.value = false;
        setTimeout(() => loadUsers(), 300);

    } catch (err) {
        const errorMsg = err.response?.data?.error || "System error occurred";
        alert("Failed: " + errorMsg);
    }
};

const handleDelete = async () => {
    try {
        await UserService.deleteUser(selectedUser.value.id);
        showDeleteModal.value = false;
        triggerSuccess('User account permanently deleted.');
        loadUsers();
    } catch (err) {
        alert("Error: " + (err.response?.data?.error || "Failed to delete user"));
    }
};

onMounted(loadUsers);
</script>

<template>
    <div class="flex min-h-screen bg-slate-50 font-sans">
        <AdminSidebar :isOpen="isSidebarOpen" @close-sidebar="isSidebarOpen = false" />
        <main class="flex-1 w-full max-w-7xl mx-auto p-4 md:p-8 space-y-6 md:space-y-8 overflow-y-auto transition-all">
            <div class="flex flex-col md:flex-row justify-between items-start md:items-end border-b border-slate-200 pb-5 gap-4">
                <div class="flex items-center gap-4">
                    <button @click="isSidebarOpen = true" class="lg:hidden p-2 -ml-2 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">
                        <i class="pi pi-bars text-xl"></i>
                    </button>
                    <div>
                        <h2 class="text-xl md:text-2xl font-black text-slate-800 uppercase tracking-tighter italic">User Management</h2>
                        <p class="text-xs md:text-sm text-slate-500 font-medium italic mt-1">Manage System Access</p>
                    </div>
                </div>
                <Button 
                    label="Add New User" 
                    icon="pi pi-user-plus" 
                    class="w-full md:w-auto bg-indigo-600 hover:bg-indigo-700 text-white font-bold shadow-lg shadow-indigo-100 px-5 py-2.5 rounded-xl transition-all hover:scale-105 active:scale-95 text-sm" 
                    @click="openAddModal" 
                />
            </div>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-5">
                <StatCard 
                    title="Total Accounts" 
                    :value="userStats.total" 
                    icon="pi-users" 
                    variant="indigo" 
                />
                <StatCard 
                    title="Administrators" 
                    :value="userStats.admins" 
                    icon="pi-shield" 
                    variant="slate" 
                />
                <StatCard 
                    title="General Users" 
                    :value="userStats.general" 
                    icon="pi-id-card" 
                    variant="blue" 
                />
            </div>
            <div class="bg-white rounded-2xl border border-slate-200 shadow-lg shadow-slate-100/50 overflow-hidden">
                
                <div class="p-4 md:p-6 border-b border-slate-100 bg-white flex flex-col sm:flex-row justify-between items-center gap-4">
                    <div class="w-full sm:w-auto">
                        <h3 class="font-black text-slate-800 text-sm md:text-base uppercase tracking-wider flex items-center gap-3">
                            <div class="p-2 bg-indigo-50 rounded-lg text-indigo-600">
                                <i class="pi pi-table"></i>
                            </div>
                            Registered Users List
                        </h3>
                    </div>
                    <Button 
                        icon="pi pi-refresh" 
                        text 
                        rounded 
                        @click="loadUsers" 
                        :loading="loading" 
                        class="text-slate-400 hover:text-indigo-600 hover:bg-indigo-50"
                    />
                </div>
                <div class="p-0 md:p-2">
                    <DataTable 
                        :value="users" 
                        :loading="loading" 
                        paginator :rows="10" 
                        responsiveLayout="scroll"
                        class="p-datatable-sm"
                        :pt="{
                            header: { class: 'bg-slate-50 text-slate-800 font-bold uppercase text-xs tracking-wider border-b border-slate-100 py-3' },
                            bodyRow: { class: 'hover:bg-slate-50 transition-colors border-b border-slate-50 last:border-0' }
                        }"
                    >
                        <template #empty>
                            <div class="text-center py-12 text-slate-400">
                                <i class="pi pi-user-minus text-5xl mb-3 block opacity-20 mx-auto"></i>
                                <p class="font-bold text-xs uppercase tracking-widest">No users found.</p>
                            </div>
                        </template>
                        <Column field="name" header="Full Name" sortable headerClass="pl-4 md:pl-6" bodyClass="pl-4 md:pl-6 text-left">
                            <template #body="{ data }">
                                <div class="flex items-center gap-3">
                                    <div class="w-8 h-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold text-xs shrink-0">
                                        {{ data.name.charAt(0).toUpperCase() }}
                                    </div>
                                    <span class="font-bold text-slate-700 text-sm whitespace-nowrap">{{ data.name }}</span>
                                </div>
                            </template>
                        </Column>
                        <Column field="email" header="Email Address" headerClass="text-left" bodyClass="text-left">
                            <template #body="{ data }">
                                <span class="text-slate-600 text-sm whitespace-nowrap">{{ data.email }}</span>
                            </template>
                        </Column>
                        <Column field="role" header="Role" headerClass="text-center" bodyClass="text-center">
                            <template #body="{ data }">
                                <Tag :value="data.role.toUpperCase()" :severity="getRoleSeverity(data.role)" class="px-2 py-1 text-[10px] font-black tracking-widest uppercase rounded shadow-sm" />
                            </template>
                        </Column>
                        <Column header="Actions" headerClass="text-center pr-4 md:pr-6" bodyClass="text-center pr-4 md:pr-6">
                            <template #body="{ data }">
                                <div class="flex justify-center gap-2">
                                    <Button icon="pi pi-pencil" severity="info" text rounded size="small" @click="openEditModal(data)" />
                                    <Button icon="pi pi-trash" severity="danger" text rounded size="small" @click="openDeleteConfirm(data)" />
                                </div>
                            </template>
                        </Column>
                    </DataTable>
                </div>
            </div>
            <BaseModal :show="showFormModal" :title="isEditing ? 'Edit User Profile' : 'Create New User'" icon="pi-user-edit" @close="showFormModal = false">
                <div class="space-y-5 py-2">
                    <div class="space-y-1.5">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Full Name</label>
                        <InputText v-model="form.name" placeholder="Ex: Dr. John Doe" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-slate-700 text-sm" />
                    </div>
                    <div class="space-y-1.5">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Email Address</label>
                        <InputText v-model="form.email" type="email" placeholder="Ex: john@iti.ac.id" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-slate-700 text-sm" />
                    </div>
                    <div class="space-y-1.5">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Role</label>
                        <div class="relative">
                            <select v-model="form.role" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-slate-700 text-sm appearance-none cursor-pointer">
                                <option value="admin">Administrator</option>
                                <option value="user">General User</option>
                            </select>
                            <i class="pi pi-chevron-down absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 text-xs pointer-events-none"></i>
                        </div>
                    </div>
                    <div class="space-y-1.5">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Password</label>
                        <InputText v-model="form.password" type="password" placeholder="••••••••" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-slate-700 text-sm" />
                        <p v-if="isEditing" class="text-[10px] text-amber-600 font-bold mt-1 ml-1 flex items-center gap-1">
                            <i class="pi pi-info-circle"></i> Leave blank to keep current password
                        </p>
                    </div>
                </div>
                <template #footer>
                    <div class="flex flex-col-reverse sm:flex-row gap-2 w-full pt-2">
                        <Button label="Cancel" severity="secondary" @click="showFormModal = false" class="flex-1 font-bold rounded-xl" text />
                        <Button :label="isEditing ? 'Save Changes' : 'Create Account'" icon="pi pi-check" @click="handleSave" class="flex-1 font-bold shadow-lg shadow-indigo-100 rounded-xl" />
                    </div>
                </template>
            </BaseModal>
            <BaseModal :show="showDeleteModal" title="Confirm Deletion" icon="pi-exclamation-triangle" @close="showDeleteModal = false">
                <div class="text-center py-4 space-y-3">
                    <p class="text-slate-600 font-medium text-sm">Are you sure you want to delete user <br><span class="font-black text-slate-800 text-lg">{{ selectedUser?.name }}</span>?</p>
                    <p class="text-[10px] text-red-500 font-bold bg-red-50 py-1.5 px-3 rounded-lg inline-block uppercase tracking-wider">This action is permanent</p>
                </div>
                <template #footer>
                    <div class="flex flex-col-reverse sm:flex-row gap-2 w-full">
                        <Button label="Cancel" severity="secondary" @click="showDeleteModal = false" class="flex-1 font-bold rounded-xl" text />
                        <Button label="Yes, Delete" severity="danger" @click="handleDelete" class="flex-1 font-bold shadow-lg shadow-red-100 rounded-xl" />
                    </div>
                </template>
            </BaseModal>
            <BaseModal :show="showSuccessModal" title="Success" icon="pi-check-circle" @close="showSuccessModal = false">
                <div class="flex flex-col items-center justify-center py-6 space-y-4">
                    <div class="w-16 h-16 bg-green-50 text-green-500 rounded-full flex items-center justify-center animate-bounce">
                        <i class="pi pi-check text-3xl"></i>
                    </div>
                    <p class="font-bold text-slate-800 text-base text-center">{{ successMessage }}</p>
                </div>
                <template #footer>
                    <Button label="OK, Continue" class="w-full font-bold rounded-xl shadow-lg" @click="showSuccessModal = false" />
                </template>
            </BaseModal>
        </main>
    </div>
</template>