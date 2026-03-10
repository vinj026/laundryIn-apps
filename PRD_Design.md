# PRD Design: LaundryIn UI/UX Architecture (Frontend)

Dokumen ini mendeskripsikan secara spesifik pedoman desain (UI/UX) dan kerangka antarmuka pengguna untuk aplikasi **LaundryIn**. Desain berfokus pada estetika minimalis modern, teknikal, dan tegas, memanfaatkan palet warna monokrom (Hitam & Abu-abu gelap) dengan aksen tekstur ruang tata letak.

---

## 1. Design Philosophy & Aesthetics

### A. Konsep Visual (Material You 3 / Minimalist)
- **Tema Utama:** *Dark Mode Exclusive* (Google Material You 3 - MD3 Dark Theme).
- **Nuansa:** Simple, modern, ergonomis, dan fokus pada visibilitas data.
- **Elemen Dekoratif:** 
  - **Background Pattern:** Menggunakan pola **Kotak-kotak (Grid Pattern)** tipis atau pola **Dotted (Titik-titik)** pada *outer background* (ruang kosong di luar container utama).
  - **Sudut (Border Radius):** Melengkung lembut ala Material You. Tombol/Chip menggunakan `rounded-full` (pill-shaped), Card/Dialog menggunakan `rounded-2xl` atau `rounded-3xl`.
  - **Elevasi & Batas:** Menggunakan permainan *Surface Colors* (elevasi MD3) atau *soft border* yang sangat tipis alih-alih shadow yang kaku.

### B. Color Palette (Material You 3 Dark Theme)
Skema warna mengikuti sistem tonal Material Design 3, namun tetap diusahakan simple dengan dominasi hitam/abu-abu gelap, ditambah satu warna aksen *Primary* yang lembut.
- **Background Utama (Outer):** Hitam pekat `bg-[#000000]` atau abu sangat gelap `bg-[#0a0a0a]` dengan grid/dotted pattern.
- **Surface / Container Utama:** `bg-neutral-900` atau `bg-[#121212]` (Material Surface).
- **Surface Container (Card):** `bg-neutral-800` (Material Surface Container Highest/High).
- **On-Surface (Primary Text):** `text-neutral-100` (Putih terang).
- **On-Surface Variant (Muted):** `text-neutral-400` atau `text-neutral-500`.
- **Primary Accent:** Menggunakan warna aksen pastel khas MD3 (contoh: Soft Indigo / Blue `text-indigo-300`, `bg-indigo-300/10`) sebagai *Call to Action* utama.

### C. Container Layout Constraint (Tablet/Mobile-Max)
Sesuai referensi visual, aplikasi MENGGUNAKAN CONSTRAINT LEBAR TETAP agar terlihat seragam di semua device.
- **Layout Utama:** Seluruh isi aplikasi, baik login, dashboard owner, atau customer, dibungkus dalam container tipe Tablet (contoh: `max-w-2xl` atau max `768px`).
- **Posisi:** Container ini diposisikan tepat di tengah layar (`mx-auto`, `min-h-screen`). 
- Akibatnya, jika dibuka di Desktop/Layar Lebar, pengguna akan melihat aplikasi mengambang di tengah (dengan bentuk mirip tablet portrait) dengan background grid/dotted menghiasi sisa layar di kiri dan kanan.

### C. Typography Rules
- **Font Utama (Smooth & Modern):** `Inter` atau `Outfit` (Sans-serif modern). Memberikan kesan yang sangat halus, elegan, dan clean ala Material You. Sangat memanjakan mata untuk bacaan lama.
- **Data & Harga (Tabular):** Jika dibutuhkan presisi saat menyejajarkan nominal uang, dapat dipadu dengan `Roboto Mono`, tapi dominasi UI tetap dipegang oleh font Sans-Serif yang smooth.

---

## 2. Struktur Halaman & Komponen (Customer Flow - Mobile First)
Pelanggan mengakses situs melalui HP. UI dirancang *responsive* menyerupai struktur aplikasi HP (Max-width: `480px` centered).

### A. Auth Pages (Login & Register)
- **Background:** Pattern Grid gelap meliputi seluruh layer.
- **Komponen:**
  - `Brand Logo Typography`: Teks "LaundryIn" monospace putih tebal.
  - `Input Field Component`: Kotak transparan dengan border `border-neutral-800`. Fokus: `border-white`. Placeholder abu gelap.
  - `Button Primary`: `bg-white text-black font-bold outline-none`.

