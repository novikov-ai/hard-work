# Формализуем понятие надёжности системы

## Ключевые свойства системы

### **Свойство 1: Доступность (Availability)**
- **Узкий диапазон**: 99.99% (≤ 5 минут простоя в месяц).
- **Широкий диапазон**: 95% (≤ 36 часов простоя в месяц при экстремальных сбоях).
- **Обоснование**: Система должна автоматически восстанавливаться после сбоев (например, перезапуск подов в Kubernetes).

### **Свойство 2: Время отклика API (Latency)**
- **Узкий диапазон**: 50–200 мс (для 95% запросов).
- **Широкий диапазон**: ≤ 500 мс (для 99% запросов при пиковой нагрузке).
- **Предусловие**: Нет DDoS-атак.  
- **Обоснование**: После пиковой нагрузки система должна стабилизироваться.

### **Свойство 3: Целостность данных (Data Integrity)**
- **Узкий диапазон**: 0 потерь данных (все транзакции подтверждены).
- **Широкий диапазон**: ≤ 0.001% потерянных данных (при крахе кластера).
- **Обоснование**: Потеря данных недопустима даже в экстремальных условиях (репликация в реальном времени).

### **Свойство 4: Пропускная способность (Throughput)**
- **Узкий диапазон**: 10,000 RPS (штатная нагрузка).
- **Широкий диапазон**: 5,000 RPS (при аварийном режиме).
- **Обоснование**: После снижения нагрузки до 5k RPS система должна вернуться в штатный режим.

### **Свойство 5: Скорость восстановления (Recovery Time)**
- **Узкий диапазон**: ≤ 1 минута (для автоматического восстановления).
- **Широкий диапазон**: ≤ 30 минут (при ручном вмешательстве).
- **Обоснование**: Сейчас восстановление возможно только вручную, но планируется автоматизация.

## Усиление инвариантов

### Доступность (RESILIENT → STRONG)
- **Проблема**: Сейчас при потере дата-центра доступность падает до 95%.
- **Решение**:  
  - Реализовать multi-region deployment с автоматическим переключением трафика (на Go через Consul API).  
  ```go
  func switchTrafficToBackupRegion() {
      consulClient := consul.NewClient()
      err := consulClient.UpdateServiceWeights("primary", 0) // Отключаем primary
      if err != nil {
          log.Fatal("Failed to switch traffic:", err)
      }
      log.Println("Traffic switched to backup region")
  }
  ```

### Скорость восстановления (WEAK → RESILIENT)
- **Проблема**: Восстановление занимает до 30 минут.
- **Решение**:  
  - Добавить health-checks и автоматический перезапуск подов.  
  ```go
  func startHealthCheck(ctx context.Context) {
      ticker := time.NewTicker(10 * time.Second)
      for {
          select {
          case <-ticker.C:
              if isUnhealthy() {
                  restartPod() // Используем Kubernetes API
              }
          case <-ctx.Done():
              return
          }
      }
  }
  ```

### Целостность данных (STRONG → с метастабильностью)
- **Проблема**: Репликация замедляется при высокой нагрузке.
- **Решение**:  
  - Добавить асинхронную очередь для гарантированной доставки данных (на основе Kafka + Go-воркеров).  
  ```go
  func asyncReplicate(data []byte) {
      producer := kafka.NewProducer()
      err := producer.Send("replication-topic", data)
      if err != nil {
          storeInDlq(data) // Dead Letter Queue для ручного восстановления
      }
  }
  ```

## Скорость восстановления как характеристика

### **Как добавить:**
1. **Интеграция с мониторингом**:  
   - Использовать Prometheus для отслеживания MTTR (Mean Time To Recovery).  
   ```go
   func trackRecoveryTime(start time.Time) {
       duration := time.Since(start)
       prometheus.RecoveryTime.Observe(duration.Seconds())
   }
   ```

2. **Автоматизация через Kubernetes Operators**:  
   - Написать оператор на Go, который реагирует на события (например, `CrashLoopBackOff`) и перезапускает зависимости.  
   ```go
   func handlePodEvent(event k8s.Event) {
       if event.Reason == "CrashLoopBackOff" {
           restartDependentServices(event.Pod)
       }
   }
   ```

3. **Тесты на скорость восстановления**:  
   - Реализовать chaos-инжиниринг через Gremlin API.  
   ```go
   func simulateNetworkPartition() {
       gremlinClient := gremlin.NewClient()
       err := gremlinClient.RunExperiment("network-partition")
       if err != nil {
           log.Fatal("Chaos experiment failed:", err)
       }
   }
   ```