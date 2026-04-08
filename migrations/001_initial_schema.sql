-- RateYourProduction initial schema
-- Matches MVP spec data model

CREATE TYPE work_type AS ENUM ('play', 'musical', 'opera');
CREATE TYPE creator_role AS ENUM ('playwright', 'composer', 'lyricist', 'librettist', 'book_writer');
CREATE TYPE credit_role AS ENUM ('actor', 'director', 'designer', 'choreographer', 'conductor', 'musician');
CREATE TYPE submission_status AS ENUM ('pending', 'approved', 'rejected');

-- Profiles
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT UNIQUE NOT NULL,
    display_name TEXT,
    avatar_url TEXT,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_profiles_username ON profiles (username);

-- Local auth accounts
CREATE TABLE auth_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL UNIQUE REFERENCES profiles(id) ON DELETE CASCADE,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_auth_accounts_email ON auth_accounts (email);

CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_profile ON sessions (profile_id);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);

-- People (creators, cast, crew)
CREATE TABLE people (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_people_slug ON people (slug);

-- Genres
CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT UNIQUE NOT NULL,
    slug TEXT UNIQUE NOT NULL
);

CREATE UNIQUE INDEX idx_genres_slug ON genres (slug);

-- Works
CREATE TABLE works (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    normalized_title TEXT NOT NULL,
    type work_type NOT NULL,
    description TEXT,
    premiere_year INT,
    average_rating NUMERIC(3,2) DEFAULT 0,
    rating_count INT DEFAULT 0,
    weighted_score NUMERIC(5,3) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_works_slug ON works (slug);
CREATE INDEX idx_works_type ON works (type);
CREATE INDEX idx_works_normalized_title ON works (normalized_title);
CREATE INDEX idx_works_weighted_score ON works (weighted_score DESC);
CREATE INDEX idx_works_rating_count ON works (rating_count DESC);

-- Work creators (playwright, composer, etc.)
CREATE TABLE work_creators (
    work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
    person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    role_type creator_role NOT NULL,
    PRIMARY KEY (work_id, person_id, role_type)
);

CREATE INDEX idx_work_creators_person ON work_creators (person_id);

-- Work genres
CREATE TABLE work_genres (
    work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
    genre_id UUID NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (work_id, genre_id)
);

CREATE INDEX idx_work_genres_genre ON work_genres (genre_id);

-- Companies
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    city TEXT,
    country TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_companies_slug ON companies (slug);

-- Venues
CREATE TABLE venues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    city TEXT,
    country TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_venues_slug ON venues (slug);

-- Productions
CREATE TABLE productions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
    slug TEXT UNIQUE NOT NULL,
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    venue_id UUID REFERENCES venues(id) ON DELETE SET NULL,
    city TEXT,
    country TEXT,
    start_date DATE,
    end_date DATE,
    production_label TEXT,
    average_rating NUMERIC(3,2) DEFAULT 0,
    rating_count INT DEFAULT 0,
    weighted_score NUMERIC(5,3) DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_productions_work ON productions (work_id);
CREATE INDEX idx_productions_company ON productions (company_id);
CREATE INDEX idx_productions_venue ON productions (venue_id);
CREATE INDEX idx_productions_city ON productions (city);
CREATE INDEX idx_productions_country ON productions (country);
CREATE INDEX idx_productions_weighted_score ON productions (weighted_score DESC);

-- Production credits (cast, director, etc.)
CREATE TABLE production_credits (
    production_id UUID NOT NULL REFERENCES productions(id) ON DELETE CASCADE,
    person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    role_type credit_role NOT NULL,
    PRIMARY KEY (production_id, person_id, role_type)
);

CREATE INDEX idx_production_credits_person ON production_credits (person_id);
CREATE INDEX idx_production_credits_role ON production_credits (role_type);

-- Logs (user attendance + rating + review)
CREATE TABLE logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
    production_id UUID REFERENCES productions(id) ON DELETE SET NULL,
    seen_date DATE,
    rating NUMERIC(2,1) CHECK (rating >= 0.5 AND rating <= 5.0),
    review_text TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_logs_user ON logs (user_id);
CREATE INDEX idx_logs_work ON logs (work_id);
CREATE INDEX idx_logs_production ON logs (production_id);
CREATE INDEX idx_logs_created ON logs (created_at DESC);

-- User-submitted productions
CREATE TABLE production_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
    submitted_by UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    company_id UUID REFERENCES companies(id) ON DELETE SET NULL,
    venue_id UUID REFERENCES venues(id) ON DELETE SET NULL,
    city TEXT,
    country TEXT,
    start_date DATE,
    end_date DATE,
    production_label TEXT,
    status submission_status NOT NULL DEFAULT 'pending',
    notes TEXT,
    reviewed_at TIMESTAMPTZ,
    reviewer_id UUID REFERENCES profiles(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_production_submissions_work ON production_submissions (work_id);
CREATE INDEX idx_production_submissions_submitted_by ON production_submissions (submitted_by);
CREATE INDEX idx_production_submissions_status ON production_submissions (status);
