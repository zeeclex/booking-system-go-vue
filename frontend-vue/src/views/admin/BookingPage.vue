<script setup>
import { ref, onMounted, computed } from 'vue';
import { BookingService } from '@/service/BookingService';
import { RoomService } from '@/service/RoomService'; 
import AdminSidebar from '../../components/admin/Sidebar.vue';
import StatCard from '../../components/ui/StatCard.vue';

// --- STATE MANAGEMENT ---
const bookings = ref([]);
const rooms = ref([]);
const loading = ref(true);

// Mobile Sidebar State
const isSidebarOpen = ref(false);

// --- HELPER: Format DateTime ---
const formatDateTime = (dateString) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('en-GB', { 
        day: '2-digit', month: 'short', year: 'numeric', 
        hour: '2-digit', minute: '2-digit', hour12: false 
    });
};

// --- HELPER: Calculate Duration ---
const calculateDuration = (start, end) => {
    const s = new Date(start);
    const e = new Date(end);
    const diffMs = e - s;
    const diffHrs = diffMs / (1000 * 60 * 60);
    return diffHrs.toFixed(1) + ' Hrs'; 
};

// --- HELPER: Lookup Room Name ---
const getRoomName = (id) => {
    const room = rooms.value.find(r => r.id === id);
    return room ? room.name : `Room #${id}`;
};

// --- COMPUTED STATS ---
const bookingStats = computed(() => {
    return {
        total: bookings.value.length,
        approved: bookings.value.filter(b => b.status === 'approved').length,
        pending: bookings.value.filter(b => b.status === 'pending').length,
        rejected: bookings.value.filter(b => b.status === 'rejected').length
    };
});

// --- SORTING LOGIC ---
const sortedBookings = computed(() => {
    return [...bookings.value].sort((a, b) => {
        const priority = { 'pending': 1, 'approved': 2, 'rejected': 3 };
        return priority[a.status] - priority[b.status];
    });
});

// --- LOAD DATA ---
const loadData = async () => {
    loading.value = true;
    try {
        const [resBookings, resRooms] = await Promise.all([
            BookingService.getBookings(),
            RoomService.getRooms()
        ]);
        bookings.value = resBookings.data;
        rooms.value = resRooms.data; 
    } catch (err) {
        console.error("Failed to load data:", err);
    } finally {
        loading.value = false;
    }
};

const handleStatusUpdate = async (id, action) => {
    try {
        await BookingService.updateStatus(id, action);
        await loadData();
    } catch (err) {
        console.error(`Failed to ${action} booking:`, err);
        alert(`Failed to ${action} booking. Please try again.`);
    }
};

onMounted(loadData);

