<script setup>
defineProps({
    show: Boolean,
    title: String,
    icon: { type: String, default: 'pi-info-circle' },
    widthClass: { type: String, default: 'max-w-md' } 
});
defineEmits(['close']);
</script>

<template>
    <Teleport to="body">
        <Transition name="modal">
            <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4 sm:p-6">
                <div class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm transition-opacity" @click="$emit('close')"></div>
                <div :class="['bg-white w-full rounded-2xl shadow-2xl transform transition-all relative flex flex-col max-h-[90vh]', widthClass]">
                    <div class="px-4 py-3 md:px-6 md:py-4 border-b border-slate-100 flex items-center justify-between bg-white rounded-t-2xl shrink-0">
                        <div class="flex items-center gap-3">
                            <div class="w-8 h-8 md:w-10 md:h-10 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center shrink-0 transition-all">
                                <i :class="['pi text-sm md:text-lg', icon]"></i>
                            </div>
                            <h3 class="font-black text-slate-800 uppercase tracking-tight text-sm md:text-base transition-all">
                                {{ title }}
                            </h3>
                        </div>
                        <button 
                            @click="$emit('close')" 
                            class="w-8 h-8 flex items-center justify-center rounded-lg text-slate-400 hover:text-slate-600 hover:bg-slate-50 transition-all"
                            title="Close"
                        >
                            <i class="pi pi-times"></i>
                        </button>
                    </div>
                    <div class="p-4 md:p-6 overflow-y-auto custom-scrollbar">
                        <slot></slot>
                    </div>
                    <div v-if="$slots.footer" class="px-4 py-3 md:px-6 md:py-4 bg-slate-50 border-t border-slate-100 flex flex-col-reverse sm:flex-row sm:justify-end gap-2 md:gap-3 rounded-b-2xl shrink-0 transition-all">
                        <slot name="footer"></slot>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.3s ease;
}
.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}
.modal-enter-active .transform,
.modal-leave-active .transform {
    transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1); 
}
.modal-enter-from .transform,
.modal-leave-to .transform {
    opacity: 0;
    transform: scale(0.95) translateY(10px);
}
.custom-scrollbar::-webkit-scrollbar {
    width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
    background-color: #e2e8f0;
    border-radius: 20px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background-color: #cbd5e1;
}
</style>