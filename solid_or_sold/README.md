# SOLID или SOLD

## Пример 1: Репозиторий с избыточными методами

**Было:**
```csharp
public interface IUserRepository
{
    User GetById(int id);
    User GetByEmail(string email);
    IEnumerable<User> GetAll();
    void Add(User user);
    void Update(User user);
    void Delete(int id);
    void ChangePassword(int userId, string newPassword);
    void UpdateLastLogin(int userId);
}

// Сервис только для чтения пользователей
public class UserQueryService
{
    private readonly IUserRepository _userRepository; // Зависит от ВСЕХ методов!
    
    public UserQueryService(IUserRepository userRepository)
    {
        _userRepository = userRepository;
    }
    
    public User GetUserProfile(int id) => _userRepository.GetById(id);
}
```

**Стало:**
```csharp
// Разделяем на читающие и пишущие интерфейсы
public interface IUserReader
{
    User GetById(int id);
    User GetByEmail(string email);
    IEnumerable<User> GetAll();
}

public interface IUserWriter
{
    void Add(User user);
    void Update(User user);
    void Delete(int id);
}

public interface IUserPasswordManager
{
    void ChangePassword(int userId, string newPassword);
    void UpdateLastLogin(int userId);
}

// Сервис зависит только от того, что использует
public class UserQueryService
{
    private readonly IUserReader _userReader; // Только чтение!
    
    public UserQueryService(IUserReader userReader)
    {
        _userReader = userReader;
    }
    
    public User GetUserProfile(int id) => _userReader.GetById(id);
}
```

## Пример 2: Сервис уведомлений с разными каналами

**Было:**
```csharp
public interface INotificationService
{
    void SendEmail(string to, string subject, string body);
    void SendSms(string phone, string message);
    void SendPush(string userId, string title, string message);
    void SendTelegram(string chatId, string message);
}

// Сервис только для email-уведомлений
public class OrderConfirmationService
{
    private readonly INotificationService _notificationService; // Зависит от всех каналов!
    
    public void SendOrderConfirmation(Order order)
    {
        _notificationService.SendEmail(order.CustomerEmail, "Order Confirmed", "...");
    }
}
```

**Стало:**
```csharp
public interface IEmailNotifier
{
    void SendEmail(string to, string subject, string body);
}

public interface ISmsNotifier
{
    void SendSms(string phone, string message);
}

public interface IPushNotifier
{
    void SendPush(string userId, string title, string message);
}

// Сервис зависит только от email
public class OrderConfirmationService
{
    private readonly IEmailNotifier _emailNotifier;
    
    public OrderConfirmationService(IEmailNotifier emailNotifier)
    {
        _emailNotifier = emailNotifier;
    }
    
    public void SendOrderConfirmation(Order order)
    {
        _emailNotifier.SendEmail(order.CustomerEmail, "Order Confirmed", "...");
    }
}
```

## Пример 3: Конфигурация с множеством настроек

**Было:**
```csharp
public interface IAppSettings
{
    string DatabaseConnectionString { get; }
    string ApiKey { get; }
    int CacheTimeoutMinutes { get; }
    string LogLevel { get; }
    string EmailFromAddress { get; }
    bool EnableFeatureX { get; }
}

// Класс использует только настройки кэша
public class CacheService
{
    private readonly IAppSettings _settings; // Видит ВСЕ настройки!
    
    public CacheService(IAppSettings settings)
    {
        _settings = settings;
    }
    
    public TimeSpan GetCacheTimeout() => TimeSpan.FromMinutes(_settings.CacheTimeoutMinutes);
}
```

**Стало:**
```csharp
public interface IDatabaseSettings
{
    string DatabaseConnectionString { get; }
}

public interface ICacheSettings
{
    int CacheTimeoutMinutes { get; }
}

public interface IEmailSettings
{
    string EmailFromAddress { get; }
    string ApiKey { get; }
}

// Класс зависит только от настроек кэша
public class CacheService
{
    private readonly ICacheSettings _cacheSettings;
    
    public CacheService(ICacheSettings cacheSettings)
    {
        _cacheSettings = cacheSettings;
    }
    
    public TimeSpan GetCacheTimeout() => TimeSpan.FromMinutes(_cacheSettings.CacheTimeoutMinutes);
}
```

