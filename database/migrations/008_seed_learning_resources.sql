-- Migration 008: Seed learning resources with curated high-quality content
-- Populates the learning resource catalog with 80+ resources covering
-- the most common technical skills for software engineering roles.

BEGIN;

-- ─────────────────────────────────────────────────────────────────────────────
-- Resource Providers
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO resource_providers (id, name, normalized_name, website_url, description) VALUES
    ('10000000-0000-0000-0000-000000000001', 'Coursera',       'coursera',       'https://www.coursera.org',       'Online learning platform with university-level courses'),
    ('10000000-0000-0000-0000-000000000002', 'Udemy',          'udemy',          'https://www.udemy.com',          'Marketplace for online courses on any topic'),
    ('10000000-0000-0000-0000-000000000003', 'edX',            'edx',            'https://www.edx.org',            'Online learning platform from Harvard and MIT'),
    ('10000000-0000-0000-0000-000000000004', 'Pluralsight',    'pluralsight',    'https://www.pluralsight.com',    'Technology skills platform for developers'),
    ('10000000-0000-0000-0000-000000000005', 'AWS Training',   'aws-training',   'https://aws.amazon.com/training','Official Amazon Web Services training and certification'),
    ('10000000-0000-0000-0000-000000000006', 'Google Cloud',   'google-cloud',   'https://cloud.google.com/learn', 'Official Google Cloud training and certification'),
    ('10000000-0000-0000-0000-000000000007', 'Microsoft Learn','microsoft-learn','https://learn.microsoft.com',    'Official Microsoft learning platform'),
    ('10000000-0000-0000-0000-000000000008', 'freeCodeCamp',   'freecodecamp',   'https://www.freecodecamp.org',   'Free coding bootcamp with certifications'),
    ('10000000-0000-0000-0000-000000000009', 'The Odin Project','the-odin-project','https://www.theodinproject.com','Free full-stack web development curriculum'),
    ('10000000-0000-0000-0000-000000000010', 'LeetCode',       'leetcode',       'https://leetcode.com',           'Coding interview preparation platform'),
    ('10000000-0000-0000-0000-000000000011', 'HackerRank',     'hackerrank',     'https://www.hackerrank.com',     'Coding challenges and technical assessments'),
    ('10000000-0000-0000-0000-000000000012', 'Exercism',       'exercism',       'https://exercism.org',           'Free coding exercises with mentorship'),
    ('10000000-0000-0000-0000-000000000013', 'YouTube',        'youtube',        'https://www.youtube.com',        'Free video content platform'),
    ('10000000-0000-0000-0000-000000000014', 'Official Docs',  'official-docs',  NULL,                             'Official language and framework documentation'),
    ('10000000-0000-0000-0000-000000000015', 'O''Reilly',      'oreilly',        'https://www.oreilly.com',        'Technical books and online learning'),
    ('10000000-0000-0000-0000-000000000016', 'Acloud.guru',    'acloud-guru',    'https://acloudguru.com',         'Cloud computing training platform'),
    ('10000000-0000-0000-0000-000000000017', 'Codecademy',     'codecademy',     'https://www.codecademy.com',     'Interactive coding courses for beginners'),
    ('10000000-0000-0000-0000-000000000018', 'Fast.ai',        'fastai',         'https://www.fast.ai',            'Free practical deep learning courses'),
    ('10000000-0000-0000-0000-000000000019', 'Kaggle',         'kaggle',         'https://www.kaggle.com',         'Data science competitions and free courses'),
    ('10000000-0000-0000-0000-000000000020', 'Linux Foundation','linux-foundation','https://training.linuxfoundation.org','Linux and open source training');

-- ─────────────────────────────────────────────────────────────────────────────
-- Learning Resources
-- ─────────────────────────────────────────────────────────────────────────────

-- ── Python ──────────────────────────────────────────────────────────────────
INSERT INTO learning_resources (id, title, slug, description, url, provider_id, resource_type, difficulty, cost_type, cost_amount, duration_hours, duration_label, has_certificate, has_hands_on, rating, rating_count, is_verified, is_featured) VALUES
    ('20000000-0000-0000-0000-000000000001',
     'Python for Everybody Specialization',
     'python-for-everybody-coursera',
     'Learn to program and analyze data with Python. Develop programs to gather, clean, analyze, and visualize data. Covers Python basics, data structures, web access, databases, and data visualization.',
     'https://www.coursera.org/specializations/python',
     '10000000-0000-0000-0000-000000000001', 'course', 'beginner', 'free_audit', 49.00, 80, '8 months', TRUE, TRUE, 4.80, 1200000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000002',
     'Complete Python Bootcamp: From Zero to Hero',
     'complete-python-bootcamp-udemy',
     'Learn Python like a professional. Start from the basics and go all the way to creating your own applications and games. Covers Python 3, OOP, decorators, generators, and more.',
     'https://www.udemy.com/course/complete-python-bootcamp/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 22, '22 hours', TRUE, TRUE, 4.60, 500000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000003',
     'Python Official Documentation',
     'python-official-docs',
     'The official Python 3 documentation including tutorial, library reference, and language reference. The definitive resource for Python programming.',
     'https://docs.python.org/3/',
     '10000000-0000-0000-0000-000000000014', 'documentation', 'all_levels', 'free', NULL, NULL, NULL, FALSE, FALSE, NULL, NULL, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000004',
     'Automate the Boring Stuff with Python',
     'automate-boring-stuff-python',
     'A practical programming book for office workers, academics, and administrators who want to improve their productivity. Free to read online.',
     'https://automatetheboringstuff.com/',
     '10000000-0000-0000-0000-000000000014', 'book', 'beginner', 'free', NULL, 20, '20 hours', FALSE, TRUE, 4.70, 50000, TRUE, FALSE),

-- ── Go (Golang) ─────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000005',
     'Go: The Complete Developer''s Guide',
     'go-complete-developers-guide-udemy',
     'Master the fundamentals and advanced features of the Go programming language. Covers goroutines, channels, interfaces, testing, and building web services.',
     'https://www.udemy.com/course/go-the-complete-developers-guide/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 9, '9 hours', TRUE, TRUE, 4.60, 45000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000006',
     'A Tour of Go',
     'a-tour-of-go',
     'An interactive introduction to Go programming language. Covers all the basics of Go with hands-on exercises directly in the browser.',
     'https://go.dev/tour/',
     '10000000-0000-0000-0000-000000000014', 'documentation', 'beginner', 'free', NULL, 4, '4 hours', FALSE, TRUE, 4.80, 100000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000007',
     'Go by Example',
     'go-by-example',
     'Hands-on introduction to Go using annotated example programs. Covers all major Go features with runnable examples.',
     'https://gobyexample.com/',
     '10000000-0000-0000-0000-000000000014', 'documentation', 'beginner', 'free', NULL, 8, '8 hours', FALSE, TRUE, 4.90, 200000, TRUE, FALSE),

-- ── JavaScript / TypeScript ──────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000008',
     'The Complete JavaScript Course 2024',
     'complete-javascript-course-udemy',
     'The modern JavaScript course for everyone. Master JavaScript with projects, challenges and theory. Covers ES6+, OOP, async/await, and modern tooling.',
     'https://www.udemy.com/course/the-complete-javascript-course/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 69, '69 hours', TRUE, TRUE, 4.70, 350000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000009',
     'JavaScript: Understanding the Weird Parts',
     'javascript-understanding-weird-parts',
     'An advanced JavaScript course that dives deep into the language internals. Covers closures, prototypal inheritance, the event loop, and more.',
     'https://www.udemy.com/course/understand-javascript/',
     '10000000-0000-0000-0000-000000000002', 'course', 'advanced', 'paid', 19.99, 12, '12 hours', TRUE, TRUE, 4.70, 200000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000010',
     'TypeScript: The Complete Developer''s Guide',
     'typescript-complete-developers-guide',
     'Master TypeScript by building real projects. Covers type system, generics, decorators, and integration with React and Node.js.',
     'https://www.udemy.com/course/typescript-the-complete-developers-guide/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 27, '27 hours', TRUE, TRUE, 4.60, 80000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000011',
     'freeCodeCamp JavaScript Algorithms and Data Structures',
     'freecodecamp-javascript-algorithms',
     'Free certification covering JavaScript fundamentals, ES6, regular expressions, debugging, data structures, and algorithm scripting.',
     'https://www.freecodecamp.org/learn/javascript-algorithms-and-data-structures/',
     '10000000-0000-0000-0000-000000000008', 'course', 'beginner', 'free', NULL, 300, '300 hours', TRUE, TRUE, 4.50, 500000, TRUE, TRUE),

