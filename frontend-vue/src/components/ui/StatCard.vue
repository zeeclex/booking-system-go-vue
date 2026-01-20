<script setup>
import { computed } from 'vue';

const props = defineProps({
  title: String,
  value: [Number, String],
  icon: String,
  variant: {
    type: String,
    default: 'indigo' 
  }
});

// Class mapping for Tailwind colors
const colorStyles = computed(() => {
  const styles = {
    indigo: {
      bar: 'bg-indigo-500',
      text: 'text-indigo-600',
      bgIcon: 'bg-indigo-50 group-hover:bg-indigo-100',
      border: 'border-indigo-200'
    },
    green: {
      bar: 'bg-green-500',
      text: 'text-green-600',
      bgIcon: 'bg-green-50 group-hover:bg-green-100',
      border: 'border-green-200'
    },
    red: {
      bar: 'bg-red-500',
      text: 'text-red-600',
      bgIcon: 'bg-red-50 group-hover:bg-red-100',
      border: 'border-red-200'
    },
    amber: {
      bar: 'bg-amber-500',
      text: 'text-amber-600',
      bgIcon: 'bg-amber-50 group-hover:bg-amber-100',
      border: 'border-amber-200'
    },
    blue: {
      bar: 'bg-blue-500',
      text: 'text-blue-600',
      bgIcon: 'bg-blue-50 group-hover:bg-blue-100',
      border: 'border-blue-200'
    },
    slate: {
      bar: 'bg-slate-500',
      text: 'text-slate-600',
      bgIcon: 'bg-slate-100 group-hover:bg-slate-200',
      border: 'border-slate-200'
    },
    emerald: {
      bar: 'bg-emerald-500',
      text: 'text-emerald-600',
      bgIcon: 'bg-emerald-50 group-hover:bg-emerald-100',
      border: 'border-emerald-200'
    }
  };
  return styles[props.variant] || styles.indigo;
});
</script>

<template>
  <div class="bg-white p-4 md:p-5 rounded-xl border border-slate-200 shadow-sm hover:shadow-md transition-all duration-300 relative overflow-hidden group">
    <div :class="['absolute left-0 top-0 bottom-0 w-1 md:w-1.5', colorStyles.bar]"></div>
    <div class="flex justify-between items-center pl-2 md:pl-3"> 
      <div>
        <p class="text-[10px] md:text-xs font-black text-slate-400 uppercase tracking-widest mb-0.5 md:mb-1">{{ title }}</p>
        <h3 class="text-xl md:text-2xl font-black text-slate-800 tracking-tight">{{ value }}</h3>
      </div>
      <div :class="['w-10 h-10 md:w-12 md:h-12 rounded-xl flex items-center justify-center transition-all duration-300 group-hover:scale-110', colorStyles.bgIcon, colorStyles.text]">
        <i :class="['pi text-lg md:text-xl', icon]"></i>
      </div>
    </div>
  </div>
</template>