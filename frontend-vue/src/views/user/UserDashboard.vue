<script setup>
import { ref, onMounted, computed } from 'vue';
import { BookingService } from '@/service/BookingService';
import { RoomService } from '@/service/RoomService';
import Navbar from '../../components/layout/Navbar.vue';
import StatCard from '../../components/ui/StatCard.vue';
import BaseModal from '../../components/ui/BaseModal.vue';
import Footer from '../../components/layout/Footer.vue';

// 1. State Management
const bookings = ref([]);
const rooms = ref([]);
const loading = ref(true);
const showBookingModal = ref(false);
const showSuccessModal = ref(false);
const showErrorModal = ref(false);
const errorMessage = ref('');
const userId = localStorage.getItem('user_id'); 

// 2. Form State
const bookingForm = ref({
    room_id: null,
    start_date: '', 
    start_time: '', 
    end_date: '',   
    end_time: '',   
    purpose: ''
});

// Combine Date + Time
const combineDateTime = (date, time) => {
    return `${date} ${time}:00`; 
};

// Format Date Display (English Locale)
const formatDateDisplay = (dateTimeString) => {
    if (!dateTimeString) return '-';
    return new Date(dateTimeString).toLocaleDateString('en-GB', {
        day: '2-digit', month: 'short', year: 'numeric',
        hour: '2-digit', minute: '2-digit'
    });
};

// Get Room Name
const getRoomName = (id) => {
    const room = rooms.value.find(r => r.id === id);
    return room ? room.name : `Room #${id}`; 
};

// 3. Trigger Error Modal
const triggerError = (msg) => {
    errorMessage.value = msg;
    showErrorModal.value = true;
};

// 4. Data Fetching
const fetchData = async () => {
    loading.value = true;
    try {
        const [resBookings, resRooms] = await Promise.all([
            BookingService.getBookings(userId),
            RoomService.getRooms()
        ]);
        bookings.value = resBookings.data;
        rooms.value = resRooms.data.filter(r => r.is_active);
    } catch (err) {
        console.error("Failed to sync data:", err);
    } finally {
        loading.value = false;
    }
};

// 5. Handle Submit
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
            room_id: f.room_id,
            user_id: parseInt(userId),
            start_time: fullStart,
            end_time: fullEnd,
            purpose: f.purpose,
            status: 'pending'
        };
        await BookingService.createBooking(payload);
        showBookingModal.value = false;
        showSuccessModal.value = true;
        bookingForm.value = { room_id: null, start_date: '', start_time: '', end_date: '', end_time: '', purpose: '' };
        await fetchData();
    } catch (err) {
        console.error("Error creating booking:", err);
        const msg = err.response?.data?.error || "Failed to create booking. Please try again.";
        triggerError(msg);
    }
};

// Auto-fill End Date
const onStartDateChange = () => {
    if (!bookingForm.value.end_date) {
        bookingForm.value.end_date = bookingForm.value.start_date;
    }
};
onMounted(fetchData);