-- ── React ────────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000012',
     'React - The Complete Guide 2024',
     'react-complete-guide-udemy',
     'Dive in and learn React.js from scratch. Learn Reactjs, Hooks, Redux, React Router, Next.js, Best Practices and way more.',
     'https://www.udemy.com/course/react-the-complete-guide-incl-redux/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 68, '68 hours', TRUE, TRUE, 4.60, 250000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000013',
     'React Official Documentation',
     'react-official-docs',
     'The official React documentation with interactive examples, tutorials, and API reference. Covers React 18 features including hooks and concurrent mode.',
     'https://react.dev/',
     '10000000-0000-0000-0000-000000000014', 'documentation', 'all_levels', 'free', NULL, NULL, NULL, FALSE, TRUE, NULL, NULL, TRUE, FALSE),

-- ── Node.js ──────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000014',
     'The Complete Node.js Developer Course',
     'complete-nodejs-developer-course',
     'Learn Node.js by building real-world applications with Node, Express, MongoDB, Jest, and more. Covers REST APIs, authentication, and deployment.',
     'https://www.udemy.com/course/the-complete-nodejs-developer-course-2/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 35, '35 hours', TRUE, TRUE, 4.60, 150000, TRUE, FALSE),

-- ── Python Data Science / ML ─────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000015',
     'Machine Learning Specialization',
     'machine-learning-specialization-coursera',
     'Andrew Ng''s updated ML course. Covers supervised learning, unsupervised learning, and best practices for ML development.',
     'https://www.coursera.org/specializations/machine-learning-introduction',
     '10000000-0000-0000-0000-000000000001', 'course', 'intermediate', 'free_audit', 49.00, 90, '3 months', TRUE, TRUE, 4.90, 500000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000016',
     'Deep Learning Specialization',
     'deep-learning-specialization-coursera',
     'Become a Deep Learning expert. Master deep neural networks, CNNs, RNNs, LSTMs, and transformers. Build and train deep neural networks.',
     'https://www.coursera.org/specializations/deep-learning',
     '10000000-0000-0000-0000-000000000001', 'course', 'advanced', 'free_audit', 49.00, 120, '5 months', TRUE, TRUE, 4.90, 400000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000017',
     'Practical Deep Learning for Coders',
     'practical-deep-learning-fastai',
     'Free course from fast.ai. Learn deep learning with PyTorch and fastai. Covers computer vision, NLP, tabular data, and collaborative filtering.',
     'https://course.fast.ai/',
     '10000000-0000-0000-0000-000000000018', 'course', 'intermediate', 'free', NULL, 30, '30 hours', FALSE, TRUE, 4.80, 100000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000018',
     'Kaggle Learn: Python',
     'kaggle-learn-python',
     'Free micro-courses on Python, pandas, machine learning, and more. Hands-on exercises with immediate feedback.',
     'https://www.kaggle.com/learn',
     '10000000-0000-0000-0000-000000000019', 'course', 'beginner', 'free', NULL, 5, '5 hours', TRUE, TRUE, 4.70, 200000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000019',
     'TensorFlow Developer Certificate',
     'tensorflow-developer-certificate',
     'Official TensorFlow certification. Demonstrates proficiency in using TensorFlow to solve deep learning and ML problems.',
     'https://www.tensorflow.org/certificate',
     '10000000-0000-0000-0000-000000000014', 'certification', 'intermediate', 'paid', 100.00, 40, '40 hours prep', TRUE, TRUE, 4.60, 20000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000020',
     'PyTorch for Deep Learning Bootcamp',
     'pytorch-deep-learning-bootcamp-udemy',
     'Learn PyTorch for deep learning. Covers tensors, neural networks, CNNs, RNNs, and transfer learning with hands-on projects.',
     'https://www.udemy.com/course/pytorch-for-deep-learning-bootcamp/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 17, '17 hours', TRUE, TRUE, 4.60, 30000, TRUE, FALSE),

-- ── SQL / Databases ──────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000021',
     'The Complete SQL Bootcamp',
     'complete-sql-bootcamp-udemy',
     'Become an expert at SQL. Learn how to read and write complex queries to a database using one of the most in-demand skills.',
     'https://www.udemy.com/course/the-complete-sql-bootcamp/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 9, '9 hours', TRUE, TRUE, 4.70, 200000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000022',
     'PostgreSQL: The Complete Developer''s Guide',
     'postgresql-complete-developers-guide',
     'Master PostgreSQL with this comprehensive course. Covers advanced queries, indexing, performance tuning, and administration.',
     'https://www.udemy.com/course/sql-and-postgresql/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 22, '22 hours', TRUE, TRUE, 4.70, 50000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000023',
     'SQLZoo Interactive SQL Tutorial',
     'sqlzoo-interactive-tutorial',
     'Free interactive SQL tutorial with exercises. Covers SELECT, INSERT, UPDATE, DELETE, and advanced SQL features.',
     'https://sqlzoo.net/',
     '10000000-0000-0000-0000-000000000014', 'documentation', 'beginner', 'free', NULL, 10, '10 hours', FALSE, TRUE, 4.50, 500000, TRUE, FALSE),

-- ── Docker & Kubernetes ──────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000024',
     'Docker and Kubernetes: The Complete Guide',
     'docker-kubernetes-complete-guide',
     'Build, test, and deploy Docker applications with Kubernetes while learning production-style development workflows.',
     'https://www.udemy.com/course/docker-and-kubernetes-the-complete-guide/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 22, '22 hours', TRUE, TRUE, 4.60, 100000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000025',
     'Docker Official Documentation',
     'docker-official-docs',
     'Official Docker documentation covering installation, getting started, guides, and reference material.',
     'https://docs.docker.com/',
     '10000000-0000-0000-0000-000000000014', 'documentation', 'all_levels', 'free', NULL, NULL, NULL, FALSE, TRUE, NULL, NULL, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000026',
     'Certified Kubernetes Administrator (CKA)',
     'certified-kubernetes-administrator',
     'The CKA certification is designed to ensure that certification holders have the skills, knowledge, and competency to perform the responsibilities of Kubernetes administrators.',
     'https://training.linuxfoundation.org/certification/certified-kubernetes-administrator-cka/',
     '10000000-0000-0000-0000-000000000020', 'certification', 'advanced', 'paid', 395.00, 60, '60 hours prep', TRUE, TRUE, 4.70, 50000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000027',
     'Kubernetes for Developers (LFD259)',
     'kubernetes-for-developers-lfd259',
     'Official Linux Foundation course for Kubernetes application developers. Covers pods, deployments, services, and CI/CD.',
     'https://training.linuxfoundation.org/training/kubernetes-for-developers/',
     '10000000-0000-0000-0000-000000000020', 'course', 'intermediate', 'paid', 299.00, 30, '30 hours', TRUE, TRUE, 4.50, 10000, TRUE, FALSE),

