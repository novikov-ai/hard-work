# Группировка в функциях и файлах

~~~go
// 1. User Processing
func ProcessUserProfile(user *User) error {
    // ------------------------- VALIDATION -------------------------
    {
        if user == nil {
            return errors.New("nil user")
        }
        
        var validationErrs []string
        if user.ID == "" {
            validationErrs = append(validationErrs, "missing user ID")
        }
        if !isValidEmail(user.Email) {
            validationErrs = append(validationErrs, "invalid email")
        }
        
        if len(validationErrs) > 0 {
            return fmt.Errorf("validation errors: %s", strings.Join(validationErrs, "; "))
        }
        metrics.Record("user.validation", 1)
    }

    // ------------------------- TRANSFORMATION -------------------------
    {
        user.Name = normalizeName(user.Name)
        user.Email = strings.ToLower(user.Email)
        
        // Обогащение метаданных
        user.Metadata = enrichUserMetadata(user)
        user.UpdatedAt = time.Now().UTC().Truncate(time.Millisecond)
    }

    // ------------------------- PERSISTENCE -------------------------
    {
        ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
        defer cancel()
        
        if err := dbClient.SaveUserWithRetry(ctx, user, 3); err != nil {
            sentry.CaptureException(err)
            return fmt.Errorf("user persistence failed: %w", err)
        }
        log.Info().Str("user_id", user.ID).Msg("User saved")
    }
    
    return nil
}

// 2. File Processing
func ProcessUploadedFile(filePath string) (FileMetadata, error) {
    var metadata FileMetadata
    
    // ------------------------- FILE HANDLING -------------------------
    file, err := os.Open(filePath)
    if err != nil {
        return metadata, fmt.Errorf("file open: %w", err)
    }
    defer file.Close()
    
    // ------------------------- CONTENT ANALYSIS -------------------------
    {
        stats, err := file.Stat()
        if err != nil {
            return metadata, fmt.Errorf("file stats: %w", err)
        }
        
        if stats.Size() > config.MaxFileSize {
            return metadata, ErrFileTooLarge
        }
        
        // Потоковая обработка
        metadata = analyzeContent(file)
    }

    // ------------------------- SECURITY CHECKS -------------------------
    {
        if err := validateContentSafety(metadata); err != nil {
            antivirus.ScanFile(filePath)
            return metadata, fmt.Errorf("security: %w", err)
        }
        
        metadata.SHA256 = computeFileHash(file)
        metadata.ProcessedAt = time.Now().UTC()
    }
    
    return metadata, nil
}

// 3. API Request
func FetchWeatherData(location string) (WeatherData, error) {
    var result WeatherData
    
    // ------------------------- REQUEST PREPARATION -------------------------
    req, err := createWeatherRequest(location)
    if err != nil {
        return result, fmt.Errorf("request prep: %w", err)
    }

    // ------------------------- EXECUTION WITH RESILIENCE -------------------------
    resp, err := executeRequestWithRetry(req, 3)
    if err != nil {
        return result, fmt.Errorf("request execution: %w", err)
    }
    defer resp.Body.Close()

    // ------------------------- RESPONSE HANDLING -------------------------
    {
        if err := decodeWeatherResponse(resp.Body, &result); err != nil {
            return result, fmt.Errorf("response decode: %w", err)
        }
        
        if !result.IsValid() {
            log.Warn().Str("location", location).Msg("Invalid weather data")
            return result, ErrInvalidWeatherData
        }
        
        result.Source = "api.example.com"
        result.TTL = calculateCacheTTL(result.Temperature)
    }
    
    return result, nil
}
~~~