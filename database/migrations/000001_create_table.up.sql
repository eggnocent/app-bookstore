CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'member')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS roles (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    identifier VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS authors (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name VARCHAR(255) NOT NULL UNIQUE,
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS publishers (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name VARCHAR(255) NOT NULL UNIQUE,
    address TEXT,
    phone VARCHAR(15),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS categories (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS books (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    title VARCHAR(255) NOT NULL,
    author_id UUID REFERENCES authors(id) ON DELETE SET NULL,
    publisher_id UUID REFERENCES publishers(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    published_year INT,
    isbn VARCHAR(20) NOT NULL UNIQUE,
    status VARCHAR(10) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'borrowed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS book_authors (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS loans (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    loan_date DATE NOT NULL,
    return_date DATE,
    status VARCHAR(10) NOT NULL DEFAULT 'borrowed' CHECK (status IN ('borrowed', 'returned')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS comments (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    updated_at TIMESTAMPTZ,
    updated_by UUID REFERENCES users(id),
    PRIMARY KEY (id)
);