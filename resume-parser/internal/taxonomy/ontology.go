// Package taxonomy – ontology.go defines the built-in skill taxonomy database.
package taxonomy

// builtinSkills is the canonical skill ontology.
// Each entry defines a skill node with its domain, category, aliases,
// prerequisites, and related skills.
//
// Aliases are stored in lowercase for case-insensitive matching.
var builtinSkills = []SkillNode{
	// ─────────────────────────────────────────────────────────────────────────
	// Programming Languages
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "go", CanonicalName: "Go",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"golang", "go lang", "go programming"},
		RelatedSkills: []string{"grpc", "gin", "echo", "fiber", "kubernetes"},
		Description:   "Statically typed, compiled language designed at Google.",
	},
	{
		ID: "python", CanonicalName: "Python",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"py", "python3", "python 3", "python2", "python 2"},
		RelatedSkills: []string{"django", "flask", "fastapi", "pandas", "numpy", "tensorflow", "pytorch"},
		Description:   "High-level, general-purpose programming language.",
	},
	{
		ID: "javascript", CanonicalName: "JavaScript",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"js", "ecmascript", "es6", "es2015", "es2016", "es2017", "es2018", "es2019", "es2020", "es2021", "es2022", "vanilla js", "vanilla javascript"},
		RelatedSkills: []string{"typescript", "react", "angular", "vue", "nodejs"},
		Description:   "Lightweight, interpreted scripting language for the web.",
	},
	{
		ID: "typescript", CanonicalName: "TypeScript",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"ts", "typescript lang"},
		Prerequisites: []string{"javascript"},
		RelatedSkills: []string{"javascript", "react", "angular", "nodejs"},
		Description:   "Strongly typed superset of JavaScript.",
	},
	{
		ID: "java", CanonicalName: "Java",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"java se", "java ee", "java 8", "java 11", "java 17", "java 21"},
		RelatedSkills: []string{"spring-boot"},
		Description:   "Object-oriented, class-based programming language.",
	},
	{
		ID: "rust", CanonicalName: "Rust",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"rust lang", "rust programming"},
		RelatedSkills: []string{"actix"},
		Description:   "Systems programming language focused on safety and performance.",
	},
	{
		ID: "csharp", CanonicalName: "C#",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"c#", "csharp", "c sharp", ".net c#", "dotnet c#"},
		RelatedSkills: []string{"dotnet"},
		Description:   "Modern, object-oriented language for the .NET platform.",
	},
	{
		ID: "cpp", CanonicalName: "C++",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"c++", "cpp", "c plus plus", "cplusplus"},
		RelatedSkills: []string{"c"},
		Description:   "General-purpose language with low-level memory manipulation.",
	},
	{
		ID: "c", CanonicalName: "C",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"c language", "c programming", "ansi c"},
		RelatedSkills: []string{"cpp"},
		Description:   "General-purpose, procedural programming language.",
	},
	{
		ID: "ruby", CanonicalName: "Ruby",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"ruby lang", "ruby programming"},
		RelatedSkills: []string{"rails"},
		Description:   "Dynamic, open-source programming language.",
	},
	{
		ID: "php", CanonicalName: "PHP",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"php7", "php8", "php 7", "php 8"},
		RelatedSkills: []string{"laravel"},
		Description:   "Server-side scripting language for web development.",
	},
	{
		ID: "swift", CanonicalName: "Swift",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"swift lang", "swift programming", "apple swift"},
		RelatedSkills: []string{"ios"},
		Description:   "General-purpose language developed by Apple.",
	},
	{
		ID: "kotlin", CanonicalName: "Kotlin",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"kotlin lang", "kotlin programming"},
		Prerequisites: []string{"java"},
		RelatedSkills: []string{"android", "spring-boot"},
		Description:   "Cross-platform, statically typed language for JVM.",
	},
	{
		ID: "scala", CanonicalName: "Scala",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"scala lang"},
		RelatedSkills: []string{"spark"},
		Description:   "Strong static type system language for JVM.",
	},
	{
		ID: "r", CanonicalName: "R",
		Domain: DomainDataScience, Category: CategoryLanguage,
		Aliases:       []string{"r language", "r programming", "r stats"},
		Description:   "Language for statistical computing and graphics.",
	},
	{
		ID: "sql", CanonicalName: "SQL",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"structured query language", "ansi sql", "t-sql", "tsql", "pl/sql", "plsql"},
		RelatedSkills: []string{"postgresql", "mysql", "sqlite"},
		Description:   "Domain-specific language for managing relational databases.",
	},
	{
		ID: "bash", CanonicalName: "Bash",
		Domain: DomainEngineering, Category: CategoryLanguage,
		Aliases:       []string{"shell", "shell scripting", "bash scripting", "sh", "zsh", "unix shell"},
		Description:   "Unix shell and command language.",
	},
	{
		ID: "html", CanonicalName: "HTML",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"html5", "html 5", "hypertext markup language"},
		RelatedSkills: []string{"css", "javascript"},
		Description:   "Standard markup language for web pages.",
	},
	{
		ID: "css", CanonicalName: "CSS",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"css3", "css 3", "cascading style sheets"},
		RelatedSkills: []string{"html", "sass", "less", "tailwind"},
		Description:   "Style sheet language for HTML documents.",
	},
	{
		ID: "graphql", CanonicalName: "GraphQL",
		Domain: DomainEngineering, Category: CategoryAPI,
		Aliases:       []string{"graph ql", "gql"},
		RelatedSkills: []string{"rest", "apollo", "hasura"},
		Description:   "Query language for APIs.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// Frontend Frameworks
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "react", CanonicalName: "React",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"react.js", "reactjs", "react js", "react native"},
		Prerequisites: []string{"javascript"},
		RelatedSkills: []string{"redux", "next.js", "typescript", "jsx"},
		Description:   "JavaScript library for building user interfaces.",
	},
	{
		ID: "angular", CanonicalName: "Angular",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"angular.js", "angularjs", "angular js", "angular 2", "angular 4", "angular 8", "angular 12", "angular 14", "angular 16"},
		Prerequisites: []string{"typescript"},
		RelatedSkills: []string{"rxjs", "typescript", "ngrx"},
		Description:   "TypeScript-based web application framework by Google.",
	},
	{
		ID: "vue", CanonicalName: "Vue.js",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"vue.js", "vuejs", "vue js", "vue 2", "vue 3", "nuxt", "nuxt.js"},
		Prerequisites: []string{"javascript"},
		RelatedSkills: []string{"vuex", "pinia", "typescript"},
		Description:   "Progressive JavaScript framework for building UIs.",
	},
	{
		ID: "nextjs", CanonicalName: "Next.js",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"next.js", "nextjs", "next js"},
		Prerequisites: []string{"react"},
		RelatedSkills: []string{"react", "typescript", "vercel"},
		Description:   "React framework for production-grade web applications.",
	},
	{
		ID: "svelte", CanonicalName: "Svelte",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"svelte.js", "sveltejs", "sveltekit"},
		Prerequisites: []string{"javascript"},
		Description:   "Compiler-based JavaScript framework.",
	},
	{
		ID: "tailwind", CanonicalName: "Tailwind CSS",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"tailwindcss", "tailwind css"},
		Prerequisites: []string{"css"},
		Description:   "Utility-first CSS framework.",
	},
	{
		ID: "sass", CanonicalName: "Sass",
		Domain: DomainEngineering, Category: CategoryFrontend,
		Aliases:       []string{"scss", "sass/scss"},
		Prerequisites: []string{"css"},
		Description:   "CSS preprocessor scripting language.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// Backend Frameworks
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "django", CanonicalName: "Django",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"django framework", "django rest framework", "drf"},
		Prerequisites: []string{"python"},
		RelatedSkills: []string{"python", "postgresql", "rest"},
		Description:   "High-level Python web framework.",
	},
	{
		ID: "flask", CanonicalName: "Flask",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"flask framework", "flask python"},
		Prerequisites: []string{"python"},
		RelatedSkills: []string{"python", "sqlalchemy"},
		Description:   "Lightweight Python web framework.",
	},
	{
		ID: "fastapi", CanonicalName: "FastAPI",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"fast api", "fastapi framework"},
		Prerequisites: []string{"python"},
		RelatedSkills: []string{"python", "pydantic", "uvicorn"},
		Description:   "Modern, fast Python web framework for building APIs.",
	},
	{
		ID: "spring-boot", CanonicalName: "Spring Boot",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"spring boot", "springboot", "spring framework", "spring"},
		Prerequisites: []string{"java"},
		RelatedSkills: []string{"java", "maven", "gradle", "hibernate"},
		Description:   "Java-based framework for building microservices.",
	},
	{
		ID: "nodejs", CanonicalName: "Node.js",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"node.js", "nodejs", "node js", "node"},
		Prerequisites: []string{"javascript"},
		RelatedSkills: []string{"express", "nestjs", "typescript"},
		Description:   "JavaScript runtime built on Chrome's V8 engine.",
	},
	{
		ID: "express", CanonicalName: "Express.js",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"express.js", "expressjs", "express js", "express framework"},
		Prerequisites: []string{"nodejs"},
		Description:   "Minimal and flexible Node.js web application framework.",
	},
	{
		ID: "nestjs", CanonicalName: "NestJS",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"nest.js", "nestjs", "nest js"},
		Prerequisites: []string{"nodejs", "typescript"},
		Description:   "Progressive Node.js framework for scalable server-side apps.",
	},
	{
		ID: "rails", CanonicalName: "Ruby on Rails",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"ruby on rails", "rails", "ror"},
		Prerequisites: []string{"ruby"},
		Description:   "Server-side web application framework written in Ruby.",
	},
	{
		ID: "laravel", CanonicalName: "Laravel",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"laravel framework", "laravel php"},
		Prerequisites: []string{"php"},
		Description:   "PHP web application framework.",
	},
	{
		ID: "gin", CanonicalName: "Gin",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"gin framework", "gin-gonic"},
		Prerequisites: []string{"go"},
		Description:   "HTTP web framework written in Go.",
	},
	{
		ID: "echo", CanonicalName: "Echo",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"echo framework", "echo go"},
		Prerequisites: []string{"go"},
		Description:   "High performance, minimalist Go web framework.",
	},
	{
		ID: "fiber", CanonicalName: "Fiber",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"fiber framework", "gofiber"},
		Prerequisites: []string{"go"},
		Description:   "Express-inspired web framework written in Go.",
	},
	{
		ID: "actix", CanonicalName: "Actix",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{"actix-web", "actix web"},
		Prerequisites: []string{"rust"},
		Description:   "Powerful, pragmatic, and extremely fast Rust web framework.",
	},
	{
		ID: "dotnet", CanonicalName: ".NET",
		Domain: DomainEngineering, Category: CategoryBackend,
		Aliases:       []string{".net", "dotnet", "asp.net", "asp.net core", "aspnet", ".net core", "dotnet core"},
		Prerequisites: []string{"csharp"},
		Description:   "Free, cross-platform, open-source developer platform.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// Mobile
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "ios", CanonicalName: "iOS Development",
		Domain: DomainEngineering, Category: CategoryMobile,
		Aliases:       []string{"ios development", "ios dev", "iphone development"},
		Prerequisites: []string{"swift"},
		RelatedSkills: []string{"swift", "xcode", "objective-c"},
		Description:   "Development for Apple iOS platform.",
	},
	{
		ID: "android", CanonicalName: "Android Development",
		Domain: DomainEngineering, Category: CategoryMobile,
		Aliases:       []string{"android development", "android dev"},
		Prerequisites: []string{"kotlin"},
		RelatedSkills: []string{"kotlin", "java", "android studio"},
		Description:   "Development for Google Android platform.",
	},
	{
		ID: "flutter", CanonicalName: "Flutter",
		Domain: DomainEngineering, Category: CategoryMobile,
		Aliases:       []string{"flutter sdk", "flutter framework"},
		RelatedSkills: []string{"dart", "ios", "android"},
		Description:   "Google's UI toolkit for cross-platform apps.",
	},
	{
		ID: "react-native", CanonicalName: "React Native",
		Domain: DomainEngineering, Category: CategoryMobile,
		Aliases:       []string{"react native", "reactnative", "rn"},
		Prerequisites: []string{"react"},
		Description:   "Framework for building native apps using React.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// Databases
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "postgresql", CanonicalName: "PostgreSQL",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"postgres", "psql", "pg", "postgresql database"},
		Prerequisites: []string{"sql"},
		Description:   "Advanced open-source relational database.",
	},
	{
		ID: "mysql", CanonicalName: "MySQL",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"mysql database", "mysql server"},
		Prerequisites: []string{"sql"},
		Description:   "Open-source relational database management system.",
	},
	{
		ID: "mongodb", CanonicalName: "MongoDB",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"mongo", "mongo db", "mongodb database"},
		Description:   "Document-oriented NoSQL database.",
	},
	{
		ID: "redis", CanonicalName: "Redis",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"redis cache", "redis db"},
		Description:   "In-memory data structure store.",
	},
	{
		ID: "elasticsearch", CanonicalName: "Elasticsearch",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"elastic search", "elastic", "es", "opensearch"},
		Description:   "Distributed, RESTful search and analytics engine.",
	},
	{
		ID: "cassandra", CanonicalName: "Cassandra",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"apache cassandra", "cassandra db"},
		Description:   "Distributed NoSQL database for high availability.",
	},
	{
		ID: "dynamodb", CanonicalName: "DynamoDB",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"dynamo db", "aws dynamodb", "amazon dynamodb"},
		Description:   "AWS managed NoSQL database service.",
	},
	{
		ID: "sqlite", CanonicalName: "SQLite",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"sqlite3", "sqlite database"},
		Prerequisites: []string{"sql"},
		Description:   "Lightweight, file-based relational database.",
	},
	{
		ID: "neo4j", CanonicalName: "Neo4j",
		Domain: DomainEngineering, Category: CategoryDatabase,
		Aliases:       []string{"neo4j database", "graph database"},
		Description:   "Graph database management system.",
	},
	{
		ID: "pinecone", CanonicalName: "Pinecone",
		Domain: DomainDataScience, Category: CategoryDatabase,
		Aliases:       []string{"pinecone db", "pinecone vector"},
		Description:   "Managed vector database for ML applications.",
	},
	{
		ID: "weaviate", CanonicalName: "Weaviate",
		Domain: DomainDataScience, Category: CategoryDatabase,
		Aliases:       []string{"weaviate db"},
		Description:   "Open-source vector database.",
	},
	{
		ID: "qdrant", CanonicalName: "Qdrant",
		Domain: DomainDataScience, Category: CategoryDatabase,
		Aliases:       []string{"qdrant db"},
		Description:   "Vector similarity search engine.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// Cloud Platforms
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "aws", CanonicalName: "AWS",
		Domain: DomainDevOps, Category: CategoryCloud,
		Aliases:       []string{"amazon web services", "amazon aws", "aws cloud"},
		RelatedSkills: []string{"ec2", "s3", "lambda", "rds", "dynamodb", "ecs", "eks"},
		Description:   "Amazon Web Services cloud platform.",
	},
	{
		ID: "azure", CanonicalName: "Azure",
		Domain: DomainDevOps, Category: CategoryCloud,
		Aliases:       []string{"microsoft azure", "azure cloud", "ms azure"},
		Description:   "Microsoft Azure cloud platform.",
	},
	{
		ID: "gcp", CanonicalName: "GCP",
		Domain: DomainDevOps, Category: CategoryCloud,
		Aliases:       []string{"google cloud", "google cloud platform", "google cloud services"},
		Description:   "Google Cloud Platform.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// DevOps & Infrastructure
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "docker", CanonicalName: "Docker",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"docker container", "docker compose", "dockerfile"},
		RelatedSkills: []string{"kubernetes", "containerization"},
		Description:   "Platform for developing, shipping, and running containers.",
	},
	{
		ID: "kubernetes", CanonicalName: "Kubernetes",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"k8s", "kube", "k8", "kubernetes orchestration"},
		Prerequisites: []string{"docker"},
		RelatedSkills: []string{"helm", "istio", "docker"},
		Description:   "Container orchestration system.",
	},
	{
		ID: "terraform", CanonicalName: "Terraform",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"terraform iac", "hashicorp terraform"},
		Description:   "Infrastructure as code tool.",
	},
	{
		ID: "ansible", CanonicalName: "Ansible",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"ansible automation", "red hat ansible"},
		Description:   "IT automation platform.",
	},
	{
		ID: "jenkins", CanonicalName: "Jenkins",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"jenkins ci", "jenkins pipeline"},
		Description:   "Open-source automation server for CI/CD.",
	},
	{
		ID: "github-actions", CanonicalName: "GitHub Actions",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"github actions", "gh actions", "github ci"},
		Description:   "CI/CD platform integrated with GitHub.",
	},
	{
		ID: "gitlab-ci", CanonicalName: "GitLab CI/CD",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"gitlab ci", "gitlab ci/cd", "gitlab pipeline"},
		Description:   "CI/CD platform integrated with GitLab.",
	},
	{
		ID: "helm", CanonicalName: "Helm",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"helm chart", "helm charts"},
		Prerequisites: []string{"kubernetes"},
		Description:   "Package manager for Kubernetes.",
	},
	{
		ID: "prometheus", CanonicalName: "Prometheus",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"prometheus monitoring"},
		RelatedSkills: []string{"grafana"},
		Description:   "Open-source monitoring and alerting toolkit.",
	},
	{
		ID: "grafana", CanonicalName: "Grafana",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"grafana dashboard"},
		RelatedSkills: []string{"prometheus"},
		Description:   "Open-source analytics and monitoring platform.",
	},
	{
		ID: "git", CanonicalName: "Git",
		Domain: DomainEngineering, Category: CategoryDevOps,
		Aliases:       []string{"git version control", "git scm"},
		RelatedSkills: []string{"github", "gitlab", "bitbucket"},
		Description:   "Distributed version control system.",
	},
	{
		ID: "cicd", CanonicalName: "CI/CD",
		Domain: DomainDevOps, Category: CategoryDevOps,
		Aliases:       []string{"ci/cd", "continuous integration", "continuous delivery", "continuous deployment", "ci cd", "cicd pipeline"},
		Description:   "Continuous integration and continuous delivery practices.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// API & Messaging
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "rest", CanonicalName: "REST",
		Domain: DomainEngineering, Category: CategoryAPI,
		Aliases:       []string{"restful", "rest api", "restful api", "rest apis", "restful apis", "rest services", "restful services"},
		Description:   "Representational State Transfer architectural style.",
	},
	{
		ID: "grpc", CanonicalName: "gRPC",
		Domain: DomainEngineering, Category: CategoryAPI,
		Aliases:       []string{"grpc", "google rpc", "protocol buffers", "protobuf"},
		Description:   "High-performance RPC framework.",
	},
	{
		ID: "kafka", CanonicalName: "Apache Kafka",
		Domain: DomainEngineering, Category: CategoryMessaging,
		Aliases:       []string{"kafka", "apache kafka", "kafka streaming"},
		Description:   "Distributed event streaming platform.",
	},
	{
		ID: "rabbitmq", CanonicalName: "RabbitMQ",
		Domain: DomainEngineering, Category: CategoryMessaging,
		Aliases:       []string{"rabbit mq", "amqp"},
		Description:   "Open-source message broker.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// ML / Data Science Frameworks
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "tensorflow", CanonicalName: "TensorFlow",
		Domain: DomainDataScience, Category: CategoryMLFramework,
		Aliases:       []string{"tensor flow", "tf", "tensorflow 2", "tensorflow2"},
		Prerequisites: []string{"python"},
		RelatedSkills: []string{"keras", "python", "deep-learning"},
		Description:   "Open-source machine learning framework by Google.",
	},
	{
		ID: "pytorch", CanonicalName: "PyTorch",
		Domain: DomainDataScience, Category: CategoryMLFramework,
		Aliases:       []string{"py torch", "torch", "pytorch framework"},
		Prerequisites: []string{"python"},
		RelatedSkills: []string{"python", "deep-learning"},
		Description:   "Open-source machine learning framework by Meta.",
	},
	{
		ID: "keras", CanonicalName: "Keras",
		Domain: DomainDataScience, Category: CategoryMLFramework,
		Aliases:       []string{"keras api"},
		Prerequisites: []string{"tensorflow"},
		Description:   "High-level neural networks API.",
	},
	{
		ID: "scikit-learn", CanonicalName: "scikit-learn",
		Domain: DomainDataScience, Category: CategoryMLFramework,
		Aliases:       []string{"sklearn", "scikit learn", "scikitlearn"},
		Prerequisites: []string{"python"},
		Description:   "Machine learning library for Python.",
	},
	{
		ID: "pandas", CanonicalName: "Pandas",
		Domain: DomainDataScience, Category: CategoryDataTools,
		Aliases:       []string{"pandas library", "pandas dataframe"},
		Prerequisites: []string{"python"},
		Description:   "Data analysis and manipulation library for Python.",
	},
	{
		ID: "numpy", CanonicalName: "NumPy",
		Domain: DomainDataScience, Category: CategoryDataTools,
		Aliases:       []string{"numpy library", "np"},
		Prerequisites: []string{"python"},
		Description:   "Fundamental package for scientific computing in Python.",
	},
	{
		ID: "spark", CanonicalName: "Apache Spark",
		Domain: DomainDataScience, Category: CategoryDataTools,
		Aliases:       []string{"apache spark", "pyspark", "spark streaming"},
		Description:   "Unified analytics engine for large-scale data processing.",
	},
	{
		ID: "machine-learning", CanonicalName: "Machine Learning",
		Domain: DomainDataScience, Category: CategoryMLConcept,
		Aliases:       []string{"ml", "machine learning", "supervised learning", "unsupervised learning"},
		Description:   "Field of AI that enables systems to learn from data.",
	},
	{
		ID: "deep-learning", CanonicalName: "Deep Learning",
		Domain: DomainDataScience, Category: CategoryMLConcept,
		Aliases:       []string{"dl", "deep learning", "neural networks", "neural network"},
		Prerequisites: []string{"machine-learning"},
		Description:   "Subset of ML using multi-layered neural networks.",
	},
	{
		ID: "nlp", CanonicalName: "NLP",
		Domain: DomainDataScience, Category: CategoryMLConcept,
		Aliases:       []string{"natural language processing", "text mining", "text analytics", "computational linguistics"},
		Prerequisites: []string{"machine-learning"},
		Description:   "AI field focused on interaction between computers and human language.",
	},
	{
		ID: "llm", CanonicalName: "LLM",
		Domain: DomainDataScience, Category: CategoryMLConcept,
		Aliases:       []string{"large language model", "large language models", "llms", "gpt", "chatgpt", "openai"},
		Prerequisites: []string{"deep-learning", "nlp"},
		Description:   "Large language models for natural language tasks.",
	},
	{
		ID: "rag", CanonicalName: "RAG",
		Domain: DomainDataScience, Category: CategoryMLConcept,
		Aliases:       []string{"retrieval augmented generation", "retrieval-augmented generation"},
		Prerequisites: []string{"llm"},
		Description:   "Retrieval-Augmented Generation for LLM applications.",
	},

	// ─────────────────────────────────────────────────────────────────────────
	// Soft Skills
	// ─────────────────────────────────────────────────────────────────────────
	{
		ID: "leadership", CanonicalName: "Leadership",
		Domain: DomainManagement, Category: CategoryLeadership,
		Aliases:       []string{"team leadership", "technical leadership", "tech lead", "engineering leadership"},
		Description:   "Ability to guide and inspire a team.",
	},
	{
		ID: "communication", CanonicalName: "Communication",
		Domain: DomainCommunication, Category: CategoryCommunication,
		Aliases:       []string{"verbal communication", "written communication", "interpersonal communication", "effective communication"},
		Description:   "Ability to convey information clearly and effectively.",
	},
	{
		ID: "teamwork", CanonicalName: "Teamwork",
		Domain: DomainCommunication, Category: CategoryCollaboration,
		Aliases:       []string{"team player", "collaboration", "collaborative", "cross-functional collaboration"},
		Description:   "Ability to work effectively in a team.",
	},
	{
		ID: "problem-solving", CanonicalName: "Problem Solving",
		Domain: DomainManagement, Category: CategoryProblemSolving,
		Aliases:       []string{"problem solving", "analytical thinking", "critical thinking", "analytical skills"},
		Description:   "Ability to identify and resolve complex problems.",
	},
	{
		ID: "project-management", CanonicalName: "Project Management",
		Domain: DomainManagement, Category: CategoryProjectMgmt,
		Aliases:       []string{"project management", "program management", "pmp", "agile project management"},
		Description:   "Planning, executing, and closing projects.",
	},
	{
		ID: "agile", CanonicalName: "Agile",
		Domain: DomainManagement, Category: CategoryProjectMgmt,
		Aliases:       []string{"agile methodology", "agile development", "scrum", "kanban", "sprint", "agile scrum"},
		Description:   "Iterative approach to project management and software development.",
	},
	{
		ID: "mentoring", CanonicalName: "Mentoring",
		Domain: DomainManagement, Category: CategoryLeadership,
		Aliases:       []string{"mentorship", "coaching", "staff development"},
		Description:   "Guiding and developing junior team members.",
	},
	{
		ID: "stakeholder-management", CanonicalName: "Stakeholder Management",
		Domain: DomainManagement, Category: CategoryProjectMgmt,
		Aliases:       []string{"stakeholder management", "stakeholder communication", "executive communication"},
		Description:   "Managing relationships with project stakeholders.",
	},
}
