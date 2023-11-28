package treblle_fiber

var Config internalConfiguration

// Configuration sets up and customizes communication with the Treblle API
type Configuration struct {
	APIKey     string
	ProjectID  string
	AdditionalFieldsToMask []string
	IgnoreExact []string
	IgnorePrefix []string
	ServerURL  string
}

// internalConfiguration is used for communication with Treblle API and contains optimizations
type internalConfiguration struct {
	Configuration
	FieldsMap    map[string]bool
	serverInfo   ServerInfo
	languageInfo LanguageInfo
	Debug        bool
}

func Configure(config Configuration) {
	if config.APIKey != "" {
		Config.APIKey = config.APIKey
	}
	if config.ProjectID != "" {
		Config.ProjectID = config.ProjectID
	}
	if len(config.AdditionalFieldsToMask) > 0 {
		Config.AdditionalFieldsToMask = config.AdditionalFieldsToMask
	}

	Config.IgnoreExact = config.IgnoreExact
	Config.IgnorePrefix = config.IgnorePrefix
	Config.FieldsMap = generateFieldsToMask(Config.AdditionalFieldsToMask)
	Config.serverInfo = getServerInfo()
	Config.languageInfo = getLanguageInfo()
}

func generateFieldsToMask(additionalFieldsToMask []string) map[string]bool {
	defaultFieldsToMask := []string{
		"password",
		"pwd",
		"secret",
		"password_confirmation",
		"passwordConfirmation",
		"cc",
		"card_number",
		"cardNumber",
		"ccv",
		"ssn",
		"credit_score",
		"creditScore",
	}

	fields := append(defaultFieldsToMask, additionalFieldsToMask...)
	fieldsToMask := make(map[string]bool)
	for _, field := range fields {
		fieldsToMask[field] = true
	}

	return fieldsToMask
}
