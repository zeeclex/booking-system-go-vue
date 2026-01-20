<script setup>
import { ref, computed, onMounted } from 'vue';
import { ReportService } from '@/service/ReportService';
import { BookingService } from '@/service/BookingService';
import { RoomService } from '@/service/RoomService'; 
import AdminSidebar from '../../components/admin/Sidebar.vue';
import BaseModal from '../../components/ui/BaseModal.vue';

// --- STATE ---
const selectedMonth = ref(new Date().getMonth() + 1);
const selectedYear = ref(new Date().getFullYear());
const loadingPreview = ref(false);
const loadingGenerate = ref(false);
const previewData = ref([]);
const rooms = ref([]);
const isDataLoaded = ref(false);

// Mobile Sidebar State
const isSidebarOpen = ref(false);

// Modal State
const showModal = ref(false);
const modalTitle = ref('');
const modalMessage = ref('');
const downloadLinks = ref({ csv: '', json: '' });

// Static Data
const months = [
    { value: 1, label: 'January' }, { value: 2, label: 'February' },
    { value: 3, label: 'March' }, { value: 4, label: 'April' },
    { value: 5, label: 'May' }, { value: 6, label: 'June' },
    { value: 7, label: 'July' }, { value: 8, label: 'August' },
    { value: 9, label: 'September' }, { value: 10, label: 'October' },
    { value: 11, label: 'November' }, { value: 12, label: 'December' }
];

const years = Array.from({ length: 5 }, (_, i) => new Date().getFullYear() - i);

// --- HELPER ---
const formatDateTime = (dateString) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString('en-GB', {
        day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit'
    });
};

const calculateDuration = (start, end) => {
    const s = new Date(start);
    const e = new Date(end);
    const diffMs = e - s;
    const diffHrs = diffMs / (1000 * 60 * 60);
    return diffHrs.toFixed(1) + ' Hrs';
};

const getRoomName = (id) => {
    const room = rooms.value.find(r => r.id === id);
    return room ? room.name : `Room #${id}`;
};

// --- LOAD INITIAL DATA ---
onMounted(async () => {
    try {
        const res = await RoomService.getRooms();
        rooms.value = res.data;
    } catch (err) {
        console.error("Failed to load rooms for report:", err);
    }
});

// --- PREVIEW DATA ---
const handlePreview = async () => {
    loadingPreview.value = true;
    previewData.value = [];
    isDataLoaded.value = false;
    try {
        const res = await BookingService.getBookings(); 
        const filtered = res.data.filter(b => {
            const date = new Date(b.start_time);
            return date.getMonth() + 1 === selectedMonth.value && 
                    date.getFullYear() === selectedYear.value &&
                    b.status === 'approved';
        });
        previewData.value = filtered;
        isDataLoaded.value = true;
    } catch (err) {
        console.error("Preview failed:", err);
        alert("Failed to load preview data.");
    } finally {
        loadingPreview.value = false;
    }
};

// --- GENERATE REPORT ---
const handleGenerate = async () => {
    loadingGenerate.value = true;
    try {
        const res = await ReportService.generateReport(selectedMonth.value, selectedYear.value);
        downloadLinks.value = {
            csv: res.data.download_url_csv,
            json: res.data.download_url_json
        };
        modalTitle.value = "Report Generated!";
        modalMessage.value = "Your report files are ready to download.";
        showModal.value = true;
    } catch (err) {
        console.error("Generate failed:", err);
        const msg = err.response?.data?.error || "Failed to generate report files.";
        alert(msg);
    } finally {
        loadingGenerate.value = false;
    }
};

const downloadFile = (url) => {
    if (url) {
        const fullUrl = url.startsWith('http') ? url : `http://localhost:8080${url}`;
        window.open(fullUrl, '_blank');
    }
};
</script>

