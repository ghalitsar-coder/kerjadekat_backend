package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kerjadekat/backend/config"
	agentrepo "kerjadekat/backend/internal/agent/repository"
	agenthttp "kerjadekat/backend/internal/agent/delivery/http"
	agentusecase "kerjadekat/backend/internal/agent/usecase"
	authusecase "kerjadekat/backend/internal/auth/usecase"
	"kerjadekat/backend/internal/httpapi"
	"kerjadekat/backend/internal/location"
	orderpkg "kerjadekat/backend/internal/order"
	"kerjadekat/backend/internal/order/delivery/ws"
	orderrepo "kerjadekat/backend/internal/order/repository"
	orderusecase "kerjadekat/backend/internal/order/usecase"
	"kerjadekat/backend/internal/platform/ocr"
	"kerjadekat/backend/internal/platform/payment"
	"kerjadekat/backend/internal/platform/sms"
	"kerjadekat/backend/internal/platform/storage"
	skillrepo "kerjadekat/backend/internal/skill/repository"
	skillusecase "kerjadekat/backend/internal/skill/usecase"
	userrepo "kerjadekat/backend/internal/user/repository"
	userusecase "kerjadekat/backend/internal/user/usecase"
	workerrepo "kerjadekat/backend/internal/worker/repository"
	workerusecase "kerjadekat/backend/internal/worker/usecase"
	"kerjadekat/backend/pkg/database"
	"kerjadekat/backend/pkg/mq"
	"kerjadekat/backend/pkg/otp"
	"kerjadekat/backend/pkg/redisx"
	"kerjadekat/backend/pkg/token"

	"github.com/gin-gonic/gin"
)

func main() {
	cfgPath := ""
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := database.OpenPostgres(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	rdb := redisx.NewClient(cfg)
	if err := redisx.Ping(ctx, rdb); err != nil {
		log.Fatalf("redis: %v", err)
	}

	mqClient, err := mq.Connect(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("rabbitmq: %v", err)
	}
	defer mqClient.Close()

	userRepository := userrepo.NewUserPostgres(db)
	workerRepository := workerrepo.NewWorkerPostgres(db)
	orderRepository := orderrepo.NewOrderPostgres(db)
	skillRepository := skillrepo.NewSkillPostgres(db)

	locationTTL := time.Duration(cfg.WorkerLocationTTLMinutes) * time.Minute
	presence := location.NewRedisPresence(rdb, locationTTL)

	offerTTL := time.Duration(cfg.MatchMaxRounds*cfg.MatchTimerSeconds) * time.Second
	if offerTTL <= 0 {
		offerTTL = 30 * time.Minute
	}
	offerTracker := orderpkg.NewRedisOfferTracker(rdb, offerTTL)

	redisPub := orderpkg.NewPublisher(rdb)
	matchPublisher := orderpkg.NewMatchPublisher(redisPub, mqClient)

	pay := payment.NewMockGateway()
	smsNotifier := sms.NewMockNotifier()
	otpStore := otp.NewRedis(rdb, 5*time.Minute)
	issuer := token.NewIssuer(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)

	authUC := authusecase.NewAuth(userRepository, workerRepository, otpStore, smsNotifier, issuer)
	usersUC := userusecase.NewUsers(userRepository)
	skillsUC := skillusecase.NewSkills(skillRepository)
	ordersUC := orderusecase.NewOrders(
		orderRepository,
		skillRepository,
		pay,
		presence,
		matchPublisher,
		offerTracker,
		orderusecase.MatchSettings{
			RadiusMeters: float64(cfg.MatchRadiusMeters),
			TimerDelayMs: cfg.MatchTimerSeconds * 1000,
			MaxRounds:    cfg.MatchMaxRounds,
		},
	)
	workersUC := workerusecase.NewWorkers(workerRepository, presence)

	agentRepository := agentrepo.NewAgentPostgres(db)
	fileStore := storage.NewMock()
	ocrService := ocr.NewMock()
	agentsUC := agentusecase.NewAgents(agentRepository, userRepository, skillRepository, fileStore, ocrService)
	agentHandler := agenthttp.NewHandler(agentsUC)

	if err := mqClient.ConsumeMatchTimer(ctx, func(c context.Context, msg mq.MatchTimerMessage) error {
		return ordersUC.HandleMatchTimer(c, msg.OrderID, msg.Round)
	}); err != nil {
		log.Fatalf("rabbitmq consumer: %v", err)
	}

	wsHub := ws.NewHub(issuer, rdb)
	wsHub.RunSubscribe(ctx)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	httpapi.Mount(r, httpapi.Deps{
		Auth:    authUC,
		Users:   usersUC,
		Skills:  skillsUC,
		Orders:  ordersUC,
		Workers: workersUC,
		Agents:  agentHandler,
		Tokens:  issuer,
		WSHub:   wsHub,
	})

	addr := ":" + cfg.ServerPort
	log.Printf("KerjaDekat API listening on %s (ws=GET /api/v1/ws, auth via first frame)", addr)
	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")
}
