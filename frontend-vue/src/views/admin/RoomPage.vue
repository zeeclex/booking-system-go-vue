<script setup>
import { ref, onMounted, computed } from 'vue';
import { RoomService } from '@/service/RoomService';
import AdminSidebar from '../../components/admin/Sidebar.vue';
import BaseModal from '../../components/ui/BaseModal.vue';
import StatCard from '../../components/ui/StatCard.vue';

// --- STATE MANAGEMENT ---
const rooms = ref([]);
const loading = ref(true);
const successMessage = ref('');
const showFormModal = ref(false);
const showDeleteModal = ref(false);
const showSuccessModal = ref(false);
const isEditing = ref(false);
const selectedRoom = ref(null);
const form = ref({
    id: null,
    name: '',
    type: 'lab',
    capacity: 30,
    is_active: true
});

// Mobile Sidebar Toggle State
const isSidebarOpen = ref(false);

// --- COMPUTED STATS ---
const roomStats = computed(() => {
    return {
        total: rooms.value.length,
        active: rooms.value.filter(r => r.is_active).length,
        inactive: rooms.value.filter(r => !r.is_active).length
    };
});

// --- DATA LOGIC ---
const loadRooms = async () => {
    loading.value = true;
    try {
        const res = await RoomService.getRooms();
        rooms.value = res.data;
    } catch (err) {
        console.error("Failed to fetch room data:", err);
    } finally {
        loading.value = false;
    }
};

// --- ACTIONS ---
const openAddModal = () => {
    isEditing.value = false;
    form.value = { name: '', type: 'lab', capacity: 30, is_active: true };
    showFormModal.value = true;
};

const openEditModal = (room) => {
    isEditing.value = true;
    form.value = { 
        ...room,
        is_active: !!room.is_active 
    };
    showFormModal.value = true;
};

const openDeleteConfirm = (room) => {
    selectedRoom.value = room;
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
        if (!form.value.name || !form.value.capacity) {
            alert("Room Name and Capacity are required!");
            return;
        }
        const payload = {
            name: form.value.name,
            type: form.value.type,
            capacity: parseInt(form.value.capacity),
            is_active: Boolean(form.value.is_active)
        };
        if (isEditing.value) {
            await RoomService.updateRoom(form.value.id, payload);
            triggerSuccess('Room updated successfully!');
        } else {
            await RoomService.createRoom(payload);
            triggerSuccess('New room added successfully!');
        }
        showFormModal.value = false;
        
        setTimeout(() => loadRooms(), 300);
    } catch (err) {
        const errorMsg = err.response?.data?.message || "System error occurred";
        console.error("Failed to save data:", errorMsg);
        alert("Failed: " + errorMsg);
    }
};

const handleDelete = async () => {
    try {
        await RoomService.deleteRoom(selectedRoom.value.id);
        showDeleteModal.value = false;
        triggerSuccess('Room permanently deleted.');
        loadRooms();
    } catch (err) {
        if (err.response && (err.response.status === 409 || err.response.status === 500)) {
            showDeleteModal.value = false;
            const confirmArchive = confirm(
                "DELETE FAILED: This room has associated booking history.\n\n" +
                "Deleting this data will corrupt historical reports.\n" +
                "Do you want to archive (Deactivate) this room instead?"
            );
            if (confirmArchive) {
                try {
                    await RoomService.updateRoom(selectedRoom.value.id, {
                        ...selectedRoom.value,
                        is_active: false
                    });
                    triggerSuccess('Room archived (Inactive) successfully.');
                    loadRooms();
                } catch (updateErr) {
                    alert("Failed to archive room.");
                }
            }
        } else {
            console.error("Failed to delete data:", err);
            alert("Error: " + (err.response?.data?.message || "Failed to delete data"));
        }
    }
};

onMounted(loadRooms);
</script>

