package main

// func main() {
// 	PORT := os.Getenv("PORT")
// 	if PORT == "" {
// 		PORT = "8000"
// 	}
//
// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != http.MethodGet {
// 			w.WriteHeader(http.StatusMethodNotAllowed)
// 			w.Write([]byte("Method not allowed"))
// 			return
// 		}
// 		w.Header().Add("Content-Type", "text/html")
// 		w.Write([]byte("Hello World"))
// 	})
//
// 	s := http.Server{
// 		Addr:        ":" + PORT,
// 		Handler:     mux,
// 		ReadTimeout: 2 * time.Second,
// 	}
//
// 	go func() {
// 		log.Println("Starting server at http://localhost:" + PORT)
// 		err := s.ListenAndServe()
// 		if err != nil {
// 			log.Panicf("error starting server")
// 		}
// 	}()
//
// 	ch := make(chan os.Signal, 1)
//
// 	signal.Notify(ch, os.Interrupt)
// 	signal.Notify(ch, os.Kill)
//
// 	sig := <-ch
// 	log.Println("Got signal ", sig)
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
//
// 	s.Shutdown(ctx)
// }