### B. Halaman Eksplorasi (Home / Outlets)
- **Layout:** Sticky Header atas, Body list, Bottom Navigation bar.
- **Komponen:**
  - `Greeting Header`: "Pilih Outlet".
  - `Search Box Component`: Kotak pencarian minimalis.
  - `Outlet Card Component`: Kotak `bg-neutral-900` dengan border garis. Memuat Nama Outlet (Geist Mono), Alamat (text muted), dan icon panah murni teks `->`.

### C. Halaman Outlet Detail & Katalog Pemesanan
- **Layout:** Header informasi Outlet, list harga/layanan scroll ke bawah, Sticky *Checkout Bar* di dasar layar.
- **Komponen:**
  - `Service List Item`: Baris layanan terbagi dua sisi (Kiri: Nama & Harga Desimal, Kanan: Input Qty / `+` `-`).
  - `Qty Input`: *Bordered input*, mendukung penulisan Float (Desimal mis: `2.5`).
  - `Floating Cart Bar`: Bar nempel di bawah berisi rangkuman "X item" dan tombol Checkout (PutihSolid).

### D. Halaman Order Tracking (Riwayat)
- **Komponen:**
  - `Order History Card`: Bordered box. Isi Teks memuat Tanggal, Daftar Servis (Snapshot Data), dan Total Bayar.
  - `Status Badge`: Chip komponen kecil di pojok. Misal: `[ PENDING ]` border abu-abu teks putih, `[ COMPLETED ]` dengan ketebalan teks berbeda.

---

## 3. Struktur Halaman & Komponen (Owner / Kasir Flow)
Meskipun diakses via Desktop oleh Kasir/Owner, seluruh antarmuka **TETAP** berada di dalam constraint ukuran Tablet di tengah layar (Max-width \~768px). Sidebar klasikal ditiadakan atau dipindah menjadi Bottom Nav / Hamburger Menu yang bersahabat untuk layar sempit.

### A. Layout Utama (Dashboard Shell)
- **Komponen:**
  - `Top App Bar`: Header Sticky ala MD3 di atas. Berisi Nama Halaman dan Menu/User Avatar.
  - `Main Content Area`: Area konten vertikal dengan *padding* standar.
  - `Bottom Navigation Bar`: Lengket di bawah untuk Owner Menu (Home, Outlets, Services, Orders). Berisi label dan icon dengan indikator aktif berupa pill-shape MD3.

### B. Halaman Analytics (Beranda Owner)
- **Komponen:**
  - `Filter Controller`: *Dropdown / segmented buttons* yang modern.
  - `Matrik Omzet Card`: Card `rounded-3xl` dengan *Surface Elevated* berisi angka uang ukuran besar (`text-5xl Geist Mono`).
  - `Order FSM Summary Cards`: Kotak berderet (grid 2x2) untuk status FSM: Pending, Process, Completed, Picked Up.
  - `Top Services List`: List terurut, dipisahkan garis tipis atau cukup ruang spasi yang lega.

### C. Halaman Outlet & Service Management
- **Komponen:**
  - `Data Table Component`: Tabel monokrom, baris head border `border-b-4`, baris isi tabel `border-b-1 text-neutral-300`.
  - `Action Column`: Tiga Teks Button berdampingan tanpa background: `[Edit]` | `[Hapus]`.
  - `Dialog / Modal Drawer`: Muncul dari tengah dengan overlay *backdrop blur* gelap. Berisi form pembuatan Outlet/Layanan.

### D. Halaman Order Management (Kasir Kasar)
- **Komponen:**
  - `Status Pipeline Board` (Opsional Kanban) atau `Expanded List Table`.
  - Tabel menampilkan siapa Customernya, Apa Snapshot Itemnya, dan *Dropdown Button* ganti state FSM.
  - `State Modifier Action`: Jika status = "Pending", muncul tombol border `[Set to Process]`. Jika diklik memanggil PATCH endpoint lalu merekayasa tabel list secara reaktif.

---

## Ringkasan Eksekusi UI/UX Developer:
Desain ini direalisasikan via **Tailwind CSS (`@nuxtjs/tailwindcss`)**. Konfigurasi utama:
1. Bungkus `app.vue` di dalam `<div class="max-w-2xl mx-auto min-h-screen bg-surface position-relative shadow-2xl">`. Terapkan outer container dengan pattern CSS background (titik/grid).
2. Terapkan paradigma ukuran **Material Design 3**, border radius besar (`rounded-2xl`, `rounded-full`), warna abu elegan (`neutral-900`, `neutral-800`), tanpa sudut tajam brutalis.
3. Tetap gunakan tipografi modern yang *smooth* seperti `Inter` dipadu dengan `Roboto Mono` untuk menjaga keterbacaan data finansial yang elegan.
