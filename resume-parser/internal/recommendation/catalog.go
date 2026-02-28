// Package recommendation – catalog.go contains the built-in learning resource
// catalog used by the recommendation engine. This is a curated, in-memory
// catalog of high-quality resources for the most common technical skills.
package recommendation

// builtinCatalog is the built-in resource catalog.
// Resources are keyed by their primary skill (normalized lowercase).
var builtinCatalog = []ResourceEntry{
	// ── Python ──────────────────────────────────────────────────────────────
	{
		ID: "python-for-everybody", Title: "Python for Everybody Specialization",
		Description:    "Learn to program and analyze data with Python. Covers Python basics, data structures, web access, databases, and data visualization.",
		URL:            "https://www.coursera.org/specializations/python",
		Provider:       "Coursera", ResourceType: "course", Difficulty: "beginner",
		CostType: "free_audit", CostUSD: 49.00, DurationHours: 80, DurationLabel: "8 months",
		Skills: []string{"python", "data analysis", "sql"}, PrimarySkill: "python",
		Rating: 4.80, RatingCount: 1200000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "complete-python-bootcamp", Title: "Complete Python Bootcamp: From Zero to Hero",
		Description:    "Learn Python like a professional. Start from the basics and go all the way to creating your own applications and games.",
		URL:            "https://www.udemy.com/course/complete-python-bootcamp/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 22, DurationLabel: "22 hours",
		Skills: []string{"python", "oop"}, PrimarySkill: "python",
		Rating: 4.60, RatingCount: 500000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "python-docs", Title: "Python Official Documentation",
		Description:    "The official Python 3 documentation including tutorial, library reference, and language reference.",
		URL:            "https://docs.python.org/3/",
		Provider:       "Python.org", ResourceType: "documentation", Difficulty: "all_levels",
		CostType: "free", CostUSD: 0, DurationHours: 0, DurationLabel: "Self-paced",
		Skills: []string{"python"}, PrimarySkill: "python",
		Rating: 4.90, RatingCount: 500000, HasCertificate: false, HasHandsOn: false, IsVerified: true,
	},
	{
		ID: "automate-boring-stuff", Title: "Automate the Boring Stuff with Python",
		Description:    "A practical programming book for office workers. Free to read online.",
		URL:            "https://automatetheboringstuff.com/",
		Provider:       "No Starch Press", ResourceType: "book", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 20, DurationLabel: "20 hours",
		Skills: []string{"python", "automation"}, PrimarySkill: "python",
		Rating: 4.70, RatingCount: 50000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},

	// ── Go ──────────────────────────────────────────────────────────────────
	{
		ID: "go-complete-guide", Title: "Go: The Complete Developer's Guide",
		Description:    "Master the fundamentals and advanced features of the Go programming language.",
		URL:            "https://www.udemy.com/course/go-the-complete-developers-guide/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 9, DurationLabel: "9 hours",
		Skills: []string{"go", "concurrency"}, PrimarySkill: "go",
		Rating: 4.60, RatingCount: 45000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "tour-of-go", Title: "A Tour of Go",
		Description:    "An interactive introduction to Go programming language with hands-on exercises.",
		URL:            "https://go.dev/tour/",
		Provider:       "Go.dev", ResourceType: "documentation", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 4, DurationLabel: "4 hours",
		Skills: []string{"go"}, PrimarySkill: "go",
		Rating: 4.80, RatingCount: 100000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "go-by-example", Title: "Go by Example",
		Description:    "Hands-on introduction to Go using annotated example programs.",
		URL:            "https://gobyexample.com/",
		Provider:       "Go by Example", ResourceType: "documentation", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 8, DurationLabel: "8 hours",
		Skills: []string{"go"}, PrimarySkill: "go",
		Rating: 4.90, RatingCount: 200000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},

	// ── JavaScript / TypeScript ──────────────────────────────────────────────
	{
		ID: "complete-javascript-course", Title: "The Complete JavaScript Course 2024",
		Description:    "The modern JavaScript course for everyone. Master JavaScript with projects, challenges and theory.",
		URL:            "https://www.udemy.com/course/the-complete-javascript-course/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 69, DurationLabel: "69 hours",
		Skills: []string{"javascript", "es6"}, PrimarySkill: "javascript",
		Rating: 4.70, RatingCount: 350000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "freecodecamp-javascript", Title: "freeCodeCamp JavaScript Algorithms and Data Structures",
		Description:    "Free certification covering JavaScript fundamentals, ES6, data structures, and algorithm scripting.",
		URL:            "https://www.freecodecamp.org/learn/javascript-algorithms-and-data-structures/",
		Provider:       "freeCodeCamp", ResourceType: "course", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 300, DurationLabel: "300 hours",
		Skills: []string{"javascript", "algorithms"}, PrimarySkill: "javascript",
		Rating: 4.50, RatingCount: 500000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "typescript-complete-guide", Title: "TypeScript: The Complete Developer's Guide",
		Description:    "Master TypeScript by building real projects. Covers type system, generics, decorators.",
		URL:            "https://www.udemy.com/course/typescript-the-complete-developers-guide/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 27, DurationLabel: "27 hours",
		Skills: []string{"typescript", "javascript"}, PrimarySkill: "typescript",
		Rating: 4.60, RatingCount: 80000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── React ────────────────────────────────────────────────────────────────
	{
		ID: "react-complete-guide", Title: "React - The Complete Guide 2024",
		Description:    "Dive in and learn React.js from scratch. Learn Reactjs, Hooks, Redux, React Router, Next.js.",
		URL:            "https://www.udemy.com/course/react-the-complete-guide-incl-redux/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 68, DurationLabel: "68 hours",
		Skills: []string{"react", "redux", "javascript"}, PrimarySkill: "react",
		Rating: 4.60, RatingCount: 250000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "react-docs", Title: "React Official Documentation",
		Description:    "The official React documentation with interactive examples, tutorials, and API reference.",
		URL:            "https://react.dev/",
		Provider:       "React.dev", ResourceType: "documentation", Difficulty: "all_levels",
		CostType: "free", CostUSD: 0, DurationHours: 0, DurationLabel: "Self-paced",
		Skills: []string{"react"}, PrimarySkill: "react",
		Rating: 4.80, RatingCount: 200000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},

	// ── Node.js ──────────────────────────────────────────────────────────────
	{
		ID: "complete-nodejs-course", Title: "The Complete Node.js Developer Course",
		Description:    "Learn Node.js by building real-world applications with Node, Express, MongoDB, Jest.",
		URL:            "https://www.udemy.com/course/the-complete-nodejs-developer-course-2/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 35, DurationLabel: "35 hours",
		Skills: []string{"node.js", "express", "mongodb"}, PrimarySkill: "node.js",
		Rating: 4.60, RatingCount: 150000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Machine Learning ─────────────────────────────────────────────────────
	{
		ID: "ml-specialization", Title: "Machine Learning Specialization",
		Description:    "Andrew Ng's updated ML course. Covers supervised learning, unsupervised learning, and best practices.",
		URL:            "https://www.coursera.org/specializations/machine-learning-introduction",
		Provider:       "Coursera", ResourceType: "course", Difficulty: "intermediate",
		CostType: "free_audit", CostUSD: 49.00, DurationHours: 90, DurationLabel: "3 months",
		Skills: []string{"machine learning", "python", "tensorflow"}, PrimarySkill: "machine learning",
		Rating: 4.90, RatingCount: 500000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "deep-learning-specialization", Title: "Deep Learning Specialization",
		Description:    "Become a Deep Learning expert. Master deep neural networks, CNNs, RNNs, LSTMs, and transformers.",
		URL:            "https://www.coursera.org/specializations/deep-learning",
		Provider:       "Coursera", ResourceType: "course", Difficulty: "advanced",
		CostType: "free_audit", CostUSD: 49.00, DurationHours: 120, DurationLabel: "5 months",
		Skills: []string{"deep learning", "tensorflow", "python"}, PrimarySkill: "deep learning",
		Rating: 4.90, RatingCount: 400000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "practical-deep-learning-fastai", Title: "Practical Deep Learning for Coders",
		Description:    "Free course from fast.ai. Learn deep learning with PyTorch and fastai.",
		URL:            "https://course.fast.ai/",
		Provider:       "fast.ai", ResourceType: "course", Difficulty: "intermediate",
		CostType: "free", CostUSD: 0, DurationHours: 30, DurationLabel: "30 hours",
		Skills: []string{"deep learning", "pytorch", "python"}, PrimarySkill: "deep learning",
		Rating: 4.80, RatingCount: 100000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "pytorch-bootcamp", Title: "PyTorch for Deep Learning Bootcamp",
		Description:    "Learn PyTorch for deep learning. Covers tensors, neural networks, CNNs, RNNs, and transfer learning.",
		URL:            "https://www.udemy.com/course/pytorch-for-deep-learning-bootcamp/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 17, DurationLabel: "17 hours",
		Skills: []string{"pytorch", "deep learning"}, PrimarySkill: "pytorch",
		Rating: 4.60, RatingCount: 30000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "tensorflow-certificate", Title: "TensorFlow Developer Certificate",
		Description:    "Official TensorFlow certification. Demonstrates proficiency in using TensorFlow for deep learning.",
		URL:            "https://www.tensorflow.org/certificate",
		Provider:       "Google", ResourceType: "certification", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 100.00, DurationHours: 40, DurationLabel: "40 hours prep",
		Skills: []string{"tensorflow", "deep learning"}, PrimarySkill: "tensorflow",
		Rating: 4.60, RatingCount: 20000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── SQL / Databases ──────────────────────────────────────────────────────
	{
		ID: "complete-sql-bootcamp", Title: "The Complete SQL Bootcamp",
		Description:    "Become an expert at SQL. Learn how to read and write complex queries to a database.",
		URL:            "https://www.udemy.com/course/the-complete-sql-bootcamp/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 9, DurationLabel: "9 hours",
		Skills: []string{"sql", "postgresql"}, PrimarySkill: "sql",
		Rating: 4.70, RatingCount: 200000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "postgresql-complete-guide", Title: "PostgreSQL: The Complete Developer's Guide",
		Description:    "Master PostgreSQL with this comprehensive course. Covers advanced queries, indexing, performance tuning.",
		URL:            "https://www.udemy.com/course/sql-and-postgresql/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 22, DurationLabel: "22 hours",
		Skills: []string{"postgresql", "sql"}, PrimarySkill: "postgresql",
		Rating: 4.70, RatingCount: 50000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "sqlzoo", Title: "SQLZoo Interactive SQL Tutorial",
		Description:    "Free interactive SQL tutorial with exercises. Covers SELECT, INSERT, UPDATE, DELETE.",
		URL:            "https://sqlzoo.net/",
		Provider:       "SQLZoo", ResourceType: "documentation", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 10, DurationLabel: "10 hours",
		Skills: []string{"sql"}, PrimarySkill: "sql",
		Rating: 4.50, RatingCount: 500000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},

	// ── Docker & Kubernetes ──────────────────────────────────────────────────
	{
		ID: "docker-kubernetes-complete", Title: "Docker and Kubernetes: The Complete Guide",
		Description:    "Build, test, and deploy Docker applications with Kubernetes while learning production-style workflows.",
		URL:            "https://www.udemy.com/course/docker-and-kubernetes-the-complete-guide/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 22, DurationLabel: "22 hours",
		Skills: []string{"docker", "kubernetes"}, PrimarySkill: "docker",
		Rating: 4.60, RatingCount: 100000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "docker-docs", Title: "Docker Official Documentation",
		Description:    "Official Docker documentation covering installation, getting started, guides, and reference material.",
		URL:            "https://docs.docker.com/",
		Provider:       "Docker", ResourceType: "documentation", Difficulty: "all_levels",
		CostType: "free", CostUSD: 0, DurationHours: 0, DurationLabel: "Self-paced",
		Skills: []string{"docker"}, PrimarySkill: "docker",
		Rating: 4.70, RatingCount: 300000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "cka-certification", Title: "Certified Kubernetes Administrator (CKA)",
		Description:    "The CKA certification ensures holders have the skills to perform Kubernetes administrator responsibilities.",
		URL:            "https://training.linuxfoundation.org/certification/certified-kubernetes-administrator-cka/",
		Provider:       "Linux Foundation", ResourceType: "certification", Difficulty: "advanced",
		CostType: "paid", CostUSD: 395.00, DurationHours: 60, DurationLabel: "60 hours prep",
		Skills: []string{"kubernetes"}, PrimarySkill: "kubernetes",
		Rating: 4.70, RatingCount: 50000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── AWS ──────────────────────────────────────────────────────────────────
	{
		ID: "aws-saa-udemy", Title: "Ultimate AWS Certified Solutions Architect Associate",
		Description:    "Pass the AWS Certified Solutions Architect Associate certification. Covers all AWS services with hands-on labs.",
		URL:            "https://www.udemy.com/course/aws-certified-solutions-architect-associate-saa-c03/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 27, DurationLabel: "27 hours",
		Skills: []string{"aws", "cloud architecture"}, PrimarySkill: "aws",
		Rating: 4.70, RatingCount: 300000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "aws-saa-cert", Title: "AWS Certified Solutions Architect – Associate",
		Description:    "The AWS SAA certification validates the ability to design and implement distributed systems on AWS.",
		URL:            "https://aws.amazon.com/certification/certified-solutions-architect-associate/",
		Provider:       "AWS", ResourceType: "certification", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 300.00, DurationHours: 80, DurationLabel: "80 hours prep",
		Skills: []string{"aws", "cloud architecture"}, PrimarySkill: "aws",
		Rating: 4.80, RatingCount: 200000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "aws-cloud-practitioner", Title: "AWS Cloud Practitioner Essentials",
		Description:    "Free foundational course for AWS Cloud Practitioner certification.",
		URL:            "https://aws.amazon.com/training/digital/aws-cloud-practitioner-essentials/",
		Provider:       "AWS", ResourceType: "course", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 6, DurationLabel: "6 hours",
		Skills: []string{"aws"}, PrimarySkill: "aws",
		Rating: 4.60, RatingCount: 500000, HasCertificate: false, HasHandsOn: false, IsVerified: true,
	},

	// ── GCP ──────────────────────────────────────────────────────────────────
	{
		ID: "gcp-ace-cert", Title: "Google Cloud Associate Cloud Engineer",
		Description:    "Validates ability to deploy applications, monitor operations, and manage enterprise solutions on GCP.",
		URL:            "https://cloud.google.com/certification/cloud-engineer",
		Provider:       "Google Cloud", ResourceType: "certification", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 200.00, DurationHours: 60, DurationLabel: "60 hours prep",
		Skills: []string{"gcp"}, PrimarySkill: "gcp",
		Rating: 4.60, RatingCount: 30000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Azure ────────────────────────────────────────────────────────────────
	{
		ID: "az-900-cert", Title: "AZ-900: Microsoft Azure Fundamentals",
		Description:    "Foundational knowledge of cloud services and how those services are provided with Microsoft Azure.",
		URL:            "https://learn.microsoft.com/en-us/certifications/azure-fundamentals/",
		Provider:       "Microsoft", ResourceType: "certification", Difficulty: "beginner",
		CostType: "paid", CostUSD: 165.00, DurationHours: 20, DurationLabel: "20 hours prep",
		Skills: []string{"azure"}, PrimarySkill: "azure",
		Rating: 4.70, RatingCount: 100000, HasCertificate: true, HasHandsOn: false, IsVerified: true,
	},

	// ── Terraform ────────────────────────────────────────────────────────────
	{
		ID: "terraform-beginner-master", Title: "Terraform: From Beginner to Master",
		Description:    "Learn Terraform from scratch. Covers infrastructure as code, AWS provisioning, modules, state management.",
		URL:            "https://www.udemy.com/course/terraform-beginner-to-advanced/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 12, DurationLabel: "12 hours",
		Skills: []string{"terraform", "aws"}, PrimarySkill: "terraform",
		Rating: 4.60, RatingCount: 40000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "terraform-associate-cert", Title: "HashiCorp Certified: Terraform Associate",
		Description:    "Validates knowledge of infrastructure automation using Terraform.",
		URL:            "https://www.hashicorp.com/certification/terraform-associate",
		Provider:       "HashiCorp", ResourceType: "certification", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 70.50, DurationHours: 40, DurationLabel: "40 hours prep",
		Skills: []string{"terraform"}, PrimarySkill: "terraform",
		Rating: 4.70, RatingCount: 30000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Git ──────────────────────────────────────────────────────────────────
	{
		ID: "git-github-crash-course", Title: "Git & GitHub Crash Course",
		Description:    "Free comprehensive Git and GitHub tutorial. Covers version control fundamentals, branching, merging.",
		URL:            "https://www.youtube.com/watch?v=RGOj5yH7evk",
		Provider:       "YouTube/freeCodeCamp", ResourceType: "video", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 1, DurationLabel: "1 hour",
		Skills: []string{"git", "github"}, PrimarySkill: "git",
		Rating: 4.80, RatingCount: 5000000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "pro-git-book", Title: "Pro Git Book",
		Description:    "The entire Pro Git book, written by Scott Chacon and Ben Straub. Free to read online.",
		URL:            "https://git-scm.com/book/en/v2",
		Provider:       "Git SCM", ResourceType: "book", Difficulty: "all_levels",
		CostType: "free", CostUSD: 0, DurationHours: 15, DurationLabel: "15 hours",
		Skills: []string{"git"}, PrimarySkill: "git",
		Rating: 4.90, RatingCount: 100000, HasCertificate: false, HasHandsOn: false, IsVerified: true,
	},

	// ── System Design ────────────────────────────────────────────────────────
	{
		ID: "system-design-primer", Title: "System Design Primer (GitHub)",
		Description:    "Free open-source guide to learning how to design large-scale systems.",
		URL:            "https://github.com/donnemartin/system-design-primer",
		Provider:       "GitHub", ResourceType: "documentation", Difficulty: "intermediate",
		CostType: "free", CostUSD: 0, DurationHours: 20, DurationLabel: "20 hours",
		Skills: []string{"system design"}, PrimarySkill: "system design",
		Rating: 4.90, RatingCount: 200000, HasCertificate: false, HasHandsOn: false, IsVerified: true,
	},
	{
		ID: "grokking-system-design", Title: "Grokking the System Design Interview",
		Description:    "Learn how to design large-scale systems. Covers load balancing, caching, databases, microservices.",
		URL:            "https://www.educative.io/courses/grokking-the-system-design-interview",
		Provider:       "Educative", ResourceType: "course", Difficulty: "intermediate",
		CostType: "subscription", CostUSD: 59.00, DurationHours: 20, DurationLabel: "20 hours",
		Skills: []string{"system design"}, PrimarySkill: "system design",
		Rating: 4.70, RatingCount: 100000, HasCertificate: false, HasHandsOn: false, IsVerified: true,
	},

	// ── Algorithms & Data Structures ─────────────────────────────────────────
	{
		ID: "leetcode", Title: "LeetCode Practice Platform",
		Description:    "The leading platform for coding interview preparation. 2000+ problems covering all major topics.",
		URL:            "https://leetcode.com/",
		Provider:       "LeetCode", ResourceType: "practice", Difficulty: "all_levels",
		CostType: "freemium", CostUSD: 35.00, DurationHours: 0, DurationLabel: "Self-paced",
		Skills: []string{"algorithms", "data structures"}, PrimarySkill: "algorithms",
		Rating: 4.70, RatingCount: 2000000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "algorithms-stanford", Title: "Algorithms Specialization (Stanford)",
		Description:    "Learn algorithms from Stanford University. Covers divide and conquer, graph algorithms, dynamic programming.",
		URL:            "https://www.coursera.org/specializations/algorithms",
		Provider:       "Coursera", ResourceType: "course", Difficulty: "advanced",
		CostType: "free_audit", CostUSD: 49.00, DurationHours: 60, DurationLabel: "4 months",
		Skills: []string{"algorithms", "data structures"}, PrimarySkill: "algorithms",
		Rating: 4.80, RatingCount: 200000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Linux / Bash ─────────────────────────────────────────────────────────
	{
		ID: "linux-command-line-book", Title: "The Linux Command Line (Book)",
		Description:    "A complete introduction to the Linux command line. Free to read online.",
		URL:            "https://linuxcommand.org/tlcl.php",
		Provider:       "LinuxCommand.org", ResourceType: "book", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 20, DurationLabel: "20 hours",
		Skills: []string{"linux", "bash"}, PrimarySkill: "linux",
		Rating: 4.80, RatingCount: 50000, HasCertificate: false, HasHandsOn: false, IsVerified: true,
	},

	// ── Java ─────────────────────────────────────────────────────────────────
	{
		ID: "java-masterclass", Title: "Java Programming Masterclass",
		Description:    "Learn Java in this complete masterclass. Covers Java 17, OOP, data structures, algorithms.",
		URL:            "https://www.udemy.com/course/java-the-complete-java-developer-course/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 80, DurationLabel: "80 hours",
		Skills: []string{"java"}, PrimarySkill: "java",
		Rating: 4.60, RatingCount: 300000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Spring Boot ──────────────────────────────────────────────────────────
	{
		ID: "spring-boot-masterclass", Title: "Spring Boot 3 & Spring Framework 6 Masterclass",
		Description:    "Master Spring Boot 3 and Spring Framework 6. Covers REST APIs, Spring Security, Spring Data JPA.",
		URL:            "https://www.udemy.com/course/spring-boot-tutorial-for-beginners/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 45, DurationLabel: "45 hours",
		Skills: []string{"spring boot", "java"}, PrimarySkill: "spring boot",
		Rating: 4.60, RatingCount: 80000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Rust ─────────────────────────────────────────────────────────────────
	{
		ID: "rust-book", Title: "The Rust Programming Language (Book)",
		Description:    "The official Rust book. Free to read online. Covers ownership, borrowing, lifetimes.",
		URL:            "https://doc.rust-lang.org/book/",
		Provider:       "Rust Foundation", ResourceType: "book", Difficulty: "intermediate",
		CostType: "free", CostUSD: 0, DurationHours: 30, DurationLabel: "30 hours",
		Skills: []string{"rust"}, PrimarySkill: "rust",
		Rating: 4.90, RatingCount: 200000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "rustlings", Title: "Rustlings",
		Description:    "Small exercises to get you used to reading and writing Rust code.",
		URL:            "https://github.com/rust-lang/rustlings",
		Provider:       "Rust Foundation", ResourceType: "practice", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 10, DurationLabel: "10 hours",
		Skills: []string{"rust"}, PrimarySkill: "rust",
		Rating: 4.80, RatingCount: 50000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},

	// ── MongoDB ──────────────────────────────────────────────────────────────
	{
		ID: "mongodb-university", Title: "MongoDB University: Introduction to MongoDB",
		Description:    "Free official MongoDB course. Learn the fundamentals of MongoDB, CRUD operations, and indexing.",
		URL:            "https://learn.mongodb.com/learning-paths/introduction-to-mongodb",
		Provider:       "MongoDB", ResourceType: "course", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 8, DurationLabel: "8 hours",
		Skills: []string{"mongodb"}, PrimarySkill: "mongodb",
		Rating: 4.70, RatingCount: 200000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Redis ────────────────────────────────────────────────────────────────
	{
		ID: "redis-university", Title: "Redis University: RU101 Introduction to Redis",
		Description:    "Free official Redis course. Learn Redis data structures, commands, and use cases.",
		URL:            "https://university.redis.com/courses/ru101/",
		Provider:       "Redis", ResourceType: "course", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 8, DurationLabel: "8 hours",
		Skills: []string{"redis"}, PrimarySkill: "redis",
		Rating: 4.60, RatingCount: 50000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Kafka ────────────────────────────────────────────────────────────────
	{
		ID: "kafka-beginners", Title: "Apache Kafka Series - Learn Apache Kafka for Beginners",
		Description:    "Learn Apache Kafka from scratch. Covers producers, consumers, topics, partitions.",
		URL:            "https://www.udemy.com/course/apache-kafka/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 8, DurationLabel: "8 hours",
		Skills: []string{"kafka"}, PrimarySkill: "kafka",
		Rating: 4.70, RatingCount: 80000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Agile / Scrum ────────────────────────────────────────────────────────
	{
		ID: "psm1-cert", Title: "Professional Scrum Master I (PSM I)",
		Description:    "The PSM I certification validates knowledge of the Scrum framework.",
		URL:            "https://www.scrum.org/assessments/professional-scrum-master-i-certification",
		Provider:       "Scrum.org", ResourceType: "certification", Difficulty: "beginner",
		CostType: "paid", CostUSD: 150.00, DurationHours: 20, DurationLabel: "20 hours prep",
		Skills: []string{"scrum", "agile"}, PrimarySkill: "scrum",
		Rating: 4.70, RatingCount: 100000, HasCertificate: true, HasHandsOn: false, IsVerified: true,
	},
	{
		ID: "agile-jira-coursera", Title: "Agile with Atlassian Jira",
		Description:    "Free course on Agile development with Jira. Covers Scrum, Kanban, sprints, backlogs.",
		URL:            "https://www.coursera.org/learn/agile-atlassian-jira",
		Provider:       "Coursera", ResourceType: "course", Difficulty: "beginner",
		CostType: "free_audit", CostUSD: 49.00, DurationHours: 6, DurationLabel: "6 hours",
		Skills: []string{"agile", "scrum"}, PrimarySkill: "agile",
		Rating: 4.50, RatingCount: 50000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Data Engineering ─────────────────────────────────────────────────────
	{
		ID: "data-engineering-zoomcamp", Title: "Data Engineering Zoomcamp",
		Description:    "Free 9-week data engineering course. Covers containerization, workflow orchestration, data warehousing.",
		URL:            "https://github.com/DataTalksClub/data-engineering-zoomcamp",
		Provider:       "DataTalks.Club", ResourceType: "course", Difficulty: "intermediate",
		CostType: "free", CostUSD: 0, DurationHours: 80, DurationLabel: "9 weeks",
		Skills: []string{"data engineering", "kafka", "docker"}, PrimarySkill: "data engineering",
		Rating: 4.80, RatingCount: 30000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},
	{
		ID: "pyspark-udemy", Title: "Apache Spark with Python - PySpark",
		Description:    "Learn Apache Spark with Python. Covers RDDs, DataFrames, Spark SQL, Spark Streaming.",
		URL:            "https://www.udemy.com/course/spark-and-python-for-big-data-with-pyspark/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 10, DurationLabel: "10 hours",
		Skills: []string{"spark", "python"}, PrimarySkill: "spark",
		Rating: 4.60, RatingCount: 50000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── CI/CD ────────────────────────────────────────────────────────────────
	{
		ID: "github-actions-complete", Title: "GitHub Actions: The Complete Guide",
		Description:    "Master GitHub Actions for CI/CD. Covers workflows, jobs, steps, actions, secrets.",
		URL:            "https://www.udemy.com/course/github-actions-the-complete-guide/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 10, DurationLabel: "10 hours",
		Skills: []string{"github actions", "ci/cd"}, PrimarySkill: "github actions",
		Rating: 4.60, RatingCount: 20000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── NLP ──────────────────────────────────────────────────────────────────
	{
		ID: "nlp-specialization", Title: "Natural Language Processing Specialization",
		Description:    "Break into NLP. Covers sentiment analysis, machine translation, question answering.",
		URL:            "https://www.coursera.org/specializations/natural-language-processing",
		Provider:       "Coursera", ResourceType: "course", Difficulty: "advanced",
		CostType: "free_audit", CostUSD: 49.00, DurationHours: 80, DurationLabel: "4 months",
		Skills: []string{"nlp", "deep learning", "python"}, PrimarySkill: "nlp",
		Rating: 4.80, RatingCount: 100000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Exercism ─────────────────────────────────────────────────────────────
	{
		ID: "exercism", Title: "Exercism: Code Practice and Mentorship",
		Description:    "Free platform for coding exercises in 60+ programming languages with mentorship.",
		URL:            "https://exercism.org/",
		Provider:       "Exercism", ResourceType: "practice", Difficulty: "all_levels",
		CostType: "free", CostUSD: 0, DurationHours: 0, DurationLabel: "Self-paced",
		Skills: []string{"algorithms", "python", "go", "javascript", "rust"}, PrimarySkill: "algorithms",
		Rating: 4.80, RatingCount: 200000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},

	// ── Full Stack ───────────────────────────────────────────────────────────
	{
		ID: "odin-project", Title: "The Odin Project: Full Stack JavaScript",
		Description:    "Free full-stack web development curriculum. Covers HTML, CSS, JavaScript, Node.js, React.",
		URL:            "https://www.theodinproject.com/paths/full-stack-javascript",
		Provider:       "The Odin Project", ResourceType: "course", Difficulty: "beginner",
		CostType: "free", CostUSD: 0, DurationHours: 1000, DurationLabel: "1000+ hours",
		Skills: []string{"javascript", "react", "node.js", "html", "css"}, PrimarySkill: "javascript",
		Rating: 4.90, RatingCount: 100000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},

	// ── Ansible ──────────────────────────────────────────────────────────────
	{
		ID: "ansible-beginner", Title: "Ansible for the Absolute Beginner",
		Description:    "Learn Ansible from scratch. Covers playbooks, roles, variables, templates.",
		URL:            "https://www.udemy.com/course/learn-ansible/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 5, DurationLabel: "5 hours",
		Skills: []string{"ansible"}, PrimarySkill: "ansible",
		Rating: 4.60, RatingCount: 50000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Elasticsearch ────────────────────────────────────────────────────────
	{
		ID: "elasticsearch-complete", Title: "Complete Guide to Elasticsearch",
		Description:    "Learn Elasticsearch from scratch. Covers indexing, searching, aggregations, mappings.",
		URL:            "https://www.udemy.com/course/elasticsearch-complete-guide/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 15, DurationLabel: "15 hours",
		Skills: []string{"elasticsearch"}, PrimarySkill: "elasticsearch",
		Rating: 4.70, RatingCount: 30000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Communication ────────────────────────────────────────────────────────
	{
		ID: "communication-skills-engineers", Title: "Communication Skills for Engineers",
		Description:    "Improve your technical communication skills. Covers written communication, presentations, code reviews.",
		URL:            "https://www.coursera.org/learn/communication-skills-engineers",
		Provider:       "Coursera", ResourceType: "course", Difficulty: "beginner",
		CostType: "free_audit", CostUSD: 49.00, DurationHours: 12, DurationLabel: "4 weeks",
		Skills: []string{"communication"}, PrimarySkill: "communication",
		Rating: 4.40, RatingCount: 20000, HasCertificate: true, HasHandsOn: false, IsVerified: true,
	},

	// ── Vue ──────────────────────────────────────────────────────────────────
	{
		ID: "vue-complete-guide", Title: "Vue - The Complete Guide",
		Description:    "Learn Vue.js from the ground up. Covers Vue 3, Composition API, Vue Router, Pinia.",
		URL:            "https://www.udemy.com/course/vuejs-2-the-complete-guide/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 32, DurationLabel: "32 hours",
		Skills: []string{"vue", "javascript"}, PrimarySkill: "vue",
		Rating: 4.70, RatingCount: 100000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Angular ──────────────────────────────────────────────────────────────
	{
		ID: "angular-complete-guide", Title: "Angular - The Complete Guide (2024 Edition)",
		Description:    "Master Angular 17. Covers components, directives, services, routing, forms, HTTP, and RxJS.",
		URL:            "https://www.udemy.com/course/the-complete-guide-to-angular-2/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 36, DurationLabel: "36 hours",
		Skills: []string{"angular", "typescript"}, PrimarySkill: "angular",
		Rating: 4.60, RatingCount: 200000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── C# ───────────────────────────────────────────────────────────────────
	{
		ID: "csharp-basics", Title: "C# Basics for Beginners: Learn C# Fundamentals",
		Description:    "Learn C# from scratch. Covers C# syntax, OOP, LINQ, async/await, and .NET fundamentals.",
		URL:            "https://www.udemy.com/course/csharp-tutorial-for-beginners/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 5, DurationLabel: "5 hours",
		Skills: []string{"c#", ".net"}, PrimarySkill: "c#",
		Rating: 4.50, RatingCount: 100000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Kotlin ───────────────────────────────────────────────────────────────
	{
		ID: "android-kotlin-masterclass", Title: "Android Kotlin Development Masterclass",
		Description:    "Learn Android development with Kotlin. Covers Jetpack Compose, MVVM, Room, Retrofit.",
		URL:            "https://www.udemy.com/course/android-oreo-kotlin-app-masterclass/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 60, DurationLabel: "60 hours",
		Skills: []string{"kotlin", "android"}, PrimarySkill: "kotlin",
		Rating: 4.60, RatingCount: 50000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Swift ────────────────────────────────────────────────────────────────
	{
		ID: "ios-swift-bootcamp", Title: "iOS & Swift - The Complete iOS App Development Bootcamp",
		Description:    "Learn iOS development with Swift. Covers UIKit, SwiftUI, Core Data, networking.",
		URL:            "https://www.udemy.com/course/ios-13-app-development-bootcamp/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "beginner",
		CostType: "paid", CostUSD: 19.99, DurationHours: 55, DurationLabel: "55 hours",
		Skills: []string{"swift", "ios"}, PrimarySkill: "swift",
		Rating: 4.80, RatingCount: 100000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── GraphQL ──────────────────────────────────────────────────────────────
	{
		ID: "graphql-react-guide", Title: "GraphQL with React: The Complete Developers Guide",
		Description:    "Learn GraphQL with React. Covers schemas, queries, mutations, subscriptions, Apollo Client.",
		URL:            "https://www.udemy.com/course/graphql-with-react-course/",
		Provider:       "Udemy", ResourceType: "course", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 19.99, DurationHours: 13, DurationLabel: "13 hours",
		Skills: []string{"graphql", "react"}, PrimarySkill: "graphql",
		Rating: 4.60, RatingCount: 30000, HasCertificate: true, HasHandsOn: true, IsVerified: true,
	},

	// ── Security ─────────────────────────────────────────────────────────────
	{
		ID: "comptia-security-plus", Title: "CompTIA Security+",
		Description:    "The CompTIA Security+ certification validates baseline cybersecurity skills.",
		URL:            "https://www.comptia.org/certifications/security",
		Provider:       "CompTIA", ResourceType: "certification", Difficulty: "intermediate",
		CostType: "paid", CostUSD: 392.00, DurationHours: 60, DurationLabel: "60 hours prep",
		Skills: []string{"security"}, PrimarySkill: "security",
		Rating: 4.70, RatingCount: 100000, HasCertificate: true, HasHandsOn: false, IsVerified: true,
	},

	// ── LangChain / LLMs ─────────────────────────────────────────────────────
	{
		ID: "langchain-llm-course", Title: "LangChain for LLM Application Development",
		Description:    "Free short course on building LLM-powered applications with LangChain. Covers chains, agents, memory, and RAG.",
		URL:            "https://www.deeplearning.ai/short-courses/langchain-for-llm-application-development/",
		Provider:       "DeepLearning.AI", ResourceType: "course", Difficulty: "intermediate",
		CostType: "free", CostUSD: 0, DurationHours: 1, DurationLabel: "1 hour",
		Skills: []string{"langchain", "python", "machine learning"}, PrimarySkill: "langchain",
		Rating: 4.70, RatingCount: 50000, HasCertificate: false, HasHandsOn: true, IsVerified: true,
	},
}

// skillAliases maps common skill aliases to their canonical names in the catalog.
var skillAliases = map[string]string{
	"golang":     "go",
	"js":         "javascript",
	"ts":         "typescript",
	"node":       "node.js",
	"nodejs":     "node.js",
	"react.js":   "react",
	"reactjs":    "react",
	"vue.js":     "vue",
	"vuejs":      "vue",
	"angular.js": "angular",
	"angularjs":  "angular",
	"postgres":   "postgresql",
	"psql":       "postgresql",
	"k8s":        "kubernetes",
	"py":         "python",
	"rb":         "ruby",
	"cpp":        "c++",
	"csharp":     "c#",
	"dotnet":     ".net",
	"net":        ".net",
	"ml":         "machine learning",
	"dl":         "deep learning",
	"tf":         "tensorflow",
	"spring":     "spring boot",
}
