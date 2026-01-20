<script setup>
import { ref, onMounted } from 'vue';
import { BookingService } from '@/service/BookingService';
import { RoomService } from '@/service/RoomService';
import { UserService } from '@/service/UserService';
import AdminSidebar from '../../components/admin/Sidebar.vue';
import LineChart from '../../components/admin/LineChart.vue'; 

// STATE
const stats = ref({
    total_rooms: 0,
    rooms_active: 0,
    rooms_inactive: 0,
    total_bookings: 0,
    total_approved: 0,
    total_rejected: 0,
    pending_requests: 0,
    total_users: 0 
});
const loading = ref(true);
const rawBookings = ref([]); 

// Mobile Sidebar State
const isSidebarOpen = ref(false);

// DATA FETCHING
const loadDashboardData = async () => {
    loading.value = true;
    try {
        const [resStats, resBookings, resRooms, resUsers] = await Promise.all([
            BookingService.getStats(),
            BookingService.getBookings(),
            RoomService.getRooms(),
            UserService.getUsers()
        ]);
        const allRooms = resRooms.data;
        const allUsers = resUsers.data;
        rawBookings.value = resBookings.data;
        stats.value = { 
            total_rooms: allRooms.length,
            rooms_active: allRooms.filter(r => r.is_active).length,
            rooms_inactive: allRooms.filter(r => !r.is_active).length,
            total_bookings: rawBookings.value.length,
            total_approved: rawBookings.value.filter(b => b.status === 'approved').length,
            total_rejected: rawBookings.value.filter(b => b.status === 'rejected').length,
            pending_requests: rawBookings.value.filter(b => b.status === 'pending').length,
            total_users: allUsers.length || 0, 
        };
    } catch (err) {
        console.error("Dashboard sync error:", err);
    } finally {
        loading.value = false;
    }
};

onMounted(loadDashboardData);
</script>

