# Прекратите вызывать throw/raise

В Go нет исключений как таковых, но есть паники, которые можно обрабатывать, поэтому примеры будут на них.

## Пример 1

https://github.com/go-task/task/blob/1f34895185fe325c464e6f5e0d80b97c34f4f00c/taskfile/snippet.go#L26

Пример:
~~~go
func init() {
	r, err := embedded.Open("themes/task.xml")
	if err != nil {
		panic(err)
	}
	style, err := chroma.NewXMLStyle(r)
	if err != nil {
		panic(err)
	}
	styles.Register(style)
}
~~~

Без исключения:
~~~go
func init() error {
	r, err := embedded.Open("themes/task.xml")
	if err != nil {
		return err
	}
	style, err := chroma.NewXMLStyle(r)
	if err != nil {
		return err
	}
	styles.Register(style)
}
~~~

Вместо паники возвращаем ошибку, которую можем корректно обработать и завершить процесс.

## Пример 2

https://github.com/go-task/task/blob/1f34895185fe325c464e6f5e0d80b97c34f4f00c/watch_test.go#L69

Пример:
~~~go
func TestFileWatch(t *testing.T) {
	t.Parallel()

	const dir = "testdata/watch"
	_ = os.RemoveAll(filepathext.SmartJoin(dir, ".task"))
	_ = os.RemoveAll(filepathext.SmartJoin(dir, "src"))

	expectedOutput := strings.TrimSpace(`
task: Started watching for tasks: default
task: [default] echo "Task running!"
Task running!
task: task "default" finished running
task: [default] echo "Task running!"
Task running!
task: task "default" finished running
	`)

	var buff bytes.Buffer
	e := task.NewExecutor(
		task.WithDir(dir),
		task.WithStdout(&buff),
		task.WithStderr(&buff),
		task.WithWatch(true),
	)

	require.NoError(t, e.Setup())
	buff.Reset()

	dirPath := filepathext.SmartJoin(dir, "src")
	filePath := filepathext.SmartJoin(dirPath, "a")

	err := os.MkdirAll(dirPath, 0o755)
	require.NoError(t, err)

	err = os.WriteFile(filePath, []byte("test"), 0o644)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := e.Run(ctx, &task.Call{Task: "default"})
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	time.Sleep(200 * time.Millisecond)
	err = os.WriteFile(filePath, []byte("test updated"), 0o644)
	require.NoError(t, err)

	time.Sleep(200 * time.Millisecond)
	cancel()
	assert.Equal(t, expectedOutput, strings.TrimSpace(buff.String()))
}
~~~

Без исключения:
~~~go
func TestFileWatch(t *testing.T) {
// ...

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := e.Run(ctx, &task.Call{Task: "default"})
				if err != nil {
					fmt.Fprintln(os.Stderr, err.Error())
					os.Exit(1)
				}
			}
		}
	}()

// ...
}
~~~

Заменили панику корректным кодом ошибки и добавили логирование.


## Пример 3

https://github.com/binh234/go-project/blob/5b371629e98ea2e18f23c22fa8500360d116192d/go-fiber-crm/main.go#L25

Пример:
~~~go
func initDatabase() {
	fmt.Println(database.SayHello())
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "leads.db")
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Connection opened to database")
	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database Migrated")
}
~~~

Без исключения:
~~~go
func initDatabase() {
	fmt.Println(database.SayHello())
	var err error
	database.DBConn, err = gorm.Open("sqlite3", "leads.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to connect database")
		os.Exit(1)
	}
	fmt.Println("Connection opened to database")
	database.DBConn.AutoMigrate(&lead.Lead{})
	fmt.Println("Database Migrated")
}
~~~

Заменили панику корректным кодом ошибки и добавили логирование.

## Пример 4

https://github.com/Onelvay/go-pet-project/blob/da5e4b7579ec5c8e5b9da576971f0e8085091ad6/main.go#L40

Пример:
~~~go
func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	config := postgres.NewConfig(viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.dbname"),
		viper.GetString("db.user"),
		viper.GetString("db.pass"),
	)
	client.InitConst(viper.GetString("payment.merchantId"), viper.GetString("payment.merchantPassword"), viper.GetString("payment.checkoutUrl"))

	mongoProductDb := mongoDb.MongoProductCollection(viper.GetString("mongoDB.host"))
	postgresDb := postgres.NewPostgresDb(*config)

	redis, err := redisClient.InitRedis(viper.GetString("redis.host"), viper.GetString("redis.password"))
	if err != nil {
		panic(err)
	}

	productDb, userDb, tokenDb, orderDb := initDbControllers(postgresDb, redis, mongoProductDb)
	hasher := service.NewHasher(viper.GetString("app.hash"))
	userContr := contr.NewUserController(userDb, tokenDb, hasher, orderDb)
	handlers := contr.NewHandlers(productDb, &userContr, orderDb, tokenDb, userDb)
	router := routes.InitRoutes(handlers)

	var PORT string
	if PORT = os.Getenv("PORT"); PORT == "" {
		PORT = "8080"
	}
	err = http.ListenAndServe(":"+PORT, router)
	if err != nil {
		fmt.Println(err.Error())
	}

}
~~~

Без исключения:
~~~go
// ...
	redis, err := redisClient.InitRedis(viper.GetString("redis.host"), viper.GetString("redis.password"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
// ...
~~~

Заменили панику корректным кодом ошибки и добавили логирование.

## Пример 5

https://github.com/Onelvay/go-pet-project/blob/da5e4b7579ec5c8e5b9da576971f0e8085091ad6/db/postgres/db.go#L33

Пример:
~~~go
func NewPostgresDb(cfg Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.HOST, cfg.PORT, cfg.USER, cfg.DB_NAME, cfg.PASS)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&domain.User{})

	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&domain.User{}, &domain.Refresh_token{})

	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&domain.User{}, &domain.Order{})

	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&domain.Order{}, &domain.FinalResponse{})

	if err != nil {
		panic(err)
	}

	return db
}
~~~

Без исключения:
~~~go
func NewPostgresDb(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.HOST, cfg.PORT, cfg.USER, cfg.DB_NAME, cfg.PASS)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&domain.User{})

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&domain.User{}, &domain.Refresh_token{})

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&domain.User{}, &domain.Order{})

	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&domain.Order{}, &domain.FinalResponse{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
~~~

Заменили панику ошибкой, чтобы корректно обработать на уровне выше.