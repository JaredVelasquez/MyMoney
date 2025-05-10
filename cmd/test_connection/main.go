package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"mi-app-backend/db/config"

	"github.com/joho/godotenv"
)

// Estructura para almacenar la respuesta de la API de Supabase
type SupabaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Determinar qué tipo de prueba realizar
	testType := "rest"
	if len(os.Args) > 1 {
		testType = os.Args[1]
	}

	switch testType {
	case "rest":
		testRESTAPI()
	case "db":
		testDatabaseConnection()
	default:
		fmt.Println("Unknown test type. Use 'rest' or 'db'")
		os.Exit(1)
	}
}

// Prueba la conexión a la API REST de Supabase
func testRESTAPI() {
	// Verificar que las variables necesarias estén configuradas
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_API_KEY")

	if supabaseURL == "" {
		log.Fatal("SUPABASE_URL environment variable must be set")
	}

	if supabaseKey == "" {
		log.Fatal("SUPABASE_API_KEY environment variable must be set")
	}

	// Intentar conectar a Supabase usando la API REST
	fmt.Println("\nTesting Supabase REST API connection...")
	fmt.Printf("Supabase URL: %s\n", supabaseURL)

	// Crear una solicitud para verificar la salud del servicio
	healthURL := fmt.Sprintf("%s/rest/v1/", supabaseURL)

	// Configurar timeout para diagnóstico
	fmt.Println("Attempting connection with 10 second timeout...")
	connectionStart := time.Now()

	// Crear un cliente HTTP con timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Crear la solicitud
	req, err := http.NewRequest("GET", healthURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Agregar los encabezados de autenticación de Supabase
	req.Header.Add("apikey", supabaseKey)
	req.Header.Add("Authorization", "Bearer "+supabaseKey)

	// Ejecutar la solicitud
	resp, err := client.Do(req)
	if err != nil {
		connectionDuration := time.Since(connectionStart)
		fmt.Printf("Connection failed after %v\n", connectionDuration)
		log.Fatalf("Error connecting to Supabase API: %v", err)
	}
	defer resp.Body.Close()

	connectionDuration := time.Since(connectionStart)
	fmt.Printf("Connection established in %v\n", connectionDuration)

	// Leer la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Verificar el código de estado
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from Supabase API: %s - %s", resp.Status, string(body))
	}

	fmt.Printf("Supabase API connection successful! Status code: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))
	fmt.Printf("Total connection test completed in %v\n", time.Since(connectionStart))

	// Ejecutar una consulta simple para verificar la funcionalidad
	fmt.Println("\nTesting a simple query to Supabase...")
	queryStart := time.Now()

	// Crear una consulta simple a una tabla (ajustar según la estructura de tu base de datos)
	// Por ejemplo, intentar obtener la primera fila de la tabla 'users'
	queryURL := fmt.Sprintf("%s/rest/v1/users?select=*&limit=1", supabaseURL)

	req, err = http.NewRequest("GET", queryURL, nil)
	if err != nil {
		log.Fatalf("Error creating query request: %v", err)
	}

	// Agregar los encabezados de autenticación de Supabase
	req.Header.Add("apikey", supabaseKey)
	req.Header.Add("Authorization", "Bearer "+supabaseKey)
	req.Header.Add("Content-Type", "application/json")

	// Ejecutar la solicitud
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Error executing query: %v", err)
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading query response: %v", err)
	}

	queryDuration := time.Since(queryStart)

	// Verificar el resultado (no importa si no hay datos, solo queremos verificar que la API responda)
	fmt.Printf("Query completed in %v\n", queryDuration)
	fmt.Printf("Status code: %d\n", resp.StatusCode)

	// Formatear la respuesta JSON para mejor legibilidad
	var prettyJSON bytes.Buffer
	if json.Valid(body) {
		err = json.Indent(&prettyJSON, body, "", "  ")
		if err != nil {
			fmt.Printf("Response: %s\n", string(body))
		} else {
			fmt.Printf("Response: %s\n", prettyJSON.String())
		}
	} else {
		fmt.Printf("Response: %s\n", string(body))
	}
}

// Prueba la conexión directa a la base de datos de Supabase
func testDatabaseConnection() {
	// Configurar variables de entorno para la conexión
	if os.Getenv("SUPABASE_HOST") == "" {
		os.Setenv("SUPABASE_HOST", "aws-0-us-east-2.pooler.supabase.com")
		fmt.Println("SUPABASE_HOST not set, using default: aws-0-us-east-2.pooler.supabase.com")
	}

	if os.Getenv("SUPABASE_PORT") == "" {
		os.Setenv("SUPABASE_PORT", "6543")
		fmt.Println("SUPABASE_PORT not set, using default: 6543")
	}

	if os.Getenv("SUPABASE_USER") == "" {
		os.Setenv("SUPABASE_USER", "postgres.eddrisdxafllaekoaduy")
		fmt.Println("SUPABASE_USER not set, using default: postgres.eddrisdxafllaekoaduy")
	}

	if os.Getenv("SUPABASE_DBNAME") == "" {
		os.Setenv("SUPABASE_DBNAME", "postgres")
		fmt.Println("SUPABASE_DBNAME not set, using default: postgres")
	}

	// Verificar que la contraseña esté configurada
	if os.Getenv("SUPABASE_PASSWORD") == "" {
		log.Fatal("SUPABASE_PASSWORD environment variable must be set")
	}

	// Intentar crear una conexión a la base de datos usando el cliente de Supabase
	fmt.Println("\nTesting Supabase database connection...")

	// Mostrar información de conexión (sin contraseña)
	host := os.Getenv("SUPABASE_HOST")
	port := os.Getenv("SUPABASE_PORT")
	user := os.Getenv("SUPABASE_USER")
	dbname := os.Getenv("SUPABASE_DBNAME")

	fmt.Printf("Connection info: host=%s port=%s user=%s dbname=%s\n",
		host, port, user, dbname)

	// Configurar timeout para diagnóstico
	fmt.Println("Attempting connection with 10 second timeout...")
	connectionStart := time.Now()

	// Intentar conexión directa con el cliente Supabase
	client, err := config.NewSupabaseClient()
	if err != nil {
		connectionDuration := time.Since(connectionStart)
		fmt.Printf("Connection failed after %v\n", connectionDuration)
		log.Fatalf("Error connecting to Supabase: %v", err)
	}
	defer client.Close()

	connectionDuration := time.Since(connectionStart)
	fmt.Printf("Connection established in %v\n", connectionDuration)

	// Verificar que la conexión funcione con una consulta simple
	queryStart := time.Now()
	result, err := client.ExecuteQuery()
	if err != nil {
		log.Fatalf("Error executing test query: %v", err)
	}
	queryDuration := time.Since(queryStart)

	fmt.Printf("Query successful in %v! Test query result: %d\n", queryDuration, result)
	fmt.Printf("Total connection test completed in %v\n", time.Since(connectionStart))
}
