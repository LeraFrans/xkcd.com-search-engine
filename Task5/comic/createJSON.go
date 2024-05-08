package comic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"task5/words"
)

// просто "временная структура", тут для удобства храним считанные с сайта данные об одном комиксе до обработки
type InputData struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

// а тут уже обработанные данные, предназначенные для записи в json
type OutputData struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

// создаёт/обновляет БД, возвращает количество новых комиксов и общее количество комиксов
func MakeJSONWithComicsData(source_url string, db_file string, parallel int) (int, int) {

	max_num := findMaxNumberOfComics(source_url)

	var wg sync.WaitGroup

	resultMap := make(map[int]OutputData) // мапа где обработанные данные
	numbers := make(chan int, max_num)    // в канале номера комиксов, подлежащие обработке
	var written []int                     // номера комиксов, которые успели записать (для кэша)
	mu := &sync.Mutex{}                   // Мьютекс для синхронизации доступа к resultMap

	// Устанавливаем обработчик сигнала SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go signalHandler(c, db_file, resultMap, &written) // ловит и обрабатывает сигнал

	// Заполняем канал numbers номерами комиксов, которые нужно обработать
	countOfNewComics := findNumbersToGet(numbers, max_num, db_file, source_url)

	// Создаем пул горутин для выполнения запросов
	for i := 0; i < parallel; i++ {
		// тут комикс достаётся с сайта и обрабатывается
		wg.Add(1)
		go func() {
			defer wg.Done()
			for num := range numbers {
				input := getDataOfOneComic(num, source_url)
				output := DataProcess(input)
				fmt.Printf("Get %d go %d\n", num, i)
				// а тут комикс записывается в результирующую мапу и слайс для кэша
				wg.Add(1)
				go func() {
					defer wg.Done()
					mu.Lock() // Захватываем мьютекс перед записью в мапу
					resultMap[input.Num] = output
					written = append(written, num) //запишем номер обработанного комикса в слайс для кэша
					fmt.Printf("Write %d\n", num)
					mu.Unlock() // Освобождаем мьютекс после записи
				}()
			}
		}()
	}

	wg.Wait()                             // Ожидаем завершения всех горутин
	writeResultInJSON(resultMap, db_file) // Записываем данные в JSON
	os.Remove("written_comics.txt")       // удаляем файл с кэшем

	return countOfNewComics, max_num
}

// Ловит сигнал об остновке программы (Ctrl+C) и безопасно завершает её
func signalHandler(c chan os.Signal, db_file string, resultMap map[int]OutputData, written *[]int) {
	<-c // Ожидаем сигнал SIGINT
	// если прервана:
	fmt.Println("\nПрограмма прервана.")
	writeResultInJSON(resultMap, db_file) //записываем в json то, что успели обработать
	createCache(written)                  //создаём кэш с номерами записанных комиксов
	os.Exit(0)                            // читала, что это антипатерн, но другое пока не смогла

}

// кладёт в канал номера комиксов для записи и возвращает их количество
func findNumbersToGet(numbers chan int, max_num int, db_file string, source_url string) int {

	countToReturn := max_num // по умолчанию база пустая

	// существует ли пустая дб
	ok1 := isDataBaseExist(db_file)
	// Проверяем, есть ли что-то в кэше
	ok2, cache := readCache("written_comics.txt")

	// если нет базы данных и нет кэша, заполняем все номера
	if !ok1 && !ok2 {
		for num := 1; num <= max_num; num++ {
			numbers <- num
		}
	} else if ok1 && ok2 {
		// если есть кэш и бд, тоже заполняем канал, но пропускаем те комиксы, которые в кэше
		for num := 1; num <= max_num; num++ {
			if bytes.Contains(cache, []byte(fmt.Sprintf("%d", num))) {
				continue
			}
			numbers <- num
		}
	} else {
		// Если есть непустая бд, но нет кэша, значит проверяем максимальное количество
		// комиксов на сайте, сравниваем с нашей ДБ и если нужно, обновляем базу
		// (кладём к канал недостающие номера)
		realMax := findMaxNumberOfComics(source_url)

		ourMax := 0
		db := readDataBase()
		for key := range db {
			if key > ourMax {
				ourMax = key
			}
		}

		if realMax > ourMax {
			for num := ourMax + 1; num <= max_num; num++ {
				numbers <- num
			}
			countToReturn = realMax - ourMax
		}

	}

	return countToReturn

}

func isDataBaseExist(db_file string) bool {
	// Если файл не существует возвращаем FALSE
	info, err := os.Stat(db_file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	// Проверяем, равен ли размер файла нулю
	if info.Size() == 0 {
		return false
	}
	return true
}

// Создаёт файл written_comics.txt в который записывает номера комиксов, которые успели записать в json
func createCache(written *[]int) {

	// Преобразование слайса чисел в слайс байт
	byteSlice := []byte(fmt.Sprint(written))

	// если уже есть старый кэш, то склеиваем старый и новый слайсы байтов
	ok, old_byteSlice := readCache("written_comics.txt")
	if ok {
		byteSlice = append(old_byteSlice, byteSlice...)
	}

	// Создаём файл кэша
	file, err := os.Create("written_comics.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
	}
	defer file.Close() // Не забываем закрыть файл после использования

	// Запись слайса байт в файл
	_, err = file.Write(byteSlice)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}
}

