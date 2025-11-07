.PHONY: run build test clean lint fmt json-server

# Variables
BINARY_NAME=api
MAIN_PATH=cmd/api/main.go
BUILD_DIR=bin

run: ## Lancer l'application
	@echo "Démarrage de l'application..."
	@go run $(MAIN_PATH)

build: ## Compiler l'application
	@echo "Compilation de l'application..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binaire créé: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Lancer les tests
	@echo "Lancement des tests..."
	@go test -v -race ./...

test-coverage: ## Lancer les tests avec couverture
	@echo "Lancement des tests avec couverture..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Rapport de couverture: coverage.html"

lint: ## Lancer le linter
	@echo "Lancement du linter..."
	@go vet ./...
	@golangci-lint run ./... 2>/dev/null || echo "golangci-lint non installé, utilisation de go vet uniquement"

fmt: ## Formater le code
	@echo "Formatage du code..."
	@go fmt ./...
	@goimports -w . 2>/dev/null || echo "goimports non installé"

clean: ## Nettoyer les fichiers générés
	@echo "Nettoyage..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f data.db
	@echo "Nettoyage terminé"

deps: ## Installer les dépendances
	@echo "Installation des dépendances..."
	@go mod download
	@go mod tidy

json-server: ## Lancer json-server
	@echo "Lancement de json-server..."
	@json-server --watch json-server/db.json --port 3000

.DEFAULT_GOAL := run
