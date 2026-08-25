# Когда интерфейсы -- плохие абстракции

## 1. LSP-нарушение / даункаст: `database/sql/driver`

В стандартной библиотеке Go (`database/sql/driver`) есть базовый интерфейс `driver.Conn` с обязательными методами (`Prepare`, `Close`, `Begin`), и рядом -- набор "опциональных" интерфейсов: `driver.Queryer`, `driver.QueryerContext`, `driver.ExecerContext`, `driver.ConnPrepareContext`, `driver.ConnBeginTx` и т.д.

```go
// упрощённо из src/database/sql/sql.go
if queryer, ok := dc.ci.(driver.QueryerContext); ok {
    // быстрый путь: драйвер поддерживает QueryContext нативно
} else if queryer, ok := dc.ci.(driver.Queryer); ok {
    // fallback на старый Queryer
} else {
    // fallback на Prepare+Exec
}
```

`driver.Conn` формально задаёт "контракт" соединения с БД, но реальное поведение зависит от того, какие ещё интерфейсы реализует конкретный драйвер (pq, mysql, sqlite3...). Пользователь стандартной библиотеки вынужден делать `type assertion`, чтобы понять, с чем он реально работает.

## 2. Заголовочный интерфейс: mockery/wire-генерируемые `XService` в enterprise Go

Очень частый паттерн в реальных Go-репозиториях (Uber, многие внутренние сервисы на github, шаблоны типа `golang-standards/project-layout`): для каждого `struct`, у которого нужен мок в тестах, через `mockery`/`moq` механически генерируется интерфейс:

```go
//go:generate mockery --name=OrderService
type OrderService interface {
    CreateOrder(ctx context.Context, in CreateOrderInput) (*Order, error)
    GetOrder(ctx context.Context, id string) (*Order, error)
    CancelOrder(ctx context.Context, id string) error
}

type orderService struct{ repo OrderRepository }
```

Это ровно "Extract Interface": единственная продакшн-реализация -- `orderService`, интерфейс существует только ради DI/мокирования в юнит-тестах, а не потому что появилась вторая содержательная реализация.

## 3. Поверхностный интерфейс: типизированные клиенты `k8s.io/client-go`

`kubernetes.Interface` из `client-go` -- формально интерфейс:

```go
type Interface interface {
    CoreV1() corev1.CoreV1Interface
    AppsV1() appsv1.AppsV1Interface
    // ... десятки других *V1Interface
}

type PodInterface interface {
    Create(ctx context.Context, pod *v1.Pod, opts metav1.CreateOptions) (*v1.Pod, error)
    Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Pod, error)
    // ...
}
```

Формы вызовов абстрагированы интерфейсом, но сигнатуры целиком построены на конкретных типах Kubernetes API (`*v1.Pod`, `metav1.CreateOptions`), сгенерированных из OpenAPI-схемы конкретного кластера. Единственная альтернативная реализация -- `client-go/kubernetes/fake`, которая эмулирует тот же самый Kubernetes API in-memory, а не какой-то принципиально другой backend. Заменить "поставщика" на что-то радикально иное (не Kubernetes) невозможно без переписывания всего кода вызывающей стороны -- классический признак поверхностного извлечения интерфейса.

## 4. Протекающая абстракция: `http.ResponseWriter` + `http.Hijacker`/`http.Flusher`

`net/http.ResponseWriter` -- нарочно минимальный интерфейс (`Header`, `Write`, `WriteHeader`). Но как только нужен SSE (Server-Sent Events) или апгрейд до WebSocket, приходится вскрывать эту "абстракцию" даункастом до платформо- и транспортно-специфичных интерфейсов:

```go
func sseHandler(w http.ResponseWriter, r *http.Request) {
    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "streaming unsupported", http.StatusInternalServerError)
        return
    }
    // ...
    flusher.Flush()
}
```

Библиотека `gorilla/websocket` в функции `Upgrade` делает то же самое с `http.Hijacker`:

```go
hj, ok := w.(http.Hijacker)
if !ok {
    return nil, errors.New("webserver doesn't support hijacking")
}
conn, brw, err := hj.Hijack()
```

`ResponseWriter` как "абстракция HTTP-ответа" протекает: реальное поведение (можно ли захватить сырое TCP-соединение, можно ли принудительно сбросить буфер) зависит от конкретной реализации `http.Server`, и об этом знании пользователь интерфейса узнаёт только через ошибку даункаста в рантайме.

## 5. Фат-интерфейс с "не поддерживается": `afero.Fs`

`github.com/spf13/afero` определяет единый интерфейс файловой системы:

```go
type Fs interface {
    Create(name string) (File, error)
    Mkdir(name string, perm os.FileMode) error
    MkdirAll(path string, perm os.FileMode) error
    Open(name string) (File, error)
    OpenFile(name string, flag int, perm os.FileMode) (File, error)
    Remove(name string) error
    RemoveAll(path string) error
    Rename(oldname, newname string) error
    // ...
}
```

И тут же в комплекте идёт `afero.ReadOnlyFs`, оборачивающий любую другую `Fs` и на все write-операции возвращающий `syscall.EPERM`:

```go
func (r *ReadOnlyFs) Remove(n string) error {
    return syscall.EPERM
}

func (r *ReadOnlyFs) Rename(o, n string) error {
    return syscall.EPERM
}
```

Интерфейс `Fs` совмещает чтение и запись (нарушение ISP), поэтому read-only реализация вынуждена симулировать поддержку методов, которые не поддерживает, через ошибку вместо честного отсутствия метода в типе.