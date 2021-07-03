package main

type Limiter struct {
	concurrentCount int
	bucket          chan int
}

func NewLimiter(cc int) *Limiter {
	return &Limiter{
		concurrentCount: cc,
		bucket:          make(chan int, cc),
	}
}

func (l *Limiter) Get() bool {
	if len(l.bucket) >= l.concurrentCount {
		return false
	}
	l.bucket <- 1
	return true
}

func (l *Limiter) Release() {
	<-l.bucket
}

// func dosome(l *Limiter, n int) bool {
// 	if !l.Get() {
// 		return false
// 	}
// 	defer l.Release()
// 	log.Println(n)
// 	return true
// }

// func main() {
// 	log.Println("start...\n=====")

// 	l := NewLimiter(10)
// 	for i := 0; i < 50; i++ {
// 		i := i
// 		go func() {
// 			for {
// 				if dosome(l, i) {
// 					break
// 				}
// 			}
// 		}()
// 	}
// 	time.Sleep(5 * time.Second)

// 	log.Println("=====\nexit.")
// }
