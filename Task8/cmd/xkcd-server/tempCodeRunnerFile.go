package main
// db, err := repository.ConnectDB()
// 	if err != nil {
// 		log.Fatalf("failed to initializing db: %s", err.Error())
// 	}
// 	defer db.Close()

// 	// создание трёх слоёв: репозитория(работа с БД), сервиса(бизнес-логика) и хендлеров(связь с клиентом)
// 	repos := repository.NewRepository(db)
// 	services := service.NewService(repos)