// Проверяет наличие кэша и если есть, достаёт оттуда данные
func readCache(filename string) (bool, []byte) {

	var bytes []byte

	// Если файл не существует возвращаем FALSE
	_, err := os.Stat("written_comics.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return false, bytes
		}
	}

	// Если кэш существует, считываем его и возвращаем TRUE и содержимое кэша
	file, _ := os.Open(filename)
	defer file.Close()
	bytes, _ = os.ReadFile("written_comics.txt")
	return true, bytes
}

// Возвращает данные об одном комиксе, полученные с сервера (данные в виде временной структуры InputData)
func getDataOfOneComic(num int, source_url string) InputData {

	//собираем url
	url := fmt.Sprint(source_url, "/", num, "/info.0.json")

	//получаем данные от сервера
	resp, err := http.Get(url)

	if err != nil {
		log.Println("Failed to make a get request to the server: ", err)
	}
	defer resp.Body.Close()

	var input InputData
	// Проверяем статус ответа
	if resp.StatusCode != 200 {
		log.Println("URL Not Found")
	} else {
		//декодируем данные во временную структуру
		json.NewDecoder(resp.Body).Decode(&input)
		// пришлось убрать обработку ошибок при декодировании, тк она не давала парсить те комиксы, где
		// некоторые поля были не заполнены
	}

	return input
}

// Нормализует данные об одном комиксе и оформляет их в структурку
func DataProcess(input InputData) OutputData {
	//обрабатываем нормализатором
	transcriptWithAlt := fmt.Sprint(input.Transcript, input.Alt) //описание и краткое описание сливаем в одну строку, чтобы для каждого нормализацию не делать отдельно
	normalizated, err := words.WordsNormalizator(transcriptWithAlt)
	if err != nil {
		log.Println(err)
	}

	//формируем элемент данных об одном коммиксе и добавляем его в результирующую мапу
	output := OutputData{
		Url:      input.Img,
		Keywords: normalizated,
	}

	return output
}

// Проверяет, существует ли файл json, если нет, то создаёт его
func createFileIfNotExist(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// Записывает результирующую мапу в json
// Если json не пустой, старые данные сохраняются, а новые добавляются рядом с ними.
func writeResultInJSON(resultMap map[int]OutputData, db_file string) {
	// Проверка наличия файла
	createFileIfNotExist(db_file)
	// Открытие файла для чтения
	file, _ := os.Open(db_file)
	defer file.Close()

	// Создание новой мапы для хранения старых данных из файла
	var oldResultMap map[int]OutputData
	json.NewDecoder(file).Decode(&oldResultMap)

	// Объединение старых и новых данных
	combinedResultMap := make(map[int]OutputData)
	for k, v := range oldResultMap {
		combinedResultMap[k] = v
	}
	for k, v := range resultMap {
		combinedResultMap[k] = v
	}

	// Подготовка результирующей мапы к записи в json
	dataBytes, _ := json.MarshalIndent(combinedResultMap, "", "	")

	// Запись
	os.WriteFile(db_file, dataBytes, 0)
}

// // возвращает путь до конфигурационного файла, обрабатывает флаг -с
// func ParsArgument_C() string {

// 	// смотрим, есть ли тут флаг -с
// 	cFlag := flag.Bool("c", false, "Флаг -c")
// 	flag.Parse()

// 	// обрабатываем оставшиеся аргументы (ищем путь до конфига)
// 	path := flag.Args()
// 	var configName string
// 	//устанавливаем max_num
// 	if *cFlag {
// 		configName = path[0] // ноль, тк путь должен идти сразу после флага -с
// 	} else {
// 		configName = "config.yaml" // если нет флага -с
// 	}
// 	return configName
// }

// проверяет, существует ли URL страница комикса
func isPageExists(num int, source_url string) bool {
	url := fmt.Sprint(source_url, "/", num, "/info.0.json") //собираем url
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
	}
	// Закрываем ответ
	defer resp.Body.Close()
	// Проверяем статус ответа
	if resp.StatusCode == 200 {
		return true
	} else {
		return false
	}
}

// Возвращает максимальный номер комикса
func findMaxNumberOfComics(source_url string) int {

	var num int = 100
	var emptyNum int = -1
	var lastNotEmptyNum int = -1

	for {
		// с шагом в 100 к номерам комиксов от 100 до бесконечности делаем запросы к сайту
		// останавливаемся на той сотке, где не найден комикс, записываем его в emptyNum
		for ; emptyNum == -1; num += 100 {
			if !isPageExists(num, source_url) {
				emptyNum = num
			}
		}
		num -= 100

		// От ненайденного сотого номера идём назад с шагом 1. Считаем количество
		// отсутствующих комиксов счётчиком counterOfEmpty. При нахождении
		// существующей страницы пишем её в lastNotEmptyNum
		var counterOfEmpty int = 0
		for ; lastNotEmptyNum == -1; num-- {
			if isPageExists(num, source_url) {
				lastNotEmptyNum = num
			} else {
				counterOfEmpty++
			}
		}
		// Проверяем счётчик пустых комиксов, если он меньше 10, то на всякий
		// случай нужно сходить ещё и наверх от ненайденного сотого, вдуг эти ненайденные
		// комиксы просто удалены
		if counterOfEmpty >= 10 {
			break //
		} else {
			flagRetry := 0
			num = emptyNum + 1
			for i := 0; i < 10; i++ {
				if isPageExists(num, source_url) {
					flagRetry = 1 // нашли комикс, идём на следующую сотку и заново всё
				}
			}
			if flagRetry == 0 {
				break
			} // если не нашли правее существующих, выходим
		}

		// до сих пор не вышли, значит идём на следующую сотку
		num = emptyNum + 100
	}

	return lastNotEmptyNum
}