<template>
    <div class="flex min-h-screen bg-slate-50 font-sans">
        
        <AdminSidebar :isOpen="isSidebarOpen" @close-sidebar="isSidebarOpen = false" />

        <main class="flex-1 w-full max-w-7xl mx-auto p-4 md:p-8 space-y-6 md:space-y-8 overflow-y-auto transition-all">
            
            <div class="flex flex-col md:flex-row justify-between items-start md:items-end border-b border-slate-200 pb-5 gap-4">
                <div class="flex items-center gap-4">
                    <button @click="isSidebarOpen = true" class="md:hidden p-2 -ml-2 text-slate-600 hover:bg-slate-100 rounded-lg">
                        <i class="pi pi-bars text-xl"></i>
                    </button>
                    <div>
                        <h2 class="text-xl md:text-2xl font-black text-slate-800 uppercase tracking-tighter italic">Room Manager</h2>
                        <p class="text-xs md:text-sm text-slate-500 font-medium italic mt-1">Manage Labs & Classrooms - ITI Informatics Engineering</p>
                    </div>
                </div>
                
                <Button 
                    label="Add New Room" 
                    icon="pi pi-plus" 
                    class="w-full md:w-auto bg-indigo-600 hover:bg-indigo-700 text-white font-bold shadow-lg shadow-indigo-100 px-5 py-2.5 rounded-xl transition-all hover:scale-105 active:scale-95 text-sm" 
                    @click="openAddModal" 
                />
            </div>

            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-5">
                <StatCard 
                    title="Total Inventory" 
                    :value="roomStats.total" 
                    icon="pi-building" 
                    variant="blue" 
                />
                <StatCard 
                    title="Active Rooms" 
                    :value="roomStats.active" 
                    icon="pi-check-circle" 
                    variant="emerald" 
                />
                <StatCard 
                    title="Inactive / Archived" 
                    :value="roomStats.inactive" 
                    icon="pi-ban" 
                    variant="slate" 
                />
            </div>

            <div class="bg-white rounded-2xl border border-slate-200 shadow-lg shadow-slate-100/50 overflow-hidden">
                
                <div class="p-4 md:p-6 border-b border-slate-100 bg-white flex flex-col sm:flex-row justify-between items-center gap-4">
                    <div class="w-full sm:w-auto">
                        <h3 class="font-black text-slate-800 text-sm md:text-base uppercase tracking-wider flex items-center gap-3">
                            <div class="p-2 bg-indigo-50 rounded-lg text-indigo-600">
                                <i class="pi pi-list"></i>
                            </div>
                            Room Inventory List
                        </h3>
                    </div>
                    
                    <Button 
                        icon="pi pi-refresh" 
                        text 
                        rounded 
                        @click="loadRooms" 
                        :loading="loading" 
                        class="text-slate-400 hover:text-indigo-600 hover:bg-indigo-50"
                    />
                </div>

                <div class="p-0 md:p-2">
                    <DataTable 
                        :value="rooms" 
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
                                <i class="pi pi-box text-5xl mb-3 block opacity-20 mx-auto"></i>
                                <p class="font-bold text-xs uppercase tracking-widest">No room data found.</p>
                            </div>
                        </template>

                        <Column field="name" header="Room Name" sortable headerClass="pl-4 md:pl-6" bodyClass="pl-4 md:pl-6 text-left">
                            <template #body="{ data }">
                                <span class="font-bold text-slate-700 text-sm whitespace-nowrap">{{ data.name }}</span>
                            </template>
                        </Column>

                        <Column field="type" header="Type" headerClass="text-center" bodyClass="text-center text-slate-600">
                            <template #body="{ data }">
                                <span class="capitalize text-xs font-bold bg-slate-100 px-2 py-1 rounded-md text-slate-500">{{ data.type }}</span>
                            </template>
                        </Column>

                        <Column field="capacity" header="Capacity" headerClass="text-center" bodyClass="text-center font-medium">
                            <template #body="{ data }"> 
                                <span class="text-sm whitespace-nowrap">{{ data.capacity }} Seats</span>
                            </template>
                        </Column>

                        <Column field="is_active" header="Status" headerClass="text-center" bodyClass="text-center">
                            <template #body="{ data }">
                                <Tag :value="data.is_active ? 'ACTIVE' : 'INACTIVE'" :severity="data.is_active ? 'success' : 'danger'" class="px-2 py-1 text-[10px] font-black tracking-widest uppercase rounded shadow-sm" />
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

            <BaseModal :show="showFormModal" :title="isEditing ? 'Edit Room' : 'Add New Room'" icon="pi-pencil" @close="showFormModal = false">
                <div class="space-y-5 py-2">
                    <div class="space-y-1.5">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Room Name</label>
                        <InputText v-model="form.name" placeholder="Ex: Network Lab" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-slate-700 text-sm" />
                    </div>

                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div class="space-y-1.5">
                            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Type</label>
                            <div class="relative">
                                <select v-model="form.type" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-slate-700 text-sm appearance-none cursor-pointer">
                                    <option value="lab">Laboratory</option>
                                    <option value="kelas">Classroom</option>
                                    <option value="aula">Hall</option>
                                </select>
                                <i class="pi pi-chevron-down absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 text-xs pointer-events-none"></i>
                            </div>
                        </div>
                        <div class="space-y-1.5">
                            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Capacity</label>
                            <InputText v-model="form.capacity" type="number" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-slate-700 text-sm" />
                        </div>
                    </div>

                    <div class="space-y-1.5">
                        <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Room Status</label>
                        
                        <div 
                            @click="form.is_active = !form.is_active" 
                            class="flex items-center justify-between bg-slate-50 p-3 rounded-xl border border-slate-200 cursor-pointer hover:bg-slate-100 transition-colors group select-none"
                        >
                            <span class="text-sm font-bold text-slate-700">
                                {{ form.is_active ? 'Active (Bookable)' : 'Inactive (Archived)' }}
                            </span>

                            <div :class="[
                                'w-12 h-6 rounded-full p-1 transition-colors duration-300 ease-in-out flex items-center',
                                form.is_active ? 'bg-indigo-500' : 'bg-slate-300'
                            ]">
                                <div :class="[
                                    'w-4 h-4 rounded-full bg-white shadow-sm transform transition-transform duration-300',
                                    form.is_active ? 'translate-x-6' : 'translate-x-0'
                                ]"></div>
                            </div>
                        </div>
                    </div>
                </div>

                <template #footer>
                    <div class="flex flex-col-reverse sm:flex-row gap-2 w-full pt-2">
                        <Button label="Cancel" severity="secondary" @click="showFormModal = false" class="flex-1 font-bold rounded-xl" text />
                        <Button :label="isEditing ? 'Save Changes' : 'Add Room'" icon="pi pi-check" @click="handleSave" class="flex-1 font-bold shadow-lg shadow-indigo-100 rounded-xl" />
                    </div>
                </template>
            </BaseModal>

            <BaseModal :show="showDeleteModal" title="Confirm Deletion" icon="pi-exclamation-triangle" @close="showDeleteModal = false">
                <div class="text-center py-4 space-y-3">
                    <p class="text-slate-600 font-medium text-sm">Are you sure you want to delete room <br><span class="font-black text-slate-800 text-lg">{{ selectedRoom?.name }}</span>?</p>
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