<template>
    <div class="flex min-h-screen bg-slate-50 font-sans">
        <AdminSidebar :isOpen="isSidebarOpen" @close-sidebar="isSidebarOpen = false" />
        <main class="flex-1 w-full max-w-7xl mx-auto p-4 md:p-8 space-y-6 md:space-y-8 overflow-y-auto">
            <div class="flex flex-col md:flex-row justify-between items-start md:items-end gap-4">
                <div class="flex items-center gap-4">
                    <button @click="isSidebarOpen = true" class="md:hidden p-2 -ml-2 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">
                        <i class="pi pi-bars text-xl"></i>
                    </button>
                    <div>
                        <h2 class="text-2xl md:text-3xl font-black text-slate-800 uppercase tracking-tighter italic">Dashboard</h2>
                        <p class="text-xs md:text-sm text-slate-500 font-medium italic mt-1">ITIDoorz Analytics & Overview</p>
                    </div>
                </div>
                <div class="hidden md:flex items-center gap-3 px-4 py-2 bg-white border border-slate-200 rounded-full shadow-sm">
                    <span class="relative flex h-2.5 w-2.5">
                        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                        <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-500"></span>
                    </span>
                    <span class="text-xs font-bold text-slate-600 uppercase tracking-widest">System Online</span>
                </div>
            </div>
            <LineChart :bookings="rawBookings" :loading="loading" />
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-6 relative z-10">
                <div class="group relative bg-white h-40 rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:border-blue-300 transition-all duration-300 overflow-hidden">
                    <div class="p-6 h-full flex flex-col justify-center group-hover:opacity-0 transition-opacity duration-300">
                        <div class="flex justify-between items-center">
                            <div>
                                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Inventory</p>
                                <h3 class="text-3xl md:text-4xl font-black text-slate-800 mt-1">{{ stats.total_rooms }}</h3>
                                <p class="text-xs text-slate-400 font-bold mt-1">Total Rooms</p>
                            </div>
                            <div class="w-12 h-12 md:w-14 md:h-14 rounded-2xl bg-blue-50 text-blue-600 flex items-center justify-center">
                                <i class="pi pi-building text-xl md:text-2xl"></i>
                            </div>
                        </div>
                    </div>
                    <div class="absolute inset-0 bg-slate-800 p-6 flex flex-col justify-center translate-y-full group-hover:translate-y-0 transition-transform duration-300 ease-out">
                        <div class="flex flex-col space-y-3">
                            <div class="flex justify-between items-center border-b border-slate-600 pb-2">
                                <span class="text-slate-400 font-bold text-xs uppercase">Total Rooms</span>
                                <span class="text-white font-black text-xl">{{ stats.total_rooms }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-emerald-400 font-bold text-xs flex gap-2 items-center"><i class="pi pi-check-circle"></i> Active</span>
                                <span class="text-white font-bold text-sm">{{ stats.rooms_active }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-slate-400 font-bold text-xs flex gap-2 items-center"><i class="pi pi-ban"></i> Inactive</span>
                                <span class="text-white font-bold text-sm">{{ stats.rooms_inactive }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="group relative bg-white h-40 rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:border-indigo-300 transition-all duration-300 overflow-hidden">
                    <div class="p-6 h-full flex flex-col justify-center group-hover:opacity-0 transition-opacity duration-300">
                        <div class="flex justify-between items-center">
                            <div>
                                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Reservations</p>
                                <h3 class="text-3xl md:text-4xl font-black text-slate-800 mt-1">{{ stats.total_bookings }}</h3>
                                <p class="text-xs text-slate-400 font-bold mt-1">Total Requests</p>
                            </div>
                            <div class="w-12 h-12 md:w-14 md:h-14 rounded-2xl bg-indigo-50 text-indigo-600 flex items-center justify-center">
                                <i class="pi pi-calendar text-xl md:text-2xl"></i>
                            </div>
                        </div>
                    </div>
                    <div class="absolute inset-0 bg-slate-800 p-6 flex flex-col justify-center translate-y-full group-hover:translate-y-0 transition-transform duration-300 ease-out">
                        <div class="flex flex-col space-y-2">
                            <div class="flex justify-between items-end border-b border-slate-600 pb-1 mb-1">
                                <span class="text-slate-400 font-bold text-xs uppercase">Total</span>
                                <span class="text-white font-black text-lg">{{ stats.total_bookings }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-emerald-400 font-bold text-xs flex gap-2 items-center"><i class="pi pi-check"></i> Approved</span>
                                <span class="text-white font-bold text-sm">{{ stats.total_approved }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-amber-400 font-bold text-xs flex gap-2 items-center"><i class="pi pi-clock"></i> Pending</span>
                                <span class="text-white font-bold text-sm">{{ stats.pending_requests }}</span>
                            </div>
                            <div class="flex justify-between items-center">
                                <span class="text-red-400 font-bold text-xs flex gap-2 items-center"><i class="pi pi-times"></i> Rejected</span>
                                <span class="text-white font-bold text-sm">{{ stats.total_rejected }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="group relative bg-white h-40 rounded-2xl border border-slate-200 shadow-sm hover:shadow-xl hover:border-slate-300 transition-all duration-300 overflow-hidden">
                    <div class="p-6 h-full flex flex-col justify-center group-hover:opacity-0 transition-opacity duration-300">
                        <div class="flex justify-between items-center">
                            <div>
                                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">User Base</p>
                                <h3 class="text-3xl md:text-4xl font-black text-slate-800 mt-1">{{ stats.total_users }}</h3>
                                <p class="text-xs text-slate-400 font-bold mt-1">Total Users</p>
                            </div>
                            <div class="w-12 h-12 md:w-14 md:h-14 rounded-2xl bg-slate-100 text-slate-600 flex items-center justify-center">
                                <i class="pi pi-users text-xl md:text-2xl"></i>
                            </div>
                        </div>
                    </div>
                    <div class="absolute inset-0 bg-slate-800 p-6 flex flex-col justify-center items-center translate-y-full group-hover:translate-y-0 transition-transform duration-300 ease-out text-center">
                        <i class="pi pi-users text-4xl md:text-5xl text-indigo-400 mb-2 md:mb-3"></i>
                        <span class="text-3xl md:text-4xl font-black text-white block">{{ stats.total_users }}</span>
                        <div class="w-10 h-1 md:w-12 bg-emerald-500 rounded-full mt-2 md:mt-3 mb-2"></div>
                        <span class="text-[10px] md:text-xs font-bold text-slate-400 uppercase tracking-widest">Registered Accounts</span>
                    </div>
                </div>
            </div>
            <div class="text-center pt-8 pb-4 opacity-50">
                <p class="text-[10px] font-black uppercase tracking-[0.3em]">
                    ITIDoorz v1.0 • Engineering
                </p>
            </div>
        </main>
    </div>
</template>