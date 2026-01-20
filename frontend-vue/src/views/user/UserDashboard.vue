<script setup>
import { ref, onMounted, computed } from 'vue';
import { BookingService } from '@/service/BookingService';
import { RoomService } from '@/service/RoomService';

// Components
import Navbar from '../../components/layout/Navbar.vue';
import Footer from '../../components/layout/Footer.vue';
import StatCard from '../../components/ui/StatCard.vue';
import BaseModal from '../../components/ui/BaseModal.vue';
import BookingCalendar from '../../components/ui/BookingCalendar.vue';

// --- STATE ---
const allBookings = ref([]);
const userBookings = ref([]);
const rooms = ref([]);
const loading = ref(true);
const viewMode = ref('list');

// Modals
const showBookingModal = ref(false);
const showSuccessModal = ref(false);
const showErrorModal = ref(false);
const errorMessage = ref('');

// User Info
const userId = parseInt(localStorage.getItem('user_id')); 

// Form
const bookingForm = ref({
    room_id: null,
    start_date: '', start_time: '', 
    end_date: '', end_time: '', 
    purpose: ''
});

// --- HELPERS ---
const combineDateTime = (date, time) => `${date} ${time}:00`;
const formatDateTime = (dateString) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('en-GB', { 
        day: '2-digit', month: 'short', year: 'numeric', 
        hour: '2-digit', minute: '2-digit' 
    });
};

// Helper for calculating duration
const calculateDuration = (start, end) => {
    const s = new Date(start);
    const e = new Date(end);
    const diffMs = e - s;
    const diffHrs = diffMs / (1000 * 60 * 60);
    return Number.isInteger(diffHrs) ? `${diffHrs}H` : `${diffHrs.toFixed(1)}H`;
};
const getRoomName = (id) => {
    const room = rooms.value.find(r => r.id === id);
    return room ? room.name : `Room #${id}`;
};
const triggerError = (msg) => {
    errorMessage.value = msg;
    showErrorModal.value = true;
};

// --- DATA FETCHING ---
const fetchData = async () => {
    loading.value = true;
    try {
        const [resAll, resUser, resRooms] = await Promise.all([
            BookingService.getBookings(),       // Ambil semua (buat kalender)
            BookingService.getBookings(userId), // Ambil punya user (buat list history)
            RoomService.getRooms()
        ]);
        allBookings.value = resAll.data;
        userBookings.value = resUser.data;
        rooms.value = resRooms.data.filter(r => r.is_active);
    } catch (err) {
        console.error("Failed to sync data:", err);
    } finally {
        loading.value = false;
    }
};

onMounted(fetchData);

// --- PRIVACY LOGIC ---
const sanitizedCalendarBookings = computed(() => {
    return allBookings.value.map(b => {
        const isMyBooking = b.user_id === userId;
        if (isMyBooking) {
            return b;
        } else {
            return {
                ...b,
                room_name: b.room_name, 
                user_name: 'Occupied',
                purpose: 'Reserved',
                status: 'occupied'
            };
        }
    });
});

// --- SUBMIT HANDLER ---
const handleCreateBooking = async () => {
    try {
        const f = bookingForm.value;
        if (!f.room_id || !f.start_date || !f.start_time || !f.end_date || !f.end_time || !f.purpose) {
            triggerError("Please fill in all required fields!");
            return;
        }
        const fullStart = combineDateTime(f.start_date, f.start_time);
        const fullEnd = combineDateTime(f.end_date, f.end_time);
        if (new Date(fullStart) >= new Date(fullEnd)) {
            triggerError("End Time must be later than Start Time!");
            return;
        }
        const payload = {
            room_id: f.room_id, user_id: userId,
            start_time: fullStart, end_time: fullEnd,
            purpose: f.purpose, status: 'pending'
        };
        await BookingService.createBooking(payload);
        showBookingModal.value = false;
        showSuccessModal.value = true;
        bookingForm.value = { room_id: null, start_date: '', start_time: '', end_date: '', end_time: '', purpose: '' };
        await fetchData();
    } catch (err) {
        const msg = err.response?.data?.error || "Failed to create booking.";
        triggerError(msg);
    }
};

