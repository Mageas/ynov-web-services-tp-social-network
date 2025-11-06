# Social Network API

Une API REST moderne pour un réseau social, construite selon les meilleures pratiques Go avec une architecture clean et une séparation claire des responsabilités.

## 🏗️ Architecture

Ce projet suit une **architecture en couches** (layered architecture) inspirée de l'architecture hexagonale et des principes DDD (Domain-Driven Design).

```
.
├── cmd/
│   └── api/                  # Point d'entrée de l'application
│       └── main.go           # Initialisation et configuration
├── internal/
│   ├── api/                  # Couche HTTP/API
│   │   ├── dto/              # Data Transfer Objects
│   │   │   ├── request.go    # Structures de requête
│   │   │   └── response.go   # Structures de réponse
│   │   ├── handler/          # Handlers HTTP
│   │   │   ├── auth.go       # Endpoints d'authentification
│   │   │   └── post.go       # Endpoints des posts
│   │   ├── middleware/       # Middlewares HTTP
│   │   │   ├── auth.go       # Middleware d'authentification JWT
│   │   │   ├── error.go      # Middleware de gestion d'erreurs
│   │   │   └── utils.go      # Utilitaires middleware
│   │   └── router/           # Configuration des routes
│   │       └── router.go
│   ├── config/               # Configuration de l'application
│   │   └── config.go         # Chargement et validation de la config
│   ├── domain/               # Couche métier (Domain Layer)
│   │   ├── post/
│   │   │   ├── post.go       # Entité Post
│   │   │   └── repository.go # Interface du repository Post
│   │   └── user/
│   │       ├── user.go       # Entité User
│   │       └── repository.go # Interface du repository User
│   ├── pkg/                  # Packages utilitaires internes
│   │   ├── apperrors/        # Gestion centralisée des erreurs
│   │   │   └── errors.go
│   │   ├── logger/           # Logger structuré
│   │   │   └── logger.go
│   │   └── validator/        # Validation des données
│   │       └── validator.go
│   ├── repository/           # Couche d'accès aux données
│   │   └── sqlite/
│   │       ├── database.go   # Connexion et migration DB
│   │       ├── models.go     # Modèles GORM
│   │       ├── post_repository.go  # Implémentation Post
│   │       └── user_repository.go  # Implémentation User
│   └── service/              # Couche de logique métier
│       ├── auth/
│       │   ├── jwt.go        # Service JWT
│       │   └── password.go   # Service de hachage
│       ├── post/
│       │   └── service.go    # Logique métier des posts
│       └── user/
│           └── service.go    # Logique métier des users
├── postman/                  # Collections Postman pour les tests
├── go.mod
└── go.sum
```

## 📋 Principes Architecturaux

### Séparation des Responsabilités

1. **Domain Layer** (`internal/domain/`)
   - Contient les entités métier pures
   - Définit les interfaces des repositories (inversion de dépendance)
   - Aucune dépendance sur les frameworks ou l'infrastructure

2. **Service Layer** (`internal/service/`)
   - Contient la logique métier
   - Orchestre les opérations entre le domain et les repositories
   - Valide les règles métier complexes
   - Indépendant de la couche HTTP

3. **Repository Layer** (`internal/repository/`)
   - Implémente les interfaces définies dans le domain
   - Gère la persistance des données
   - Isole la logique de la base de données

4. **API Layer** (`internal/api/`)
   - Gère les requêtes/réponses HTTP
   - Transforme les données (DTOs)
   - Applique les middlewares (auth, logging, erreurs)

5. **Infrastructure** (`internal/pkg/`, `internal/config/`)
   - Packages réutilisables (logger, validator, errors)
   - Configuration centralisée

### Avantages de cette Architecture

- ✅ **Testabilité**: Chaque couche peut être testée indépendamment
- ✅ **Maintenabilité**: Code organisé et facile à comprendre
- ✅ **Scalabilité**: Facile d'ajouter de nouvelles fonctionnalités
- ✅ **Flexibilité**: Changement de DB/framework sans affecter le métier
- ✅ **Réutilisabilité**: Services métier réutilisables

