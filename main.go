package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type App struct {
	RedisClient        *redis.Client
	HttpClient         *http.Client
	SqsSvc             *sqs.SQS
	SqsQueueURL        string
	FlagServiceURL     string
	TargetingServiceURL string
}

func main() {
	// ============================
	// 1. Variáveis de ambiente
	// ============================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8005"
	}

	flagURL := os.Getenv("FLAG_SERVICE_URL")
	targetURL := os.Getenv("TARGETING_SERVICE_URL")

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL deve ser definida")
	}

	sqsURL := os.Getenv("SQS_QUEUE_URL")
	awsRegion := os.Getenv("AWS_REGION")
	awsKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsToken := os.Getenv("AWS_SESSION_TOKEN")

	// ============================
	// 2. Iniciar Redis
	// ============================
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Erro no Redis: %v", err)
	}

	redisClient := redis.NewClient(opt)
	log.Println("Redis conectado.")

	// ============================
	// 3. Iniciar SQS (ElasticMQ)
	// ============================
	var sqsSvc *sqs.SQS

	if sqsURL != "" && awsKey != "" && awsSecret != "" {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(awsRegion),
			Endpoint:    aws.String("http://sqs-local:9324"),
			Credentials: credentials.NewStaticCredentials(awsKey, awsSecret, awsToken),
		})

		if err != nil {
			log.Fatalf("Erro iniciando sessão AWS: %v", err)
		}

		sqsSvc = sqs.New(sess)
		log.Println("Cliente SQS inicializado.")
	} else {
		log.Println("⚠️ SQS DESABILITADO — variáveis de ambiente ausentes.")
	}

	// ============================
	// 4. Criar App
	// ============================
	app := &App{
		RedisClient:        redisClient,
		HttpClient:         &http.Client{Timeout: 5 * time.Second},
		SqsSvc:             sqsSvc,
		SqsQueueURL:        sqsURL,
		FlagServiceURL:     flagURL,
		TargetingServiceURL: targetURL,
	}

	// ============================
	// 5. Rotas HTTP
	// ============================
	http.HandleFunc("/health", app.healthHandler)
	http.HandleFunc("/evaluate", app.evaluationHandler)

	log.Printf("Evaluation-service rodando na porta %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
