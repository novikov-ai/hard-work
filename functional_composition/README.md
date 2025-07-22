# Применяем функциональную композицию правильно

### Пример 1: Валидация данных  
**Проблема**: Глубокие иерархии валидаторов через наследование  
**Решение**: Композиция простых предикатов  

~~~csharp
// до (ООП с наследованием)
public abstract class Validator {
    public abstract bool Validate(string value);
}

public class EmailValidator : Validator {
    public override bool Validate(string value) 
        => value.Contains("@");
}

public class CompositeValidator : Validator {
    private List<Validator> _validators = new();
    public void Add(Validator v) => _validators.Add(v);
    
    public override bool Validate(string value) 
        => _validators.All(v => v.Validate(value));
}

// Использование:
var validator = new CompositeValidator();
validator.Add(new EmailValidator());
bool isValid = validator.Validate("test@example.com");
~~~

~~~csharp
// после (функциональная композиция)
// Плоские функции-предикаты
bool IsNotEmpty(string s) => !string.IsNullOrEmpty(s);
bool IsEmail(string s) => s.Contains("@");
bool IsStrongPassword(string s) => s.Length >= 8;

// Композиция через лямбда-функцию
Func<string, bool> ValidateUserInput = s => 
    IsNotEmpty(s) && 
    IsEmail(s) && 
    IsStrongPassword(s);

// Использование:
bool isValid = ValidateUserInput("test@example.com");
~~~

**Комментарии**:  
- Убрали 3 класса + иерархию наследования  
- Зависимости стали плоскими (нет вложенных объектов)  
- Нет tight coupling: функции независимы и переиспользуемы  

---

### Пример 2: Обработка заказов  
**Проблема**: Шаблонный метод в базовом классе  
**Решение**: Последовательная композиция функций  

~~~csharp
// до (Template Method pattern)
public abstract class OrderProcessor {
    public void Process(Order order) {
        Validate(order);
        CalculateTotal(order);
        ProcessPayment(order);
    }
    protected abstract void Validate(Order order);
    protected abstract void CalculateTotal(Order order);
    protected abstract void ProcessPayment(Order order);
}
~~~

~~~csharp
// после (функциональный pipeline)
// Чистые функции обработки
Order Validate(Order o) { /*...*/ return o; }
Order CalculateTotal(Order o) { /*...*/ return o; }
Order ApplyDiscount(Order o) { /*...*/ return o; }
Order ProcessPayment(Order o) { /*...*/ return o; }

// Композиция через вызовы
Order ProcessOrder(Order order) => 
    ProcessPayment(
        ApplyDiscount(
            CalculateTotal(
                Validate(order))));

// Добавление нового шага без изменений
Order AddGiftWrap(Order o) { /*...*/ return o; }

Order ProcessOrderWithGift(Order order) => 
    ProcessOrder(AddGiftWrap(order));
~~~

**Комментарии**:  
- Убрали абстрактный класс и принудительное наследование  
- Каждый шаг - самостоятельная тестируемая функция  
- Новые требования = новые комбинации, а не изменение иерархии  

---

### Пример 3: Трансформация данных  
**Проблема**: Декораторы для добавления поведения  
**Решение**: Цепочка преобразований через композицию  

~~~csharp
// до (Паттерн Декоратор)
public interface ITextTransformer {
    string Transform(string text);
}

public class UpperCaseTransformer : ITextTransformer {
    private ITextTransformer _inner;
    public UpperCaseTransformer(ITextTransformer inner) => _inner = inner;
    public string Transform(string text) => _inner.Transform(text).ToUpper();
}
~~~

~~~csharp
// после (композиция функций)
// Базовые трансформации
string ToUpper(string s) => s.ToUpper();
string AddExclamation(string s) => s + "!!!";
string AddGreeting(string s) => $"Hello, {s}";

// Композиция через агрегацию
Func<string, string> ProcessText = s => 
    AddExclamation(ToUpper(AddGreeting(s)));

// Использование:
var result = ProcessText("World"); // "HELLO, WORLD!!!"
~~~

**Комментарии**:  
- Убрали 4 класса (интерфейс + реализации)  
- Зависимости визуализируются как поток данных  
- Легко менять порядок операций: `ToUpper(AddGreeting(s))`  

---

### Пример 4: Фильтрация данных  
**Проблема**: Стратегии через наследование  
**Решение**: Комбинация предикатов  

~~~csharp
// до (Паттерн Стратегия)
public abstract class FilterStrategy {
    public abstract bool IsMatch(Product p);
}

public class PriceFilter : FilterStrategy {
    public override bool IsMatch(Product p) => p.Price < 100;
}
~~~

~~~csharp
// после (композиция условий)
// Чистые функции-фильтры
bool IsUnderPrice(Product p, decimal max) => p.Price < max;
bool InCategory(Product p, string category) => p.Category == category;
bool IsInStock(Product p) => p.StockCount > 0;

// Динамическая композиция фильтров
Func<Product, bool> BuildFilter(
    decimal maxPrice, 
    string category
) => p => 
    IsUnderPrice(p, maxPrice) && 
    InCategory(p, category) && 
    IsInStock(p);

// Использование:
var filter = BuildFilter(maxPrice: 100, category: "Electronics");
var results = products.Where(filter);
~~~

**Комментарии**:  
- Нет классов-стратегий и сложных зависимостей  
- Параметры фильтрации явные и контролируемые  
- Функции можно переиспользовать в других комбинациях