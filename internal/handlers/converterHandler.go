package handlers

import (
	"TestTaskBeelineDev/internal/config"
	"TestTaskBeelineDev/internal/models"
	"TestTaskBeelineDev/internal/services"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ConverterHandler обрабатывает HTTP-запросы на преобразование данных
type ConverterHandler struct {
	Cfg *config.Config
}

// XmlConverterHandler обрабатывает POST-запросы с XML-данными пользователей
// и преобразует их в JSON для отправки во внешний сервис
func (c *ConverterHandler) XmlConverterHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("XmlConverterHandler serving the request")

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/xml") {
		log.Printf("Wrong contentType: %s", contentType)
		http.Error(w, "Content-Type not supported", http.StatusUnsupportedMediaType)
		return
	}

	if r.Method != http.MethodPost {
		log.Printf("Wrong method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var users models.UsersInput
	if err := xml.NewDecoder(r.Body).Decode(&users); err != nil {
		log.Println("Error while parsing XML body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(users.Users) == 0 {
		log.Println("Error while parsing XML body: empty users list")
		http.Error(w, "No users found", http.StatusBadRequest)
		return
	}

	converter := services.Converter{}

	var convertedUsers []models.UserOutput
	var mx sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, c.Cfg.NumOfWorkers) // Семафор для ограничения горутин

	log.Println("Processing users")
	// Обработка каждого пользователя в отдельной горутине
	for _, user := range users.Users {
		wg.Add(1)
		sem <- struct{}{}
		go func(u models.UserInput) {
			defer wg.Done()
			defer func() { <-sem }()

			res, err := converter.ConvertUser(u)
			if err != nil {
				log.Println("Error while converting user:", err)
				return
			}

			mx.Lock()
			convertedUsers = append(convertedUsers, res)
			mx.Unlock()
		}(user)
	}

	wg.Wait()
	log.Println("Finished users processing")

	// Маршалинг преобразованных пользователей в JSON
	result, err := json.Marshal(convertedUsers)
	if err != nil {
		log.Println("Error while marshalling users:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Creating the new request")
	client := http.Client{
		Timeout: c.Cfg.Timeout,
	}

	// Подготовка запроса к внешнему сервису
	req, err := http.NewRequest("POST", c.Cfg.OutputServiceURL, bytes.NewBuffer(result))
	if err != nil {
		log.Println("Error while creating request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.Cfg.AuthKey)

	req = req.WithContext(r.Context())

	log.Println(string(result))
	// Попытки отправки с задержкой (3 попытки)
	log.Println("Sending the request")
	var resp *http.Response
	for i := 0; i < 3; i++ {
		//Проверка, не закрыл ли клиент запрос
		select {
		case <-r.Context().Done():
			log.Println("Request cancelled by client")
			return
		default:
		}

		// Клонирование запроса с новым телом для повторных попыток
		newReq := req.Clone(req.Context())
		newReq.Body = io.NopCloser(bytes.NewReader(result))

		resp, err = client.Do(newReq)
		if err == nil {
			break
		}

		// Задержка перед повторной попыткой
		log.Println("The external service is not responding: retrying...")
		time.Sleep(time.Duration(i*5) * time.Second)
	}

	if err != nil {
		log.Println("Error while sending request:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Println("Checking response from the external service")
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("External service error: Status - %d, %s", resp.StatusCode, string(body))
		http.Error(w, "External service error", resp.StatusCode)
		return
	}

	log.Println("Successfully completed the request")
	w.WriteHeader(http.StatusAccepted)
}
