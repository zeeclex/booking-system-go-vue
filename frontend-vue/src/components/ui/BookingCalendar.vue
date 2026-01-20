<script setup>
import { computed, ref } from 'vue';
import FullCalendar from '@fullcalendar/vue3';
import dayGridPlugin from '@fullcalendar/daygrid';
import timeGridPlugin from '@fullcalendar/timegrid';
import listPlugin from '@fullcalendar/list';
import interactionPlugin from '@fullcalendar/interaction';
import BaseModal from './BaseModal.vue'; 

// --- PROPS ---
const props = defineProps({
  bookings: {
    type: Array,
    default: () => []
  }
});

// --- STATE ---
const calendarRef = ref(null);
const showModal = ref(false);
const selectedEvent = ref(null);

// --- HELPER: Hitung Durasi (Jam) ---
const getDurationText = (start, end) => {
    const diff = new Date(end) - new Date(start);
    const hours = diff / (1000 * 60 * 60);
    return Number.isInteger(hours) ? `${hours}H` : `${hours.toFixed(1)}H`;
};

// --- DATA TRANSFORMATION ---
const calendarEvents = computed(() => {
  return props.bookings
    .filter(b => b.status !== 'rejected')
    .map(b => {
      const duration = getDurationText(b.start_time, b.end_time);
      const baseTitle = b.room_name || 'Room ' + b.room_id;
      let statusClass = 'status-pending'; 
      if (b.status === 'approved') statusClass = 'status-approved';
      if (b.status === 'occupied') statusClass = 'status-occupied';
      return {
        id: b.id,
        title: `[${duration}] ${baseTitle}`, 
        start: b.start_time,
        end: b.end_time,
        classNames: [statusClass], 
        extendedProps: {
          status: b.status,
          purpose: b.purpose,
          user: b.user_name || 'Unknown User',
          room: b.room_name || 'Room ' + b.room_id,
          duration: duration
        }
      };
    });
});

// --- CALENDAR CONFIG ---
const calendarOptions = ref({
  plugins: [dayGridPlugin, timeGridPlugin, listPlugin, interactionPlugin],
  initialView: window.innerWidth < 768 ? 'listWeek' : 'dayGridMonth',
  headerToolbar: {
    left: 'prev,next today',
    center: 'title',
    right: 'dayGridMonth,timeGridWeek,listWeek'
  },
  
  // UX Fixes:
  eventDisplay: 'block',
  eventTimeFormat: {
    hour: '2-digit',
    minute: '2-digit',
    meridiem: false,
    hour12: false 
  },
  events: calendarEvents,
  height: 'auto',
  contentHeight: 600,
  slotMinTime: '07:00:00',
  slotMaxTime: '21:00:00',
  navLinks: true,
  eventClick: handleEventClick,
  windowResize: (arg) => {
    const api = calendarRef.value.getApi();
    api.changeView(window.innerWidth < 768 ? 'listWeek' : 'dayGridMonth');
  },
});

// --- HANDLERS ---
function handleEventClick(info) {
  selectedEvent.value = {
    title: info.event.title,
    start: info.event.start,
    end: info.event.end,
    status: info.event.extendedProps.status,
    purpose: info.event.extendedProps.purpose,
    user: info.event.extendedProps.user,
    room: info.event.extendedProps.room,
    duration: info.event.extendedProps.duration
  };
  showModal.value = true;
}

function formatDateTime(date) {
  if (!date) return '-';
  return new Date(date).toLocaleString('en-GB', {
    weekday: 'short', month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit'
  });
}

// Badge Color Helper
function getStatusBadgeClass(status) {
    if (status === 'approved') return 'bg-emerald-100 text-emerald-700 border-emerald-200';
    if (status === 'occupied') return 'bg-slate-100 text-slate-600 border-slate-200';
    return 'bg-amber-100 text-amber-700 border-amber-200';
}
</script>

