import { createRouter, createWebHistory } from 'vue-router';

import LoginPage from '../views/auth/Login.vue';
import AdminDashboard from '../views/admin/Dashboard.vue';
import AdminRoomPage from '../views/admin/RoomPage.vue';
import AdminBookingPage from '../views/admin/BookingPage.vue';
import AdminReportPage from '../views/admin/ReportPage.vue';
import AdminUserPage from '../views/admin/UserPage.vue';
import UserDashboard from '../views/user/UserDashboard.vue';

const routes = [
  // --- LOGIN ---
  { 
    path: '/login', 
    name: 'Login', 
    component: LoginPage,
    meta: { guestOnly: true } 
  },

  // --- USER ROUTES (Role: 'user') ---
  { 
    path: '/user/dashboard', 
    name: 'UserDashboard', 
    component: UserDashboard, 
    meta: { requiresAuth: true, allowedRoles: ['user'] } 
  },
  { 
    path: '/user', 
    redirect: '/user/dashboard' 
  },

  // --- ADMIN ROUTES (Role: 'admin') ---
  { 
    path: '/admin', 
    redirect: '/admin/dashboard'
  },
  { 
    path: '/admin/dashboard', 
    name: 'AdminDashboard', 
    component: AdminDashboard, 
    meta: { requiresAuth: true, allowedRoles: ['admin'] } 
  },
  { 
    path: '/admin/rooms', 
    name: 'AdminRooms', 
    component: AdminRoomPage, 
    meta: { requiresAuth: true, allowedRoles: ['admin'] } 
  },
  { 
    path: '/admin/bookings', 
    name: 'AdminBookings', 
    component: AdminBookingPage, 
    meta: { requiresAuth: true, allowedRoles: ['admin'] } 
  },
  { 
    path: '/admin/reports', 
    name: 'AdminReports', 
    component: AdminReportPage, 
    meta: { requiresAuth: true, allowedRoles: ['admin'] } 
  },
  { 
    path: '/admin/users',
    name: 'AdminUsers', 
    component: AdminUserPage, 
    meta: { requiresAuth: true, allowedRoles: ['admin'] } 
  },
  { path: '/', redirect: '/login' },
  { path: '/:pathMatch(.*)*', redirect: '/login' }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// --- NAVIGATION GUARD ---
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  const role = localStorage.getItem('role');

  // 1. Guest Only Guard (Login Page)
  // If user is logged in, redirect them to their respective dashboard
  if (to.meta.guestOnly && token) {
    if (role === 'admin') return next('/admin/dashboard');
    return next('/user/dashboard');
  }

  // 2. Auth Guard
  // If route requires auth and no token is present, send to login
  if (to.meta.requiresAuth && !token) {
    return next('/login');
  }

  // 3. Role Based Access Control (RBAC)
  if (token && to.meta.allowedRoles) {
    if (!to.meta.allowedRoles.includes(role)) {
      if (role === 'admin') {
        return next('/admin/dashboard');
      } else {
        return next('/user/dashboard');
      }
    }
  }
  next();
});

export default router;