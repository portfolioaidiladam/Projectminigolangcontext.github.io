package belajar_golang_context

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

// TestContext mendemonstrasikan pembuatan context dasar
// Background() dan TODO() adalah dua context kosong yang dapat digunakan
func TestContext(t *testing.T) {
	// context.Background() adalah context dasar, root dari semua context
	background := context.Background()
	fmt.Println(background)

	// context.TODO() digunakan ketika belum jelas context apa yang akan digunakan
	todo := context.TODO()
	fmt.Println(todo)
}

// TestContextWithValue mendemonstrasikan penggunaan context dengan nilai (key-value)
// dan bagaimana nilai tersebut diwariskan ke context turunannya
func TestContextWithValue(t *testing.T) {
	// Membuat context induk
	contextA := context.Background()

	// Membuat context turunan dengan nilai
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	// Membuat context turunan dari contextB
	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	// Membuat context turunan dari contextC
	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextF, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)

	fmt.Println(contextF.Value("f"))
	fmt.Println(contextF.Value("c"))
	fmt.Println(contextF.Value("b"))

	fmt.Println(contextA.Value("b"))
}

// CreateCounter membuat channel yang menghasilkan angka berurutan
// dan dapat dibatalkan menggunakan context
func CreateCounter(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			select {
			case <-ctx.Done(): // Memeriksa apakah context sudah dibatalkan
				return
			default:
				destination <- counter
				counter++
				time.Sleep(1 * time.Second) // Simulasi proses yang lambat
			}
		}
	}()

	return destination
}

// TestContextWithCancel mendemonstrasikan penggunaan context.WithCancel
// untuk membatalkan operasi secara manual
func TestContextWithCancel(t *testing.T) {
	// Mencetak jumlah goroutine sebelum operasi
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	// Membuat context yang dapat dibatalkan
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounter(ctx)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
		if n == 10 {
			break
		}
	}

	cancel() // mengirim sinyal cancel ke context

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

// TestContextWithTimeout mendemonstrasikan penggunaan context.WithTimeout
// untuk membatalkan operasi setelah durasi tertentu
func TestContextWithTimeout(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	// Membuat context yang akan kedaluwarsa setelah 5 detik
	ctx, cancel := context.WithTimeout(parent, 5 * time.Second)
	defer cancel() // Pastikan untuk selalu memanggil cancel untuk membersihkan resources

	destination := CreateCounter(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}

// TestContextWithDeadline mendemonstrasikan penggunaan context.WithDeadline
// untuk membatalkan operasi pada waktu tertentu
func TestContextWithDeadline(t *testing.T) {
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	parent := context.Background()
	// Membuat context yang akan kedaluwarsa 5 detik dari sekarang
	ctx, cancel := context.WithDeadline(parent, time.Now().Add(5 * time.Second))
	defer cancel() // Pastikan untuk selalu memanggil cancel untuk membersihkan resources

	destination := CreateCounter(ctx)
	fmt.Println("Total Goroutine", runtime.NumGoroutine())

	for n := range destination {
		fmt.Println("Counter", n)
	}

	time.Sleep(2 * time.Second)

	fmt.Println("Total Goroutine", runtime.NumGoroutine())
}