const getStatusSeverity = (status) => {
    switch (status) {
        case 'approved': return 'success';
        case 'pending': return 'warn';
        case 'rejected': return 'danger';
        default: return 'info';
    }
};
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
                        <h2 class="text-xl md:text-2xl font-black text-slate-800 uppercase tracking-tighter italic">Booking Manager</h2>
                        <p class="text-xs md:text-sm text-slate-500 font-medium italic mt-1">Manage & Approve Room Requests</p>
                    </div>
                </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-5">
                <StatCard title="Total Requests" :value="bookingStats.total" icon="pi-list" variant="indigo" />
                <StatCard title="Approved" :value="bookingStats.approved" icon="pi-check-circle" variant="green" />
                <StatCard title="Pending Review" :value="bookingStats.pending" icon="pi-clock" variant="amber" />
                <StatCard title="Rejected" :value="bookingStats.rejected" icon="pi-times-circle" variant="red" />
            </div>
            <div class="bg-white rounded-2xl border border-slate-200 shadow-lg shadow-slate-100/50 overflow-hidden">
                <div class="p-4 md:p-6 border-b border-slate-100 bg-white flex flex-col sm:flex-row justify-between items-center gap-4">
                    <div class="w-full sm:w-auto">
                        <h3 class="font-black text-slate-800 text-sm md:text-base uppercase tracking-wider flex items-center gap-3">
                            <div class="p-2 bg-indigo-50 rounded-lg text-indigo-600"><i class="pi pi-list"></i></div>
                            Approval Queue
                        </h3>
                        <p class="text-[10px] md:text-xs text-slate-400 font-bold mt-1 ml-0 sm:ml-11 uppercase tracking-widest">Sorted by Priority (Pending First)</p>
                    </div>
                    <Button 
                        icon="pi pi-refresh" 
                        text 
                        rounded 
                        @click="loadData" 
                        :loading="loading" 
                        class="text-slate-400 hover:text-indigo-600 hover:bg-indigo-50"
                    />
                </div>
                <div class="p-0 md:p-2">
                    <DataTable 
                        :value="sortedBookings" 
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
                                <i class="pi pi-inbox text-5xl mb-3 block opacity-20 mx-auto"></i>
                                <p class="font-bold text-xs uppercase tracking-widest">No booking requests found.</p>
                            </div>
                        </template>
                        <Column header="Requester" headerClass="text-left pl-4 md:pl-6" bodyClass="text-left pl-4 md:pl-6">
                            <template #body="{ data }">
                                <div class="flex flex-col">
                                    <span class="font-black text-slate-700 text-sm whitespace-nowrap">User #{{ data.user_id }}</span>
                                    <span class="text-[10px] text-slate-400 font-bold uppercase tracking-wider">Lecturer</span>
                                </div>
                            </template>
                        </Column>
                        <Column header="Room & Duration" headerClass="text-left" bodyClass="text-left">
                            <template #body="{ data }">
                                <div class="flex flex-col">
                                    <span class="font-bold text-indigo-600 text-xs whitespace-nowrap">{{ getRoomName(data.room_id) }}</span>
                                    <span class="text-[10px] font-mono text-slate-500 bg-slate-100 px-1.5 py-0.5 rounded w-fit mt-1 whitespace-nowrap">
                                        Dur: {{ calculateDuration(data.start_time, data.end_time) }}
                                    </span>
                                </div>
                            </template>
                        </Column>
                        <Column header="Schedule" headerClass="text-left" bodyClass="text-left">
                            <template #body="{ data }">
                                <div class="flex flex-col gap-1">
                                    <span class="text-xs font-bold text-slate-700 whitespace-nowrap">
                                        <i class="pi pi-calendar-plus mr-1 text-green-500"></i>
                                        {{ formatDateTime(data.start_time) }}
                                    </span>
                                    <span class="text-xs font-bold text-slate-500 whitespace-nowrap">
                                        <i class="pi pi-calendar-minus mr-1 text-red-400"></i>
                                        {{ formatDateTime(data.end_time) }}
                                    </span>
                                </div>
                            </template>
                        </Column>
                        <Column field="purpose" header="Purpose" headerClass="text-left" bodyClass="text-left">
                            <template #body="{ data }">
                                <span class="font-medium text-slate-600 text-xs block truncate max-w-37.5 md:max-w-xs cursor-help border-b border-dashed border-slate-300 pb-0.5" v-tooltip.top="data.purpose">
                                    {{ data.purpose }}
                                </span>
                            </template>
                        </Column>
                        <Column field="status" header="Status" headerClass="text-center" bodyClass="text-center">
                            <template #body="{ data }">
                                <Tag 
                                    :value="data.status.toUpperCase()" 
                                    :severity="getStatusSeverity(data.status)" 
                                    class="px-2 py-1 text-[10px] font-black tracking-widest uppercase rounded shadow-sm"
                                />
                            </template>
                        </Column>
                        <Column header="Actions" headerClass="text-center" bodyClass="text-center">
                            <template #body="{ data }">
                                <div class="flex justify-center gap-2" v-if="data.status === 'pending'">
                                    <Button icon="pi pi-check" severity="success" rounded outlined size="small" class="w-8 h-8 p-0! shadow-sm hover:bg-green-600 hover:text-white transition-all" v-tooltip.top="'Approve'" @click="handleStatusUpdate(data.id, 'approve')" />
                                    <Button icon="pi pi-times" severity="danger" rounded outlined size="small" class="w-8 h-8 p-0! shadow-sm hover:bg-red-600 hover:text-white transition-all" v-tooltip.top="'Reject'" @click="handleStatusUpdate(data.id, 'reject')" />
                                </div>
                                <div v-else class="opacity-40 select-none">
                                    <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest bg-slate-100 px-2 py-1 rounded">Done</span>
                                </div>
                            </template>
                        </Column>
                    </DataTable>
                </div>
            </div>
        </main>
    </div>
</template>