-- ── AWS ──────────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000028',
     'AWS Certified Solutions Architect – Associate',
     'aws-certified-solutions-architect-associate',
     'The AWS Certified Solutions Architect – Associate certification validates the ability to design and implement distributed systems on AWS.',
     'https://aws.amazon.com/certification/certified-solutions-architect-associate/',
     '10000000-0000-0000-0000-000000000005', 'certification', 'intermediate', 'paid', 300.00, 80, '80 hours prep', TRUE, TRUE, 4.80, 200000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000029',
     'Ultimate AWS Certified Solutions Architect Associate',
     'ultimate-aws-saa-udemy',
     'Pass the AWS Certified Solutions Architect Associate certification. Covers all AWS services with hands-on labs.',
     'https://www.udemy.com/course/aws-certified-solutions-architect-associate-saa-c03/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 27, '27 hours', TRUE, TRUE, 4.70, 300000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000030',
     'AWS Certified Developer – Associate',
     'aws-certified-developer-associate',
     'Validates technical expertise in developing and maintaining applications on the AWS platform.',
     'https://aws.amazon.com/certification/certified-developer-associate/',
     '10000000-0000-0000-0000-000000000005', 'certification', 'intermediate', 'paid', 300.00, 60, '60 hours prep', TRUE, TRUE, 4.70, 100000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000031',
     'AWS Cloud Practitioner Essentials',
     'aws-cloud-practitioner-essentials',
     'Free foundational course for AWS Cloud Practitioner certification. Covers cloud concepts, AWS services, security, and pricing.',
     'https://aws.amazon.com/training/digital/aws-cloud-practitioner-essentials/',
     '10000000-0000-0000-0000-000000000005', 'course', 'beginner', 'free', NULL, 6, '6 hours', FALSE, FALSE, 4.60, 500000, TRUE, FALSE),

-- ── Google Cloud ─────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000032',
     'Google Cloud Professional Cloud Architect',
     'gcp-professional-cloud-architect',
     'The Google Cloud Professional Cloud Architect certification validates expertise in designing, developing, and managing robust, secure, scalable, highly available, and dynamic solutions.',
     'https://cloud.google.com/certification/cloud-architect',
     '10000000-0000-0000-0000-000000000006', 'certification', 'advanced', 'paid', 200.00, 100, '100 hours prep', TRUE, TRUE, 4.70, 50000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000033',
     'Google Cloud Associate Cloud Engineer',
     'gcp-associate-cloud-engineer',
     'Validates ability to deploy applications, monitor operations, and manage enterprise solutions on Google Cloud.',
     'https://cloud.google.com/certification/cloud-engineer',
     '10000000-0000-0000-0000-000000000006', 'certification', 'intermediate', 'paid', 200.00, 60, '60 hours prep', TRUE, TRUE, 4.60, 30000, TRUE, FALSE),

-- ── Azure ────────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000034',
     'AZ-900: Microsoft Azure Fundamentals',
     'az-900-azure-fundamentals',
     'Foundational knowledge of cloud services and how those services are provided with Microsoft Azure.',
     'https://learn.microsoft.com/en-us/certifications/azure-fundamentals/',
     '10000000-0000-0000-0000-000000000007', 'certification', 'beginner', 'paid', 165.00, 20, '20 hours prep', TRUE, FALSE, 4.70, 100000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000035',
     'AZ-204: Developing Solutions for Microsoft Azure',
     'az-204-developing-azure-solutions',
     'Validates expertise in designing, building, testing, and maintaining cloud applications and services on Microsoft Azure.',
     'https://learn.microsoft.com/en-us/certifications/azure-developer/',
     '10000000-0000-0000-0000-000000000007', 'certification', 'intermediate', 'paid', 165.00, 60, '60 hours prep', TRUE, TRUE, 4.60, 40000, TRUE, FALSE),

-- ── Terraform / Infrastructure as Code ──────────────────────────────────────
    ('20000000-0000-0000-0000-000000000036',
     'HashiCorp Certified: Terraform Associate',
     'hashicorp-terraform-associate',
     'Validates knowledge of infrastructure automation using Terraform. Covers HCL, state management, modules, and providers.',
     'https://www.hashicorp.com/certification/terraform-associate',
     '10000000-0000-0000-0000-000000000014', 'certification', 'intermediate', 'paid', 70.50, 40, '40 hours prep', TRUE, TRUE, 4.70, 30000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000037',
     'Terraform: From Beginner to Master',
     'terraform-beginner-to-master-udemy',
     'Learn Terraform from scratch. Covers infrastructure as code, AWS provisioning, modules, state management, and CI/CD integration.',
     'https://www.udemy.com/course/terraform-beginner-to-advanced/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 12, '12 hours', TRUE, TRUE, 4.60, 40000, TRUE, FALSE),

-- ── Git / Version Control ────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000038',
     'Git & GitHub Crash Course',
     'git-github-crash-course-freecodecamp',
     'Free comprehensive Git and GitHub tutorial. Covers version control fundamentals, branching, merging, pull requests, and collaboration workflows.',
     'https://www.youtube.com/watch?v=RGOj5yH7evk',
     '10000000-0000-0000-0000-000000000013', 'video', 'beginner', 'free', NULL, 1, '1 hour', FALSE, TRUE, 4.80, 5000000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000039',
     'Pro Git Book',
     'pro-git-book',
     'The entire Pro Git book, written by Scott Chacon and Ben Straub. Free to read online. The definitive guide to Git.',
     'https://git-scm.com/book/en/v2',
     '10000000-0000-0000-0000-000000000014', 'book', 'all_levels', 'free', NULL, 15, '15 hours', FALSE, FALSE, 4.90, 100000, TRUE, FALSE),

-- ── System Design ────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000040',
     'Grokking the System Design Interview',
     'grokking-system-design-interview',
     'Learn how to design large-scale systems. Covers load balancing, caching, databases, microservices, and real-world system design examples.',
     'https://www.educative.io/courses/grokking-the-system-design-interview',
     '10000000-0000-0000-0000-000000000014', 'course', 'intermediate', 'subscription', 59.00, 20, '20 hours', FALSE, FALSE, 4.70, 100000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000041',
     'System Design Primer (GitHub)',
     'system-design-primer-github',
     'Free open-source guide to learning how to design large-scale systems. Covers scalability, availability, consistency, and common patterns.',
     'https://github.com/donnemartin/system-design-primer',
     '10000000-0000-0000-0000-000000000014', 'documentation', 'intermediate', 'free', NULL, 20, '20 hours', FALSE, FALSE, 4.90, 200000, TRUE, TRUE),

-- ── Algorithms & Data Structures ─────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000042',
     'LeetCode Practice Platform',
     'leetcode-practice',
     'The leading platform for coding interview preparation. 2000+ problems covering arrays, strings, trees, graphs, dynamic programming, and more.',
     'https://leetcode.com/',
     '10000000-0000-0000-0000-000000000010', 'practice', 'all_levels', 'freemium', 35.00, NULL, 'Self-paced', FALSE, TRUE, 4.70, 2000000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000043',
     'HackerRank Problem Solving',
     'hackerrank-problem-solving',
     'Practice coding challenges in algorithms, data structures, mathematics, and more. Used by companies for technical screening.',
     'https://www.hackerrank.com/domains/algorithms',
     '10000000-0000-0000-0000-000000000011', 'practice', 'all_levels', 'free', NULL, NULL, 'Self-paced', FALSE, TRUE, 4.40, 1000000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000044',
     'Algorithms Specialization (Stanford)',
     'algorithms-specialization-stanford-coursera',
     'Learn algorithms from Stanford University. Covers divide and conquer, graph algorithms, greedy algorithms, and dynamic programming.',
     'https://www.coursera.org/specializations/algorithms',
     '10000000-0000-0000-0000-000000000001', 'course', 'advanced', 'free_audit', 49.00, 60, '4 months', TRUE, TRUE, 4.80, 200000, TRUE, FALSE),

