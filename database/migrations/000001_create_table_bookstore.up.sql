-- Tabel users
CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel roles
CREATE TABLE IF NOT EXISTS roles (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    identifier VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel user_roles (hubungan banyak-ke-banyak antara users dan roles)
CREATE TABLE IF NOT EXISTS user_roles (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel user_requests
CREATE TABLE IF NOT EXISTS user_requests (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    requested_role_id UUID NOT NULL REFERENCES roles(id),
    status VARCHAR(10) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel authors
CREATE TABLE IF NOT EXISTS authors (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name VARCHAR(255) NOT NULL UNIQUE,
    bio TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel publishers
CREATE TABLE IF NOT EXISTS publishers (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name VARCHAR(255) NOT NULL UNIQUE,
    address TEXT,
    phone VARCHAR(15),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel categories
CREATE TABLE IF NOT EXISTS categories (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel books
CREATE TABLE IF NOT EXISTS books (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    title VARCHAR(255) NOT NULL,
    author_id UUID REFERENCES authors(id) ON DELETE SET NULL,
    publisher_id UUID REFERENCES publishers(id) ON DELETE SET NULL,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    published_year INT,
    isbn VARCHAR(20) NOT NULL UNIQUE,
    status VARCHAR(10) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'borrowed')),
    access_level VARCHAR(20) NOT NULL DEFAULT 'public' CHECK (access_level IN ('public', 'member-only', 'admin-only')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel loans
CREATE TABLE IF NOT EXISTS loans (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    member_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    loan_date DATE NOT NULL,
    return_date DATE,
    status VARCHAR(10) NOT NULL DEFAULT 'borrowed' CHECK (status IN ('borrowed', 'returned')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel audit_logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(255) NOT NULL,
    details JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id)
);

-- Tabel ratings
CREATE TABLE IF NOT EXISTS ratings (
    id UUID NOT NULL DEFAULT GEN_RANDOM_UUID(),
    book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    review TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMPTZ,
    updated_by UUID,
    PRIMARY KEY (id),
    UNIQUE (book_id, user_id)
);
