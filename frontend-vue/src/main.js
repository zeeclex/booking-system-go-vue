import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura';
import ToastService from 'primevue/toastservice';
import './assets/css/main.css';
import 'primeicons/primeicons.css';

const app = createApp(App);

app.use(router);
app.use(PrimeVue, { 
    theme: { preset: Aura, options: { darkModeSelector: 'none' } } 
});
app.use(ToastService);

import Button from 'primevue/button';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Tag from 'primevue/tag';
import InputText from 'primevue/inputtext';
import Card from 'primevue/card';
import Toast from 'primevue/toast';

app.component('Button', Button);
app.component('DataTable', DataTable);
app.component('Column', Column);
app.component('Tag', Tag);
app.component('InputText', InputText);
app.component('Card', Card);
app.component('Toast', Toast);

// Daftarkan Directive PrimeVue
import Tooltip from 'primevue/tooltip';
app.directive('tooltip', Tooltip);
app.mount('#app');