const onStartDateChange = () => {
    if (!bookingForm.value.end_date) bookingForm.value.end_date = bookingForm.value.start_date;
};

// Stats
const stats = computed(() => {
    return {
        total: userBookings.value.length,
        approved: userBookings.value.filter(b => b.status === 'approved').length,
        pending: userBookings.value.filter(b => b.status === 'pending').length,
        rejected: userBookings.value.filter(b => b.status === 'rejected').length
    };
});

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
    <div class="min-h-screen flex flex-col bg-slate-50 font-sans">
        <Navbar />
        
        <main class="flex-1 w-full max-w-7xl mx-auto p-4 md:p-8 space-y-8">
            
            <div class="bg-indigo-600 rounded-3xl p-6 md:p-10 text-white shadow-xl shadow-indigo-200 flex flex-col md:flex-row justify-between items-center gap-6 relative overflow-hidden">
                <div class="absolute top-0 right-0 w-64 h-64 bg-white/10 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2 pointer-events-none"></div>
                <div class="relative z-10 text-center md:text-left">
                    <h2 class="text-2xl md:text-3xl font-black uppercase tracking-tighter">Lecturer Dashboard</h2>
                    <p class="text-indigo-100 font-medium mt-1 text-sm md:text-base">Manage reservations & check availability.</p>
                </div>
                <button @click="showBookingModal = true" class="relative z-10 bg-white text-indigo-700 hover:bg-indigo-50 font-bold px-6 py-3 rounded-xl shadow-lg transition-transform active:scale-95 flex items-center gap-2 group w-full md:w-auto justify-center">
                    <i class="pi pi-plus-circle text-xl group-hover:rotate-90 transition-transform"></i>
                    <span>New Booking</span>
                </button>
            </div>
            <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4">
                <StatCard title="My Requests" :value="stats.total" icon="pi-briefcase" variant="indigo" />
                <StatCard title="Approved" :value="stats.approved" icon="pi-check-circle" variant="green" />
                <StatCard title="Pending" :value="stats.pending" icon="pi-clock" variant="amber" />
                <StatCard title="Rejected" :value="stats.rejected" icon="pi-times-circle" variant="red" />
            </div>
            <div class="bg-white rounded-2xl border border-slate-200 shadow-lg shadow-slate-100/50 overflow-hidden min-h-125">
                
                <div class="p-4 md:p-6 border-b border-slate-100 bg-white flex flex-col sm:flex-row justify-between items-center gap-4">
                    <div class="w-full sm:w-auto">
                        <h3 class="font-black text-slate-800 text-sm md:text-base uppercase tracking-wider flex items-center gap-3">
                            <div class="p-2 bg-indigo-50 rounded-lg text-indigo-600">
                                <i :class="viewMode === 'list' ? 'pi pi-list' : 'pi pi-calendar'"></i>
                            </div>
                            {{ viewMode === 'list' ? 'My Booking History' : 'Room Availability' }}
                        </h3>
                    </div>
                    <div class="w-full md:w-auto bg-slate-200 p-1 rounded-xl flex">
                        <button @click="viewMode = 'list'" :class="['flex-1 md:flex-none px-6 py-2 rounded-lg text-xs font-bold uppercase tracking-wider flex items-center justify-center gap-2 transition-all', viewMode === 'list' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700']">
                            <i class="pi pi-list"></i> My History
                        </button>
                        <button @click="viewMode = 'calendar'" :class="['flex-1 md:flex-none px-6 py-2 rounded-lg text-xs font-bold uppercase tracking-wider flex items-center justify-center gap-2 transition-all', viewMode === 'calendar' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700']">
                            <i class="pi pi-calendar"></i> Calendar
                        </button>
                    </div>
                </div>
                <div v-if="viewMode === 'list'" class="animate-fade-in p-2">
                    <DataTable :value="userBookings" :loading="loading" paginator :rows="10" responsiveLayout="scroll" class="p-datatable-sm" :pt="{ header: { class: 'bg-slate-50 text-slate-800 font-bold uppercase text-xs tracking-wider border-b border-slate-100 py-3' } }">
                        <template #empty>
                            <div class="text-center py-12 text-slate-400">
                                <i class="pi pi-folder-open text-5xl mb-3 block opacity-20 mx-auto"></i>
                                <p class="font-bold text-xs uppercase tracking-widest">No history found.</p>
                            </div>
                        </template>
                        <Column header="Room Name" field="room_name" sortable>
                            <template #body="{ data }">
                                <span class="font-black text-slate-700">{{ data.room_name || getRoomName(data.room_id) }}</span>
                            </template>
                        </Column>
                        
                        <Column header="Duration" headerClass="text-center" bodyClass="text-center">
                            <template #body="{ data }">
                                <span class="font-bold text-slate-700 bg-slate-100 px-2 py-1 rounded text-xs border border-slate-200">
                                    {{ calculateDuration(data.start_time, data.end_time) }}
                                </span>
                            </template>
                        </Column>
                        <Column header="Start Date" field="start_time" sortable>
                            <template #body="{ data }">
                                <span class="font-medium text-slate-600">{{ formatDateTime(data.start_time) }}</span>
                            </template>
                        </Column>
                        <Column header="End Date" field="end_time" sortable>
                            <template #body="{ data }">
                                <span class="font-medium text-slate-600">{{ formatDateTime(data.end_time) }}</span>
                            </template>
                        </Column>
                        <Column header="Purpose" field="purpose">
                            <template #body="{ data }">
                                <span class="italic text-slate-500 text-sm" v-tooltip="data.purpose">{{ data.purpose }}</span>
                            </template>
                        </Column>
                        <Column header="Status" field="status" class="text-center">
                            <template #body="{ data }">
                                <Tag :value="data.status.toUpperCase()" :severity="getStatusSeverity(data.status)" class="px-2 py-1 text-[10px] font-black tracking-widest uppercase" />
                            </template>
                        </Column>
                    </DataTable>
                </div>
                <div v-else class="p-4 animate-fade-in">
                    <div class="bg-blue-50 border border-blue-100 p-3 rounded-lg mb-4 flex gap-3 text-xs text-blue-800">
                        <i class="pi pi-info-circle text-lg"></i>
                        <p><strong>Privacy Mode:</strong> You can only see details of your own bookings. Other reservations are marked as "Occupied".</p>
                    </div>
                    <BookingCalendar :bookings="sanitizedCalendarBookings" />
                </div>
            </div>
        </main>
        <BaseModal :show="showBookingModal" title="Room Booking Form" icon="pi-calendar-plus" @close="showBookingModal = false">
            <div class="space-y-6 py-2">
                <div class="space-y-1.5">
                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Select Room</label>
                    <div class="relative group">
                        <i class="pi pi-building absolute left-4 top-1/2 -translate-y-1/2 text-slate-400 group-hover:text-indigo-500 transition-colors"></i>
                        <select v-model="bookingForm.room_id" class="w-full pl-10 pr-4 py-3 bg-slate-50 border border-slate-200 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-slate-700 text-sm appearance-none cursor-pointer hover:bg-slate-100 transition-colors">
                            <option :value="null" disabled>-- Select a Room --</option>
                            <option v-for="room in rooms" :key="room.id" :value="room.id">
                                {{ room.name }} (Capacity: {{ room.capacity }})
                            </option>
                        </select>
                    </div>
                </div>
                <div class="bg-slate-50 p-4 rounded-xl border border-slate-100 space-y-4">
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div>
                            <label class="text-[10px] font-bold text-slate-400 mb-1 block">Start Date</label>
                            <input type="date" v-model="bookingForm.start_date" @change="onStartDateChange" class="w-full p-2.5 border border-slate-200 rounded-lg text-sm font-bold text-slate-700 outline-none focus:ring-2 focus:ring-indigo-500" />
                        </div>
                        <div>
                            <label class="text-[10px] font-bold text-slate-400 mb-1 block">Start Time</label>
                            <input type="time" v-model="bookingForm.start_time" class="w-full p-2.5 border border-slate-200 rounded-lg text-sm font-bold text-slate-700 outline-none focus:ring-2 focus:ring-indigo-500" />
                        </div>
                    </div>
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div>
                            <label class="text-[10px] font-bold text-slate-400 mb-1 block">End Date</label>
                            <input type="date" v-model="bookingForm.end_date" class="w-full p-2.5 border border-slate-200 rounded-lg text-sm font-bold text-slate-700 outline-none focus:ring-2 focus:ring-indigo-500" />
                        </div>
                        <div>
                            <label class="text-[10px] font-bold text-slate-400 mb-1 block">End Time</label>
                            <input type="time" v-model="bookingForm.end_time" class="w-full p-2.5 border border-slate-200 rounded-lg text-sm font-bold text-slate-700 outline-none focus:ring-2 focus:ring-indigo-500" />
                        </div>
                    </div>
                </div>
                <div class="space-y-1.5">
                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Booking Purpose</label>
                    <textarea v-model="bookingForm.purpose" rows="3" placeholder="e.g. Web Programming Lab Session..." class="w-full p-4 bg-slate-50 border border-slate-200 rounded-lg outline-none focus:ring-2 focus:ring-indigo-500 font-medium text-slate-700 text-sm resize-none"></textarea>
                </div>
            </div>
            <template #footer>
                <div class="flex flex-col-reverse sm:flex-row gap-3 pt-2 w-full">
                    <Button label="Cancel" text severity="secondary" @click="showBookingModal = false" class="flex-1 font-bold rounded-lg" />
                    <Button label="Submit Request" icon="pi pi-send" @click="handleCreateBooking" class="flex-1 bg-indigo-600 hover:bg-indigo-700 border-none font-bold shadow-lg shadow-indigo-100 rounded-lg" />
                </div>
            </template>
        </BaseModal>
        <BaseModal :show="showSuccessModal" title="Success" icon="pi-check-circle" @close="showSuccessModal = false">
            <div class="flex flex-col items-center py-6 space-y-4">
                <div class="w-16 h-16 bg-green-50 text-green-500 rounded-full flex items-center justify-center border-4 border-green-100 animate-bounce">
                    <i class="pi pi-check text-3xl"></i>
                </div>
                <div class="text-center space-y-1">
                    <h4 class="text-lg font-black text-slate-800">Request Submitted!</h4>
                    <p class="text-sm text-slate-500 font-medium px-4">Your booking request is being processed.</p>
                </div>
            </div>
            <template #footer>
                <Button label="Done" class="w-full bg-slate-800 hover:bg-slate-900 border-none font-bold py-3 shadow-lg rounded-lg" @click="showSuccessModal = false" />
            </template>
        </BaseModal>
        <BaseModal :show="showErrorModal" title="Failed" icon="pi-exclamation-triangle" @close="showErrorModal = false">
            <div class="flex flex-col items-center py-6 space-y-4">
                <div class="w-16 h-16 bg-red-50 text-red-500 rounded-full flex items-center justify-center border-4 border-red-100 animate-pulse">
                    <i class="pi pi-times text-3xl font-bold"></i>
                </div>
                <div class="text-center space-y-1">
                    <h4 class="text-lg font-black text-slate-800">Submission Failed</h4>
                    <p class="text-sm text-slate-600 font-bold px-4 leading-relaxed">{{ errorMessage }}</p>
                </div>
            </div>
            <template #footer>
                <Button label="Close" severity="danger" class="w-full font-bold py-3 shadow-lg shadow-red-100 rounded-lg" @click="showErrorModal = false" />
            </template>
        </BaseModal>
        <Footer />
    </div>
</template>

<style scoped>
.animate-fade-in {
    animation: fadeIn 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}
@keyframes fadeIn {
    from { opacity: 0; transform: translateY(5px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>