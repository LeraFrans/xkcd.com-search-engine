package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"task8/config"
	"task8/pkg/repository"
)

type UpdateService struct {
	repo repository.Update
}

func NewUpdateService(repo repository.Update) *UpdateService {
	return &UpdateService{repo: repo}
}

// отчёт по обновлению таблицы комиксов
type UpdateResponse struct {
	Total int `json:"total"`
	New   int `json:"new"`
}

// просто "временная структура", тут для удобства храним считанные с сайта данные об одном комиксе до обработки
type DataFromXksdCom struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

// вся логика обновления таблицы с комиксами
// возвращает в структуре ответа для клиента количество новых комиксов и общее количество комиксов в бд
func (s *UpdateService) UpdateComicTable() (UpdateResponse, error) {

	max_num := findMaxNumberOfComics() //сколько всего комиксов на сайте xkcd.com в данный момент

	var wg sync.WaitGroup

	var resultComicsSlice []repository.Comic // слайс всех обработанных комиксов
	numbers := make(chan int, max_num)       // в канале номера комиксов, подлежащие обработке
	var written []int                        // номера комиксов, которые успели записать (для кэша)
	mu := &sync.Mutex{}                      // Мьютекс для синхронизации доступа к resultComicsSlice

	// Устанавливаем обработчик сигнала SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	go s.signalHandler(c, resultComicsSlice, &written) // ловит и обрабатывает сигнал

	// Заполняем канал numbers номерами комиксов, которые нужно обработать
	countOfNewComics, err := s.findNumbersToGet(numbers, max_num)
	if err != nil {
		ErrResponse := UpdateResponse{}
		return ErrResponse, err
	}

	// Создаем пул горутин для выполнения запросов
	for i := 0; i < config.ReadConfig().Parallel; i++ {
		// тут комикс достаётся с сайта и обрабатывается
		wg.Add(1)
		go func() {
			defer wg.Done()
			for num := range numbers {
				oneComicFromXKCD := getDataOfOneComic(num)         // забирает данные об одном комиксе с сайта
				normalizatedComic := DataProcess(oneComicFromXKCD) // обрабатывает их до нужной нам формы
				fmt.Printf("Get %d go %d\n", num, i)
				// а тут комикс записывается в результирующий слайс и слайс для кэша
				wg.Add(1)
				go func() {
					defer wg.Done()
					mu.Lock() // Захватываем мьютекс перед записью в результирующий слайс

					resultComicsSlice = append(resultComicsSlice, normalizatedComic)
					written = append(written, num) //запишем номер обработанного комикса в слайс для кэша
					fmt.Printf("Write %d\n", num)
					mu.Unlock() // Освобождаем мьютекс после записи
				}()
			}
		}()
	}

	wg.Wait()                                                 // Ожидаем завершения всех горутин
	err = s.repo.WriteResultInKeywordTable(resultComicsSlice) // Записываем данные в JSON
	os.Remove("written_comics.txt")                           // удаляем файл с кэшем

	// формируем ответ (потом выше пойдёт клиенту)
	response := UpdateResponse{
		Total: max_num,
		New:   countOfNewComics,
	}
	return response, err
}

// Ловит сигнал об остновке программы (Ctrl+C) и безопасно завершает её
func (s *UpdateService) signalHandler(c chan os.Signal, resultComicsSlice []repository.Comic, written *[]int) {
	<-c // Ожидаем сигнал SIGINT
	// если прервана:
	fmt.Println("\nПрограмма прервана.")
	err := s.repo.WriteResultInKeywordTable(resultComicsSlice) //записываем в json то, что успели обработать
	if err != nil {
		log.Println(err)
	}
	createCache(written) //создаём кэш с номерами записанных комиксов
	os.Exit(0)           // читала, что это антипатерн, но другое пока не смогла

}