// 6. Stats Logic
const stats = computed(() => {
    return {
        total: bookings.value.length,
        approved: bookings.value.filter(b => b.status === 'approved').length,
        pending: bookings.value.filter(b => b.status === 'pending').length,
        rejected: bookings.value.filter(b => b.status === 'rejected').length
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
        <main class="flex-1 w-full max-w-7xl mx-auto p-4 md:p-8 space-y-6 md:space-y-8">
            <div class="flex flex-col md:flex-row justify-between items-start md:items-end border-b border-slate-200 pb-5 gap-4">
                <div>
                    <h2 class="text-xl md:text-2xl font-black text-slate-800 uppercase tracking-tighter italic">Lecturer Dashboard</h2>
                    <p class="text-xs md:text-sm text-slate-500 font-medium italic mt-1">Smart Door Reservation System - ITI Informatics Engineering</p>
                </div>
                <div class="hidden md:block">
                    <span class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-emerald-50 border border-emerald-100 text-emerald-600 text-[10px] font-black uppercase tracking-widest shadow-sm">
                        <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                        System Online
                    </span>
                </div>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-5">
                <StatCard title="My Total Bookings" :value="stats.total" icon="pi-calendar" variant="indigo" />
                <StatCard title="Approved Requests" :value="stats.approved" icon="pi-check-circle" variant="green" />
                <StatCard title="Pending Approval" :value="stats.pending" icon="pi-clock" variant="amber" />
                <StatCard title="Rejected Requests" :value="stats.rejected" icon="pi-times-circle" variant="red" />
            </div>
            <div class="bg-white rounded-2xl border border-slate-200 shadow-lg shadow-slate-100/50 overflow-hidden">
                <div class="p-4 md:p-6 border-b border-slate-100 bg-white flex flex-col sm:flex-row justify-between items-center gap-4">
                    <div class="w-full sm:w-auto">
                        <h3 class="font-black text-slate-800 text-sm md:text-base uppercase tracking-wider flex items-center gap-3">
                            <div class="p-2 bg-indigo-50 rounded-lg text-indigo-600"><i class="pi pi-list"></i></div>
                            Booking History
                        </h3>
                    </div>
                    <Button 
                        label="New Booking" icon="pi pi-plus" 
                        class="w-full sm:w-auto bg-indigo-600 hover:bg-indigo-700 text-white font-bold shadow-md shadow-indigo-100 px-5 py-2.5 rounded-lg transition-all text-xs" 
                        @click="showBookingModal = true" 
                    />
                </div>
                <div class="p-2">
                    <DataTable 
                        :value="bookings" :loading="loading" paginator :rows="5" responsiveLayout="scroll" class="p-datatable-sm"
                        :pt="{ header: { class: 'bg-slate-50 text-slate-800 font-bold uppercase text-xs tracking-wider border-b border-slate-100' } }"
                    >
                        <template #empty>
                            <div class="text-center py-12 text-slate-400">
                                <i class="pi pi-folder-open text-5xl mb-3 block opacity-20 mx-auto"></i>
                                <p class="font-bold text-xs uppercase tracking-widest">No booking data found.</p>
                            </div>
                        </template>
                        <Column header="Room Name" field="room_id" headerClass="text-left pl-4 md:pl-6" bodyClass="text-left pl-4 md:pl-6">
                            <template #body="{ data }">
                                <div class="flex items-center gap-2">
                                    <div class="w-8 h-8 rounded-lg bg-indigo-50 text-indigo-600 flex items-center justify-center shrink-0">
                                        <i class="pi pi-building text-xs"></i>
                                    </div>
                                    <span class="font-black text-slate-700 text-sm whitespace-nowrap">{{ getRoomName(data.room_id) }}</span>
                                </div>
                            </template>
                        </Column>
                        <Column header="Start Time" sortable field="start_time">
                            <template #body="{ data }">
                                <span class="font-bold text-slate-600 text-xs block whitespace-nowrap">{{ formatDateDisplay(data.start_time) }}</span>
                            </template>
                        </Column>
                        <Column header="End Time" field="end_time">
                            <template #body="{ data }">
                                <span class="font-medium text-slate-500 text-xs block whitespace-nowrap">{{ formatDateDisplay(data.end_time) }}</span>
                            </template>
                        </Column>
                        <Column field="purpose" header="Purpose">
                            <template #body="{ data }">
                                <span class="font-medium text-slate-600 text-sm block truncate max-w-37.5 md:max-w-xs" v-tooltip="data.purpose">{{ data.purpose }}</span>
                            </template>
                        </Column>
                        <Column field="status" header="Status" headerClass="text-center" bodyClass="text-center">
                            <template #body="{ data }">
                                <Tag :value="data.status.toUpperCase()" :severity="getStatusSeverity(data.status)" class="px-2 py-1 text-[10px] font-black tracking-widest uppercase rounded shadow-sm" />
                            </template>
                        </Column>
                    </DataTable>
                </div>
            </div>
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
                    <div class="bg-slate-50 p-4 rounded-xl border border-slate-100">
                        <p class="text-xs font-black text-indigo-600 uppercase mb-3 flex items-center gap-2"><i class="pi pi-clock"></i> Start Schedule</p>
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
                    </div>
                    <div class="bg-slate-50 p-4 rounded-xl border border-slate-100">
                        <p class="text-xs font-black text-slate-500 uppercase mb-3 flex items-center gap-2"><i class="pi pi-stopwatch"></i> End Schedule</p>
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
                        <Button label="Submit Booking" icon="pi pi-send" @click="handleCreateBooking" class="flex-1 bg-indigo-600 hover:bg-indigo-700 border-none font-bold shadow-lg shadow-indigo-100 rounded-lg" />
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
                        <p class="text-sm text-slate-500 font-medium px-4">Your booking request is being processed by Admin. Please check status periodically.</p>
                    </div>
                </div>
                <template #footer>
                    <Button label="Got it" class="w-full bg-slate-800 hover:bg-slate-900 border-none font-bold py-3 shadow-lg rounded-lg" @click="showSuccessModal = false" />
                </template>
            </BaseModal>
            <BaseModal :show="showErrorModal" title="Action Failed" icon="pi-exclamation-triangle" @close="showErrorModal = false">
                <div class="flex flex-col items-center py-6 space-y-4">
                    <div class="w-16 h-16 bg-red-50 text-red-500 rounded-full flex items-center justify-center border-4 border-red-100 animate-pulse">
                        <i class="pi pi-times text-3xl font-bold"></i>
                    </div>
                    <div class="text-center space-y-1">
                        <h4 class="text-lg font-black text-slate-800">Request Rejected</h4>
                        <p class="text-sm text-slate-600 font-bold px-4 leading-relaxed">
                            {{ errorMessage }}
                        </p>
                    </div>
                </div>
                <template #footer>
                    <Button label="Close & Fix" severity="danger" class="w-full font-bold py-3 shadow-lg shadow-red-100 rounded-lg" @click="showErrorModal = false" />
                </template>
            </BaseModal>
        </main>
        <Footer />
    </div>
</template>