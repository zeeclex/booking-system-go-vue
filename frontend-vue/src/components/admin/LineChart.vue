<script setup>
import { ref, computed, watch, onMounted } from 'vue';
import { Line } from 'vue-chartjs';
import { 
    Chart as ChartJS, 
    Title, Tooltip, Legend, LineElement,    
    PointElement, CategoryScale, LinearScale, Filler          
} from 'chart.js';

ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale, Filler);

// PROPS: Receive raw booking data from parent
const props = defineProps({
    bookings: {
        type: Array,
        default: () => []
    },
    loading: {
        type: Boolean,
        default: false
    }
});

// STATE FILTERS
const currentYear = new Date().getFullYear();
const selectedMonth = ref(new Date().getMonth());
const selectedYear = ref(currentYear);
const months = [
    { value: 0, label: 'January' }, { value: 1, label: 'February' }, { value: 2, label: 'March' },
    { value: 3, label: 'April' }, { value: 4, label: 'May' }, { value: 5, label: 'June' },
    { value: 6, label: 'July' }, { value: 7, label: 'August' }, { value: 8, label: 'September' },
    { value: 9, label: 'October' }, { value: 10, label: 'November' }, { value: 11, label: 'December' }
];
const years = [currentYear - 1, currentYear, currentYear + 1];

// CHART REACTIVE DATA
const chartLabels = ref([]);
const chartValues = ref([]);

// --- LOGIC: PROCESS DATA ---
const updateChartData = () => {
    // 1. Calculate number of days in the selected month
    const daysInMonth = new Date(selectedYear.value, selectedMonth.value + 1, 0).getDate();
    
    // 2. Prepare empty labels & data arrays
    const labels = [];
    const dataCounts = new Array(daysInMonth).fill(0);
    for (let i = 1; i <= daysInMonth; i++) {
        labels.push(`${i}`);
    }

    // 3. Filter Data from Props
    props.bookings.forEach(booking => {
        const dateStr = booking.start_time || booking.date; 
        const bDate = new Date(dateStr);
        if (bDate.getMonth() === selectedMonth.value && 
            bDate.getFullYear() === selectedYear.value &&
            booking.status !== 'rejected') { 
            
            const dayIndex = bDate.getDate() - 1; 
            if (dayIndex >= 0 && dayIndex < daysInMonth) {
                dataCounts[dayIndex]++;
            }
        }
    });
    chartLabels.value = labels;
    chartValues.value = dataCounts;
};

// --- WATCHERS ---
watch([selectedMonth, selectedYear, () => props.bookings], () => {
    updateChartData();
}, { deep: true });

// Initial Load
onMounted(updateChartData);

// --- CHART CONFIG ---
const chartData = computed(() => ({
    labels: chartLabels.value,
    datasets: [{
        label: 'Daily Usage',
        borderColor: '#6366f1',
        borderWidth: 2,
        backgroundColor: (context) => {
            const ctx = context.chart.ctx;
            const gradient = ctx.createLinearGradient(0, 0, 0, 300);
            gradient.addColorStop(0, 'rgba(99, 102, 241, 0.4)');
            gradient.addColorStop(1, 'rgba(99, 102, 241, 0.0)');
            return gradient;
        },
        data: chartValues.value,
        fill: true,
        tension: 0.3,
        pointBackgroundColor: '#fff',
        pointBorderColor: '#6366f1',
        pointBorderWidth: 2,
        pointRadius: 4,
        pointHoverRadius: 6
    }]
}));

const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: { legend: { display: false } },
    scales: {
        y: { 
            beginAtZero: true, 
            suggestedMax: 10,
            grid: { display: true, color: '#f3f4f6', borderDash: [5, 5] }, 
            ticks: { stepSize: 1, font: { size: 10 } } 
        },
        x: { 
            grid: { display: false }, 
            ticks: { font: { size: 10 }, maxRotation: 0 } 
        }
    },
    interaction: { intersect: false, mode: 'index' },
};
</script>

<template>
    <div class="bg-white rounded-2xl md:rounded-3xl border border-slate-200 shadow-xl shadow-slate-100/50 overflow-hidden relative z-0 transition-all">
        <div class="absolute top-0 left-0 w-full h-1.5 bg-linear-to-r from-indigo-500 via-purple-500 to-pink-500"></div>
        <div class="p-5 md:p-8 pb-0 flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
            <div class="w-full md:w-auto">
                <h3 class="text-lg md:text-xl font-black text-slate-800 tracking-tight flex items-center gap-2">
                    <i class="pi pi-chart-line text-indigo-500"></i> Usage Statistics
                </h3>
                <p class="text-xs md:text-sm text-slate-400 font-medium mt-1">Daily reservation breakdown</p>
            </div>
            <div class="flex items-center gap-3 bg-slate-50 p-1.5 rounded-xl border border-slate-200 w-full md:w-auto">
                <select 
                    v-model="selectedMonth" 
                    class="bg-white border border-slate-200 text-slate-600 text-xs font-bold rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 block w-full p-2 outline-none cursor-pointer hover:bg-slate-50 transition-colors"
                >
                    <option v-for="m in months" :key="m.value" :value="m.value">{{ m.label }}</option>
                </select>
                <select 
                    v-model="selectedYear" 
                    class="bg-white border border-slate-200 text-slate-600 text-xs font-bold rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 block w-24 p-2 outline-none cursor-pointer hover:bg-slate-50 transition-colors"
                >
                    <option v-for="y in years" :key="y" :value="y">{{ y }}</option>
                </select>
            </div>
        </div>
        <div class="h-64 md:h-80 w-full px-4 md:px-8 pb-4 md:pb-6 mt-6">
            <div v-if="props.loading" class="h-full flex items-center justify-center">
                <i class="pi pi-spin pi-spinner text-3xl md:text-4xl text-indigo-500 opacity-50"></i>
            </div>
            <Line v-else :data="chartData" :options="chartOptions" />
        </div>
    </div>
</template>