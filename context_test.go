package belajar_golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// TestContext adalah fungsi pengujian yang mendemonstrasikan dua jenis context dasar di Go:
// context.Background() dan context.TODO().
// Best practice: Selalu sertakan dokumentasi fungsi yang menjelaskan tujuan dan penggunaannya
func TestContext(t *testing.T) {
	// context.Background() digunakan sebagai root context dan merupakan pilihan default
	// untuk aplikasi level tinggi atau ketika jelas bahwa context ini adalah yang tertinggi
	// dalam hierarki.
	// Best practice: Gunakan Background() untuk inisialisasi context di level tertinggi
	background := context.Background()
	fmt.Println(background)

	// context.TODO() digunakan sebagai placeholder ketika tidak yakin context apa yang
	// seharusnya digunakan atau ketika refactoring kode yang belum menggunakan context.
	// Best practice: Gunakan TODO() hanya sementara selama development, hindari di production
	todo := context.TODO()
	fmt.Println(todo)
}

// TestContextWithValue mendemonstrasikan penggunaan context dengan nilai (key-value)
// dan menunjukkan hierarki pewarisan nilai antar context
func TestContextWithValue(t *testing.T) {
	// Membuat context induk (root context) yang akan menjadi dasar untuk context lainnya
	// Best practice: Selalu mulai dengan context.Background() untuk root context
	contextA := context.Background()

	// Membuat context turunan level pertama dari contextA
	// Best practice: Gunakan tipe yang spesifik untuk key, hindari string
	// Best practice: Dokumentasikan struktur key-value yang digunakan
	contextB := context.WithValue(contextA, "b", "B")  // contextB mewarisi contextA
	contextC := context.WithValue(contextA, "c", "C")  // contextC mewarisi contextA

	// Membuat context turunan level kedua dari contextB
	// Mendemonstrasikan bahwa context bisa memiliki multiple children
	contextD := context.WithValue(contextB, "d", "D")  // contextD mewarisi contextB dan contextA
	contextE := context.WithValue(contextB, "e", "E")  // contextE mewarisi contextB dan contextA

	// Membuat context turunan berjenjang dari contextC
	// Mendemonstrasikan rantai pewarisan yang lebih dalam
	contextF := context.WithValue(contextC, "f", "F")  // contextF mewarisi contextC dan contextA
	contextG := context.WithValue(contextF, "g", "G")  // contextG mewarisi contextF, contextC, dan contextA

	// Mencetak representasi string dari setiap context
	// Berguna untuk debugging dan memahami struktur context
	fmt.Println(contextA)  // Menampilkan context induk
	fmt.Println(contextB)  // Menampilkan context dengan nilai "b"
	fmt.Println(contextC)  // Menampilkan context dengan nilai "c"
	fmt.Println(contextD)  // Menampilkan context dengan nilai "b" dan "d"
	fmt.Println(contextE)  // Menampilkan context dengan nilai "b" dan "e"
	fmt.Println(contextF)  // Menampilkan context dengan nilai "c" dan "f"
	fmt.Println(contextG)  // Menampilkan context dengan nilai "c", "f", dan "g"

	// Mendemonstrasikan cara mengakses nilai dalam context
	// Best practice: Selalu periksa apakah nilai yang diambil sesuai dengan tipe yang diharapkan
	fmt.Println(contextF.Value("f"))  // Akan mengembalikan "F" karena ada di contextF
	fmt.Println(contextF.Value("c"))  // Akan mengembalikan "C" karena diwarisi dari contextC
	fmt.Println(contextF.Value("b"))  // Akan mengembalikan nil karena "b" tidak ada di rantai contextF

	// Mendemonstrasikan bahwa context induk tidak dapat mengakses nilai dari context turunan
	fmt.Println(contextA.Value("b"))  // Akan mengembalikan nil karena contextA tidak memiliki nilai
}

