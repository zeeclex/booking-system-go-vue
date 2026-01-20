# 🚪 ITIDoorz - Room & Lab Booking System

**ITIDoorz** is a modern web application designed to streamline classroom and laboratory booking management at **Institut Teknologi Indonesia (ITI)**. This system enables students to easily reserve rooms and assists administrators in managing schedules and facility usage reports efficiently.

---

## 🛠️ Tech Stack

This project is built using modern technologies for optimal performance and ease of development.

### Backend (Server-side)

- **Language:** [Go (Golang)](https://go.dev/)
- **Framework:** Gin Gonic
- **Database:** MySQL
- **ORM/Driver:** SQLX
- **Auth:** JWT & Bcrypt
- **Docs:** Swagger (Swaggo)

### Frontend (Client-side)

- **Framework:** [Vue.js 3](https://vuejs.org/) (Composition API)
- **Styling:** Tailwind CSS (Fully Responsive)
- **UI Components:** PrimeIcons
- **Charts:** Chart.js (Dashboard Statistics)

---

## ✨ Key Features

### 👤 User Panel

- **Real-time Availability:** Check room availability instantly.
- **Booking Form:** Submit booking requests for classrooms or laboratories.

### 🛡️ Admin Dashboard

- **Statistics:** Visual monthly room usage charts.
- **Inventory Management:** Add, edit, and delete room data.
- **User Management:** Full CRUD capabilities to manage system users (Admins, Lecturers, Students).
- **Approval System:** Feature to **Approve** or **Reject** booking requests.
- **Report Generator:** Export monthly statistical reports to **CSV** & **JSON** formats.

### 🎨 UI & UX Highlights

- **📱 Mobile-First Experience:** The Dashboard and User Panel layouts automatically adapt to Mobile, Tablet, and Desktop screens for a seamless experience.
- **🍔 Adaptive Sidebar:** The Admin sidebar now utilizes a Drawer/Overlay system on mobile with a Hamburger Menu, while remaining sticky on desktop.
- **📊 Smart Grids:** Statistic cards (StatCards) intelligently stack vertically (1 column) on mobile devices and expand horizontally (3-4 columns) on desktops.
- **↔️ Responsive Tables:** Data tables for Bookings and Rooms now support horizontal scrolling to ensure data integrity on smaller screens.
- **👆 Optimized Components:** Navbars, Modals, and Form Inputs have been fine-tuned with appropriate padding and sizing for optimal touch interaction.

---

## 📸 Screenshots

|                  Admin Dashboard                   |                  User Dashboard                  |
| :------------------------------------------------: | :----------------------------------------------: |
| ![Admin Dashboard](docs/image/admin-dashboard.png) | ![User Dashboard](docs/image/user-dashboard.png) |

### Optimized for all devices

This application has been optimized for various devices, ensuring a seamless user experience on mobile phones, tablets, and desktops.

---

## 🚀 Installation & Setup

### 1. Clone Repository

```bash
git clone https://github.com/MuhammadZaidan1/booking-system-go-vue.git
cd booking-system-go-vue
```

### 2. Database Setup

Import the `database.sql` file into your local MySQL database.

### 3. Backend Setup

```bash
cd backend
go mod tidy
go run main.go
```

> The Backend will run on port **8080**.

### 4. Frontend Setup

```bash
cd front-end
npm install
npm run dev
```

> The Frontend will run on port **5173**.

---

## 📖 API Documentation

Complete and interactive API documentation is available via Swagger UI.
Access at: **http://localhost:8080/swagger/index.html**