-- ── Linux / Bash ─────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000045',
     'Linux Command Line Basics',
     'linux-command-line-basics-udacity',
     'Learn the Linux command line. Covers file system navigation, file manipulation, permissions, processes, and shell scripting.',
     'https://www.udacity.com/course/linux-command-line-basics--ud595',
     '10000000-0000-0000-0000-000000000014', 'course', 'beginner', 'free', NULL, 5, '5 hours', FALSE, TRUE, 4.50, 100000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000046',
     'The Linux Command Line (Book)',
     'the-linux-command-line-book',
     'A complete introduction to the Linux command line. Free to read online. Covers shell basics, file system, processes, and shell scripting.',
     'https://linuxcommand.org/tlcl.php',
     '10000000-0000-0000-0000-000000000014', 'book', 'beginner', 'free', NULL, 20, '20 hours', FALSE, FALSE, 4.80, 50000, TRUE, FALSE),

-- ── Java ─────────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000047',
     'Java Programming Masterclass',
     'java-programming-masterclass-udemy',
     'Learn Java in this complete masterclass. Covers Java 17, OOP, data structures, algorithms, JavaFX, and more.',
     'https://www.udemy.com/course/java-the-complete-java-developer-course/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 80, '80 hours', TRUE, TRUE, 4.60, 300000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000048',
     'Oracle Certified Professional Java SE 17',
     'oracle-certified-java-se-17',
     'The OCP Java SE 17 Developer certification validates expertise in Java programming. Covers Java fundamentals, OOP, generics, collections, and more.',
     'https://education.oracle.com/java-se-17-developer/pexam_1Z0-829',
     '10000000-0000-0000-0000-000000000014', 'certification', 'advanced', 'paid', 245.00, 80, '80 hours prep', TRUE, TRUE, 4.60, 20000, TRUE, FALSE),

-- ── Spring Boot ──────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000049',
     'Spring Boot 3 & Spring Framework 6 Masterclass',
     'spring-boot-3-masterclass-udemy',
     'Master Spring Boot 3 and Spring Framework 6. Covers REST APIs, Spring Security, Spring Data JPA, microservices, and testing.',
     'https://www.udemy.com/course/spring-boot-tutorial-for-beginners/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 45, '45 hours', TRUE, TRUE, 4.60, 80000, TRUE, FALSE),

-- ── Rust ─────────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000050',
     'The Rust Programming Language (Book)',
     'the-rust-programming-language-book',
     'The official Rust book, affectionately nicknamed "the book". Free to read online. Covers ownership, borrowing, lifetimes, and more.',
     'https://doc.rust-lang.org/book/',
     '10000000-0000-0000-0000-000000000014', 'book', 'intermediate', 'free', NULL, 30, '30 hours', FALSE, TRUE, 4.90, 200000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000051',
     'Rustlings',
     'rustlings-exercises',
     'Small exercises to get you used to reading and writing Rust code. Covers all major Rust concepts with hands-on practice.',
     'https://github.com/rust-lang/rustlings',
     '10000000-0000-0000-0000-000000000014', 'practice', 'beginner', 'free', NULL, 10, '10 hours', FALSE, TRUE, 4.80, 50000, TRUE, FALSE),

-- ── MongoDB ──────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000052',
     'MongoDB University: M001 MongoDB Basics',
     'mongodb-university-m001',
     'Free official MongoDB course. Learn the fundamentals of MongoDB, including CRUD operations, aggregation, and indexing.',
     'https://learn.mongodb.com/learning-paths/introduction-to-mongodb',
     '10000000-0000-0000-0000-000000000014', 'course', 'beginner', 'free', NULL, 8, '8 hours', TRUE, TRUE, 4.70, 200000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000053',
     'MongoDB Certified Developer Associate',
     'mongodb-certified-developer-associate',
     'Validates expertise in building applications with MongoDB. Covers data modeling, CRUD, aggregation, indexing, and replication.',
     'https://www.mongodb.com/products/certifications/developer',
     '10000000-0000-0000-0000-000000000014', 'certification', 'intermediate', 'paid', 150.00, 40, '40 hours prep', TRUE, TRUE, 4.60, 10000, TRUE, FALSE),

-- ── Redis ────────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000054',
     'Redis University: RU101 Introduction to Redis',
     'redis-university-ru101',
     'Free official Redis course. Learn Redis data structures, commands, and use cases for caching, sessions, and pub/sub.',
     'https://university.redis.com/courses/ru101/',
     '10000000-0000-0000-0000-000000000014', 'course', 'beginner', 'free', NULL, 8, '8 hours', TRUE, TRUE, 4.60, 50000, TRUE, FALSE),

-- ── Kafka ────────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000055',
     'Apache Kafka Series - Learn Apache Kafka for Beginners',
     'apache-kafka-beginners-udemy',
     'Learn Apache Kafka from scratch. Covers producers, consumers, topics, partitions, Kafka Streams, and Kafka Connect.',
     'https://www.udemy.com/course/apache-kafka/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 8, '8 hours', TRUE, TRUE, 4.70, 80000, TRUE, FALSE),

-- ── Agile / Scrum ────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000056',
     'Professional Scrum Master I (PSM I)',
     'professional-scrum-master-psm1',
     'The PSM I certification validates knowledge of the Scrum framework. Covers Scrum theory, events, artifacts, and roles.',
     'https://www.scrum.org/assessments/professional-scrum-master-i-certification',
     '10000000-0000-0000-0000-000000000014', 'certification', 'beginner', 'paid', 150.00, 20, '20 hours prep', TRUE, FALSE, 4.70, 100000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000057',
     'Agile with Atlassian Jira',
     'agile-atlassian-jira-coursera',
     'Free course on Agile development with Jira. Covers Scrum, Kanban, sprints, backlogs, and Jira workflows.',
     'https://www.coursera.org/learn/agile-atlassian-jira',
     '10000000-0000-0000-0000-000000000001', 'course', 'beginner', 'free_audit', 49.00, 6, '6 hours', TRUE, TRUE, 4.50, 50000, TRUE, FALSE),

-- ── Data Engineering ─────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000058',
     'Data Engineering Zoomcamp',
     'data-engineering-zoomcamp',
     'Free 9-week data engineering course. Covers containerization, workflow orchestration, data warehousing, batch processing, and streaming.',
     'https://github.com/DataTalksClub/data-engineering-zoomcamp',
     '10000000-0000-0000-0000-000000000014', 'course', 'intermediate', 'free', NULL, 80, '9 weeks', TRUE, TRUE, 4.80, 30000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000059',
     'Apache Spark with Python - PySpark',
     'apache-spark-pyspark-udemy',
     'Learn Apache Spark with Python. Covers RDDs, DataFrames, Spark SQL, Spark Streaming, and MLlib.',
     'https://www.udemy.com/course/spark-and-python-for-big-data-with-pyspark/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 10, '10 hours', TRUE, TRUE, 4.60, 50000, TRUE, FALSE),

-- ── Security ─────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000060',
     'CompTIA Security+',
     'comptia-security-plus',
     'The CompTIA Security+ certification validates baseline cybersecurity skills. Covers threats, vulnerabilities, cryptography, and network security.',
     'https://www.comptia.org/certifications/security',
     '10000000-0000-0000-0000-000000000014', 'certification', 'intermediate', 'paid', 392.00, 60, '60 hours prep', TRUE, FALSE, 4.70, 100000, TRUE, FALSE),

-- ── API Design ───────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000061',
     'REST API Design, Development & Management',
     'rest-api-design-udemy',
     'Learn REST API design best practices. Covers HTTP methods, status codes, authentication, versioning, and documentation with Swagger/OpenAPI.',
     'https://www.udemy.com/course/rest-api/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 10, '10 hours', TRUE, TRUE, 4.50, 30000, TRUE, FALSE),

