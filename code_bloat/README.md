# Про раздутость кода

### Пример 1

Сущестовал метод, который подготавливал поля для отрисовки страницы. При этом процессинг каждого из полей происходил в цикле и, если появлялась необходимость что-то поменять в каком-либо поле, то добавлялся еще один if внутри `Process`, который изменял его. 

Из-за такого подхода `Process` очень сильно раздувался, так как изменять поля нужно, но как обращаться к каждому непосредственно без итерации по всем в текущей архитектуре - непонятно. 

Было:
~~~go
func Prepare(){
    // ...

    for _, ff := range data.Fields() {
		if field, ok := ff.(models.Field); ok {
			Process(ctx, field, data)
		}
	}
}

func Process(ctx context.Context, field models.Field, data models.Data){
    if field.ID() == carID{
        updateWithCarInfo(data)
    }

    if field.Name == "cart"{
        DoRedesignV2(data)
    }

    // ...
}
~~~

Стало:
~~~go
func Prepare(){
    // ...

    // data.Fields(): []interface{} => map[int64]interface{}
    
    for ff := range data.Fields() {
		if field, ok := ff.(models.Field); ok {
			Process(ctx, field, data)
		}
	}

    modifyFieldByID(carID, data, updateWithCarInfo)
    modifyFieldByID(servicesModelID, data, DoRedesignV2)
}

func Process(ctx context.Context, field models.Field, data models.Data){
    // ...
}

func modifyFieldByID(id int64, data models.Data, modifier func(data models.Data)){
    field := data[id]
    value, ok := field.(models.Field)
    if !ok {
        return
    }

    modifier(data)
}
~~~

### Пример 2 

Каждое поле конфигурации вручную загружалось и валидировалось, что приводило к избыточности.

Было:
~~~go
type Config struct {
	Endpoint string
	Timeout  int
	Retry    int
}

func loadConfig() Config {
	config := Config{}
	endpoint, err := loadString("ENDPOINT")
	if err != nil {
		log.Fatal(err)
	}
	config.Endpoint = endpoint

	timeout, err := loadInt("TIMEOUT")
	if err != nil {
		log.Fatal(err)
	}
	config.Timeout = timeout

	retry, err := loadInt("RETRY")
	if err != nil {
		log.Fatal(err)
	}
	config.Retry = retry
	return config
}
~~~

Стало:
~~~go
func loadConfig() Config {
	var config Config
	err := loadEnvConfig(&config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
~~~

### Пример 3 

В проекте существовало множество условных выражений, возвращающие различные типы бэкендов (это были строгие структуры), что не позволяло гибко работать с бэкендами и раздувало код при большом количестве типов.

Было:
~~~go
if backendType == "cassandra" {
	return newCassandraBackend()
} else if backendType == "elasticsearch" {
	return newElasticBackend()
}

~~~

Стало:
~~~go
type Backend interface {
	Init() error
}

func getBackend(backendType string) Backend {
	switch backendType {
	case "cassandra":
		return newCassandraBackend()
	case "elasticsearch":
		return newElasticBackend()
	default:
		return nil
	}
}
~~~