// CreateCounter membuat dan mengembalikan channel yang menghasilkan angka berurutan.
// Parameter ctx digunakan untuk mengontrol lifecycle dari goroutine yang dijalankan.
// Channel yang dikembalikan akan ditutup ketika context dibatalkan atau terjadi error.
func CreateCounter(ctx context.Context) chan int {
	// Membuat channel unbuffered untuk mengirim nilai counter
	// Best practice: Gunakan unbuffered channel untuk sinkronisasi yang lebih baik
	destination := make(chan int)

	// Menjalankan goroutine untuk menghasilkan nilai counter secara asynchronous
	// Best practice: Selalu gunakan goroutine terpisah untuk operasi yang blocking
	go func() {
		// Memastikan channel selalu ditutup ketika fungsi selesai
		// Best practice: Gunakan defer untuk mencegah resource leak
		defer close(destination)

		// Inisialisasi counter dimulai dari 1
		counter := 1

		// Loop tak terbatas untuk menghasilkan nilai counter
		// Best practice: Gunakan select untuk handling pembatalan context
		for {
			select {
			case <-ctx.Done():
				// Menghentikan goroutine ketika context dibatalkan
				// Best practice: Selalu handle pembatalan context
				return
			default:
				// Mengirim nilai counter ke channel
				// Operasi ini akan blocking jika tidak ada yang menerima
				destination <- counter
				counter++

				// Simulasi proses yang memakan waktu
				// Note: Dalam kode produksi, hindari time.Sleep
				// Best practice: Gunakan mekanisme rate limiting yang proper
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Mengembalikan channel yang akan digunakan oleh consumer
	// Best practice: Channel producer hanya bertanggung jawab untuk menutup channel
	return destination
}

// TestContextWithCancel adalah fungsi pengujian yang mendemonstrasikan penggunaan context.WithCancel
// untuk mengelola dan membatalkan goroutine secara aman
func TestContextWithCancel(t *testing.T) {
	// Mencetak jumlah goroutine yang sedang berjalan sebelum memulai operasi
	// Berguna untuk memastikan tidak ada kebocoran goroutine
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	// Membuat context induk yang kosong sebagai root context
	parent := context.Background()
	
	// Membuat context turunan yang dapat dibatalkan
	// ctx: context yang dapat dibatalkan
	// cancel: fungsi yang akan digunakan untuk membatalkan operasi
	ctx, cancel := context.WithCancel(parent)

	// Membuat channel counter yang dapat dibatalkan menggunakan context
	// CreateCounter akan menjalankan goroutine baru
	destination := CreateCounter(ctx)

	// Mencetak jumlah goroutine setelah membuat counter
	// Seharusnya bertambah 1 dari jumlah awal karena CreateCounter membuat goroutine baru
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	// Melakukan iterasi nilai dari channel destination
	// Loop akan berhenti jika channel ditutup atau nilai mencapai 10
	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break // Keluar dari loop saat counter mencapai 10
		}
	}

	// Memanggil fungsi cancel untuk membatalkan context
	// Ini akan mengirim sinyal pembatalan ke semua goroutine yang menggunakan context ini
	cancel()

	// Memberikan waktu 2 detik untuk memastikan goroutine telah selesai dibersihkan
	// Best practice: Dalam produksi, lebih baik menggunakan WaitGroup daripada sleep
	time.Sleep(2 * time.Second)

	// Mencetak jumlah goroutine di akhir
	// Seharusnya kembali ke jumlah awal, menunjukkan tidak ada kebocoran goroutine
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

// TestContextWithTimeout menguji penggunaan context dengan timeout.
// Fungsi ini mendemonstrasikan cara yang benar untuk:
// - Menggunakan context.WithTimeout untuk membatasi waktu eksekusi
// - Menangani pembersihan resources dengan defer cancel
// - Memantau jumlah goroutine untuk mencegah kebocoran
func TestContextWithTimeout(t *testing.T) {
	// Mencetak jumlah goroutine awal sebagai baseline
	// Best practice: Selalu monitor jumlah goroutine sebelum operasi untuk deteksi kebocoran
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	// Membuat context induk yang akan menjadi parent
	// Best practice: Selalu gunakan Background() sebagai root context
	parent := context.Background()

	// Membuat context dengan timeout 5 detik
	// Best practice: Selalu tentukan timeout yang masuk akal sesuai kebutuhan operasi
	// Best practice: Simpan fungsi cancel untuk dibersihkan nanti
	ctx, cancel := context.WithTimeout(parent, 5 * time.Second)

	// Menjamin pembersihan resources dengan memanggil cancel
	// Best practice: Selalu gunakan defer cancel() segera setelah WithTimeout/WithDeadline
	defer cancel()

	// Membuat counter yang akan dibatalkan oleh timeout
	// Best practice: Gunakan context untuk mengontrol lifecycle goroutine
	destination := CreateCounter(ctx)

	// Mencetak jumlah goroutine setelah membuat counter
	// Best practice: Monitor perubahan jumlah goroutine untuk memastikan creation berhasil
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	// Membaca nilai dari channel sampai channel ditutup (karena timeout atau selesai)
	// Best practice: Gunakan range untuk membaca channel sampai ditutup
	for n := range destination {
		fmt.Println("Counter", n)
	}

	// Memberikan waktu untuk cleanup
	// Best practice: Dalam production, lebih baik menggunakan WaitGroup
	// Note: time.Sleep sebaiknya dihindari dalam kode production
	time.Sleep(2 * time.Second)

	// Memeriksa jumlah goroutine di akhir
	// Best practice: Pastikan jumlah goroutine kembali ke nilai awal
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

// TestContextWithDeadline mendemonstrasikan penggunaan context.WithDeadline
// untuk membatalkan operasi pada waktu tertentu di masa depan.
// Best practice: Dokumentasikan tujuan utama fungsi di awal
func TestContextWithDeadline(t *testing.T) {
	// Mencetak jumlah goroutine awal sebagai baseline
	// Best practice: Monitor jumlah goroutine untuk mendeteksi kebocoran
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	// Membuat context induk sebagai root context
	// Best practice: Selalu gunakan Background() sebagai parent context
	parent := context.Background()

	// Membuat context dengan deadline 5 detik dari waktu sekarang
	// Best practice: Tentukan deadline yang masuk akal dan sesuai kebutuhan operasi
	// Best practice: Simpan fungsi cancel untuk pembersihan resources
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5 * time.Second))

	// Memastikan resources dibersihkan ketika fungsi selesai
	// Best practice: Selalu panggil cancel dengan defer segera setelah WithDeadline
	defer cancel()

	// Membuat counter yang akan dibatalkan ketika deadline tercapai
	// Best practice: Gunakan context untuk mengontrol lifecycle goroutine
	destination := CreateCounter(ctx)

	// Memonitor perubahan jumlah goroutine setelah membuat counter
	// Best practice: Pastikan goroutine creation berhasil dengan memeriksa jumlahnya
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	// Membaca nilai dari channel sampai channel ditutup (karena deadline atau pembatalan)
	// Best practice: Gunakan range untuk membaca channel sampai selesai
	for n := range destination {
		fmt.Println("Counter", n)
	}

	// Memberikan waktu untuk proses cleanup
	// Best practice: Dalam production, lebih baik menggunakan WaitGroup
	// Note: Hindari penggunaan time.Sleep dalam kode production
	time.Sleep(2 * time.Second)

	// Memeriksa jumlah goroutine di akhir eksekusi
	// Best practice: Pastikan tidak ada kebocoran goroutine dengan membandingkan
	// jumlah akhir dengan jumlah awal
	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