<template>
    <div class="flex min-h-screen bg-slate-50 font-sans">
        
        <AdminSidebar :isOpen="isSidebarOpen" @close-sidebar="isSidebarOpen = false" />

        <main class="flex-1 w-full max-w-7xl mx-auto p-4 md:p-8 space-y-6 flex flex-col lg:h-screen lg:overflow-hidden">
            
            <div class="flex-none pb-5 border-b border-slate-200 flex items-center gap-4">
                <button @click="isSidebarOpen = true" class="lg:hidden p-2 -ml-2 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors">
                    <i class="pi pi-bars text-xl"></i>
                </button>
                
                <div>
                    <h2 class="text-xl md:text-2xl font-black text-slate-800 uppercase tracking-tighter italic">Report Center</h2>
                    <p class="text-xs md:text-sm text-slate-500 font-medium italic mt-1">Monthly Usage Recap & Export</p>
                </div>
            </div>

            <div class="flex-1 grid grid-cols-1 lg:grid-cols-4 gap-6 lg:overflow-hidden">
                
                <div class="lg:col-span-1 flex flex-col gap-6 lg:overflow-y-auto custom-scrollbar pr-1">
                    <div class="bg-white p-5 md:p-6 rounded-2xl border border-slate-200 shadow-lg shadow-slate-100/50">
                        <h3 class="text-xs font-black text-slate-400 uppercase tracking-widest mb-4">Parameters</h3>
                        
                        <div class="space-y-4">
                            <div class="space-y-1.5">
                                <label class="text-[10px] font-bold text-slate-700 uppercase">Month</label>
                                <div class="relative">
                                    <select v-model="selectedMonth" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-sm appearance-none cursor-pointer hover:bg-slate-100 transition-colors">
                                        <option v-for="m in months" :key="m.value" :value="m.value">{{ m.label }}</option>
                                    </select>
                                    <i class="pi pi-chevron-down absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 text-xs pointer-events-none"></i>
                                </div>
                            </div>

                            <div class="space-y-1.5">
                                <label class="text-[10px] font-bold text-slate-700 uppercase">Year</label>
                                <div class="relative">
                                    <select v-model="selectedYear" class="w-full p-3 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:ring-2 focus:ring-indigo-500 font-bold text-sm appearance-none cursor-pointer hover:bg-slate-100 transition-colors">
                                        <option v-for="y in years" :key="y" :value="y">{{ y }}</option>
                                    </select>
                                    <i class="pi pi-calendar absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 text-xs pointer-events-none"></i>
                                </div>
                            </div>

                            <Button 
                                label="Preview Data" 
                                icon="pi pi-search" 
                                class="w-full mt-2 bg-slate-800 hover:bg-slate-900 text-white font-bold py-3 rounded-xl transition-all shadow-lg shadow-slate-200"
                                :loading="loadingPreview"
                                @click="handlePreview"
                            />
                        </div>
                    </div>

                    <div class="bg-indigo-50 p-5 rounded-2xl border border-indigo-100">
                        <div class="flex gap-3">
                            <i class="pi pi-info-circle text-indigo-500 mt-0.5"></i>
                            <p class="text-xs text-indigo-800 font-medium leading-relaxed">
                                Preview shows only <strong>Approved</strong> bookings for the selected period.
                            </p>
                        </div>
                    </div>
                </div>

                <div class="lg:col-span-3 flex flex-col bg-white rounded-2xl border border-slate-200 shadow-lg shadow-slate-100/50 overflow-hidden min-h-125 lg:min-h-0">
                    
                    <div class="p-4 md:p-5 border-b border-slate-100 flex flex-col sm:flex-row justify-between items-center gap-4 bg-white z-10">
                        <div class="text-center sm:text-left">
                            <h3 class="font-black text-slate-800 text-base flex items-center justify-center sm:justify-start gap-2">
                                <i class="pi pi-table text-indigo-500"></i> Data Preview
                            </h3>
                            <p class="text-xs text-slate-400 font-bold mt-0.5" v-if="isDataLoaded">
                                Found {{ previewData.length }} records for {{ months[selectedMonth-1].label }} {{ selectedYear }}
                            </p>
                        </div>
                        
                        <Button 
                            v-if="isDataLoaded && previewData.length > 0"
                            label="Generate Report" 
                            icon="pi pi-download" 
                            class="w-full sm:w-auto bg-indigo-600 hover:bg-indigo-700 text-white font-bold px-6 py-2 rounded-lg text-sm shadow-md shadow-indigo-200 animate-pulse"
                            :loading="loadingGenerate"
                            @click="handleGenerate"
                        />
                    </div>

                    <div class="flex-1 overflow-y-auto p-0 relative custom-scrollbar">
                        
                        <div v-if="!isDataLoaded" class="h-full flex flex-col items-center justify-center text-slate-300 space-y-3 min-h-75">
                            <i class="pi pi-search text-4xl opacity-50"></i>
                            <p class="text-sm font-bold uppercase tracking-widest text-center px-4">Select parameters & click Preview Data</p>
                        </div>

                        <div v-else-if="previewData.length === 0" class="h-full flex flex-col items-center justify-center text-slate-400 space-y-3 min-h-75">
                            <i class="pi pi-folder-open text-4xl opacity-50"></i>
                            <p class="text-sm font-bold uppercase tracking-widest text-center px-4">No approved data found for this period.</p>
                        </div>

                        <DataTable 
                            v-else
                            :value="previewData" 
                            class="p-datatable-sm w-full"
                            stripedRows
                            responsiveLayout="scroll"
                            :pt="{
                                header: { class: 'bg-slate-50 text-slate-700 font-bold uppercase text-xs tracking-wider sticky top-0 z-10 shadow-sm' },
                                bodyRow: { class: 'text-sm text-slate-600' }
                            }"
                        >
                            <Column header="Date" class="min-w-30 pl-4 md:pl-6">
                                <template #body="{ data }">
                                    <span class="font-bold text-slate-800 whitespace-nowrap">{{ new Date(data.start_time).toLocaleDateString() }}</span>
                                </template>
                            </Column>
                            
                            <Column header="Time Slot" class="min-w-37.5">
                                <template #body="{ data }">
                                    <div class="flex flex-col">
                                        <span class="font-bold text-xs text-indigo-600 whitespace-nowrap">
                                            {{ formatDateTime(data.start_time).split(',')[1] }} - {{ formatDateTime(data.end_time).split(',')[1] }}
                                        </span>
                                        <span class="text-[10px] text-slate-400 font-bold bg-slate-100 px-1 rounded w-fit mt-0.5 whitespace-nowrap">
                                            {{ calculateDuration(data.start_time, data.end_time) }}
                                        </span>
                                    </div>
                                </template>
                            </Column>
                            
                            <Column header="Room Name" class="min-w-35">
                                <template #body="{ data }">
                                    <span class="bg-indigo-50 text-indigo-700 px-2 py-1 rounded font-black text-xs whitespace-nowrap">{{ getRoomName(data.room_id) }}</span>
                                </template>
                            </Column>
                            
                            <Column field="purpose" header="Purpose" class="min-w-50">
                                <template #body="{ data }">
                                    <span class="truncate block max-w-50" :title="data.purpose">{{ data.purpose }}</span>
                                </template>
                            </Column>
                        </DataTable>
                    </div>
                </div>
            </div>
        </main>

        <BaseModal :show="showModal" :title="modalTitle" icon="pi-check-circle" @close="showModal = false">
            <div class="text-center py-4 space-y-6">
                <p class="text-slate-600 text-sm font-medium">{{ modalMessage }}</p>
                <div class="flex flex-col gap-3">
                    <button @click="downloadFile(downloadLinks.csv)" class="flex items-center justify-between px-4 py-3 bg-emerald-50 text-emerald-700 rounded-xl border border-emerald-200 hover:bg-emerald-100 transition-all group cursor-pointer active:scale-95">
                        <div class="flex items-center gap-3">
                            <i class="pi pi-file-excel text-xl"></i>
                            <div class="text-left">
                                <p class="text-xs font-black uppercase">CSV Summary</p>
                                <p class="text-[10px]">Total hours calculation</p>
                            </div>
                        </div>
                        <i class="pi pi-download group-hover:animate-bounce"></i>
                    </button>
                    <button @click="downloadFile(downloadLinks.json)" class="flex items-center justify-between px-4 py-3 bg-blue-50 text-blue-700 rounded-xl border border-blue-200 hover:bg-blue-100 transition-all group cursor-pointer active:scale-95">
                        <div class="flex items-center gap-3">
                            <i class="pi pi-code text-xl"></i>
                            <div class="text-left">
                                <p class="text-xs font-black uppercase">JSON Data</p>
                                <p class="text-[10px]">Full raw dataset</p>
                            </div>
                        </div>
                        <i class="pi pi-download group-hover:animate-bounce"></i>
                    </button>
                </div>
            </div>
            <template #footer>
                <Button label="Close" class="w-full font-bold rounded-xl" severity="secondary" @click="showModal = false" />
            </template>
        </BaseModal>
    </div>
</template>

<style scoped>
/* Custom Scrollbar */
.custom-scrollbar::-webkit-scrollbar {
    width: 6px;
    height: 6px; /* For horizontal scroll if needed */
}
.custom-scrollbar::-webkit-scrollbar-track {
    background: #f1f5f9;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
    background-color: #cbd5e1;
    border-radius: 20px;
}
</style>