-- ── CI/CD ────────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000062',
     'GitHub Actions: The Complete Guide',
     'github-actions-complete-guide-udemy',
     'Master GitHub Actions for CI/CD. Covers workflows, jobs, steps, actions, secrets, and deployment to cloud platforms.',
     'https://www.udemy.com/course/github-actions-the-complete-guide/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 10, '10 hours', TRUE, TRUE, 4.60, 20000, TRUE, FALSE),

-- ── Microservices ────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000063',
     'Microservices with Node JS and React',
     'microservices-nodejs-react-udemy',
     'Build a microservices app with Node, React, Docker, and Kubernetes. Covers event-driven architecture, NATS Streaming, and CI/CD.',
     'https://www.udemy.com/course/microservices-with-node-js-and-react/',
     '10000000-0000-0000-0000-000000000002', 'course', 'advanced', 'paid', 19.99, 54, '54 hours', TRUE, TRUE, 4.60, 50000, TRUE, FALSE),

-- ── Soft Skills ──────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000064',
     'Communication Skills for Engineers',
     'communication-skills-engineers-coursera',
     'Improve your technical communication skills. Covers written communication, presentations, code reviews, and cross-functional collaboration.',
     'https://www.coursera.org/learn/communication-skills-engineers',
     '10000000-0000-0000-0000-000000000001', 'course', 'beginner', 'free_audit', 49.00, 12, '4 weeks', TRUE, FALSE, 4.40, 20000, TRUE, FALSE),

-- ── C# / .NET ────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000065',
     'C# Basics for Beginners: Learn C# Fundamentals',
     'csharp-basics-beginners-udemy',
     'Learn C# from scratch. Covers C# syntax, OOP, LINQ, async/await, and .NET fundamentals.',
     'https://www.udemy.com/course/csharp-tutorial-for-beginners/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 5, '5 hours', TRUE, TRUE, 4.50, 100000, TRUE, FALSE),

-- ── Vue.js ───────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000066',
     'Vue - The Complete Guide (incl. Router & Composition API)',
     'vue-complete-guide-udemy',
     'Learn Vue.js from the ground up. Covers Vue 3, Composition API, Vue Router, Vuex/Pinia, and building real-world applications.',
     'https://www.udemy.com/course/vuejs-2-the-complete-guide/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 32, '32 hours', TRUE, TRUE, 4.70, 100000, TRUE, FALSE),

-- ── Angular ──────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000067',
     'Angular - The Complete Guide (2024 Edition)',
     'angular-complete-guide-udemy',
     'Master Angular 17. Covers components, directives, services, routing, forms, HTTP, and RxJS.',
     'https://www.udemy.com/course/the-complete-guide-to-angular-2/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 36, '36 hours', TRUE, TRUE, 4.60, 200000, TRUE, FALSE),

-- ── Elasticsearch ────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000068',
     'Complete Guide to Elasticsearch',
     'complete-guide-elasticsearch-udemy',
     'Learn Elasticsearch from scratch. Covers indexing, searching, aggregations, mappings, and integration with Kibana.',
     'https://www.udemy.com/course/elasticsearch-complete-guide/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 15, '15 hours', TRUE, TRUE, 4.70, 30000, TRUE, FALSE),

-- ── NLP / LLMs ───────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000069',
     'Natural Language Processing Specialization',
     'nlp-specialization-coursera',
     'Break into NLP. Covers sentiment analysis, machine translation, question answering, and text summarization using deep learning.',
     'https://www.coursera.org/specializations/natural-language-processing',
     '10000000-0000-0000-0000-000000000001', 'course', 'advanced', 'free_audit', 49.00, 80, '4 months', TRUE, TRUE, 4.80, 100000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000070',
     'LangChain for LLM Application Development',
     'langchain-llm-application-development',
     'Free short course on building LLM-powered applications with LangChain. Covers chains, agents, memory, and RAG.',
     'https://www.deeplearning.ai/short-courses/langchain-for-llm-application-development/',
     '10000000-0000-0000-0000-000000000014', 'course', 'intermediate', 'free', NULL, 1, '1 hour', FALSE, TRUE, 4.70, 50000, TRUE, TRUE),

-- ── Exercism ─────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000071',
     'Exercism: Code Practice and Mentorship',
     'exercism-code-practice',
     'Free platform for coding exercises in 60+ programming languages with mentorship. Great for learning new languages through practice.',
     'https://exercism.org/',
     '10000000-0000-0000-0000-000000000012', 'practice', 'all_levels', 'free', NULL, NULL, 'Self-paced', FALSE, TRUE, 4.80, 200000, TRUE, FALSE),

-- ── Full Stack ───────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000072',
     'The Odin Project: Full Stack JavaScript',
     'the-odin-project-full-stack',
     'Free full-stack web development curriculum. Covers HTML, CSS, JavaScript, Node.js, React, and databases.',
     'https://www.theodinproject.com/paths/full-stack-javascript',
     '10000000-0000-0000-0000-000000000009', 'course', 'beginner', 'free', NULL, 1000, '1000+ hours', FALSE, TRUE, 4.90, 100000, TRUE, TRUE),

    ('20000000-0000-0000-0000-000000000073',
     'freeCodeCamp: Responsive Web Design',
     'freecodecamp-responsive-web-design',
     'Free certification covering HTML, CSS, flexbox, grid, and responsive design principles.',
     'https://www.freecodecamp.org/learn/2022/responsive-web-design/',
     '10000000-0000-0000-0000-000000000008', 'course', 'beginner', 'free', NULL, 300, '300 hours', TRUE, TRUE, 4.60, 500000, TRUE, FALSE),

-- ── Pluralsight Paths ────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000074',
     'Go Path on Pluralsight',
     'go-path-pluralsight',
     'Comprehensive Go learning path on Pluralsight. Covers Go fundamentals, concurrency, testing, and building web services.',
     'https://www.pluralsight.com/paths/go',
     '10000000-0000-0000-0000-000000000004', 'course', 'intermediate', 'subscription', 29.00, 20, '20 hours', TRUE, TRUE, 4.50, 10000, TRUE, FALSE),

    ('20000000-0000-0000-0000-000000000075',
     'Python Path on Pluralsight',
     'python-path-pluralsight',
     'Comprehensive Python learning path. Covers Python fundamentals, OOP, testing, data science, and web development.',
     'https://www.pluralsight.com/paths/python',
     '10000000-0000-0000-0000-000000000004', 'course', 'beginner', 'subscription', 29.00, 30, '30 hours', TRUE, TRUE, 4.50, 20000, TRUE, FALSE),

-- ── Ansible ──────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000076',
     'Ansible for the Absolute Beginner',
     'ansible-absolute-beginner-udemy',
     'Learn Ansible from scratch. Covers playbooks, roles, variables, templates, and automating infrastructure configuration.',
     'https://www.udemy.com/course/learn-ansible/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 5, '5 hours', TRUE, TRUE, 4.60, 50000, TRUE, FALSE),

-- ── GraphQL ──────────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000077',
     'GraphQL with React: The Complete Developers Guide',
     'graphql-react-complete-guide-udemy',
     'Learn GraphQL with React. Covers schemas, queries, mutations, subscriptions, Apollo Client, and building full-stack apps.',
     'https://www.udemy.com/course/graphql-with-react-course/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 13, '13 hours', TRUE, TRUE, 4.60, 30000, TRUE, FALSE),

-- ── Prometheus / Observability ───────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000078',
     'Prometheus and Grafana: The Complete Guide',
     'prometheus-grafana-complete-guide-udemy',
     'Learn monitoring and observability with Prometheus and Grafana. Covers metrics, alerting, dashboards, and Kubernetes monitoring.',
     'https://www.udemy.com/course/prometheus-course/',
     '10000000-0000-0000-0000-000000000002', 'course', 'intermediate', 'paid', 19.99, 8, '8 hours', TRUE, TRUE, 4.60, 15000, TRUE, FALSE),