## Пример 4: Сервис работы с файлами

**Было:**
```csharp
public interface IFileService
{
    Stream ReadFile(string path);
    void WriteFile(string path, Stream content);
    void DeleteFile(string path);
    bool FileExists(string path);
    IEnumerable<string> ListFiles(string directory);
    void CreateDirectory(string path);
}

// Компонент только для чтения файлов
public class TemplateLoader
{
    private readonly IFileService _fileService; // Может удалять файлы!
    
    public TemplateLoader(IFileService fileService)
    {
        _fileService = fileService;
    }
    
    public string LoadTemplate(string path)
    {
        using var stream = _fileService.ReadFile(path);
        return new StreamReader(stream).ReadToEnd();
    }
}
```

**Стало:**
```csharp
public interface IFileReader
{
    Stream ReadFile(string path);
    bool FileExists(string path);
}

public interface IFileWriter
{
    void WriteFile(string path, Stream content);
    void DeleteFile(string path);
    void CreateDirectory(string path);
}

public interface IFileLister
{
    IEnumerable<string> ListFiles(string directory);
}

// Безопасный компонент только для чтения
public class TemplateLoader
{
    private readonly IFileReader _fileReader; // Может только читать!
    
    public TemplateLoader(IFileReader fileReader)
    {
        _fileReader = fileReader;
    }
    
    public string LoadTemplate(string path)
    {
        using var stream = _fileReader.ReadFile(path);
        return new StreamReader(stream).ReadToEnd();
    }
}
```

## Пример 5: Конфликт ISP и IoC контейнера

**Было (проблема с IoC):**
```csharp
// Большой интерфейс
public interface IDataProcessor
{
    void ValidateData(string data);
    void TransformData(string data);
    void SaveData(string data);
    void LogOperation(string operation);
}

// Регистрация в IoC контейнере
services.AddScoped<IDataProcessor, DataProcessor>();

// Два разных сервиса используют разные методы
public class DataValidator
{
    private readonly IDataProcessor _processor; // Использует только ValidateData
    
    public DataValidator(IDataProcessor processor) => _processor = processor;
    
    public bool Validate(string data)
    {
        _processor.ValidateData(data);
        return true;
    }
}

public class DataSaver
{
    private readonly IDataProcessor _processor; // Использует только SaveData
    
    public DataSaver(IDataProcessor processor) => _processor = processor;
    
    public void Save(string data) => _processor.SaveData(data);
}
```

**Стало (решение конфликта):**
```csharp
// Разделяем интерфейсы
public interface IDataValidator
{
    void ValidateData(string data);
}

public interface IDataTransformer
{
    void TransformData(string data);
}

public interface IDataSaver
{
    void SaveData(string data);
}

public interface IOperationLogger
{
    void LogOperation(string operation);
}

// Класс реализует все интерфейсы
public class DataProcessor : IDataValidator, IDataTransformer, IDataSaver, IOperationLogger
{
    public void ValidateData(string data) { /* ... */ }
    public void TransformData(string data) { /* ... */ }
    public void SaveData(string data) { /* ... */ }
    public void LogOperation(string operation) { /* ... */ }
}

// Регистрируем отдельно в IoC
services.AddScoped<IDataValidator, DataProcessor>();
services.AddScoped<IDataTransformer, DataProcessor>();
services.AddScoped<IDataSaver, DataProcessor>();
services.AddScoped<IOperationLogger, DataProcessor>();

// Теперь сервисы зависят только от нужного
public class DataValidator
{
    private readonly IDataValidator _validator; // Только валидация!
    
    public DataValidator(IDataValidator validator) => _validator = validator;
    
    public bool Validate(string data)
    {
        _validator.ValidateData(data);
        return true;
    }
}

public class DataSaver
{
    private readonly IDataSaver _saver; // Только сохранение!
    
    public DataSaver(IDataSaver saver) => _saver = saver;
    
    public void Save(string data) => _saver.SaveData(data);
}
```