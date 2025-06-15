# 🌐 AyroBase – Privacy-First Web Hosting Platform

AyroBase is a self-hosted, all-in-one web hosting solution built with **Go**, focused on **simplicity**, **privacy**, and **performance**. Designed for creators, developers, and small businesses, it provides static website hosting, cookie-free analytics, integrated SQLite database management, and static file storage — all in a lightweight package that runs on your own infrastructure.

---

## 🔧 Features

- 🚀 **One-Click Hosting**  
  Upload your static website as a ZIP file and instantly host it.

- 📊 **Privacy-First Analytics**  
  Track page visits and referrers without using cookies or third-party scripts — fully GDPR/HIPAA compliant.

- 🗂️ **Object Manager**  
  Upload, organize, and serve images, PDFs, and other static files.

- 📥 **Form Integration with SQLite**  
  Collect data from static site forms and store it directly in an embedded database.

- 🛡️ **On-Premise & Dual-Server Architecture**  
  Split between a secure local admin server and a global public server.

---

## 🖼️ Use Cases

- Hosting landing pages, portfolios, or microsites  
- Collecting contact forms or surveys without writing backend code  
- Serving documents, videos, or media files  
- Viewing site traffic and referral sources without third-party tools

---

## 🏗️ Architecture

```bash

\[ Admin (Local Server) ] <-- SSH Only
├── Form DB (SQLite)
├── Analytics DB (SQLite)
├── Admin Dashboard
└── Config Controls

\[ User-Facing (Global Server) ]
├── Static Site Hosting (/Hosting)
├── File Server (/Content)
├── Form Handler (/api/forms/insert)
└── Analytics Middleware

```

---

## 🚀 Quick Start

> Requirements: Go 1.20+, Linux/Unix system, SQLite3

### 1. Clone the Repo

```bash
git clone https://github.com/Adi-ayro/ayrobase.git
cd ayrobase
````

### 2. Run the Servers

#### Local (Admin) Server

```bash
go run localserver.go
# Accessible at http://127.0.0.1:3000
```

#### Global (User) Server

```bash
go run globalserver.go
# Accessible at http://0.0.0.0:4000
```

---

## 📁 Endpoints

### Hosting

* `POST /host` – Replace website ZIP
* `POST /newhost` – Deploy website

### Content Manager

* `POST /content/upload` – Upload file
* `GET /content/get` – List uploaded files
* `DELETE /content/delete/:filename` – Delete a file

### Forms

* `POST /forms/create` – Create form table
* `POST /api/forms/insert/:formname` – Submit form data

### Admin Panels

* `/admin/hosting`
* `/admin/content`
* `/admin/analytics`
* `/admin/form/:tablename?`

---

## 🔒 Privacy & Compliance

* **No Cookies or Trackers**
* **Complies with**: GDPR, HIPAA, India DPDP Act
* **Data is stored locally** — nothing leaves your server

---

## 📦 Tech Stack

* Backend: [Go](https://golang.org/)
* Database: [SQLite](https://www.sqlite.org/index.html)
* Frontend: HTML, CSS

---

## 📄 License

MIT License. See `LICENSE` for more information.

---

Developed by [Aditya Bharadwaj](https://github.com/Adi-ayro) as part of a major project focused on self-hosted, privacy-conscious infrastructure tools.