<template>
  <div class="bg-white p-2 md:p-4 rounded-xl shadow-sm border border-slate-200 calendar-wrapper font-sans relative">
    <FullCalendar ref="calendarRef" :options="calendarOptions" />
    <BaseModal 
        :show="showModal" 
        title="Booking Details" 
        icon="pi-calendar-clock" 
        @close="showModal = false"
    >
        <div v-if="selectedEvent" class="space-y-5">
            
            <div class="flex justify-between items-start border-b border-slate-100 pb-4">
                <div>
                    <p class="text-xs text-slate-400 uppercase font-bold tracking-wider mb-1">Room Name</p>
                    <p class="text-xl font-black text-slate-800 leading-none">{{ selectedEvent.room }}</p>
                </div>
                <div class="text-right">
                    <span :class="['px-3 py-1 rounded-full text-xs font-bold uppercase border block mb-1', getStatusBadgeClass(selectedEvent.status)]">
                        {{ selectedEvent.status }}
                    </span>
                    <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ selectedEvent.duration }} Duration</span>
                </div>
            </div>
            <div class="grid grid-cols-1 gap-4 text-sm">
                <div class="bg-slate-50 p-3 rounded-xl border border-slate-100 flex items-center gap-3">
                    <div class="w-10 h-10 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-600">
                        <i class="pi pi-user text-lg"></i>
                    </div>
                    <div>
                        <p class="text-[10px] text-slate-400 font-bold uppercase">Booked By</p>
                        <p class="text-slate-700 font-bold text-base">{{ selectedEvent.user }}</p>
                    </div>
                </div>
                <div class="bg-slate-50 p-3 rounded-xl border border-slate-100">
                    <p class="text-[10px] text-slate-400 font-bold uppercase mb-2">Time Schedule</p>
                    <div class="space-y-2 text-slate-700">
                        <div class="flex items-center gap-3">
                            <i class="pi pi-clock text-emerald-500"></i>
                            <span class="font-medium">Start: <b>{{ formatDateTime(selectedEvent.start) }}</b></span>
                        </div>
                        <div class="flex items-center gap-3">
                            <i class="pi pi-clock text-rose-500"></i>
                            <span class="font-medium">End: &nbsp;&nbsp;<b>{{ formatDateTime(selectedEvent.end) }}</b></span>
                        </div>
                    </div>
                </div>
                <div v-if="selectedEvent.status !== 'occupied'">
                    <p class="text-xs text-slate-400 font-bold uppercase mb-2">Purpose / Activity</p>
                    <div class="bg-indigo-50/50 p-3 rounded-lg border border-indigo-100">
                        <p class="text-slate-600 italic leading-relaxed">"{{ selectedEvent.purpose }}"</p>
                    </div>
                </div>
                <div v-else>
                    <div class="bg-slate-50 p-3 rounded-lg border border-slate-100 text-center">
                        <p class="text-slate-400 text-xs italic">Details are hidden for privacy.</p>
                    </div>
                </div>
            </div>
        </div>
        <template #footer>
            <button 
                @click="showModal = false" 
                class="w-full sm:w-auto px-5 py-2.5 bg-slate-100 hover:bg-slate-200 text-slate-600 font-bold rounded-xl transition-colors text-sm flex items-center justify-center gap-2"
            >
                Close Details
            </button>
        </template>
    </BaseModal>
  </div>
</template>

<style>
.fc-event {
    box-shadow: 0 2px 4px rgba(0,0,0,0.02) !important;
    border: none !important; /* Hapus border default */
    border-radius: 4px !important;
    margin-bottom: 2px !important;
    cursor: pointer;
    transition: transform 0.1s;
}
.fc-event:hover {
    transform: scale(1.01);
    filter: brightness(0.97);
}
.fc-event-main {
    padding: 2px 4px !important;
    white-space: normal !important; /* FIX: Biar teks panjang turun ke bawah */
    overflow: hidden;
    line-height: 1.2;
}
.fc-event-title { 
    font-weight: 700 !important; 
    font-size: 0.7rem !important; 
}
.fc-event-time {
    font-weight: 600 !important;
    font-size: 0.65rem !important;
    margin-right: 4px;
    opacity: 0.8;
}
.status-approved {
    background-color: #ecfdf5 !important; /* Emerald-50 */
    border-left: 4px solid #10b981 !important; /* Emerald-500 */
    color: #065f46 !important; /* Emerald-800 */
}
.status-approved .fc-event-title, .status-approved .fc-event-time { color: #065f46 !important; }
.status-pending {
    background-color: #fffbeb !important; /* Amber-50 */
    border-left: 4px solid #f59e0b !important; /* Amber-500 */
    color: #92400e !important; /* Amber-800 */
}
.status-pending .fc-event-title, .status-pending .fc-event-time { color: #92400e !important; }
.status-occupied {
    background-color: #f1f5f9 !important; /* Slate-100 */
    border-left: 4px solid #94a3b8 !important; /* Slate-400 */
    color: #475569 !important; /* Slate-600 */
}
.status-occupied .fc-event-title, .status-occupied .fc-event-time { color: #475569 !important; }
.fc-toolbar-title { font-size: 1.1rem !important; font-weight: 800; color: #1e293b; text-transform: uppercase; }
.fc .fc-button-primary { background-color: #4f46e5; border-color: #4f46e5; font-weight: 600; text-transform: capitalize; padding: 0.4rem 0.8rem; }
.fc .fc-button-primary:hover { background-color: #4338ca; border-color: #4338ca; }
.fc .fc-button-primary:disabled { background-color: #a5b4fc; border-color: #a5b4fc; }
@media (max-width: 640px) {
    .fc-toolbar-title { font-size: 0.9rem !important; }
    .fc .fc-button { font-size: 0.7rem; padding: 0.3rem 0.5rem; }
}
</style>