## 🚀 Démarrage

### Prérequis

- Go 1.22+
- SQLite

### Installation

```bash
# Cloner le repository
git clone https://github.com/Mageas/ynov-web-services-tp-social-network
cd ynov-web-services-tp-social-network

# Installer les dépendances
go mod download

# Créer le fichier .env
cat > .env << EOF
JWT_SECRET=your-secret-key-here
PORT=8080
DB_PATH=data.db
EOF
```

### Lancement

```bash
# Compiler et lancer l'application
go run cmd/api/main.go

# Ou compiler d'abord
go build -o bin/api cmd/api/main.go
./bin/api
```

L'API sera disponible sur `http://localhost:8080`

## 📡 API Endpoints

### Authentification

- **POST** `/signup` - Créer un nouveau compte
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

- **POST** `/login` - Se connecter
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
  Retourne: `{"token": "jwt-token"}`

### Posts (Authentification requise)

- **GET** `/posts?page=1&limit=10&beforeTs=<timestamp>` - Lister les posts
- **POST** `/posts` - Créer un post
  ```json
  {
    "content": "Mon premier post!"
  }
  ```

- **POST** `/posts/{id}/like` - Liker un post
- **DELETE** `/posts/{id}/unlike` - Unliker un post

### Authentification

Toutes les routes protégées nécessitent un header:
```
Authorization: Bearer <jwt-token>
```

## 🔧 Configuration

Variables d'environnement:

| Variable | Description | Défaut |
|----------|-------------|--------|
| JWT_SECRET | Secret pour signer les tokens JWT | **Obligatoire** |
| PORT | Port du serveur HTTP | 8080 |
| DB_PATH | Chemin de la base SQLite | data.db |

## 🏛️ Patterns Utilisés

### Dependency Injection
Les dépendances sont injectées via les constructeurs, facilitant les tests et la flexibilité.

### Repository Pattern
Abstraction de l'accès aux données via des interfaces, permettant de changer facilement d'implémentation.

### Service Layer Pattern
Centralisation de la logique métier, séparée de la couche HTTP.

### DTO Pattern
Séparation entre les modèles de domaine et les structures API.

### Middleware Pattern
Traitement en chaîne des requêtes HTTP (auth, logging, erreurs).

## 🧪 Tests

```bash
# Lancer tous les tests
go test ./...

# Tests avec couverture
go test -cover ./...

# Tests verbeux
go test -v ./...
```

## 📦 Dépendances

- **GORM**: ORM pour Go
- **golang-jwt/jwt**: Gestion des tokens JWT
- **godotenv**: Chargement des variables d'environnement

## 🔐 Sécurité

- Mots de passe hashés avec SHA256 + salt
- Authentification JWT avec expiration (24h)
- Validation des entrées utilisateur
- Protection contre les injections SQL (via GORM)

## 📝 Bonnes Pratiques Implémentées

1. **Clean Architecture**: Séparation claire des couches
2. **SOLID Principles**: Notamment l'inversion de dépendance
3. **Error Handling**: Gestion centralisée des erreurs
4. **Logging**: Logger structuré pour le debugging
5. **Validation**: Validation des entrées utilisateur
6. **Graceful Shutdown**: Arrêt propre du serveur
7. **Configuration**: Gestion centralisée de la config
8. **Standards Go**: Respect des conventions de nommage et structure

## 📚 Pour Aller Plus Loin

Pour améliorer encore cette API:

- [ ] Ajouter des tests unitaires et d'intégration
- [ ] Implémenter le tracing et les métriques (OpenTelemetry)
- [ ] Ajouter la pagination cursor-based complète
- [ ] Implémenter le rate limiting
- [ ] Ajouter la documentation OpenAPI/Swagger
- [ ] Configurer CI/CD
- [ ] Conteneuriser avec Docker
- [ ] Ajouter la migration de base de données versionnée

## 🤝 Contribution

Les contributions sont les bienvenues! N'hésitez pas à ouvrir une issue ou une pull request.

## 📄 Licence

MIT