-- ── Swift / iOS ──────────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000079',
     'iOS & Swift - The Complete iOS App Development Bootcamp',
     'ios-swift-complete-bootcamp-udemy',
     'Learn iOS development with Swift. Covers UIKit, SwiftUI, Core Data, networking, and publishing to the App Store.',
     'https://www.udemy.com/course/ios-13-app-development-bootcamp/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 55, '55 hours', TRUE, TRUE, 4.80, 100000, TRUE, FALSE),

-- ── Kotlin / Android ─────────────────────────────────────────────────────────
    ('20000000-0000-0000-0000-000000000080',
     'Android Kotlin Development Masterclass',
     'android-kotlin-masterclass-udemy',
     'Learn Android development with Kotlin. Covers Jetpack Compose, MVVM, Room, Retrofit, and publishing to Google Play.',
     'https://www.udemy.com/course/android-oreo-kotlin-app-masterclass/',
     '10000000-0000-0000-0000-000000000002', 'course', 'beginner', 'paid', 19.99, 60, '60 hours', TRUE, TRUE, 4.60, 50000, TRUE, FALSE);

-- ─────────────────────────────────────────────────────────────────────────────
-- Resource Skills Mappings
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO resource_skills (resource_id, skill_name, normalized_name, is_primary, coverage_level) VALUES
    -- Python for Everybody
    ('20000000-0000-0000-0000-000000000001', 'Python', 'python', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000001', 'Data Analysis', 'data analysis', FALSE, 'beginner'),
    ('20000000-0000-0000-0000-000000000001', 'SQL', 'sql', FALSE, 'beginner'),
    -- Complete Python Bootcamp
    ('20000000-0000-0000-0000-000000000002', 'Python', 'python', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000002', 'OOP', 'oop', FALSE, 'intermediate'),
    -- Python Docs
    ('20000000-0000-0000-0000-000000000003', 'Python', 'python', TRUE, 'all_levels'),
    -- Automate the Boring Stuff
    ('20000000-0000-0000-0000-000000000004', 'Python', 'python', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000004', 'Automation', 'automation', FALSE, 'beginner'),
    -- Go Complete Guide
    ('20000000-0000-0000-0000-000000000005', 'Go', 'go', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000005', 'Concurrency', 'concurrency', FALSE, 'intermediate'),
    -- Tour of Go
    ('20000000-0000-0000-0000-000000000006', 'Go', 'go', TRUE, 'beginner'),
    -- Go by Example
    ('20000000-0000-0000-0000-000000000007', 'Go', 'go', TRUE, 'beginner'),
    -- Complete JavaScript Course
    ('20000000-0000-0000-0000-000000000008', 'JavaScript', 'javascript', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000008', 'ES6', 'es6', FALSE, 'intermediate'),
    -- JavaScript Weird Parts
    ('20000000-0000-0000-0000-000000000009', 'JavaScript', 'javascript', TRUE, 'advanced'),
    -- TypeScript Complete Guide
    ('20000000-0000-0000-0000-000000000010', 'TypeScript', 'typescript', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000010', 'JavaScript', 'javascript', FALSE, 'intermediate'),
    -- freeCodeCamp JS
    ('20000000-0000-0000-0000-000000000011', 'JavaScript', 'javascript', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000011', 'Algorithms', 'algorithms', FALSE, 'beginner'),
    -- React Complete Guide
    ('20000000-0000-0000-0000-000000000012', 'React', 'react', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000012', 'Redux', 'redux', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000012', 'JavaScript', 'javascript', FALSE, 'intermediate'),
    -- React Docs
    ('20000000-0000-0000-0000-000000000013', 'React', 'react', TRUE, 'all_levels'),
    -- Node.js Complete Course
    ('20000000-0000-0000-0000-000000000014', 'Node.js', 'node.js', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000014', 'Express', 'express', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000014', 'MongoDB', 'mongodb', FALSE, 'beginner'),
    -- Machine Learning Specialization
    ('20000000-0000-0000-0000-000000000015', 'Machine Learning', 'machine learning', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000015', 'Python', 'python', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000015', 'TensorFlow', 'tensorflow', FALSE, 'beginner'),
    -- Deep Learning Specialization
    ('20000000-0000-0000-0000-000000000016', 'Deep Learning', 'deep learning', TRUE, 'advanced'),
    ('20000000-0000-0000-0000-000000000016', 'TensorFlow', 'tensorflow', FALSE, 'advanced'),
    ('20000000-0000-0000-0000-000000000016', 'Python', 'python', FALSE, 'advanced'),
    -- Practical Deep Learning fastai
    ('20000000-0000-0000-0000-000000000017', 'Deep Learning', 'deep learning', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000017', 'PyTorch', 'pytorch', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000017', 'Python', 'python', FALSE, 'intermediate'),
    -- Kaggle Learn
    ('20000000-0000-0000-0000-000000000018', 'Python', 'python', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000018', 'Pandas', 'pandas', FALSE, 'beginner'),
    -- TensorFlow Certificate
    ('20000000-0000-0000-0000-000000000019', 'TensorFlow', 'tensorflow', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000019', 'Deep Learning', 'deep learning', FALSE, 'intermediate'),
    -- PyTorch Bootcamp
    ('20000000-0000-0000-0000-000000000020', 'PyTorch', 'pytorch', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000020', 'Deep Learning', 'deep learning', FALSE, 'intermediate'),
    -- SQL Bootcamp
    ('20000000-0000-0000-0000-000000000021', 'SQL', 'sql', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000021', 'PostgreSQL', 'postgresql', FALSE, 'beginner'),
    -- PostgreSQL Complete Guide
    ('20000000-0000-0000-0000-000000000022', 'PostgreSQL', 'postgresql', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000022', 'SQL', 'sql', FALSE, 'advanced'),
    -- SQLZoo
    ('20000000-0000-0000-0000-000000000023', 'SQL', 'sql', TRUE, 'beginner'),
    -- Docker & Kubernetes
    ('20000000-0000-0000-0000-000000000024', 'Docker', 'docker', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000024', 'Kubernetes', 'kubernetes', TRUE, 'intermediate'),
    -- Docker Docs
    ('20000000-0000-0000-0000-000000000025', 'Docker', 'docker', TRUE, 'all_levels'),
    -- CKA
    ('20000000-0000-0000-0000-000000000026', 'Kubernetes', 'kubernetes', TRUE, 'advanced'),
    -- Kubernetes for Developers
    ('20000000-0000-0000-0000-000000000027', 'Kubernetes', 'kubernetes', TRUE, 'intermediate'),
    -- AWS SAA Cert
    ('20000000-0000-0000-0000-000000000028', 'AWS', 'aws', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000028', 'Cloud Architecture', 'cloud architecture', FALSE, 'intermediate'),
    -- AWS SAA Udemy
    ('20000000-0000-0000-0000-000000000029', 'AWS', 'aws', TRUE, 'intermediate'),
    -- AWS Developer Cert
    ('20000000-0000-0000-0000-000000000030', 'AWS', 'aws', TRUE, 'intermediate'),
    -- AWS Cloud Practitioner
    ('20000000-0000-0000-0000-000000000031', 'AWS', 'aws', TRUE, 'beginner'),
    -- GCP Architect
    ('20000000-0000-0000-0000-000000000032', 'GCP', 'gcp', TRUE, 'advanced'),
    ('20000000-0000-0000-0000-000000000032', 'Cloud Architecture', 'cloud architecture', FALSE, 'advanced'),
    -- GCP Engineer
    ('20000000-0000-0000-0000-000000000033', 'GCP', 'gcp', TRUE, 'intermediate'),
    -- AZ-900
    ('20000000-0000-0000-0000-000000000034', 'Azure', 'azure', TRUE, 'beginner'),
    -- AZ-204
    ('20000000-0000-0000-0000-000000000035', 'Azure', 'azure', TRUE, 'intermediate'),
    -- Terraform Cert
    ('20000000-0000-0000-0000-000000000036', 'Terraform', 'terraform', TRUE, 'intermediate'),
    -- Terraform Course
    ('20000000-0000-0000-0000-000000000037', 'Terraform', 'terraform', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000037', 'AWS', 'aws', FALSE, 'intermediate'),
    -- Git Crash Course
    ('20000000-0000-0000-0000-000000000038', 'Git', 'git', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000038', 'GitHub', 'github', FALSE, 'beginner'),
    -- Pro Git Book
    ('20000000-0000-0000-0000-000000000039', 'Git', 'git', TRUE, 'all_levels'),
    -- System Design Interview
    ('20000000-0000-0000-0000-000000000040', 'System Design', 'system design', TRUE, 'intermediate'),
    -- System Design Primer
    ('20000000-0000-0000-0000-000000000041', 'System Design', 'system design', TRUE, 'intermediate'),
    -- LeetCode
    ('20000000-0000-0000-0000-000000000042', 'Algorithms', 'algorithms', TRUE, 'all_levels'),
    ('20000000-0000-0000-0000-000000000042', 'Data Structures', 'data structures', TRUE, 'all_levels'),
    -- HackerRank
    ('20000000-0000-0000-0000-000000000043', 'Algorithms', 'algorithms', TRUE, 'all_levels'),
    -- Algorithms Specialization
    ('20000000-0000-0000-0000-000000000044', 'Algorithms', 'algorithms', TRUE, 'advanced'),
    ('20000000-0000-0000-0000-000000000044', 'Data Structures', 'data structures', FALSE, 'advanced'),
    -- Linux Command Line
    ('20000000-0000-0000-0000-000000000045', 'Linux', 'linux', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000045', 'Bash', 'bash', FALSE, 'beginner'),
    -- Linux Command Line Book
    ('20000000-0000-0000-0000-000000000046', 'Linux', 'linux', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000046', 'Bash', 'bash', FALSE, 'beginner'),
    -- Java Masterclass
    ('20000000-0000-0000-0000-000000000047', 'Java', 'java', TRUE, 'intermediate'),
    -- Oracle Java Cert
    ('20000000-0000-0000-0000-000000000048', 'Java', 'java', TRUE, 'advanced'),
    -- Spring Boot
    ('20000000-0000-0000-0000-000000000049', 'Spring Boot', 'spring boot', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000049', 'Java', 'java', FALSE, 'intermediate'),
    -- Rust Book
    ('20000000-0000-0000-0000-000000000050', 'Rust', 'rust', TRUE, 'intermediate'),
    -- Rustlings
    ('20000000-0000-0000-0000-000000000051', 'Rust', 'rust', TRUE, 'beginner'),
    -- MongoDB University
    ('20000000-0000-0000-0000-000000000052', 'MongoDB', 'mongodb', TRUE, 'beginner'),
    -- MongoDB Cert
    ('20000000-0000-0000-0000-000000000053', 'MongoDB', 'mongodb', TRUE, 'intermediate'),
    -- Redis University
    ('20000000-0000-0000-0000-000000000054', 'Redis', 'redis', TRUE, 'beginner'),
    -- Kafka Course
    ('20000000-0000-0000-0000-000000000055', 'Kafka', 'kafka', TRUE, 'beginner'),
    -- PSM I
    ('20000000-0000-0000-0000-000000000056', 'Scrum', 'scrum', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000056', 'Agile', 'agile', FALSE, 'beginner'),
    -- Agile with Jira
    ('20000000-0000-0000-0000-000000000057', 'Agile', 'agile', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000057', 'Scrum', 'scrum', FALSE, 'beginner'),
    -- Data Engineering Zoomcamp
    ('20000000-0000-0000-0000-000000000058', 'Data Engineering', 'data engineering', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000058', 'Kafka', 'kafka', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000058', 'Docker', 'docker', FALSE, 'intermediate'),
    -- PySpark
    ('20000000-0000-0000-0000-000000000059', 'Spark', 'spark', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000059', 'Python', 'python', FALSE, 'intermediate'),
    -- CompTIA Security+
    ('20000000-0000-0000-0000-000000000060', 'Security', 'security', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000060', 'Networking', 'networking', FALSE, 'intermediate'),
    -- REST API Design
    ('20000000-0000-0000-0000-000000000061', 'REST API', 'rest api', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000061', 'API Design', 'api design', FALSE, 'intermediate'),
    -- GitHub Actions
    ('20000000-0000-0000-0000-000000000062', 'GitHub Actions', 'github actions', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000062', 'CI/CD', 'ci/cd', FALSE, 'intermediate'),
    -- Microservices
    ('20000000-0000-0000-0000-000000000063', 'Microservices', 'microservices', TRUE, 'advanced'),
    ('20000000-0000-0000-0000-000000000063', 'Node.js', 'node.js', FALSE, 'advanced'),
    ('20000000-0000-0000-0000-000000000063', 'Docker', 'docker', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000063', 'Kubernetes', 'kubernetes', FALSE, 'intermediate'),
    -- Communication Skills
    ('20000000-0000-0000-0000-000000000064', 'Communication', 'communication', TRUE, 'intermediate'),
    -- C# Basics
    ('20000000-0000-0000-0000-000000000065', 'C#', 'c#', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000065', '.NET', '.net', FALSE, 'beginner'),
    -- Vue Complete Guide
    ('20000000-0000-0000-0000-000000000066', 'Vue', 'vue', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000066', 'JavaScript', 'javascript', FALSE, 'intermediate'),
    -- Angular Complete Guide
    ('20000000-0000-0000-0000-000000000067', 'Angular', 'angular', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000067', 'TypeScript', 'typescript', FALSE, 'intermediate'),
    -- Elasticsearch
    ('20000000-0000-0000-0000-000000000068', 'Elasticsearch', 'elasticsearch', TRUE, 'intermediate'),
    -- NLP Specialization
    ('20000000-0000-0000-0000-000000000069', 'NLP', 'nlp', TRUE, 'advanced'),
    ('20000000-0000-0000-0000-000000000069', 'Deep Learning', 'deep learning', FALSE, 'advanced'),
    ('20000000-0000-0000-0000-000000000069', 'Python', 'python', FALSE, 'advanced'),
    -- LangChain
    ('20000000-0000-0000-0000-000000000070', 'LangChain', 'langchain', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000070', 'Python', 'python', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000070', 'Machine Learning', 'machine learning', FALSE, 'intermediate'),
    -- Exercism
    ('20000000-0000-0000-0000-000000000071', 'Algorithms', 'algorithms', TRUE, 'all_levels'),
    -- The Odin Project
    ('20000000-0000-0000-0000-000000000072', 'JavaScript', 'javascript', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000072', 'React', 'react', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000072', 'Node.js', 'node.js', FALSE, 'intermediate'),
    -- freeCodeCamp Web Design
    ('20000000-0000-0000-0000-000000000073', 'HTML', 'html', TRUE, 'beginner'),
    ('20000000-0000-0000-0000-000000000073', 'CSS', 'css', TRUE, 'beginner'),
    -- Go Pluralsight
    ('20000000-0000-0000-0000-000000000074', 'Go', 'go', TRUE, 'intermediate'),
    -- Python Pluralsight
    ('20000000-0000-0000-0000-000000000075', 'Python', 'python', TRUE, 'intermediate'),
    -- Ansible
    ('20000000-0000-0000-0000-000000000076', 'Ansible', 'ansible', TRUE, 'beginner'),
    -- GraphQL
    ('20000000-0000-0000-0000-000000000077', 'GraphQL', 'graphql', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000077', 'React', 'react', FALSE, 'intermediate'),
    -- Prometheus/Grafana
    ('20000000-0000-0000-0000-000000000078', 'Prometheus', 'prometheus', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000078', 'Grafana', 'grafana', FALSE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000078', 'Kubernetes', 'kubernetes', FALSE, 'intermediate'),
    -- iOS Swift
    ('20000000-0000-0000-0000-000000000079', 'Swift', 'swift', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000079', 'iOS', 'ios', FALSE, 'intermediate'),
    -- Android Kotlin
    ('20000000-0000-0000-0000-000000000080', 'Kotlin', 'kotlin', TRUE, 'intermediate'),
    ('20000000-0000-0000-0000-000000000080', 'Android', 'android', FALSE, 'intermediate');

-- ─────────────────────────────────────────────────────────────────────────────
-- Learning Paths
-- ─────────────────────────────────────────────────────────────────────────────
INSERT INTO learning_paths (id, title, slug, description, target_role, target_skill, difficulty, estimated_hours, is_featured) VALUES
    ('30000000-0000-0000-0000-000000000001',
     'Python Developer Path',
     'python-developer-path',
     'Go from Python beginner to job-ready Python developer. Covers Python fundamentals, OOP, data structures, and web development.',
     'Python Developer', 'Python', 'beginner', 150, TRUE),

    ('30000000-0000-0000-0000-000000000002',
     'Machine Learning Engineer Path',
     'machine-learning-engineer-path',
     'Become a machine learning engineer. Covers Python, statistics, ML algorithms, deep learning, and MLOps.',
     'Machine Learning Engineer', 'Machine Learning', 'advanced', 300, TRUE),

    ('30000000-0000-0000-0000-000000000003',
     'AWS Cloud Engineer Path',
     'aws-cloud-engineer-path',
     'Prepare for AWS certifications and cloud engineering roles. Covers AWS fundamentals, architecture, and DevOps.',
     'Cloud Engineer', 'AWS', 'intermediate', 120, TRUE),

    ('30000000-0000-0000-0000-000000000004',
     'Full Stack JavaScript Developer Path',
     'full-stack-javascript-path',
     'Become a full stack JavaScript developer. Covers HTML, CSS, JavaScript, React, Node.js, and databases.',
     'Full Stack Developer', 'JavaScript', 'beginner', 400, TRUE),

    ('30000000-0000-0000-0000-000000000005',
     'DevOps Engineer Path',
     'devops-engineer-path',
     'Master DevOps practices. Covers Docker, Kubernetes, CI/CD, Terraform, and cloud platforms.',
     'DevOps Engineer', 'Kubernetes', 'intermediate', 200, TRUE),

    ('30000000-0000-0000-0000-000000000006',
     'Go Backend Developer Path',
     'go-backend-developer-path',
     'Build high-performance backend services with Go. Covers Go fundamentals, concurrency, REST APIs, and databases.',
     'Backend Engineer', 'Go', 'intermediate', 120, FALSE),

    ('30000000-0000-0000-0000-000000000007',
     'Data Engineering Path',
     'data-engineering-path',
     'Become a data engineer. Covers Python, SQL, Spark, Kafka, and cloud data platforms.',
     'Data Engineer', 'Data Engineering', 'intermediate', 250, FALSE);

-- ─────────────────────────────────────────────────────────────────────────────
-- Learning Path Resources (ordered steps)
-- ─────────────────────────────────────────────────────────────────────────────

-- Python Developer Path
INSERT INTO learning_path_resources (path_id, resource_id, step_order, is_required, notes) VALUES
    ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 1, TRUE, 'Start with Python for Everybody for a solid foundation'),
    ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000004', 2, FALSE, 'Practical Python projects to reinforce learning'),
    ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000002', 3, TRUE, 'Deep dive into Python with the complete bootcamp'),
    ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000021', 4, TRUE, 'Learn SQL for database interactions'),
    ('30000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000042', 5, FALSE, 'Practice algorithms and data structures');

-- Machine Learning Engineer Path
INSERT INTO learning_path_resources (path_id, resource_id, step_order, is_required, notes) VALUES
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000001', 1, TRUE, 'Python fundamentals are essential for ML'),
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000018', 2, TRUE, 'Kaggle micro-courses for quick Python and pandas skills'),
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000015', 3, TRUE, 'Core ML algorithms with Andrew Ng'),
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000016', 4, TRUE, 'Deep learning specialization'),
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000017', 5, FALSE, 'Practical deep learning with fastai'),
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000020', 6, FALSE, 'PyTorch for production ML'),
    ('30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000019', 7, FALSE, 'TensorFlow certification to validate skills');

-- AWS Cloud Engineer Path
INSERT INTO learning_path_resources (path_id, resource_id, step_order, is_required, notes) VALUES
    ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000031', 1, TRUE, 'Start with AWS Cloud Practitioner fundamentals'),
    ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000029', 2, TRUE, 'Comprehensive SAA course on Udemy'),
    ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000028', 3, TRUE, 'AWS Solutions Architect Associate certification'),
    ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000024', 4, FALSE, 'Docker and Kubernetes for cloud deployments'),
    ('30000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000037', 5, FALSE, 'Terraform for infrastructure as code');

-- Full Stack JavaScript Path
INSERT INTO learning_path_resources (path_id, resource_id, step_order, is_required, notes) VALUES
    ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000073', 1, TRUE, 'HTML and CSS fundamentals'),
    ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000011', 2, TRUE, 'JavaScript fundamentals with freeCodeCamp'),
    ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000008', 3, TRUE, 'Complete JavaScript course'),
    ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000012', 4, TRUE, 'React for frontend development'),
    ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000014', 5, TRUE, 'Node.js for backend development'),
    ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000021', 6, TRUE, 'SQL for database management'),
    ('30000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000072', 7, FALSE, 'The Odin Project for additional practice');

-- DevOps Engineer Path
INSERT INTO learning_path_resources (path_id, resource_id, step_order, is_required, notes) VALUES
    ('30000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000038', 1, TRUE, 'Git fundamentals for version control'),
    ('30000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000045', 2, TRUE, 'Linux command line basics'),
    ('30000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000024', 3, TRUE, 'Docker and Kubernetes fundamentals'),
    ('30000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000062', 4, TRUE, 'CI/CD with GitHub Actions'),
    ('30000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000037', 5, TRUE, 'Infrastructure as code with Terraform'),
    ('30000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000026', 6, FALSE, 'CKA certification for Kubernetes expertise'),
    ('30000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000078', 7, FALSE, 'Monitoring with Prometheus and Grafana');

-- Go Backend Developer Path
INSERT INTO learning_path_resources (path_id, resource_id, step_order, is_required, notes) VALUES
    ('30000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000006', 1, TRUE, 'Start with the official Tour of Go'),
    ('30000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000007', 2, TRUE, 'Go by Example for practical patterns'),
    ('30000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000005', 3, TRUE, 'Complete Go developer course'),
    ('30000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000022', 4, TRUE, 'PostgreSQL for database work'),
    ('30000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000024', 5, FALSE, 'Docker for containerization');

-- Data Engineering Path
INSERT INTO learning_path_resources (path_id, resource_id, step_order, is_required, notes) VALUES
    ('30000000-0000-0000-0000-000000000007', '20000000-0000-0000-0000-000000000001', 1, TRUE, 'Python fundamentals'),
    ('30000000-0000-0000-0000-000000000007', '20000000-0000-0000-0000-000000000021', 2, TRUE, 'SQL for data querying'),
    ('30000000-0000-0000-0000-000000000007', '20000000-0000-0000-0000-000000000022', 3, TRUE, 'PostgreSQL for advanced database work'),
    ('30000000-0000-0000-0000-000000000007', '20000000-0000-0000-0000-000000000059', 4, TRUE, 'Apache Spark with PySpark'),
    ('30000000-0000-0000-0000-000000000007', '20000000-0000-0000-0000-000000000055', 5, TRUE, 'Apache Kafka for streaming'),
    ('30000000-0000-0000-0000-000000000007', '20000000-0000-0000-0000-000000000058', 6, TRUE, 'Data Engineering Zoomcamp for end-to-end skills');

COMMIT;