// кладёт в канал номера комиксов для записи и возвращает их количество
func (s *UpdateService) findNumbersToGet(numbers chan int, max_num int) (int, error) {

	countToReturn := max_num // по умолчанию база пустая

	// существует ли пустая дб
	ok1 := s.repo.IsTableHasAnyRows("keyword")
	// Проверяем, есть ли что-то в кэше
	ok2, cache := readCache("written_comics.txt")

	// если нет базы данных и нет кэша, заполняем все номера
	if !ok1 {
		for num := 1; num <= max_num; num++ {
			numbers <- num
		}
	} else if ok1 && ok2 {
		// если есть кэш и бд, тоже заполняем канал, но пропускаем те комиксы, которые в кэше
		for num := 1; num <= max_num; num++ {
			if bytes.Contains(cache, []byte(fmt.Sprintf("%d", num))) {
				countToReturn--
				continue
			}
			numbers <- num
		}
	} else {
		// Если есть непустая бд, но нет кэша, значит проверяем максимальное количество
		// комиксов на сайте, сравниваем с нашей ДБ и если нужно, обновляем базу
		// (кладём к канал недостающие номера)
		realMax := max_num
		ourMax, err := s.repo.FindOurMaxNumberOfComics()
		if err != nil {
			return 0, errors.New("error finding our max number of comics")
		}

		if realMax > ourMax {
			for num := ourMax + 1; num <= max_num; num++ {
				numbers <- num
			}
			countToReturn = realMax - ourMax
		}
		if realMax == ourMax {
			countToReturn = 0 // если обновлять не нужно, то количество новых комиксов будет 0
		}
	}

	close(numbers)

	return countToReturn, nil

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
		log.Println("Ошибка при создании файла:", err)
	}
	defer file.Close()

	// Запись слайса байт в файл
	_, err = file.Write(byteSlice)
	if err != nil {
		log.Println("Ошибка при записи в файл:", err)
	}
}

// Проверяет наличие кэша и если есть, достаёт оттуда данные
func readCache(filename string) (bool, []byte) {

	var bytes []byte

	// Если файл не существует возвращаем FALSE
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, bytes
		}
	}

	// Если кэш существует, считываем его и возвращаем TRUE и содержимое кэша
	file, _ := os.Open(filename)
	defer file.Close()
	bytes, _ = os.ReadFile(filename)

	return true, bytes
}

// Возвращает данные об одном комиксе, полученные с сайта
func getDataOfOneComic(num int) DataFromXksdCom {

	source_url := config.ReadConfig().Source_url

	//собираем url
	url := fmt.Sprint(source_url, "/", num, "/info.0.json")

	//получаем данные от сервера
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to make a get request to the server: ", err)
	}
	defer resp.Body.Close()

	var oneComicFromXKCD DataFromXksdCom
	// Проверяем статус ответа
	if resp.StatusCode != 200 {
		log.Println("URL Not Found")
	} else {
		//декодируем данные во временную структуру
		json.NewDecoder(resp.Body).Decode(&oneComicFromXKCD)
		// пришлось убрать обработку ошибок при декодировании, тк она не давала парсить те комиксы, где
		// некоторые поля были не заполнены
	}

	return oneComicFromXKCD
}

// Нормализует данные об одном комиксе и оформляет их
func DataProcess(oneComicFromXKCD DataFromXksdCom) repository.Comic {
	//обрабатываем нормализатором
	transcriptWithAlt := fmt.Sprint(oneComicFromXKCD.Transcript, oneComicFromXKCD.Alt) //описание и краткое описание сливаем в одну строку, чтобы для каждого нормализацию не делать отдельно
	normalizated, err := WordsNormalizator(transcriptWithAlt)
	if err != nil {
		log.Println(err)
	}

	//формируем элемент данных об одном коммиксе
	output := repository.Comic{
		Num:      oneComicFromXKCD.Num,
		Url:      oneComicFromXKCD.Img,
		Keywords: strings.Join(normalizated, " "),
	}

	return output
}

// проверяет, существует ли URL страница комикса
func IsPageExists(num int) bool {

	source_url := config.ReadConfig().Source_url

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

// Возвращает максимальный номер комикса на сайте xkcd.com
func findMaxNumberOfComics() int {

	//return 2936

	var num int = 100
	var emptyNum int = -1
	var lastNotEmptyNum int = -1

	fmt.Println("I'm find max numbers of comics xkcd.com\nPlease, wait...")

	for {
		// с шагом в 100 к номерам комиксов от 100 до бесконечности делаем запросы к сайту
		// останавливаемся на той сотке, где не найден комикс, записываем его в emptyNum
		for ; emptyNum == -1; num += 100 {
			if !IsPageExists(num) {
				emptyNum = num
			}
		}
		num -= 100

		// От ненайденного сотого номера идём назад с шагом 1. Считаем количество
		// отсутствующих комиксов счётчиком counterOfEmpty. При нахождении
		// существующей страницы пишем её в lastNotEmptyNum
		var counterOfEmpty int = 0
		for ; lastNotEmptyNum == -1; num-- {
			if IsPageExists(num) {
				lastNotEmptyNum = num
			} else {
				counterOfEmpty++
			}
		}
		// Проверяем счётчик пустых комиксов, если он меньше 10, то на всякий
		// случай нужно сходить ещё и наверх от ненайденного сотого, вдуг эти ненайденные
		// комиксы просто удалены
		if counterOfEmpty >= 10 {
			break
		} else {
			flagRetry := 0
			num = emptyNum + 1
			for i := 0; i < 10; i++ {
				if IsPageExists(num) {
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

	fmt.Println("Max number: ", lastNotEmptyNum)

	return lastNotEmptyNum
}
