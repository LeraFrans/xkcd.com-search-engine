// сделать флаг -с для поиска конфигурационного файла
// сделать отдельный счётчик максимального числа комиксов на сайте
// подомумать о том, какие способы “транзакционности” можно придумать

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
	"task3/words"

	"gopkg.in/yaml.v2"
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

func MakeJSONwithComicsData() {

	source_url, db_file, parallel := readConfig()
	max_num := findMaxNumberOfComics()

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
	findNumbersToGet(numbers, max_num)

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
					written = append(written, num) //запишем номер обработанного комикса в слайс для кэш
					fmt.Printf("Write %d\n", num)
					mu.Unlock() // Освобождаем мьютекс после записи
				}()
			}
		}()
	}

	wg.Wait()                             // Ожидаем завершения всех горутин
	writeResultInJSON(resultMap, db_file) // Записываем данные в JSON
	os.Remove("written_comics.txt")       // удаляем файл с кэшем
}

// Ловит сигнал об остновке программы (Ctrl+C) и безопасно завершает её
func signalHandler(c chan os.Signal, db_file string, resultMap map[int]OutputData, written *[]int) {
	<-c // Ожидаем сигнал SIGINT
	// если прервана:
	fmt.Println("\nПрограмма прервана.")
	writeResultInJSON(resultMap, db_file) //записываем в json то, что успели обработать
	createCashe(written)                  //создаём кэш с номерами записанных комиксов
	os.Exit(0)                            // читала, что это антипатерн, но другое пока не смогла
}

func findNumbersToGet(numbers chan int, max_num int) {
	// Проверяем, есть ли что-то в кэше
	ok, cashe := readCashe("written_comics")
	// если кэш пуст, заполняем канал номерами всех комиксов
	if !ok {
		for num := 1; num <= max_num; num++ {
			numbers <- num
		}
	} else {
		// если кэш есть, тоже заполняем канал, но пропускаем те комиксы, которые в кэше
		for num := 1; num <= max_num; num++ {
			if bytes.Contains(cashe, []byte(fmt.Sprintf("%d", num))) {
				continue
			}
			numbers <- num
		}
	}
	close(numbers)
}

// Создаёт файл written_comics.txt в который записывает номера комиксов, которые успели записать в json
func createCashe(written *[]int) {

	// Преобразование слайса чисел в слайс байт
	byteSlice := []byte(fmt.Sprint(written))

	// если уже есть старый кэш, то склеиваем старый и новый слайсы байтов
	ok, old_byteSlice := readCashe("written_comics.txt")
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
func readCashe(filename string) (bool, []byte) {

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

// Возвращает переменные из "config.yaml"
func readConfig() (string, string, int) {
	// Чтение файла
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	// Разбор YAML-контента
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(content, &yamlMap)
	if err != nil {
		panic(err)
	}

	source_url, ok := yamlMap["source_url"].(string)
	if !ok {
		panic("Неверный тип данных для 'source_url'")
	}
	db_file, ok := yamlMap["db_file"].(string)
	if !ok {
		panic("Неверный тип данных для 'db_file'")
	}
	parallel, ok := yamlMap["parallel"].(int)
	if !ok {
		panic("Неверный тип данных для 'parallel'")
	}

	return source_url, db_file, parallel
}

// Возвращает данные об одном комиксе, полученные с сервера (данные в виде временной структуры InputData)
func getDataOfOneComic(num int, source_url string) InputData {

	//собираем url
	url := fmt.Sprint(source_url, "/", num, "/info.0.json")

	//получаем данные от сервера
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to make a get request to the server: ", err)
	}
	defer resp.Body.Close()

	//декодируем данные во временную структуру
	var input InputData
	json.NewDecoder(resp.Body).Decode(&input)
	// пришлось убрать обработку ошибок при декодировании, тк она не давала парсить те комиксы, где
	// некоторые поля были не заполнены

	return input
}

// Нормализует данные об одном комиксе и оформляет их в структурку
func DataProcess(input InputData) OutputData {
	//обрабатываем нормализатором
	transcriptWithAlt := fmt.Sprint(input.Transcript, input.Alt) //описание и краткое описание сливаем в одну строку, чтобы для каждого нормализацию не делать отдельно
	normalizated, err := words.WorldsNormalizator(transcriptWithAlt)
	if err != nil {
		log.Fatal(err)
	}

	//формируем элемент данных об одном коммиксе и добавляем его в результирующую мапу
	output := OutputData{
		Url:      input.Img,
		Keywords: normalizated,
	}

	return output
}

// Проверяет, существует ли файл json, если нет, то создаёт его
func checkFileIsExist(filename string) error {
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
	checkFileIsExist(db_file)
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

// Доделаю позже
func findMaxNumberOfComics() int {
	return 2916
}
