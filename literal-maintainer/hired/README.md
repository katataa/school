
# Literal Maintainer 📋

An introductory cross-platform mobile application to help keep track of home maintenance jobs — built with Flutter.

---

## 📖 Overview

Literal Maintainer is a simple job-tracking app designed to help individuals like Maya, a home-maintenance engineer, manage their work more efficiently than with paper spreadsheets. Instead of using generic to-do apps, Maya needed something specific to her workflow — a way to track jobs, clients, payments, and job status. This app is tailored to her needs.

---

## 🚀 Features

- ✅ Create, edit, and delete jobs
- 🧾 Store important details like contact, address, costs, and job status
- 🖼 Upload up to 5 images per job (JPEG/PNG)
- 🔐 Encrypt sensitive client data (contact & address)
- 📋 Copy full job info to clipboard
- 🔍 Filter jobs by status (Todo, In Progress, Completed, Debtors, Cancelled)
- 💾 Persistent local storage (even after app restarts)
- ☀️ Dark/light mode with animated toggle

---

## 🛠 Setup & Installation

### ✅ Requirements

- Flutter SDK: [Install Flutter](https://docs.flutter.dev/get-started/install)
- Dart SDK (comes with Flutter)
- X Code for IOS testing (macOS only)
- Android Studio for androit testing: [Install Android studio](https://developer.android.com/studio)

### 💻 Clone & Run

```bash
git clone https://gitea.kood.tech/katriinsartakov/literal-maintainer.git
cd literal-maintainer
cd hired
flutter pub get
flutter run -d <emulator-id>
```

> **Note:** iOS testing is only possible on macOS using the iOS Simulator or a physical iPhone.

---

## 📱 Usage Guide

1. Launch the app
2. Tap the ➕ button to create a new job
3. Fill out required fields (Title, Status, Contact, Address, Cost, Amount Paid)
4. Add images (optional, max 5)
5. Save job — it will appear on the main list
6. Use the filter dropdown to view jobs by status
7. Tap a job to edit it, delete it, or copy its details

---

## 🔐 Security & Encryption

Client **contact** and **address** fields are encrypted using [flutter_secure_storage](https://pub.dev/packages/flutter_secure_storage). These values are not stored in the raw SQLite database.

---

## 🧠 State Management

The app uses the **Provider** package for clean, scalable state management.

- `JobProvider` handles:
  - fetching job data
  - updating job lists
  - notifying listeners after inserts, updates, or deletes

---

## 🧩 Custom Widgets

- `JobCard`: Reusable widget used to display job summary with title, status, and quick access to delete/view
- `AnimatedIconButton`: Animated toggle for dark/light mode

---

## 📁 App Structure

```
lib/
├── models/                # Job model
├── providers/            # State management
├── services/             # Database & secure storage
├── screens/              # UI screens
├── widgets/              # Custom widgets (e.g., job card)
└── main.dart             # App entry point
```

---

## 🌟 Bonus Features

- 🎨 Dark mode toggle with animation
- 📷 Secure in-app photo picker
- ✨ UI animations for job entries & transitions

---

## 🧪 Testing Checklist

- [x] App builds and runs successfully
- [x] Jobs persist after restart
- [x] Filtering works across all status types
- [x] Encrypted values are not shown in DB logs
- [x] Clipboard works with formatted output
- [x] No app crashes on invalid input or image load

---

## 🙋 FAQ

**Q: Can I run this on iPhone?**  
A: Yes, but only if you're using macOS with Xcode installed.

**Q: How do I reset the database?**  
A: Delete the `jobs.db` file from the emulator device or modify `initializeDatabase()` temporarily to delete the existing DB.

---