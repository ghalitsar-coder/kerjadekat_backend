package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "kerjadekat/backend/config"
    "kerjadekat/backend/internal/domain"
    agenthttp "kerjadekat/backend/internal/agent/delivery/http"
    agentrepo "kerjadekat/backend/internal/agent/repository"
    agentusecase "kerjadekat/backend/internal/agent/usecase"
    authusecase "kerjadekat/backend/internal/auth/usecase"
    "kerjadekat/backend/internal/httpapi"
    kelurahanrepo "kerjadekat/backend/internal/kelurahan/repository"
    kelurahanusecase "kerjadekat/backend/internal/kelurahan/usecase"
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
    walletrepo "kerjadekat/backend/internal/wallet/repository"
    workerrepo "kerjadekat/backend/internal/worker/repository"
    workerusecase "kerjadekat/backend/internal/worker/usecase"
    "kerjadekat/backend/pkg/database"
    "kerjadekat/backend/pkg/mq"
    "kerjadekat/backend/pkg/otp"
    "kerjadekat/backend/pkg/redisx"
    "kerjadekat/backend/pkg/token"

    "github.com/gin-contrib/cors"
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

    if err := db.AutoMigrate(
        &domain.Kelurahan{},
        &domain.SkillCategory{},
        &domain.User{},
        &domain.WorkerProfile{},
        &domain.WorkerSkill{},
        &domain.AgentTerritory{},
        &domain.Order{},
        &domain.OrderStatusLog{},
        &domain.OrderRating{},
        &domain.IncomeRecord{},
        &domain.Wallet{},
        &domain.WalletTransaction{},
    ); err != nil {
        log.Fatalf("migrate: %v", err)
    }
    log.Println("database migration completed")

    rdb := redisx.NewClient(cfg)
    if err := redisx.Ping(ctx, rdb); err != nil {
        log.Fatalf("redis: %v", err)
    }

    mqClient, err := mq.Connect(cfg.RabbitMQURL)
    if err != nil {
        log.Fatalf("rabbitmq: %v", err)
    }

    userRepository := userrepo.NewUserPostgres(db)
    workerRepository := workerrepo.NewWorkerPostgres(db)
    orderRepository := orderrepo.NewOrderPostgres(db)
    skillRepository := skillrepo.NewSkillPostgres(db)
    walletRepository := walletrepo.NewWalletPostgres(db)
    kelurahanRepository := kelurahanrepo.NewKelurahanPostgres(db)

    locationTTL := time.Duration(cfg.WorkerLocationTTLMinutes) * time.Minute
    presence := location.NewRedisPresence(rdb, locationTTL)

    offerTTL := time.Duration(cfg.MatchMaxRounds*cfg.MatchTimerSeconds) * time.Second
    if offerTTL <= 0 {
        offerTTL = 30 * time.Minute
    }
    offerTracker := orderpkg.NewRedisOfferTracker(rdb, offerTTL)

    redisPub := orderpkg.NewPublisher(rdb)
    matchPublisher := orderpkg.NewMatchPublisher(redisPub, mqClient)

    pay := payment.NewFromConfig(cfg)
    smsNotifier := sms.NewMockNotifier(cfg.OTPLogFile)
    otpStore := otp.NewRedis(rdb, 5*time.Minute)
    issuer := token.NewIssuer(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)

    authUC := authusecase.NewAuth(userRepository, workerRepository, otpStore, smsNotifier, issuer)
    usersUC := userusecase.NewUsers(userRepository)
    skillsUC := skillusecase.NewSkills(skillRepository)
    kelurahansUC := kelurahanusecase.NewKelurahans(kelurahanRepository)
    ordersUC := orderusecase.NewOrders(
        orderRepository,
        skillRepository,
        pay,
        presence,
        matchPublisher,
        offerTracker,
        workerRepository,
        walletRepository,
        orderusecase.MatchSettings{
            RadiusMeters: float64(cfg.MatchRadiusMeters),
            TimerDelayMs: cfg.MatchTimerSeconds * 1000,
            MaxRounds:    cfg.MatchMaxRounds,
        },
    )
    workersUC := workerusecase.NewWorkers(workerRepository, presence, float64(cfg.MatchRadiusMeters))
    if err := workersUC.BootstrapPresence(ctx); err != nil {
        log.Printf("worker presence bootstrap: %v", err)
    } else {
        log.Println("worker presence bootstrap: synced online workers to Redis")
    }

    agentRepository := agentrepo.NewAgentPostgres(db)

    var fileStore domain.FileStorage
    if cfg.StorageBackend == "s3" {
        s3Store, err := storage.NewS3(cfg.S3Endpoint, cfg.S3AccessKey, cfg.S3SecretKey, cfg.S3Region, cfg.S3UseSSL)
        if err != nil {
            log.Fatalf("s3 storage: %v", err)
        }
        fileStore = s3Store
        log.Println("storage backend: S3")
    } else {
        fileStore = storage.NewMock()
        log.Println("storage backend: mock (files not persisted)")
    }

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
    r.SetTrustedProxies(nil)
    r.Use(gin.Logger(), gin.Recovery())

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", "http://kerjadekat.local"},
        AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           86400,
    }))

    httpapi.Mount(r, httpapi.Deps{
		Auth:                authUC,
		Users:               usersUC,
		Skills:              skillsUC,
		Kelurahans:          kelurahansUC,
		Orders:              ordersUC,
		Workers:             workersUC,
		Agents:              agentHandler,
		Tokens:              issuer,
		WSHub:               wsHub,
		FileStorage:         fileStore,
		XenditCallbackToken: cfg.XenditCallbackToken,
	})

    addr := ":" + cfg.ServerPort
    srv := &http.Server{
        Addr:    addr,
        Handler: r,
    }
    go func() {
        log.Printf("KerjaDekat API listening on %s (ws=GET /api/v1/ws, auth via first frame)", addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("listen: %v", err)
        }
    }()

    // Menahan (block) eksekusi sampai ada sinyal mati (Ctrl+C / SIGTERM)
    <-ctx.Done()
    log.Println("Sinyal shutdown diterima. Mulai mematikan server...")

    // Beri waktu maksimal 5 detik untuk menyelesaikan HTTP request yang sedang berjalan
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Matikan HTTP server (berhenti menerima request baru)
    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Fatalf("Server dipaksa mati: %v", err)
    }

    log.Println("Menutup koneksi ke layanan eksternal...")

    // 1. Tutup koneksi RabbitMQ
    if err := mqClient.Close(); err != nil {
        log.Printf("Error menutup RabbitMQ: %v", err)
    } else {
        log.Println("RabbitMQ tertutup.")
    }

    // 2. Tutup koneksi Redis
    if err := rdb.Close(); err != nil {
        log.Printf("Error menutup Redis: %v", err)
    } else {
        log.Println("Redis tertutup.")
    }

    // 3. Tutup koneksi Postgres (GORM)
    if sqlDB, err := db.DB(); err == nil {
        if err := sqlDB.Close(); err != nil {
            log.Printf("Error menutup Postgres: %v", err)
        } else {
            log.Println("Postgres tertutup.")
        }
    } else {
        log.Printf("Gagal mendapatkan instance SQL dari GORM: %v", err)
    }

    log.Println("Server berhasil dimatikan dengan aman.")
}